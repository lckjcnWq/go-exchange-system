package controller

import (
	v1 "backend/api/v1"
	"backend/internal/logic"
	"backend/internal/model"
	"backend/internal/service/ethereum/contracts"
	"context"
	"math/big"
)

var Trade = cTrade{
	tradeLogic: logic.NewTradeLogic(),
}

type cTrade struct {
	tradeLogic *logic.TradeLogic
}

// CreateTrade 创建交易
func (c *cTrade) CreateTrade(ctx context.Context, req *v1.CreateTradeReq) (res *v1.CreateTradeRes, err error) {
	// 初始化UniswapService
	_, err = contracts.NewUniswapService()
	if err != nil {
		return nil, err
	}
	//解析参数
	amountIn := new(big.Int)
	amountIn.SetString(req.AmountIn, 10)
	amountOutMin := new(big.Int)
	amountOutMin.SetString(req.AmountOutMin, 10)

	//获取用户地址
	userAddress := ctx.Value("userAddress").(string)

	// 创建交易记录
	trade := &model.Trade{
		UserAddress: userAddress,
		TokenIn:     req.TokenIn,
		TokenOut:    req.TokenOut,
		AmountIn:    req.AmountIn,
		Status:      model.TradeStatusPending,
	}

	// 保存交易记录
	if err := c.tradeLogic.CreateTrade(ctx, trade); err != nil {
		return nil, err
	}

	return &v1.CreateTradeRes{
		TxHash: trade.TxHash,
	}, nil
}

// GetTrades 获取交易列表
func (c *cTrade) GetTrades(ctx context.Context, req *v1.GetTradesReq) (res *v1.GetTradesRes, err error) {
	// 获取用户地址
	userAddress := ctx.Value("userAddress").(string)

	// 获取交易列表
	trades, err := c.tradeLogic.GetUserTrades(ctx, userAddress, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var tradeInfos []*v1.TradeInfo
	for _, trade := range trades {
		tradeInfos = append(tradeInfos, &v1.TradeInfo{
			TxHash:        trade.TxHash,
			TokenIn:       trade.TokenIn,
			TokenOut:      trade.TokenOut,
			AmountIn:      trade.AmountIn,
			AmountOut:     trade.AmountOut,
			Status:        string(trade.Status),
			Confirmations: trade.Confirmations,
			CreatedAt:     trade.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.GetTradesRes{
		List:  tradeInfos,
		Total: len(trades),
	}, nil
}
