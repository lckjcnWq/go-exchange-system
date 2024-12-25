package model

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TradeParams struct {
	UserAddress  string
	TokenIn      common.Address
	TokenOut     common.Address
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Deadline     *big.Int
}
