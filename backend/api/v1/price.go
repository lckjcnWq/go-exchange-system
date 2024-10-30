package v1

import "github.com/gogf/gf/v2/frame/g"

type GetPriceReq struct {
	g.Meta   `path:"/price" method:"get"`
	TokenIn  string `v:"required"`
	TokenOut string `v:"required"`
	AmountIn string `v:"required"`
}

type GetPriceRes struct {
	AmountOut string `json:"amountOut"`
	Rate      string `json:"rate"`
}
