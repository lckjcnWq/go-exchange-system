package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/your-project/abix"
)

// 区块链监听服务
type BlockListener struct {
	blockChan chan *types.Header
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewBlockListener() *BlockListener {
	ctx, cancel := context.WithCancel(context.Background())
	return &BlockListener{
		blockChan: make(chan *types.Header),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (l *BlockListener) Start() error {
	client := GetWSClient()
	sub, err := client.SubscribeNewHead(l.ctx, l.blockChan)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("block listener error: %v", err)
				return
			case header := <-l.blockChan:
				l.processNewBlock(header)
			case <-l.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// processNewBlock 处理新区块
func (l *BlockListener) processNewBlock(header *types.Header) {
	client := GetHTTPClient()
	block, err := client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		log.Printf("Error getting block: %v", err)
		return
	}

	//处理区块中的交易
	for _, tx := range block.Transactions() {
		//处理交易
		l.processTransaction(tx)
	}

	g.Log().Info(context.Background(), "New block processed",
		"number", block.Number().Uint64(),
		"hash", block.Hash().Hex(),
		"transactions", len(block.Transactions()),
	)
}

// processTransaction 处理交易
func (l *BlockListener) processTransaction(tx *types.Transaction) {
	ctx := context.Background()

	// 1. 检查是否是关注的合约交易
	if tx.To() == nil {
		// 合约创建交易,跳过
		return
	}

	// 获取关注的合约地址列表
	watchedContracts := []string{
		"0x86", // UniswapV2合约地址
		"0x64", // 其他关注的合约地址
	}

	isWatched := false
	for _, contract := range watchedContracts {
		if tx.To().Hex() == contract {
			isWatched = true
			break
		}
	}

	if !isWatched {
		return
	}

	// 2. 解析交易数据
	input := tx.Data()
	if len(input) < 4 {
		return
	}

	// 获取方法签名(前4字节)
	methodID := input[:4]

	// 根据方法签名处理不同的合约调用
	switch string(methodID) {
	case "0x86": // swap方法的签名
		// 解析swap参数
		swapData, err := parseSwapData(input)
		if err != nil {
			g.Log().Error(ctx, "Failed to parse swap data",
				"tx", tx.Hash().Hex(),
				"error", err)
			return
		}

		// 3. 更新数据库
		err = l.saveSwapToDB(ctx, tx, swapData)
		if err != nil {
			g.Log().Error(ctx, "Failed to save swap to db",
				"tx", tx.Hash().Hex(),
				"error", err)
			return
		}

		// 4. 通知相关用户
		l.notifyUsers(ctx, tx, swapData)

	case "0x64": // 其他方法的签名
		// 处理其他方法调用
	}

	// 记录处理日志
	g.Log().Debug(ctx, "Transaction processed",
		"hash", tx.Hash().Hex(),
		"contract", tx.To().Hex(),
		"method", methodID)
}

// SwapEvent 表示Swap事件的数据结构
type SwapEvent struct {
	TxHash      string
	BlockNumber uint64
	Sender      string
	Amount0In   *big.Int
	Amount1In   *big.Int
	Amount0Out  *big.Int
	Amount1Out  *big.Int
	To          string
	Timestamp   time.Time
}

// parseSwapData 解析swap交易数据
func parseSwapData(input []byte) (*SwapEvent, error) {
	// 获取UniswapV2合约ABI
	uniswapABI, err := abix.GetUniswapV2ABI()
	if err != nil {
		return nil, fmt.Errorf("failed to get UniswapV2 ABI: %v", err)
	}

	// 解析方法调用数据
	method, err := uniswapABI.MethodById(input[:4])
	if err != nil {
		return nil, fmt.Errorf("failed to get method: %v", err)
	}

	// 解析参数
	args := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(args, input[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack inputs: %v", err)
	}

	// 构造SwapEvent
	event := &SwapEvent{
		Amount0In:  args["amount0In"].(*big.Int),
		Amount1In:  args["amount1In"].(*big.Int),
		Amount0Out: args["amount0Out"].(*big.Int),
		Amount1Out: args["amount1Out"].(*big.Int),
		To:         args["to"].(string),
		Timestamp:  time.Now(),
	}

	return event, nil
}

// saveSwapToDB 保存swap交易到数据库
func (l *BlockListener) saveSwapToDB(ctx context.Context, tx *types.Transaction, data *SwapEvent) error {
	// 补充交易相关信息
	data.TxHash = tx.Hash().Hex()
	data.Sender = tx.From().Hex()

	// 获取区块信息
	client := GetHTTPClient()
	block, err := client.BlockByHash(ctx, tx.BlockHash())
	if err != nil {
		return fmt.Errorf("failed to get block: %v", err)
	}
	data.BlockNumber = block.NumberU64()

	// 使用gorm保存到数据库
	db := g.DB()
	result := db.Create(data)
	if result.Error != nil {
		return fmt.Errorf("failed to save to database: %v", result.Error)
	}

	return nil
}

// notifyUsers 通知相关用户
func (l *BlockListener) notifyUsers(ctx context.Context, tx *types.Transaction, data *SwapEvent) {
	// 查询订阅了该合约的用户
	var subscribers []string
	db := g.DB()
	err := db.Table("contract_subscribers").
		Where("contract_address = ?", tx.To().Hex()).
		Pluck("user_address", &subscribers).Error

	if err != nil {
		g.Log().Error(ctx, "Failed to get subscribers",
			"contract", tx.To().Hex(),
			"error", err)
		return
	}

	// 构造通知消息
	message := map[string]interface{}{
		"type": "swap",
		"data": data,
	}

	// 向每个订阅者发送通知
	for _, subscriber := range subscribers {
		// 这里可以使用websocket或消息队列发送通知
		err := sendNotification(subscriber, message)
		if err != nil {
			g.Log().Error(ctx, "Failed to notify subscriber",
				"subscriber", subscriber,
				"error", err)
			continue
		}
	}
}

// sendNotification 发送通知
func sendNotification(subscriber string, message interface{}) error {
	// 实现具体的通知发送逻辑
	// 可以使用websocket或消息队列
	return nil
}

// Stop 停止监听
func (l *BlockListener) Stop() {
	l.cancel()
}
