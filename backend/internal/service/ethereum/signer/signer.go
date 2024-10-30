package signer

import (
	"backend/internal/service/ethereum"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogf/gf/v2/frame/g"
	"math/big"
)

// 交易签名服务
type TransactionSigner struct {
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}

func NewTransactionSigner(ctx context.Context) (*TransactionSigner, error) {
	cfg := g.Cfg()

	// 从配置或环境变量获取私钥
	privateKeyHex := cfg.MustGet(ctx, "ethereum.signer.privateKey").String()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)

	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	chainID := big.NewInt(int64(cfg.MustGet(ctx, "ethereum.chainID").Int()))

	return &TransactionSigner{
		privateKey: privateKey,
		chainID:    chainID,
	}, nil
}

func (s *TransactionSigner) SignAndSendTransaction(ctx context.Context, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	client := ethereum.GetHTTPClient()

	//1.获取nonce
	publicKey := s.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	//2.获取最新的gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}
	//3.创建交易
	tx := types.NewTransaction(nonce, to, value, 210000, gasPrice, data)
	//4.签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainID), s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	//5.发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}
	return signedTx, nil
}

// GetAddress 获取签名者地址
func (s *TransactionSigner) GetAddress() common.Address {
	publicKey := s.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
