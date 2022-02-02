package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"identity-token-relayer/log"
)

type Project struct {
	ContractAddress  string `firestore:"contract_address"`
	MappingAddress   string `firestore:"mapping_address,omitempty"`
	Name             string `firestore:"name"`
	Symbol           string `firestore:"symbol"`
	BaseURI          string `firestore:"base_uri"`
	Status           string `firestore:"status"` // created, initialed
	LastUpdateHeight int64  `firestore:"last_update_height"`
	Enable           bool   `firestore:"enable"`
}

var (
	projects map[string]Project
)

func GetAllProjects() map[string]Project {
	return projects
}

func SyncAllEnableProjects() error {
	ctx := context.Background()
	newProjects := make(map[string]Project, 0)
	iter := GetDbClient().Collection("projects").Where("enable", "==", true).Documents(ctx)
	for {
		if doc, iterErr := iter.Next(); iterErr != nil {
			if iterErr == iterator.Done {
				break
			} else {
				log.GetLogger().Error("sync project data failed.", zap.String("error", iterErr.Error()))
				return iterErr
			}
		} else {
			project := Project{}
			dataErr := doc.DataTo(&project)
			if dataErr != nil {
				log.GetLogger().Error("parse project data failed.", zap.String("error", dataErr.Error()))
				continue
			}

			newProjects[project.ContractAddress] = project
		}
	}
	projects = newProjects

	return nil
}

func UpdateProjectLastHeight(address string, height int64) error {
	if project, ok := projects[address]; !ok {
		return errors.New("project not found")
	} else {
		project.LastUpdateHeight = height

		_, updateErr := GetDbClient().Collection("projects").Doc(project.ContractAddress).Update(context.Background(), []firestore.Update{
			{
				Path:  "last_update_height",
				Value: height,
			},
		})
		if updateErr == nil {
			projects[address] = project
		}
		return updateErr
	}
}
