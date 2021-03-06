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
	"identity-token-relayer/contract/token721"
	"identity-token-relayer/lib"
	"identity-token-relayer/log"
	"math/big"
	"os"
)

var (
	hmyClient                *Client
	hmyPrivateKey            string
	OwnershipValidatorClient *ownership_validator.OwnershipValidator
	Token721ClientSet        = make(map[string]*token721.Token721, 0)
)

type Client struct {
	client *ethclient.Client
}

func InitClient() error {
	// load encrypted wallet private key
	hmyPrivateKeyRaw, err := os.ReadFile(config.Get().Hmy.PrivateKeyPath)
	if err != nil {
		return err
	}

	if config.Get().Hmy.OpenKMS {
		// decrypt the data
		result, decErr := lib.GetKmsClient().Decrypt(&kms.DecryptInput{CiphertextBlob: hmyPrivateKeyRaw})
		if decErr != nil {
			return decErr
		}

		hmyPrivateKey = string(result.Plaintext)
	} else {
		hmyPrivateKey = string(hmyPrivateKeyRaw)
	}

	// check private key
	_, err = NewPrivateKey(hmyPrivateKey)
	if err != nil {
		return err
	}

	return nil
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

func GetOwnershipValidatorClient() *ownership_validator.OwnershipValidator {
	if OwnershipValidatorClient == nil {
		address := common.HexToAddress(config.Get().Hmy.OwnershipValidatorAddress)
		OwnershipValidatorClient, _ = ownership_validator.NewOwnershipValidator(address, GetHmyClient().client)
	}
	return OwnershipValidatorClient
}

func GetToken721Client(address string) *token721.Token721 {
	if _, ok := Token721ClientSet[address]; !ok {
		Token721ClientSet[address], _ = token721.NewToken721(common.HexToAddress(address), GetHmyClient().client)
	}
	return Token721ClientSet[address]
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
