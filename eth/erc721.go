package eth

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"identity-token-relayer/contract/token721"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"math/big"
	"time"
)

type NftArrayItem struct {
	Id    int64
	Owner string
}

type NftArrayItemWithHeight struct {
	Id     int64
	Owner  string
	Height *big.Int
}

func GetAllErc721OwnersInOrderOnChain(addressRaw string, syncBlockHeight int64, isLatest bool, indexOffset int64) (allOwners []NftArrayItemWithHeight, err error) {
	address := common.HexToAddress(addressRaw)
	blockHeight := big.NewInt(syncBlockHeight)

	token721Client, err := token721.NewToken721(address, GetEthClient().EthClient)
	if err != nil {
		return nil, err
	}

	callOpts := &bind.CallOpts{}
	if !isLatest {
		callOpts.BlockNumber = blockHeight
	}

	// get total supply
	supplyRaw, err := token721Client.TotalSupply(callOpts)
	if err != nil {
		return nil, err
	}
	supply := int(supplyRaw.Int64())

	latestHeader, err := GetEthClient().EthClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	latestHeight := latestHeader.Number

	// get nft owners
	nftOwnerSet := make([]NftArrayItemWithHeight, 0)
	for i := 0 + indexOffset; i < int64(supply)+indexOffset; i++ {
		owner, err := token721Client.OwnerOf(callOpts, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}

		itemHeight := latestHeight
		if !isLatest {
			itemHeight = big.NewInt(syncBlockHeight)
		}

		nftOwnerSet = append(nftOwnerSet, NftArrayItemWithHeight{
			Id:     i,
			Owner:  owner.String(),
			Height: itemHeight,
		})

		fmt.Printf("get owner %s of %d\n", owner.String(), i)
		time.Sleep(100 * time.Millisecond)
	}

	if supply != len(nftOwnerSet) {
		return nil, errors.New("total supply not correct")
	}

	return nftOwnerSet, nil
}

func SyncOneErc721TokenOnChain(addressRaw string, tokenId int64, syncBlockHeight int64, isLatest bool) error {
	address := common.HexToAddress(addressRaw)
	blockHeight := big.NewInt(syncBlockHeight)

	token721Client, err := token721.NewToken721(address, GetEthClient().EthClient)
	if err != nil {
		return err
	}

	callOpts := &bind.CallOpts{}
	if !isLatest {
		callOpts.BlockNumber = blockHeight
	}

	// get nft owners
	nftOwnerSet := make([]NftArrayItem, 0)
	owner, err := token721Client.OwnerOf(callOpts, big.NewInt(tokenId))
	if err != nil {
		return err
	}
	nftOwnerSet = append(nftOwnerSet, NftArrayItem{
		Id:    tokenId,
		Owner: owner.String(),
	})

	// update nft owners
	syncErr := SyncAllErc721TokenFromArray(addressRaw, syncBlockHeight, nftOwnerSet)
	if syncErr != nil {
		return syncErr
	}

	return nil
}

func SyncAllErc721TokenOnChain(addressRaw string, syncBlockHeight int64, isLatest bool) error {
	address := common.HexToAddress(addressRaw)
	blockHeight := big.NewInt(syncBlockHeight)

	token721Client, err := token721.NewToken721(address, GetEthClient().EthClient)
	if err != nil {
		return err
	}

	callOpts := &bind.CallOpts{}
	if !isLatest {
		callOpts.BlockNumber = blockHeight
	}

	// get total supply
	supplyRaw, err := token721Client.TotalSupply(callOpts)
	if err != nil {
		return err
	}
	supply := int(supplyRaw.Int64())

	// get nft owners
	nftOwnerSet := make([]NftArrayItem, 0)
	for i := 0; i < supply; i++ {
		owner, err := token721Client.OwnerOf(callOpts, big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		nftOwnerSet = append(nftOwnerSet, NftArrayItem{
			Id:    int64(i),
			Owner: owner.String(),
		})
		log.GetLogger().Info("found token and synced.", zap.Int("token_id", i), zap.String("owner", owner.String()))
		time.Sleep(100 * time.Millisecond)
	}

	// update nft owners
	syncErr := SyncAllErc721TokenFromArray(addressRaw, syncBlockHeight, nftOwnerSet)
	if syncErr != nil {
		return syncErr
	}

	return nil
}

func SyncAllErc721TokenFromArray(addressRaw string, syncBlockHeight int64, nfts []NftArrayItem) error {
	if len(nfts) == 0 {
		return errors.New("nft not found")
	}

	pendingNftSet := make([]model.Nft, 0)
	for _, nft := range nfts {
		pendingNftSet = append(pendingNftSet, model.Nft{
			ContractAddress:  addressRaw,
			TokenId:          nft.Id,
			OwnerAddress:     nft.Owner,
			LastUpdateHeight: syncBlockHeight,
		})

		if len(pendingNftSet) >= 499 {
			_, batchErr := model.BatchCreateNfts(pendingNftSet)
			if batchErr != nil {
				return batchErr
			} else {
				log.GetLogger().Info(fmt.Sprintf("created NFT from %d to %d", pendingNftSet[0].TokenId, pendingNftSet[498].TokenId))
				pendingNftSet = make([]model.Nft, 0)
				time.Sleep(1 * time.Second)
			}
		}
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
