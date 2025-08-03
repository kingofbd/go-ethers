package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	sendTransfer()
}

// 查询区块
func searchBlock() {
	// 获取到ethers client
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/b15360719eb64c42adc56e6b5cc6042d")
	if err != nil {
		log.Fatal(err)
	}
	// 获取到最新的区块头，打印区块号
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	blockNumber := header.Number
	fmt.Println(blockNumber) // 8900564
	// 查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	// 区块哈希
	fmt.Println(block.Hash().Hex()) // 0x90595d296233b8a68ce65d52ca3651f11aab67d19f4d48d11b1d92522ce19a09
	// 区块时间戳
	fmt.Println(block.Time()) // 1754185752
	// 区块交易数量
	fmt.Println(len(block.Transactions())) // 164
}

// 发送交易
func sendTransfer() {
	// 获取到ethers client
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/b15360719eb64c42adc56e6b5cc6042d")
	if err != nil {
		log.Fatal(err)
	}
	// 加载私钥（提交代码中私钥进行了保护）
	privateKey, err := crypto.HexToECDSA("xxxxx")
	if err != nil {
		log.Fatal(err)
	}
	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 获取帐户交易的随机数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 设置转账金额(0.01eth)
	value := big.NewInt(10000000000000000)
	// 设置gasLimit
	gasLimit := uint64(21000)
	// 获得平均燃气价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 接收方的地址
	toAddress := common.HexToAddress("0xDb572320840Bfc710e04FD1823A07408399c83c5")
	// 生成交易事务
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	// 对交易事务进行签名
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // 0xdd15c8313ee9c7eb977c465ad8cac095d02ea445f7b6e4016167343df8604fa0
	// https://sepolia.etherscan.io/tx/0xdd15c8313ee9c7eb977c465ad8cac095d02ea445f7b6e4016167343df8604fa0
}
