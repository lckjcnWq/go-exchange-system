package contracts

import (
	ethereum1 "backend/internal/service/ethereum"
	"backend/internal/service/ethereum/abix"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	client := ethereum1.GetHTTPClient()

	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.routerAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}
	amounts, err := s.routerABI.Unpack("getAmountsOut", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}
	return amounts[0].([]*big.Int), nil
}

// GetPairReserves 获取交易对储备金
func (s *UniswapService) GetPairReserves(ctx context.Context, pairAddr common.Address) (reserve0, reserve1 *big.Int, err error) {
	data, err := s.pairABI.Pack("getReserves")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack data: %v", err)
	}
	client := ethereum1.GetHTTPClient()

	unpacked, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &pairAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call contract: %v", err)
	}
	reserve0 = new(big.Int).SetUint64(uint64(unpacked[0]))
	reserve1 = new(big.Int).SetUint64(uint64(unpacked[1]))
	return
}

// SwapExactETHForTokens ETH换代币
func (s *UniswapService) SwapExactETHForTokens(
	ctx context.Context,
	amountOutMin *big.Int,
	path []common.Address,
	to common.Address,
	deadline *big.Int,
	value *big.Int,
) (*types.Transaction, error) {
	_, err := s.routerABI.Pack("swapExactETHForTokens", amountOutMin, path, to, deadline)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %v", err)
	}
	// 这里需要实现交易签名和发送
	// 在实际应用中，你需要添加签名逻辑
	return nil, fmt.Errorf("not implemented")
}
