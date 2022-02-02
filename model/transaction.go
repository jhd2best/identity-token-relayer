package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
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
	Status          string `firestore:"status"` // created, mapping, success, failed, skipped
	CreatedAt       string `firestore:"created_at"`
	UpdatedAt       string `firestore:"updated_at"`
}

func GetTransactionByHash(hash string) (trans Transaction, err error) {
	data, err := GetDbClient().Collection("transactions").Doc(hash).Get(context.Background())
	if err != nil {
		return
	}

	_ = data.DataTo(&trans)
	return
}

func GetTransactionByStatus(status string) ([]Transaction, error) {
	transSet := make([]Transaction, 0)
	iter := GetDbClient().Collection("transactions").Where("status", "==", status).Documents(context.Background())
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
	ref := GetDbClient().Collection("transactions").Doc(trans.TxHash)
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
		docRef := GetDbClient().Collection("transactions").Doc(tran.TxHash)
		batch.Set(docRef, tran)
	}

	_, commitErr := batch.Commit(context.Background())
	if commitErr != nil {
		return nil, commitErr
	}
	return addedTxHash, nil
}

func SetTransactionStatusMapping(txHash string, mappingTxHash string) error {
	_, updateErr := GetDbClient().Collection("transactions").Doc(txHash).Update(context.Background(), []firestore.Update{
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

func SetTransactionStatus(txHash string, status string) error {
	_, updateErr := GetDbClient().Collection("transactions").Doc(txHash).Update(context.Background(), []firestore.Update{
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
