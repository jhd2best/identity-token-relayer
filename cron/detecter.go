package cron

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"identity-token-relayer/eth"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"math/big"
)

const (
	MinimumCheckingInterval = 5
	SecurityIntervalBlock   = 6
	TransferTopic           = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

var (
	isGetEthTransaction = false
)

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
		log.GetLogger().Error("get transaction logs failed.", zap.String("error", transErr.Error()))
		return
	}

	if len(transLogs) > 0 {
		for index, transLog := range transLogs {
			// TODO
			if index >= 10 {
				break
			}

			if int64(transLog.BlockNumber) <= projects[transLog.Address.String()].LastUpdateHeight {
				continue
			}

			trans := model.Transaction{
				TxHash:          transLog.TxHash.String(),
				BlockHeight:     int64(transLog.BlockNumber),
				MappingTxHash:   "",
				ContractAddress: transLog.Address.String(),
				TokenId:         transLog.Topics[3].Big().Int64(),
				FromAddress:     common.HexToAddress(transLog.Topics[1].Hex()).String(),
				ToAddress:       common.HexToAddress(transLog.Topics[2].Hex()).String(),
			}

			createErr := model.CreateTransaction(trans)
			if createErr != nil {
				log.GetLogger().Error("create transaction logs failed.", zap.String("error", createErr.Error()), zap.String("tx_hash", trans.TxHash))
				continue
			}
			log.GetLogger().Info("create transaction success", zap.String("tx_hash", trans.TxHash))
		}

		// TODO: update last block height
	}

}
