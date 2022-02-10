package model

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"identity-token-relayer/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.GetLogger().Info("testing init...")

	err := InitDb()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestAddProject(t *testing.T) {
	address := common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D").String()
	project := Project{
		ContractAddress:  address,
		MappingAddress:   "",
		Name:             "BoredApeYachtClub",
		Symbol:           "BAYC",
		BaseURI:          "ipfs://QmeSjSinHpPnmXmspMjwiXyN6zS4E9zccariGR3jxcaWtq/",
		Status:           "created",
		LastUpdateHeight: 14118753,
		Enable:           true,
	}

	_, err := GetDbClient().Collection("projects").Doc(project.ContractAddress).Set(context.Background(), project)
	if err != nil {
		log.GetLogger().Error("add project data failed.", zap.String("error", err.Error()))
	}

	log.GetLogger().Info("add project data success", zap.String("name", project.Name))
}
