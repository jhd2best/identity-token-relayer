package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
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
	Status          string `firestore:"status"` // created, mapping, success, failed
	CreatedAt       string `firestore:"created_at"`
	UpdatedAt       string `firestore:"updated_at"`
}

func GetTransactionByHash(hash string) (trans Transaction, err error) {
	data, err := getDbClient().Collection("transactions").Doc(hash).Get(context.Background())
	if err != nil {
		return
	}

	_ = data.DataTo(&trans)
	return
}

func CreateTransaction(trans Transaction) error {
	trans.Status = "created"
	trans.CreatedAt = time.Now().String()
	ref := getDbClient().Collection("transactions").Doc(trans.TxHash)
	err := getDbClient().RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		_, getErr := tx.Get(ref)
		if getErr == nil {
			return errors.New("transaction already exist")
		}

		return tx.Set(ref, trans)
	})
	return err
}
