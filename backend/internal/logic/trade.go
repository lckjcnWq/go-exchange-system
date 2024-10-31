package logic

import (
	"backend/internal/dao"
	"backend/internal/model"
	ethereum1 "backend/internal/service/ethereum"
	"context"
	"github.com/ethereum/go-ethereum/common"
)

type TradeLogic struct {
	tradeDao *dao.TradeDao
}

func NewTradeLogic() *TradeLogic {
	return &TradeLogic{
		tradeDao: dao.NewTradeDao(),
	}
}

func (l *TradeLogic) CreateTrade(ctx context.Context, trade *model.Trade) error {
	return l.tradeDao.Create(ctx, trade)
}

func (l *TradeLogic) UpdateTradeStatus(ctx context.Context, txHash string) error {
	// 获取交易
	trade, err := l.tradeDao.GetByTxHash(ctx, txHash)
	if err != nil {
		return err
	}

	// 获取交易收据
	client := ethereum1.GetHTTPClient()
	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}

	// 更新状态
	if receipt.Status == 1 {
		trade.Status = model.TradeStatusConfirmed
	} else {
		trade.Status = model.TradeStatusFailed
	}
	trade.BlockNumber = receipt.BlockNumber.Uint64()

	// 获取当前区块
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return err
	}
	trade.Confirmations = currentBlock - receipt.BlockNumber.Uint64()

	return l.tradeDao.Update(ctx, trade)
}

func (l *TradeLogic) GetUserTrades(ctx context.Context, userAddress string, page, size int) ([]*model.Trade, error) {
	return l.tradeDao.ListByUser(ctx, userAddress, page, size)
}
