package v1

// WSMessage WebSocket消息结构
type WSMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Channel string      `json:"channel,omitempty"`
}

// TradeUpdate 交易更新消息
type TradeUpdate struct {
	TxHash        string `json:"txHash"`
	Status        string `json:"status"`
	Confirmations uint64 `json:"confirmations"`
}

// PriceUpdate 价格更新消息
type PriceUpdate struct {
	Pair  string `json:"pair"`
	Price string `json:"price"`
}
