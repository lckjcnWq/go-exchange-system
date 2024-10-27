package v1

import "github.com/gogf/gf/v2/frame/g"

type WalletCheckReq struct {
	g.Meta  `path:"/wallet/check" method:"post"`
	Address string `json:"address" v:"required"`
	ChainId int64  `json:"chainId" v:"required"`
}

type WalletCheckRes struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}
