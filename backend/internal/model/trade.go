package model

import "time"

type TradeStatus string

const (
	TradeStatusPending   TradeStatus = "pending"
	TradeStatusConfirmed TradeStatus = "confirmed"
	TradeStatusFailed    TradeStatus = "failed"
)

type Trade struct {
	Id            uint64      `json:"id"`
	TxHash        string      `json:"txHash"`
	UserAddress   string      `json:"userAddress"`
	TokenIn       string      `json:"tokenIn"`
	TokenOut      string      `json:"tokenOut"`
	AmountIn      string      `json:"amountIn"`
	AmountOut     string      `json:"amountOut"`
	Status        TradeStatus `json:"status"`
	BlockNumber   uint64      `json:"blockNumber"`
	Confirmations uint64      `json:"confirmations"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}
