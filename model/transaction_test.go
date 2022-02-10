package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"identity-token-relayer/log"
	"testing"
	"time"
)

func TestBatchUpdateTransactionStatus(t *testing.T) {
	trans := make([]Transaction, 0)
	iter := GetDbClient().Collection("transactions").Where("status", "==", "failed").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}

		tran := Transaction{}
		_ = doc.DataTo(&tran)

		trans = append(trans, tran)
	}

	log.GetLogger().Info(fmt.Sprintf("total found transaction: %d", len(trans)))

	pendingTransactionSet := make([]Transaction, 0)
	updateData := []firestore.Update{
		{
			Path:  "status",
			Value: "error",
		},
		{
			Path:  "retry_times",
			Value: 0,
		},
	}
	for index, tran := range trans {
		pendingTransactionSet = append(pendingTransactionSet, tran)

		if len(pendingTransactionSet) >= 499 {
			_, batchErr := BatchUpdateTransactions(pendingTransactionSet, updateData)
			if batchErr != nil {
				panic(batchErr)
			} else {
				log.GetLogger().Info(fmt.Sprintf("updated transaction to %d", index))
				pendingTransactionSet = make([]Transaction, 0)
				time.Sleep(1 * time.Second)
			}
		}
	}

	if len(pendingTransactionSet) > 0 {
		_, batchErr := BatchUpdateTransactions(pendingTransactionSet, updateData)
		if batchErr != nil {
			panic(batchErr)
		}
	}
	log.GetLogger().Info("batch update transaction success.")
}
