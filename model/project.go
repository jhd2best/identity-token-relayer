package model

import (
	"context"
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
}

var (
	projects map[string]Project
)

func GetAllProjects() map[string]Project {
	return projects
}

func SyncAllProjects() error {
	ctx := context.Background()
	newProjects := make(map[string]Project, 0)
	iter := getDbClient().Collection("projects").Documents(ctx)
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

func AddProject() {
	// TODO: only for test
	ctx := context.Background()
	project := Project{
		ContractAddress:  "0x8a90cab2b38dba80c64b7734e58ee1db38b8992e",
		MappingAddress:   "",
		Name:             "Doodles",
		Symbol:           "DOODLE",
		BaseURI:          "ipfs://QmPMc4tcBsMqLRuCQtPmPe84bpSjrC3Ky7t3JWuHXYB4aS/",
		Status:           "created",
		LastUpdateHeight: 14108230,
	}

	_, err := getDbClient().Collection("projects").Doc(project.ContractAddress).Set(ctx, project)
	if err != nil {
		log.GetLogger().Error("add project data failed.", zap.String("error", err.Error()))
	}

	log.GetLogger().Info("add project data success", zap.String("name", project.Name))
}
