package tools

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"identity-token-relayer/contract/token721"
	"identity-token-relayer/eth"
	"identity-token-relayer/hmy"
	"identity-token-relayer/log"
	"math/big"
)

func GetMappingTotalSupply(mapping string) {
	supply, callErr := hmy.GetToken721Client(mapping).TotalSupply(new(bind.CallOpts))
	if callErr != nil {
		log.GetLogger().Fatal("get mapping total supply failed.", zap.String("error", callErr.Error()))
	}
	log.GetLogger().Info("get mapping total supply success.", zap.Int64("supply", supply.Int64()))
}

func GetOriginOwner(origin string, tokenId int64) {
	address := common.HexToAddress(origin)

	token721Client, err := token721.NewToken721(address, eth.GetEthClient().EthClient)
	if err != nil {
		log.GetLogger().Fatal("get eth client failed.", zap.String("error", err.Error()))
	}

	callOpts := &bind.CallOpts{}
	owner, err := token721Client.OwnerOf(callOpts, big.NewInt(tokenId))
	if err != nil {
		log.GetLogger().Fatal("get nft owner failed.", zap.String("error", err.Error()))
	}

	log.GetLogger().Info("get nft owner success.", zap.String("owner", owner.String()))
}
