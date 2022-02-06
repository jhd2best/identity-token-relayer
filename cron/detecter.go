package cron

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	. "identity-token-relayer/common"
	"identity-token-relayer/eth"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"math/big"
	"strings"
	"time"
)

const (
	MinimumCheckingInterval = 5
	SecurityIntervalBlock   = 6
	TransferTopic           = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

var (
	isGetEthTransaction = false
	splitTimes          = 1
	pendingTransaction  = make([]model.Transaction, 0)
)

func GetEnableProject() {
	syncErr := model.SyncAllEnableProjects()
	if syncErr != nil {
		log.GetLogger().Error("sync projects failed.", zap.String("error", syncErr.Error()))
	}
	log.GetLogger().Info("sync projects success")
}

func GetEthTransaction() {
	if isGetEthTransaction {
		return
	}
	isGetEthTransaction = true

	defer func() {
		isGetEthTransaction = false
	}()

	// get latest block height
	latestHeader, headerErr := eth.GetEthClient().HeaderByNumber(context.Background(), nil)
	if headerErr != nil {
		log.GetLogger().Error("get latest block height failed.", zap.String("error", headerErr.Error()))
		return
	}
	needCheckBlockNum := latestHeader.Number.Int64() - SecurityIntervalBlock

	minHeight := needCheckBlockNum
	projects := model.GetAllProjects()
	projectAddressSet := make([]common.Address, len(projects))

	for address, project := range projects {
		if project.LastUpdateHeight < minHeight {
			minHeight = project.LastUpdateHeight
		}
		projectAddressSet = append(projectAddressSet, common.HexToAddress(address))
	}

	// check split block range
	if splitTimes > 1 {
		newRange := (needCheckBlockNum - minHeight) / int64(splitTimes)
		needCheckBlockNum = minHeight + newRange
	}

	if needCheckBlockNum-minHeight < MinimumCheckingInterval {
		return
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(minHeight),
		ToBlock:   big.NewInt(needCheckBlockNum),
		Addresses: projectAddressSet,
		Topics: [][]common.Hash{
			{common.HexToHash(TransferTopic)},
		},
	}

	transLogs, transErr := eth.GetEthClient().FilterLogs(context.Background(), query)
	if transErr != nil {
		log.GetLogger().Error("get transaction logs failed. auto split block range", zap.String("error", transErr.Error()))
		splitTimes++
		return
	} else {
		splitTimes = 1
	}
	if len(transLogs) > 0 {
		for _, transLog := range transLogs {
			contractAddress := transLog.Address.String()

			// check project
			if _, ok := projects[contractAddress]; !ok {
				continue
			}

			// check project last update height
			if int64(transLog.BlockNumber) <= projects[contractAddress].LastUpdateHeight {
				continue
			}

			trans := model.Transaction{
				TxHash:          transLog.TxHash.String(),
				BlockHeight:     int64(transLog.BlockNumber),
				MappingTxHash:   "",
				ContractAddress: contractAddress,
				TokenId:         transLog.Topics[3].Big().Int64(),
				FromAddress:     common.HexToAddress(transLog.Topics[1].Hex()).String(),
				ToAddress:       common.HexToAddress(transLog.Topics[2].Hex()).String(),
				Status:          "created",
				CreatedAt:       time.Now().Format(TimeLayout),
			}

			createErr := model.CreateTransaction(trans)
			if createErr != nil {
				// add to pending list
				pendingTransaction = append(pendingTransaction, trans)
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetContext("data", map[string]interface{}{
						"tx_hash": trans.TxHash,
					})
					scope.SetLevel(sentry.LevelError)
					sentry.CaptureMessage("create transaction logs failed.")
				})
				log.GetLogger().Warn("create transaction logs failed.", zap.String("error", createErr.Error()))
				continue
			}
			log.GetLogger().Info("create transaction success", zap.String("tx_hash", trans.TxHash))

			time.Sleep(100 * time.Millisecond)
		}
	}

	// update project last block height
	for _, project := range projects {
		updateErr := model.UpdateProjectLastHeight(project.ContractAddress, needCheckBlockNum)
		if updateErr != nil {
			log.GetLogger().Error("update project last block height failed.", zap.String("error", updateErr.Error()))
			continue
		}
		log.GetLogger().Info("update project last block height success", zap.String("project", project.Name), zap.Int64("new_height", needCheckBlockNum))
	}
}

func HandlePendingTransaction() {
	if len(pendingTransaction) == 0 {
		return
	}

	_, createErr := model.BatchCreateTransactions(pendingTransaction)
	if createErr != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetContext("data", map[string]interface{}{
				"count": len(pendingTransaction),
			})
			scope.SetLevel(sentry.LevelError)
			sentry.CaptureMessage("batch re-create pending transaction logs failed.")
		})

		// get failed transaction hash
		failedTransactionHash := make([]string, 0)
		for _, tran := range pendingTransaction {
			failedTransactionHash = append(failedTransactionHash, tran.TxHash)
		}
		log.GetLogger().Warn("batch re-create pending transaction logs failed.",
			zap.String("error", createErr.Error()),
			zap.String("failed_tx_hash", strings.Join(failedTransactionHash, " / ")),
		)
		return
	}
	log.GetLogger().Info("batch re-create pending transaction logs success.")
}
