package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"sync"
	"time"
)

type ConnectionMode int

const (
	ModeHTTP ConnectionMode = iota
	ModeWebSocket
)

type ClientManager struct {
	ctx          context.Context
	cancel       context.CancelFunc
	httpClient   *ethclient.Client
	wsClient     *ethclient.Client
	rpcClient    *rpc.Client
	wsRpcClient  *rpc.Client
	wsConnected  bool
	mutex        sync.RWMutex
	reconnecting bool
}

var (
	clientManager *ClientManager
	once          sync.Once
)

// InitClients 初始化以太坊客户端
func InitClients(ctx context.Context) error {
	var err error
	once.Do(func() {
		ctx, cancel := context.WithCancel(ctx)
		clientManager = &ClientManager{
			ctx:    ctx,
			cancel: cancel,
		}
		err = clientManager.initialize(ctx)
	})
	return err
}

func (cm *ClientManager) initialize(ctx context.Context) error {
	// 初始化 HTTP 客户端
	if err := cm.initHTTPClient(ctx); err != nil {
		return err
	}

	// 初始化 WebSocket 客户端
	if err := cm.initWSClient(ctx); err != nil {
		g.Log().Warning(ctx, "WebSocket connection failed, will retry in background:", err)
		// 启动后台重连
		go cm.wsReconnectLoop()
	}

	return nil
}

func (cm *ClientManager) initHTTPClient(ctx context.Context) error {
	cfg := g.Cfg()
	httpUrl := cfg.MustGet(ctx, "ethereum.network.sepolia.httpUrl").String()

	rpcClient, err := rpc.Dial(httpUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to ethereum node via HTTP: %v", err)
	}

	cm.mutex.Lock()
	cm.httpClient = ethclient.NewClient(rpcClient)
	cm.rpcClient = rpcClient
	cm.mutex.Unlock()

	return nil
}

func (cm *ClientManager) initWSClient(ctx context.Context) error {
	cfg := g.Cfg()
	wsUrl := cfg.MustGet(ctx, "ethereum.network.sepolia.wsUrl").String()

	wsRpcClient, err := rpc.Dial(wsUrl)
	if err != nil {
		return err
	}

	cm.mutex.Lock()
	cm.wsClient = ethclient.NewClient(wsRpcClient)
	cm.wsRpcClient = wsRpcClient
	cm.wsConnected = true
	cm.mutex.Unlock()

	// 启动心跳检测
	go cm.heartbeat()

	return nil
}

func (cm *ClientManager) wsReconnectLoop() {
	for {
		if cm.reconnecting {
			continue
		}

		cm.mutex.Lock()
		cm.reconnecting = true
		cm.mutex.Unlock()

		ctx := gctx.New()
		err := cm.initWSClient(ctx)

		cm.mutex.Lock()
		cm.reconnecting = false
		cm.mutex.Unlock()

		if err == nil {
			g.Log().Info(ctx, "WebSocket reconnected successfully")
			return
		}

		g.Log().Warning(ctx, "WebSocket reconnection failed, will retry:", err)
		time.Sleep(5 * time.Second)
	}
}

func (cm *ClientManager) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if cm.wsClient == nil {
				continue
			}

			// 尝试获取最新区块号来检测连接状态
			_, err := cm.wsClient.BlockNumber(cm.ctx)
			if err != nil {
				g.Log().Warning(cm.ctx, "WebSocket heartbeat failed:", err)
				cm.mutex.Lock()
				cm.wsConnected = false
				cm.mutex.Unlock()

				// 触发重连
				go cm.wsReconnectLoop()
			}
		case <-cm.ctx.Done():
			return
		}
	}
}

// GetClient 智能获取客户端
func (cm *ClientManager) GetClient(mode ConnectionMode) *ethclient.Client {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	switch mode {
	case ModeWebSocket:
		if cm.wsConnected && cm.wsClient != nil {
			return cm.wsClient
		}
		return cm.httpClient
	default:
		return cm.httpClient
	}
}

// IsWSConnected 检查WebSocket连接状态
func (cm *ClientManager) IsWSConnected() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.wsConnected
}

// Close 关闭所有连接
func (cm *ClientManager) Close() {
	cm.cancel()

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.httpClient != nil {
		cm.httpClient.Close()
	}
	if cm.wsClient != nil {
		cm.wsClient.Close()
	}
	if cm.rpcClient != nil {
		cm.rpcClient.Close()
	}
	if cm.wsRpcClient != nil {
		cm.wsRpcClient.Close()
	}
}

// 提供给外部的便捷方法
func GetHTTPClient() *ethclient.Client {
	return clientManager.GetClient(ModeHTTP)
}

func GetWSClient() *ethclient.Client {
	return clientManager.GetClient(ModeWebSocket)
}

func IsWSConnected() bool {
	return clientManager.IsWSConnected()
}

func Close() {
	if clientManager != nil {
		clientManager.Close()
	}
}
