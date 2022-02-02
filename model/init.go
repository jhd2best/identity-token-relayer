package model

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"identity-token-relayer/config"
)

var (
	dbClient *firestore.Client
)

func InitDb() error {
	dbConfig := config.Get().Db
	sa := option.WithCredentialsFile(dbConfig.ServiceAccountPath)

	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return err
	}

	dbClient, err = app.Firestore(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func GetDbClient() *firestore.Client {
	return dbClient
}

func CloseDb() error {
	return dbClient.Close()
}
