package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	. "identity-token-relayer/common"
	"identity-token-relayer/log"
	"time"
)

type Transaction struct {
	TxHash          string `firestore:"tx_hash"`
	BlockHeight     int64  `firestore:"block_height"`
	MappingTxHash   string `firestore:"mapping_tx_hash"`
	ContractAddress string `firestore:"contract_address"`
	TokenId         int64  `firestore:"token_id"`
	FromAddress     string `firestore:"from_address"`
	ToAddress       string `firestore:"to_address"`
	Status          string `firestore:"status"` // created, mapping, error, success, failed, skipped
	RetryTimes      int64  `firestore:"retry_times"`
	CreatedAt       string `firestore:"created_at"`
	UpdatedAt       string `firestore:"updated_at"`
}

func GetOneTransaction(hash string, address string, tokenId int64) (trans Transaction, err error) {
	docName := fmt.Sprintf("%s-%s-%d", hash, address[34:41], tokenId)
	data, err := GetDbClient().Collection("transactions").Doc(docName).Get(context.Background())
	if err != nil {
		return
	}

	_ = data.DataTo(&trans)
	return
}

func GetTransactionByStatus(status string, limit int) ([]Transaction, error) {
	transSet := make([]Transaction, 0)
	iter := GetDbClient().Collection("transactions").Where("status", "==", status).Limit(limit).Documents(context.Background())
	for {
		if doc, iterErr := iter.Next(); iterErr != nil {
			if iterErr == iterator.Done {
				break
			} else {
				log.GetLogger().Error("get transactions failed.", zap.String("error", iterErr.Error()))
				return nil, iterErr
			}
		} else {
			trans := Transaction{}
			dataErr := doc.DataTo(&trans)
			if dataErr != nil {
				log.GetLogger().Error("parse transactions data failed.", zap.String("error", dataErr.Error()))
				continue
			}
			transSet = append(transSet, trans)
		}
	}
	return transSet, nil
}

func CreateTransaction(trans Transaction) error {
	docName := fmt.Sprintf("%s-%s-%d", trans.TxHash, trans.ContractAddress[34:41], trans.TokenId)
	ref := GetDbClient().Collection("transactions").Doc(docName)
	err := GetDbClient().RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		_, getErr := tx.Get(ref)
		if getErr == nil {
			return errors.New("transaction already exist")
		}

		return tx.Set(ref, trans)
	})
	return err
}

func BatchCreateTransactions(trans []Transaction) ([]string, error) {
	if len(trans) > 500 {
		return nil, errors.New("batch create can not more 500 items")
	}

	addedTxHash := make([]string, 0)
	batch := GetDbClient().Batch()

	for _, tran := range trans {
		addedTxHash = append(addedTxHash, tran.TxHash)
		docName := fmt.Sprintf("%s-%s-%d", tran.TxHash, tran.ContractAddress[34:41], tran.TokenId)
		docRef := GetDbClient().Collection("transactions").Doc(docName)
		batch.Set(docRef, tran)
	}

	_, commitErr := batch.Commit(context.Background())
	if commitErr != nil {
		return nil, commitErr
	}
	return addedTxHash, nil
}

func BatchUpdateTransactions(trans []Transaction, data []firestore.Update) ([]string, error) {
	if len(trans) > 500 {
		return nil, errors.New("batch update can not more 500 items")
	}

	updatedTransactionTxHash := make([]string, 0)
	batch := GetDbClient().Batch()

	for _, tran := range trans {
		updatedTransactionTxHash = append(updatedTransactionTxHash, tran.TxHash)
		docName := fmt.Sprintf("%s-%s-%d", tran.TxHash, tran.ContractAddress[34:41], tran.TokenId)
		sfRef := GetDbClient().Collection("transactions").Doc(docName)
		batch.Update(sfRef, data)
	}

	_, commitErr := batch.Commit(context.Background())
	if commitErr != nil {
		return nil, commitErr
	}
	return updatedTransactionTxHash, nil
}

func SetTransactionStatusMapping(txHash string, address string, tokenId int64, mappingTxHash string) error {
	docName := fmt.Sprintf("%s-%s-%d", txHash, address[34:41], tokenId)
	_, updateErr := GetDbClient().Collection("transactions").Doc(docName).Update(context.Background(), []firestore.Update{
		{
			Path:  "status",
			Value: "mapping",
		}, {
			Path:  "mapping_tx_hash",
			Value: mappingTxHash,
		}, {
			Path:  "updated_at",
			Value: time.Now().Format(TimeLayout),
		},
	})
	return updateErr
}

func SetTransactionStatusError(txHash string, address string, tokenId int64, retryTimes int64) error {
	docName := fmt.Sprintf("%s-%s-%d", txHash, address[34:41], tokenId)
	_, updateErr := GetDbClient().Collection("transactions").Doc(docName).Update(context.Background(), []firestore.Update{
		{
			Path:  "status",
			Value: "error",
		}, {
			Path:  "retry_times",
			Value: retryTimes,
		}, {
			Path:  "updated_at",
			Value: time.Now().Format(TimeLayout),
		},
	})
	return updateErr
}

func SetTransactionStatus(txHash string, address string, tokenId int64, status string) error {
	docName := fmt.Sprintf("%s-%s-%d", txHash, address[34:41], tokenId)
	_, updateErr := GetDbClient().Collection("transactions").Doc(docName).Update(context.Background(), []firestore.Update{
		{
			Path:  "status",
			Value: status,
		}, {
			Path:  "updated_at",
			Value: time.Now().Format(TimeLayout),
		},
	})
	return updateErr
}
