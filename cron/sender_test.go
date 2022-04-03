package cron

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"identity-token-relayer/hmy"
	"identity-token-relayer/log"
	"math/big"
	"strings"
	"testing"
)

func TestGetOwnerOf(t *testing.T) {
	opt := &bind.CallOpts{}
	owner, callErr := hmy.GetOwnershipValidatorClient().OwnerOf(opt, "BC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", big.NewInt(10001))
	if callErr != nil {
		if strings.Index(callErr.Error(), "nonexistent") >= 0 {
			log.GetLogger().Error("token non-existent")
		} else {
			panic(callErr)
		}
	} else {
		log.GetLogger().Info(fmt.Sprintf("owner:%s", owner))
	}
}

func TestUpdateOrInit(t *testing.T) {
	ethAddress := "BC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"
	tokenId := big.NewInt(16667)
	newOwner := common.HexToAddress("0xaf14a532b76d6812e8d036a08be92c6dd6839a48")

	transOpt := new(bind.TransactOpts)
	opt := new(bind.CallOpts)

	privateKey, err := crypto.HexToECDSA(hmy.GetHmyPrivateKey())
	if err != nil {
		panic(err)
	}

	gasPrice, err := hmy.GetHmyClient().SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	if gasPrice.Int64() < 40000000000 {
		gasPrice = big.NewInt(40000000000)
	}

	chainId, err := hmy.GetHmyClient().NetworkID(context.Background())
	if err != nil {
		panic(err)
	}

	transOpt, err = bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		panic(err)
	}

	transOpt.Value = big.NewInt(0)
	transOpt.GasLimit = uint64(1000000)
	transOpt.GasPrice = gasPrice

	nonce, err := hmy.GetHmyClient().Nonce(context.Background(), transOpt.From)
	if err != nil {
		panic(err)
	}
	transOpt.Nonce = big.NewInt(int64(nonce))

	owner, callErr := hmy.GetOwnershipValidatorClient().OwnerOf(opt, ethAddress, tokenId)
	if callErr != nil {
		if strings.Index(callErr.Error(), "nonexistent") >= 0 {
			log.GetLogger().Info("token non-existent. will auto-init it")

			execTrans, execErr := hmy.GetOwnershipValidatorClient().Initialize(transOpt, ethAddress, []common.Address{newOwner}, []*big.Int{tokenId})
			if execErr != nil {
				panic(execErr)
			}

			log.GetLogger().Info(fmt.Sprintf("send init tx success. hash:%s", execTrans.Hash().Hex()))
		} else {
			panic(callErr)
		}
	} else {
		log.GetLogger().Info(fmt.Sprintf("got token. owner:%s", owner))

		execTrans, execErr := hmy.GetOwnershipValidatorClient().UpdateOwnership(transOpt, ethAddress, []common.Address{newOwner}, []*big.Int{tokenId})
		if execErr != nil {
			panic(execErr)
		}

		log.GetLogger().Info(fmt.Sprintf("send update tx success. hash:%s", execTrans.Hash().Hex()))
	}
}
