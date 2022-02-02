package hmy

import (
	"context"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"identity-token-relayer/config"
	"identity-token-relayer/contract/ownership_validator"
	"identity-token-relayer/lib"
	"identity-token-relayer/log"
	"math/big"
	"os"
)

var (
	hmyClient     *Client
	hmyPrivateKey string
	OwnershipValidatorClient *ownership_validator.Contract
)

type Client struct {
	client *ethclient.Client
}

func init() {
	// load encrypted wallet private key
	hmyPrivateKeyRaw, err := os.ReadFile(config.Get().Hmy.PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	if config.Get().Hmy.OpenKMS {
		// decrypt the data
		result, err := lib.GetKmsClient().Decrypt(&kms.DecryptInput{CiphertextBlob: hmyPrivateKeyRaw})
		if err != nil {
			panic(err)
		}

		hmyPrivateKey = string(result.Plaintext)
	} else {
		hmyPrivateKey = string(hmyPrivateKeyRaw)
	}

	// check private key
	_, err = NewPrivateKey(hmyPrivateKey)
	if err != nil {
		panic(err)
	}

	log.GetLogger().Info("load hmy private key succeed")
}

func GetHmyClient() *Client {
	if hmyClient == nil {
		hmyConfig := config.Get().Hmy
		if hmyConfig.RpcEndpoints == "" {
			log.GetLogger().Fatal("harmony RPC endpoints not found")
		}

		client, err := ethclient.Dial(hmyConfig.RpcEndpoints)
		if err != nil {
			log.GetLogger().Fatal("harmony RPC endpoints init failed", zap.String("error", err.Error()))
		}
		hmyClient = &Client{client: client}
	}

	return hmyClient
}

func GetOwnershipValidatorClient() *ownership_validator.Contract {
	if OwnershipValidatorClient == nil {
		address := common.HexToAddress(config.Get().Hmy.OwnershipValidatorAddress)
		OwnershipValidatorClient, _ = ownership_validator.NewContract(address, GetHmyClient().client)
	}
	return OwnershipValidatorClient
}

func GetHmyPrivateKey() string {
	return hmyPrivateKey
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.client.SuggestGasPrice(ctx)
}

func (c *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.client.HeaderByNumber(ctx, number)
}

func (c *Client) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return c.client.FilterLogs(ctx, query)
}

func (c *Client) Nonce(ctx context.Context, fromAddress common.Address) (uint64, error) {
	return c.client.PendingNonceAt(context.Background(), fromAddress)
}

func (c *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client.NetworkID(context.Background())
}

func (c *Client) TransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(context.Background(), txHash)
}

func (c *Client) TransferringONE(ctx context.Context, pk *PrivateKey, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	nonce, err := c.Nonce(ctx, pk.Address())
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := c.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	var data []byte
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amount,
		Gas:      uint64(21000),
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk.key)
	if err != nil {
		return nil, err
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (c *Client) SendRAWTransaction(ctx context.Context, pk *PrivateKey, rawByte []byte) (*types.Transaction, error) {
	tx := &types.Transaction{}

	err := rlp.DecodeBytes(rawByte, &tx)
	if err != nil {
		return nil, err
	}

	err = c.client.SendTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
