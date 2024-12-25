package controller

import (
	"backend/api/v1"
	"backend/internal/service/ethereum/contracts"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var Price = cPrice{}

type cPrice struct{}

// GetPrice 获取代币价格
func (c *cPrice) GetPrice(ctx context.Context, req *v1.GetPriceReq) (res *v1.GetPriceRes, err error) {
	// 初始化UniswapService
	uniswapService, err := contracts.NewUniswapService()
	if err != nil {
		return nil, err
	}

	// 将输入金额转换为big.Int
	amountIn := new(big.Int)
	amountIn.SetString(req.AmountIn, 10)

	// 构建交易路径
	path := []common.Address{
		common.HexToAddress(req.TokenIn),
		common.HexToAddress(req.TokenOut),
	}

	// 获取输出金额
	amounts, err := uniswapService.GetAmountOut(ctx, amountIn, path)
	if err != nil {
		return nil, err
	}

	// 计算汇率
	rate := new(big.Float).Quo(
		new(big.Float).SetInt(amounts[1]),
		new(big.Float).SetInt(amounts[0]),
	)

	return &v1.GetPriceRes{
		AmountOut: amounts[1].String(),
		Rate:      rate.String(),
	}, nil
}
