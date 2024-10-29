package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogf/gf/v2/frame/g"
	"math/big"
	"time"
)

type EventListener struct {
	ctx         context.Context
	cancel      context.CancelFunc
	routerAddr  common.Address
	factoryAddr common.Address
}

func NewEventListener() *EventListener {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := g.Cfg()

	return &EventListener{
		ctx:         ctx,
		cancel:      cancel,
		routerAddr:  common.HexToAddress(cfg.MustGet(ctx, "ethereum.contracts.uniswap.router").String()),
		factoryAddr: common.HexToAddress(cfg.MustGet(ctx, "ethereum.contracts.uniswap.factory").String()),
	}
}

// Start 开始监听合约事件
func (l *EventListener) Start() error {
	// 启动 WebSocket 监听
	go l.startWSListener()
	// 同时启动轮询作为备份
	go l.startPolling()
	return nil
}

// startWSListener 启动WebSocket监听
func (l *EventListener) startWSListener() {
	for {
		if !IsWSConnected() {
			time.Sleep(5 * time.Second)
			continue
		}

		query := ethereum.FilterQuery{
			Addresses: []common.Address{l.routerAddr, l.factoryAddr},
		}

		logs := make(chan types.Log)
		client := GetWSClient()
		sub, err := client.SubscribeFilterLogs(l.ctx, query, logs)
		if err != nil {
			g.Log().Warning(l.ctx, "Failed to subscribe to logs:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			select {
			case err := <-sub.Err():
				g.Log().Warning(l.ctx, "Subscription error:", err)
				break
			case log := <-logs:
				l.processEvent(&log)
			case <-l.ctx.Done():
				return
			}
		}
	}
}

// startPolling 启动轮询监听
func (l *EventListener) startPolling() {
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	var lastBlock uint64

	for {
		select {
		case <-ticker.C:
			// 如果WebSocket连接正常，不需要轮询
			if IsWSConnected() {
				continue
			}

			client := GetHTTPClient()
			currentBlock, err := client.BlockNumber(l.ctx)
			if err != nil {
				g.Log().Warning(l.ctx, "Error getting block number:", err)
				continue
			}

			if currentBlock > lastBlock {
				query := ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(lastBlock + 1)),
					ToBlock:   big.NewInt(int64(currentBlock)),
					Addresses: []common.Address{l.routerAddr, l.factoryAddr},
				}

				logs, err := client.FilterLogs(l.ctx, query)
				if err != nil {
					g.Log().Warning(l.ctx, "Error filtering logs:", err)
					continue
				}

				for _, log := range logs {
					l.processEvent(&log)
				}

				lastBlock = currentBlock
			}
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *EventListener) processEvent(vLog *types.Log) {
	// Uniswap V2 事件的方法ID（前8位）
	const (
		SwapMethodID        = "0xd78ad95f" // Swap(address,uint256,uint256,uint256,uint256,address)
		SyncMethodID        = "0x1c411e9a" // Sync(uint112,uint112)
		PairCreatedMethodID = "0x0d3648bd" // PairCreated(address,address,address,uint256)
	)

	// 提取事件的方法ID
	methodID := common.Bytes2Hex(vLog.Topics[0].Bytes())[:8]

	// 根据合约地址和方法ID处理不同的事件
	switch vLog.Address {
	case l.routerAddr:
		g.Log().Info(l.ctx, "Router event received",
			"contractAddress", vLog.Address.Hex(),
			"methodID", methodID,
			"blockNumber", vLog.BlockNumber,
			"txHash", vLog.TxHash.Hex(),
			"blockHash", vLog.BlockHash.Hex(),
			"logIndex", vLog.Index,
		)

	case l.factoryAddr:
		if methodID == PairCreatedMethodID {
			// PairCreated 事件有四个参数：token0, token1, pair, allPairsLength
			g.Log().Info(l.ctx, "New pair created",
				"contractAddress", vLog.Address.Hex(),
				"token0", common.HexToAddress(vLog.Topics[1].Hex()).Hex(),
				"token1", common.HexToAddress(vLog.Topics[2].Hex()).Hex(),
				"pairAddress", common.HexToAddress(vLog.Topics[3].Hex()).Hex(),
				"blockNumber", vLog.BlockNumber,
				"txHash", vLog.TxHash.Hex(),
			)
		}

	default:
		// 可能是交易对合约的事件
		switch methodID {
		case SwapMethodID:
			g.Log().Info(l.ctx, "Swap event received",
				"pairAddress", vLog.Address.Hex(),
				"blockNumber", vLog.BlockNumber,
				"txHash", vLog.TxHash.Hex(),
				"sender", common.HexToAddress(vLog.Topics[1].Hex()).Hex(),
				"recipient", common.HexToAddress(vLog.Topics[2].Hex()).Hex(),
				"data", common.Bytes2Hex(vLog.Data),
			)

		case SyncMethodID:
			g.Log().Info(l.ctx, "Sync event received",
				"pairAddress", vLog.Address.Hex(),
				"blockNumber", vLog.BlockNumber,
				"txHash", vLog.TxHash.Hex(),
				"data", common.Bytes2Hex(vLog.Data), // reserve0 和 reserve1 数据
			)

		default:
			g.Log().Debug(l.ctx, "Other contract event",
				"contractAddress", vLog.Address.Hex(),
				"methodID", methodID,
				"blockNumber", vLog.BlockNumber,
				"txHash", vLog.TxHash.Hex(),
				"topics", formatTopics(vLog.Topics),
				"data", common.Bytes2Hex(vLog.Data),
			)
		}
	}
}

// formatTopics 格式化事件的 topics
func formatTopics(topics []common.Hash) string {
	result := "["
	for i, topic := range topics {
		if i > 0 {
			result += ", "
		}
		result += topic.Hex()
	}
	result += "]"
	return result
}

// 可选：添加事件解析辅助函数
func parseSwapEvent(vLog *types.Log) (sender, recipient common.Address, amount0In, amount1In, amount0Out, amount1Out *big.Int) {
	sender = common.HexToAddress(vLog.Topics[1].Hex())
	recipient = common.HexToAddress(vLog.Topics[2].Hex())

	// 解析 Data 字段中的数值
	data := vLog.Data
	amount0In = new(big.Int).SetBytes(data[:32])
	amount1In = new(big.Int).SetBytes(data[32:64])
	amount0Out = new(big.Int).SetBytes(data[64:96])
	amount1Out = new(big.Int).SetBytes(data[96:128])

	return
}

// 可选：添加事件解析辅助函数
func parseSyncEvent(vLog *types.Log) (reserve0, reserve1 *big.Int) {
	data := vLog.Data
	reserve0 = new(big.Int).SetBytes(data[:32])
	reserve1 = new(big.Int).SetBytes(data[32:64])
	return
}

func (l *EventListener) Stop() {
	l.cancel()
}
