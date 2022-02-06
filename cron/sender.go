package cron

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	. "identity-token-relayer/common"
	"identity-token-relayer/eth"
	"identity-token-relayer/hmy"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"math/big"
	"strings"
	"time"
)

var (
	isSendMappingTransaction  = false
	isRetryErrorTransaction   = false
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
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"contract_address": tran.ContractAddress,
						"token_id":         tran.TokenId,
					})
					scope.SetLevel(sentry.LevelWarning)
					sentry.CaptureMessage("nft not found. try to auto-create it.")
				})

				log.GetLogger().Warn("nft not found. try to auto-create it.",
					zap.String("error", nftErr.Error()),
				)

				// sync nft from chain
				syncErr := eth.SyncOneErc721TokenOnChain(tran.ContractAddress, tran.TokenId, tran.BlockHeight, true)
				if syncErr != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetContext("data", map[string]interface{}{
							"contract_address": tran.ContractAddress,
							"token_id":         tran.TokenId,
						})
						scope.SetLevel(sentry.LevelError)
						sentry.CaptureMessage("sync nft from chain failed.")
					})

					log.GetLogger().Warn("sync nft from chain failed.",
						zap.String("error", syncErr.Error()),
					)
					continue
				}

				// skip the transaction
				skipErr := model.SetTransactionStatus(tran.TxHash, tran.ContractAddress, tran.TokenId, "skipped")
				if skipErr != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetContext("data", map[string]interface{}{
							"tx_hash": tran.TxHash,
						})
						scope.SetLevel(sentry.LevelError)
						sentry.CaptureMessage("skip trans failed.")
					})

					log.GetLogger().Warn("skip trans failed.",
						zap.String("error", skipErr.Error()),
					)
				}
				continue
			}

			if nft.LastUpdateHeight >= tran.BlockHeight {
				skipErr := model.SetTransactionStatus(tran.TxHash, tran.ContractAddress, tran.TokenId, "skipped")
				if skipErr != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetContext("data", map[string]interface{}{
							"tx_hash": tran.TxHash,
						})
						scope.SetLevel(sentry.LevelError)
						sentry.CaptureMessage("skip trans failed.")
					})

					log.GetLogger().Warn("skip trans failed.",
						zap.String("error", skipErr.Error()),
					)
				}
				continue
			}

			// exec transaction on harmony
			mappingTxHash, execErr := execOwnerUpdateOnHarmony(tran)
			if execErr != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"tx_hash":       tran.TxHash,
						"harmony_error": execErr.Error(),
					})
					scope.SetLevel(sentry.LevelError)
					sentry.CaptureMessage("exec transaction on harmony failed.")
				})

				log.GetLogger().Warn("exec transaction on harmony failed.",
					zap.String("error", execErr.Error()),
				)
				continue
			}

			// update transaction status
			setErr := model.SetTransactionStatusMapping(tran.TxHash, tran.ContractAddress, tran.TokenId, mappingTxHash)
			if setErr != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"tx_hash": tran.TxHash,
					})
					scope.SetLevel(sentry.LevelError)
					sentry.CaptureMessage("update transaction status to mapping failed.")
				})

				log.GetLogger().Warn("update transaction status to mapping failed.",
					zap.String("error", setErr.Error()),
				)
				continue
			}

			log.GetLogger().Info("update transaction status to mapping success", zap.String("tx_hash", tran.TxHash))
		}
	}
}

