package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogf/gf/v2/frame/g"
	"log"
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
	// 在这里处理交易，例如：
	// 1. 检查是否是我们关注的合约交易
	// 2. 解析交易数据
	// 3. 更新数据库
	// 4. 通知相关用户
	g.Log().Debug(context.Background(), "Processing transaction",
		"hash", tx.Hash().Hex(),
		"value", tx.Value(),
	)
}

// Stop 停止监听
func (l *BlockListener) Stop() {
	l.cancel()
}
