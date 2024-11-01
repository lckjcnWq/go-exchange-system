package v1

import "github.com/gogf/gf/v2/frame/g"

// CreateTradeReq 创建交易请求
type CreateTradeReq struct {
	g.Meta       `path:"/trade" method:"post"`
	TokenIn      string `json:"tokenIn" v:"required"`
	TokenOut     string `json:"tokenOut" v:"required"`
	AmountIn     string `json:"amountIn" v:"required"`
	AmountOutMin string `json:"amountOutMin" v:"required"`
	Deadline     uint64 `json:"deadline" v:"required"`
}

type CreateTradeRes struct {
	TxHash string `json:"txHash"`
}

// GetTradesReq 获取用户交易列表请求
type GetTradesReq struct {
	g.Meta `path:"/trades" method:"get"`
	Page   int `json:"page" v:"required"`
	Size   int `json:"size" v:"required"`
}

type GetTradesRes struct {
	List  []*TradeInfo `json:"list"`
	Total int          `json:"total"`
}

// TradeInfo 交易信息
type TradeInfo struct {
	TxHash        string `json:"txHash"`
	TokenIn       string `json:"tokenIn"`
	TokenOut      string `json:"tokenOut"`
	AmountIn      string `json:"amountIn"`
	AmountOut     string `json:"amountOut"`
	Status        string `json:"status"`
	Confirmations uint64 `json:"confirmations"`
	CreatedAt     string `json:"createdAt"`
}