func RetryErrorTransaction() {
	if isRetryErrorTransaction {
		return
	}
	isRetryErrorTransaction = true

	defer func() {
		isRetryErrorTransaction = false
	}()

	// get created transactions
	trans, getErr := model.GetTransactionByStatus("error")
	if getErr != nil {
		log.GetLogger().Error("get error transactions failed.", zap.String("error", getErr.Error()))
		return
	}

	if len(trans) > 0 {
		for _, tran := range trans {
			// check nft last block height
			nft, nftErr := model.GetProjectNftByTokenId(tran.ContractAddress, tran.TokenId)
			if nftErr != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"contract_address": tran.ContractAddress,
						"token_id":         tran.TokenId,
					})
					scope.SetLevel(sentry.LevelWarning)
					sentry.CaptureMessage("nft not found.")
				})

				log.GetLogger().Warn("nft not found.",
					zap.String("error", nftErr.Error()),
				)
				continue
			}

			if nft.LastUpdateHeight >= tran.BlockHeight {
				skipErr := model.SetTransactionStatus(tran.TxHash, tran.ContractAddress, tran.TokenId, "skipped")
				if skipErr != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetContext("data", map[string]interface{}{
							"tx_hash": tran.TxHash,
						})
						scope.SetLevel(sentry.LevelError)
						sentry.CaptureMessage("skip trans failed.")
					})

					log.GetLogger().Warn("skip trans failed.",
						zap.String("error", skipErr.Error()),
					)
				}
				continue
			}

			// exec transaction on harmony
			mappingTxHash, execErr := execOwnerUpdateOnHarmony(tran)
			if execErr != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"tx_hash":       tran.TxHash,
						"harmony_error": execErr.Error(),
					})
					scope.SetLevel(sentry.LevelError)
					sentry.CaptureMessage("exec transaction on harmony failed.")
				})

				log.GetLogger().Warn("exec transaction on harmony failed.",
					zap.String("error", execErr.Error()),
				)
				continue
			}

			// update transaction status
			setErr := model.SetTransactionStatusMapping(tran.TxHash, tran.ContractAddress, tran.TokenId, mappingTxHash)
			if setErr != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"tx_hash": tran.TxHash,
					})
					scope.SetLevel(sentry.LevelError)
					sentry.CaptureMessage("update transaction status to mapping failed.")
				})

				log.GetLogger().Warn("update transaction status to mapping failed.",
					zap.String("error", setErr.Error()),
				)
				continue
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

			if time.Now().Sub(tranUpdatedTime) > time.Minute {
				// get mapping trans receipt
				receipt, recErr := hmy.GetHmyClient().TransactionReceipt(common.HexToHash(tran.MappingTxHash))
				if recErr != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetContext("data", map[string]interface{}{
							"mapping_hash": tran.MappingTxHash,
						})
						scope.SetLevel(sentry.LevelError)
						sentry.CaptureMessage("get mapping trans receipt failed.")
					})

					log.GetLogger().Warn("get mapping trans receipt failed.", zap.String("error", recErr.Error()))
					continue
				}

				// set trans final status
				if receipt.Status == 1 {
					setErr := model.SetTransactionStatus(tran.TxHash, tran.ContractAddress, tran.TokenId, "success")
					if setErr != nil {
						sentry.WithScope(func(scope *sentry.Scope) {
							scope.SetContext("data", map[string]interface{}{
								"tx_hash": tran.TxHash,
							})
							scope.SetLevel(sentry.LevelError)
							sentry.CaptureMessage("set trans final status failed.")
						})

						log.GetLogger().Warn("set trans final status failed.", zap.String("error", setErr.Error()))
						continue
					}

					// update nft new owner
					updateNftErr := model.UpdateProjectNftOwner(tran.ContractAddress, tran.TokenId, tran.ToAddress, tran.BlockHeight)
					if updateNftErr != nil {
						sentry.WithScope(func(scope *sentry.Scope) {
							scope.SetContext("data", map[string]interface{}{
								"tx_hash": tran.TxHash,
							})
							scope.SetLevel(sentry.LevelError)
							sentry.CaptureMessage("update project nft owner failed.")
						})

						log.GetLogger().Warn("update project nft owner failed.",
							zap.String("error", updateNftErr.Error()),
						)
					}
				} else {
					if tran.RetryTimes <= 3 {
						tran.RetryTimes++
						setErr := model.SetTransactionStatusError(tran.TxHash, tran.ContractAddress, tran.TokenId, tran.RetryTimes)
						if setErr != nil {
							sentry.WithScope(func(scope *sentry.Scope) {
								scope.SetContext("data", map[string]interface{}{
									"tx_hash": tran.TxHash,
								})
								scope.SetLevel(sentry.LevelError)
								sentry.CaptureMessage("set trans error status failed.")
							})

							log.GetLogger().Warn("set trans error status failed.", zap.String("error", setErr.Error()))
							continue
						}
					} else {
						setErr := model.SetTransactionStatus(tran.TxHash, tran.ContractAddress, tran.TokenId, "failed")
						if setErr != nil {
							sentry.WithScope(func(scope *sentry.Scope) {
								scope.SetContext("data", map[string]interface{}{
									"tx_hash": tran.TxHash,
								})
								scope.SetLevel(sentry.LevelError)
								sentry.CaptureMessage("set trans final status failed.")
							})

							log.GetLogger().Warn("set trans final status failed.", zap.String("error", setErr.Error()))
							continue
						}
					}

					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetContext("data", map[string]interface{}{
							"tx_hash": tran.TxHash,
							"retry_times": tran.RetryTimes,
						})
						scope.SetLevel(sentry.LevelWarning)
						sentry.CaptureMessage("found error trans.")
					})
				}
				log.GetLogger().Info("set trans final status success.", zap.String("hash", tran.TxHash))
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

	execTrans, execErr := hmy.GetOwnershipValidatorClient().UpdateOwnership(transOpt, ethContractAddress, []common.Address{newOwnerAddress}, []*big.Int{tokenIdBig})
	if execErr != nil {
		return "", execErr
	}

	return execTrans.Hash().Hex(), nil
}
