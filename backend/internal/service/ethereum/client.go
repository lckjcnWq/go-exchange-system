package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ClientManager struct {
	httpClient  *ethclient.Client
	wsClient    *ethclient.Client
	rpcClient   *rpc.Client
	wsRpcClient *rpc.Client
}

var (
	clientManager *ClientManager
)
