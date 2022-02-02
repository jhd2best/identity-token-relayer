package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"go.uber.org/zap"
	"identity-token-relayer/log"
	"strconv"
)

type Nft struct {
	ContractAddress  string `firestore:"contract_address"`
	TokenId          int64  `firestore:"token_id"`
	OwnerAddress     string `firestore:"owner_address"`
	LastUpdateHeight int64  `firestore:"last_update_height"`
}

func GetProjectNftByTokenId(address string, tokenId int64) (Nft, error) {
	tokenIdStr := strconv.Itoa(int(tokenId))
	var nft Nft
	data, getErr := GetDbClient().Collection("projects").Doc(address).Collection("nfts").Doc(tokenIdStr).Get(context.Background())
	if getErr != nil {
		return nft, getErr
	}

	dataErr := data.DataTo(&nft)
	if dataErr != nil {
		log.GetLogger().Error("parse project nft data failed.", zap.String("error", dataErr.Error()))
		return nft, dataErr
	}

	return nft, nil
}

func BatchCreateNfts(nfts []Nft) ([]int64, error) {
	if len(nfts) > 500 {
		return nil, errors.New("batch create can not more 500 items")
	}

	addedNftTokenId := make([]int64, 0)
	batch := GetDbClient().Batch()

	for _, nft := range nfts {
		addedNftTokenId = append(addedNftTokenId, nft.TokenId)
		docName := strconv.Itoa(int(nft.TokenId))
		docRef := GetDbClient().Collection("projects").Doc(nft.ContractAddress).Collection("nfts").Doc(docName)
		batch.Set(docRef, nft)
	}

	_, commitErr := batch.Commit(context.Background())
	if commitErr != nil {
		return nil, commitErr
	}
	return addedNftTokenId, nil
}

func UpdateProjectNftOwner(address string, tokenId int64, newOwner string, height int64) error {
	tokenIdStr := strconv.Itoa(int(tokenId))
	_, updateErr := GetDbClient().Collection("projects").Doc(address).Collection("nfts").Doc(tokenIdStr).Update(context.Background(), []firestore.Update{
		{
			Path:  "owner_address",
			Value: newOwner,
		},
		{
			Path:  "last_update_height",
			Value: height,
		},
	})
	return updateErr
}