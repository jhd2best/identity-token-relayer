package eth

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"identity-token-relayer/config"
	"identity-token-relayer/log"
	"math/big"
)

var (
	ethClient *Client
)

type Client struct {
	EthClient *ethclient.Client
}

func GetEthClient() *Client {
	if ethClient == nil {
		ethConfig := config.Get().Eth
		if ethConfig.RpcEndpoints == "" {
			log.GetLogger().Fatal("ethereum RPC endpoints not found")
		}

		client, err := ethclient.Dial(ethConfig.RpcEndpoints)
		if err != nil {
			log.GetLogger().Fatal("ethereum RPC endpoints init failed", zap.String("error", err.Error()))
		}
		ethClient = &Client{EthClient: client}
	}

	return ethClient
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.EthClient.SuggestGasPrice(ctx)
}

func (c *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.EthClient.HeaderByNumber(ctx, number)
}

func (c *Client) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return c.EthClient.FilterLogs(ctx, query)
}

func (c *Client) Nonce(ctx context.Context, fromAddress common.Address) (uint64, error) {
	return c.EthClient.PendingNonceAt(context.Background(), fromAddress)
}

func (c *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.EthClient.NetworkID(context.Background())
}

func (c *Client) TransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return c.EthClient.TransactionReceipt(context.Background(), txHash)
}
