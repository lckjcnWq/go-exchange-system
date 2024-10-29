package ethereum

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	blockListener *BlockListener
	eventListener *EventListener
)

func InitBlockchainService(ctx context.Context) error {
	//1.初始化以太坊客户端
	if err := InitClients(ctx); err != nil {
		return err
	}
	//2.初始化区块监听服务
	blockListener = NewBlockListener()
	if err := blockListener.Start(); err != nil {
		return err
	}
	//3.初始化并启动事件监听
	eventListener = NewEventListener()
	if err := eventListener.Start(); err != nil {
		return err
	}

	g.Log().Info(ctx, "Blockchain services initialized successfully")
	return nil
}

// CloseBlockchainServices 关闭所有区块链相关服务
func CloseBlockchainServices() {
	if blockListener != nil {
		blockListener.Stop()
	}
	if eventListener != nil {
		eventListener.Stop()
	}
	Close()
}
