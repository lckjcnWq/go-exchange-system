package dao

import (
	"backend/internal/model"
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type TradeDao struct{}

func NewTradeDao() *TradeDao {
	return &TradeDao{}
}

func (t *TradeDao) Create(ctx context.Context, trade *model.Trade) error {
	_, err := g.DB().Model("trades").Ctx(ctx).Data(trade).Insert()
	return err
}

func (t *TradeDao) Update(ctx context.Context, trade *model.Trade) error {
	_, err := g.DB().Model("trades").
		Ctx(ctx).
		Where("id", trade.Id).
		Data(trade).
		Update()

	return err
}

func (t *TradeDao) GetByTxHash(ctx context.Context, txHash string) (*model.Trade, error) {
	var trade model.Trade
	err := g.DB().Model("trades").
		Ctx(ctx).
		Where("tx_hash", txHash).
		Scan(&trade)
	return &trade, err
}

func (d *TradeDao) ListByUser(ctx context.Context, userAddress string, page, size int) ([]*model.Trade, error) {
	var trades []*model.Trade
	err := g.DB().Model("trades").Ctx(ctx).
		Where("user_address", userAddress).
		Order("created_at DESC").
		Page(page, size).
		Scan(&trades)
	return trades, err
}
