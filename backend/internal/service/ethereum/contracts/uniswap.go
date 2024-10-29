package contracts

import (
	"backend/internal/service/ethereum/abix"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"math/big"
	"strings"
)

type UniswapService struct {
	routerAddr  common.Address
	factoryAddr common.Address
	routerABI   abi.ABI
	pairABI     abi.ABI
}

func NewUniswapService() (*UniswapService, error) {
	routerABI, err := abi.JSON(strings.NewReader(abix.UniswapV2Router02ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse router ABI: %v", err)
	}

	pairABI, err := abi.JSON(strings.NewReader(abix.UniswapV2PairABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pair ABI: %v", err)
	}

	cfg := g.Cfg()
	ctx := context.Background()

	return &UniswapService{
		routerAddr:  common.HexToAddress(cfg.MustGet(ctx, "ethereum.contracts.uniswap.router").String()),
		factoryAddr: common.HexToAddress(cfg.MustGet(ctx, "ethereum.contracts.uniswap.factory").String()),
		routerABI:   routerABI,
		pairABI:     pairABI,
	}, nil
}

// GetAmountOut 计算交易输出金额
func (s *UniswapService) GetAmountOut(ctx context.Context, amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	data, err := s.routerABI.Pack("getAmountsOut", amountIn, path)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %v", err)
	}
}
