package cron

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	. "identity-token-relayer/common"
	"identity-token-relayer/hmy"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"math/big"
	"strings"
	"time"
)

var (
	isSendMappingTransaction  = false
	isCheckMappingTransaction = false
	transOpt                  *bind.TransactOpts
)

func SendMappingTransaction() {
	if isSendMappingTransaction {
		return
	}
	isSendMappingTransaction = true

	defer func() {
		isSendMappingTransaction = false
	}()

	// get created transactions
	trans, getErr := model.GetTransactionByStatus("created")
	if getErr != nil {
		log.GetLogger().Error("get created transactions failed.", zap.String("error", getErr.Error()))
		return
	}

	if len(trans) > 0 {
		for _, tran := range trans {
			// check nft last block height
			nft, nftErr := model.GetProjectNftByTokenId(tran.ContractAddress, tran.TokenId)
			if nftErr != nil {
				log.GetLogger().Error("check nft last block height failed.", zap.String("error", nftErr.Error()))
				continue
			}

			if nft.LastUpdateHeight >= tran.BlockHeight {
				skipErr := model.SetTransactionStatus(tran.TxHash, "skipped")
				if skipErr != nil {
					log.GetLogger().Error("skip trans failed.", zap.String("error", skipErr.Error()))
				}
				continue
			}

			// exec transaction on harmony
			mappingTxHash, execErr := execOwnerUpdateOnHarmony(tran)
			if execErr != nil {
				log.GetLogger().Error("exec transaction on harmony failed.", zap.String("error", execErr.Error()))
				continue
			}

			// update transaction status
			setErr := model.SetTransactionStatusMapping(tran.TxHash, mappingTxHash)
			if setErr != nil {
				log.GetLogger().Error("update transaction status to mapping failed.", zap.String("error", setErr.Error()))
				continue
			}

			// update nft new owner
			updateNftErr := model.UpdateProjectNftOwner(tran.ContractAddress, tran.TokenId, tran.ToAddress, tran.BlockHeight)
			if updateNftErr != nil {
				log.GetLogger().Error("update project nft owner failed.", zap.String("error", updateNftErr.Error()))
			}

			log.GetLogger().Info("update transaction status to mapping success", zap.String("tx_hash", tran.TxHash))
		}
	}
}

func CheckMappingTransaction() {
	if isCheckMappingTransaction {
		return
	}
	isCheckMappingTransaction = true

	defer func() {
		isCheckMappingTransaction = false
	}()

	// get mapping transactions
	trans, getErr := model.GetTransactionByStatus("mapping")
	if getErr != nil {
		log.GetLogger().Error("get mapping transactions failed.", zap.String("error", getErr.Error()))
		return
	}

	if len(trans) > 0 {
		for _, tran := range trans {
			// check status after 1 min
			tranUpdatedTime, _ := time.Parse(TimeLayout, tran.UpdatedAt)
			println(fmt.Sprintf("space: %v\n", time.Now().Sub(tranUpdatedTime)))

			if time.Now().Sub(tranUpdatedTime) > time.Minute {
				// get mapping trans receipt
				receipt, recErr := hmy.GetHmyClient().TransactionReceipt(common.HexToHash(tran.MappingTxHash))
				if recErr != nil {
					log.GetLogger().Error("get mapping trans receipt failed.", zap.String("error", recErr.Error()), zap.String("mapping_hash", tran.MappingTxHash))
					continue
				}

				// set trans final status
				if receipt.Status == 1 {
					setErr := model.SetTransactionStatus(tran.TxHash, "success")
					if setErr != nil {
						log.GetLogger().Error("set trans final status failed.", zap.String("error", setErr.Error()), zap.String("hash", tran.TxHash))
						continue
					}
				} else {
					setErr := model.SetTransactionStatus(tran.TxHash, "failed")
					if setErr != nil {
						log.GetLogger().Error("set trans final status failed.", zap.String("error", setErr.Error()), zap.String("hash", tran.TxHash))
						continue
					}
					log.GetLogger().Error("found failed trans.", zap.String("hash", tran.TxHash))
				}
			}
		}
	}
}

func execOwnerUpdateOnHarmony(tran model.Transaction) (hash string, err error) {
	if transOpt == nil {
		privateKey, err := crypto.HexToECDSA(hmy.GetHmyPrivateKey())
		if err != nil {
			return "", err
		}

		gasPrice, err := hmy.GetHmyClient().SuggestGasPrice(context.Background())
		if err != nil {
			return "", err
		}

		chainId, err := hmy.GetHmyClient().NetworkID(context.Background())
		if err != nil {
			return "", err
		}

		transOpt, err = bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			return "", err
		}

		transOpt.Value = big.NewInt(0)
		transOpt.GasLimit = uint64(1000000)
		transOpt.GasPrice = gasPrice
	}

	nonce, err := hmy.GetHmyClient().Nonce(context.Background(), transOpt.From)
	if err != nil {
		return "", err
	}
	transOpt.Nonce = big.NewInt(int64(nonce))

	// convert params
	ethContractAddress := common.HexToAddress(tran.ContractAddress).String()
	ethContractAddress = strings.Replace(ethContractAddress, "0x", "", 1)
	newOwnerAddress := common.HexToAddress(tran.ToAddress)
	tokenIdBig := big.NewInt(tran.TokenId)
	blockHeightBig := big.NewInt(tran.BlockHeight)

	execTrans, execErr := hmy.GetOwnershipValidatorClient().UpdateOwnership(transOpt, ethContractAddress, []common.Address{newOwnerAddress}, []*big.Int{tokenIdBig}, blockHeightBig)
	if execErr != nil {
		return "", execErr
	}

	return execTrans.Hash().Hex(), nil
}
