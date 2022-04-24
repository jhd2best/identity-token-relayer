package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	"identity-token-relayer/config"
	"identity-token-relayer/hmy"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	syncLimit = 500
	initLimit = 100
)

var (
	transOpt *bind.TransactOpts
	deepMode = false
)

func AutoImportFromChain(origin string, mapping string, name string, symbol string, baseUrl string, mode string) {
	// check params
	if !common.IsHexAddress(origin) || !common.IsHexAddress(mapping) {
		log.GetLogger().Fatal("origin or mapping address not valid.")
	}

	if name == "" || symbol == "" || baseUrl == "" {
		log.GetLogger().Fatal("origin project info not valid.")
	}

	originAddress := common.HexToAddress(origin).String()
	mappingAddress := common.HexToAddress(mapping).String()

	if mode == "deep" {
		deepMode = true
	}

	// get all current owners
	syncOffset := 0
	syncCursor := ""
	maxBlockHeight := 0
	totalSupply := 0
	httpClient := http.Client{}
	allOwners := make([]OwnerItem, 0)
	ownerChecker := make(map[string]bool)

	for {
		requestUrl := fmt.Sprintf("https://deep-index.moralis.io/api/v2/nft/%s/owners?limit=%d&cursor=%s", originAddress, syncLimit, syncCursor)
		req, _ := http.NewRequest(http.MethodGet, requestUrl, nil)
		req.Header.Add("accept", "application/json")
		req.Header.Add("X-API-Key", config.Get().Eth.MoralisApiKey)

		getRes, getErr := httpClient.Do(req)
		if getErr != nil {
			log.GetLogger().Error(fmt.Sprintf("get current owners failed. will try again. offset: %d", syncOffset))
			time.Sleep(3 * time.Second)
			continue
		}
		resRaw, _ := ioutil.ReadAll(getRes.Body)

		// parse the response
		res := new(NftOwnerResponse)
		unmarshalErr := json.Unmarshal(resRaw, res)
		if unmarshalErr != nil {
			log.GetLogger().Fatal(fmt.Sprintf("parse the response failed. stoped. offset: %d", syncOffset))
		}

		if res.Status != "SYNCED" {
			log.GetLogger().Fatal(fmt.Sprintf("owner data not complete syncing. stoped. offset: %d page: %d", syncOffset, res.Page))
		}

		if totalSupply == 0 {
			totalSupply = res.Total
		}

		if len(res.Result) > 0 {
			for _, item := range res.Result {
				itemBlockHeight, _ := strconv.Atoi(item.BlockNumber)
				// update max block height
				if maxBlockHeight < itemBlockHeight {
					maxBlockHeight = itemBlockHeight
				}

				if _, ok := ownerChecker[item.TokenId]; !ok {
					allOwners = append(allOwners, item)
					ownerChecker[item.TokenId] = true
				} else {
					log.GetLogger().Error(fmt.Sprintf("found duplicate owner item. skipped. tokenId: %s", item.TokenId))
				}
			}
			log.GetLogger().Info(fmt.Sprintf("get owner data success. %d/%d", syncOffset+len(res.Result), totalSupply))
		} else {
			break
		}

		_ = getRes.Body.Close()

		syncOffset += len(res.Result)
		syncCursor = res.Cursor
		time.Sleep(2 * time.Second)
	}

	// check items count
	if totalSupply != len(allOwners) {
		log.GetLogger().Fatal(fmt.Sprintf("check items count failed. stoped. total: %d allCount: %d", totalSupply, len(allOwners)))
	}

	// check project
	_, getErr := model.GetDbClient().Collection("projects").Doc(originAddress).Get(context.Background())
	if getErr != nil {
		if strings.Index(getErr.Error(), "NotFound") >= 0 {
			log.GetLogger().Info("project non-found. will auto create it")

			// create project
			project := model.Project{
				ContractAddress:  originAddress,
				MappingAddress:   mappingAddress,
				Name:             name,
				Symbol:           symbol,
				BaseURI:          baseUrl,
				Status:           "created",
				LastUpdateHeight: int64(maxBlockHeight),
				Enable:           false,
			}

			if !deepMode {
				_, createErr := model.GetDbClient().Collection("projects").Doc(project.ContractAddress).Set(context.Background(), project)
				if createErr != nil {
					log.GetLogger().Error("create project failed.", zap.String("error", createErr.Error()))
				}

				log.GetLogger().Info("create project success", zap.String("name", project.Name))
			} else {
				log.GetLogger().Info("skip create project in deep mode")
			}
		} else {
			panic(getErr)
		}
	}

	// init NFT items
	syncedItem := 0
	for {
		var subOwnerSet []OwnerItem
		if initLimit < len(allOwners) {
			subOwnerSet = allOwners[0:initLimit]
			allOwners = allOwners[initLimit:]
		} else {
			subOwnerSet = allOwners[:]
			allOwners = nil
		}

		// analysis params
		ethAddress := strings.Replace(originAddress, "0x", "", 1)
		newOwnerSet := make([]common.Address, 0)
		tokenIdSet := make([]*big.Int, 0)
		for _, owner := range subOwnerSet {
			newTokenId, _ := strconv.Atoi(owner.TokenId)
			if deepMode {
				for {
					_, callErr := hmy.GetOwnershipValidatorClient().OwnerOf(new(bind.CallOpts), ethAddress, big.NewInt(int64(newTokenId)))
					if callErr != nil {
						if strings.Index(callErr.Error(), "nonexistent") >= 0 {
							newOwnerSet = append(newOwnerSet, common.HexToAddress(owner.OwnerOf))
							tokenIdSet = append(tokenIdSet, big.NewInt(int64(newTokenId)))
							log.GetLogger().Info("found nonexistent token id.", zap.String("tokenId", owner.TokenId))
							break
						} else {
							log.GetLogger().Error("get owner failed. try again.", zap.String("error", callErr.Error()))
						}
					} else {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			} else {
				newOwnerSet = append(newOwnerSet, common.HexToAddress(owner.OwnerOf))
				tokenIdSet = append(tokenIdSet, big.NewInt(int64(newTokenId)))
			}
		}

		// init mapping contract on Harmony
		if len(newOwnerSet) > 0 {
			for {
				opt, transOptErr := getTransOpt()
				if transOptErr == nil {
					execTrans, execErr := hmy.GetOwnershipValidatorClient().Initialize(opt, ethAddress, newOwnerSet, tokenIdSet)
					if execErr != nil {
						log.GetLogger().Error("submit init trans failed.", zap.String("error", execErr.Error()))
					} else {
						log.GetLogger().Info("submit init trans success.", zap.String("hash", execTrans.Hash().Hex()))
						break
					}
				} else {
					log.GetLogger().Error("get trans opt failed.", zap.String("error", transOptErr.Error()))
				}
				time.Sleep(5 * time.Second)
			}

			// update database
			syncErr := syncDatabaseFromArray(originAddress, subOwnerSet)
			if syncErr != nil {
				log.GetLogger().Error("update database failed.", zap.String("error", syncErr.Error()))
			}
		}

		syncedItem += len(subOwnerSet)
		log.GetLogger().Info(fmt.Sprintf("init NFT items success. %d/%d", syncedItem, totalSupply))

		if allOwners == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func getTransOpt() (*bind.TransactOpts, error) {
	if transOpt == nil {
		privateKey, err := crypto.HexToECDSA(hmy.GetHmyPrivateKey())
		if err != nil {
			return nil, err
		}

		chainId, err := hmy.GetHmyClient().NetworkID(context.Background())
		if err != nil {
			return nil, err
		}

		transOpt, err = bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			return nil, err
		}

		transOpt.Value = big.NewInt(0)
		transOpt.GasLimit = uint64(20000000)
	}

	gasPrice, err := hmy.GetHmyClient().SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	if gasPrice.Int64() < 40000000000 {
		gasPrice = big.NewInt(40000000000)
	}
	transOpt.GasPrice = gasPrice

	nonce, err := hmy.GetHmyClient().Nonce(context.Background(), transOpt.From)
	if err != nil {
		return nil, err
	}
	transOpt.Nonce = big.NewInt(int64(nonce))

	return transOpt, nil
}

func syncDatabaseFromArray(addressRaw string, nfts []OwnerItem) error {
	if len(nfts) == 0 {
		return errors.New("nft not found")
	}

	pendingNftSet := make([]model.Nft, 0)
	for _, nft := range nfts {
		newTokenId, _ := strconv.Atoi(nft.TokenId)
		syncBlockHeight, _ := strconv.Atoi(nft.BlockNumber)
		pendingNftSet = append(pendingNftSet, model.Nft{
			ContractAddress:  addressRaw,
			TokenId:          int64(newTokenId),
			OwnerAddress:     common.HexToAddress(nft.OwnerOf).String(),
			LastUpdateHeight: int64(syncBlockHeight),
		})
	}

	if len(pendingNftSet) > 0 {
		_, batchErr := model.BatchCreateNfts(pendingNftSet)
		if batchErr != nil {
			return batchErr
		}
		log.GetLogger().Info(fmt.Sprintf("created NFT from %d to %d", pendingNftSet[0].TokenId, pendingNftSet[len(pendingNftSet)-1].TokenId))
	}
	log.GetLogger().Info("batch create NFT success.")

	return nil
}
