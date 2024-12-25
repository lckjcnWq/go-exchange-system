package trade

import (
	v1 "backend/api/v1"
	"backend/internal/logic"
	"backend/internal/model"
	"backend/internal/service/ethereum/contracts"
	"backend/internal/service/ws"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogf/gf/v2/frame/g"
	"math/big"
	"time"
)

type TradeService struct {
	uniswap *contracts.UniswapService
	logic   *logic.TradeLogic
}

func NewTradeService() (*TradeService, error) {
	uniswap, err := contracts.NewUniswapService()
	if err != nil {
		return nil, err
	}

	return &TradeService{
		uniswap: uniswap,
		logic:   logic.NewTradeLogic(),
	}, nil
}

// CreateTrade 创建交易
func (s *TradeService) CreateTrade(ctx context.Context, params *model.TradeParams) (*model.Trade, error) {
	// 1. 检查余额
	if err := s.checkBalance(ctx, params); err != nil {
		return nil, err
	}

	// 2. 获取价格报价
	amountOut, err := s.getQuote(ctx, params)
	if err != nil {
		return nil, err
	}

	// 3. 检查滑点
	if err := s.checkSlippage(params.AmountOutMin, amountOut); err != nil {
		return nil, err
	}

	// 4. 发送交易
	tx, err := s.sendTransaction(ctx, params)
	if err != nil {
		return nil, err
	}

	// 5. 创建交易记录
	trade := &model.Trade{
		TxHash:      tx.Hash().String(),
		UserAddress: params.UserAddress,
		TokenIn:     params.TokenIn.String(),
		TokenOut:    params.TokenOut.String(),
		AmountIn:    params.AmountIn.String(),
		AmountOut:   amountOut.String(),
		Status:      model.TradeStatusPending,
	}

	// 6. 保存交易记录
	if err := s.logic.CreateTrade(ctx, trade); err != nil {
		return nil, err
	}

	// 7. 启动交易监控
	go s.monitorTransaction(trade.TxHash)

	return trade, nil
}

// monitorTransaction 监控交易状态
func (s *TradeService) monitorTransaction(txHash string) {
	ctx := context.Background()
	trade, err := s.logic.GetTradeByTxHash(ctx, txHash)
	if err != nil {
		g.Log().Error(ctx, "Failed to get trade:", err)
		return
	}

	// 等待交易确认
	receipt, err := s.waitForTransaction(ctx, common.HexToHash(txHash))
	if err != nil {
		g.Log().Error(ctx, "Transaction failed:", err)
		trade.Status = model.TradeStatusFailed
		s.logic.UpdateTrade(ctx, trade)
		return
	}

	// 更新交易状态
	trade.Status = model.TradeStatusConfirmed
	trade.BlockNumber = receipt.BlockNumber.Uint64()
	s.logic.UpdateTrade(ctx, trade)

	// 发送WebSocket通知
	s.notifyTradeUpdate(trade)
}

// waitForTransaction 等待交易确认
func (s *TradeService) waitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := s.uniswap.GetTransactionReceipt(ctx, txHash)
		if err != nil {
			return nil, err
		}

		if receipt != nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			continue
		}
	}
}

// notifyTradeUpdate 通过WebSocket发送交易更新通知
func (s *TradeService) notifyTradeUpdate(trade *model.Trade) {
	update := &v1.TradeUpdate{
		TxHash:        trade.TxHash,
		Status:        string(trade.Status),
		Confirmations: trade.Confirmations,
	}

	message := &v1.WSMessage{
		Type: "trade_update",
		Data: update,
	}

	ws.Get().BroadcastToUser(trade.UserAddress, message)
}

// checkBalance 检查余额
func (s *TradeService) checkBalance(ctx context.Context, params *model.TradeParams) error {
	balance, err := s.uniswap.GetTokenBalance(ctx, params.TokenIn, params.UserAddress)
	if err != nil {
		return err
	}

	if balance.Cmp(params.AmountIn) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	return nil
}

// getQuote 获取报价
func (s *TradeService) getQuote(ctx context.Context, params *model.TradeParams) (*big.Int, error) {
	path := []common.Address{params.TokenIn, params.TokenOut}
	amounts, err := s.uniswap.GetAmountOut(ctx, params.AmountIn, path)
	if err != nil {
		return nil, err
	}

	return amounts[1], nil
}

// checkSlippage 检查滑点
func (s *TradeService) checkSlippage(minAmount, actualAmount *big.Int) error {
	if actualAmount.Cmp(minAmount) < 0 {
		return fmt.Errorf("price impact too high")
	}
	return nil
}

// sendTransaction 发送交易
func (s *TradeService) sendTransaction(ctx context.Context, params *model.TradeParams) (*types.Transaction, error) {
	path := []common.Address{params.TokenIn, params.TokenOut}
	return s.uniswap.SwapExactTokensForTokens(
		ctx,
		params.AmountIn,
		params.AmountOutMin,
		path,
		params.UserAddress,
		params.Deadline,
	)
}
