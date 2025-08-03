package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-ethers/counter"
	"log"
	"math/big"
)

func main() {
	// 首先编写合约，见同级sol目录下的Counter合约
	// 然后使用命令获得合约的字节码文件和abi文件
	// 1.solcjs --bin sol/Counter.so
	// 2.solcjs --abi sol/Counter.sol
	// 使用 abigen 工具根据 bin 文件和 abi 文件，生成 go 代码  abigen --bin=sol_Counter_sol_Counter.bin --abi=sol_Counter_sol_Counter.abi --pkg=counter --out=counter.go

	deployContract()
	doContract()
}

// 部署合约
func deployContract() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/b15360719eb64c42adc56e6b5cc6042d")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("xxxx")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address.Hex())   // 部署的合约地址 0xb30c7063797c9A21bCD9D01b8B9ea1B86F5680f3
	fmt.Println(tx.Hash().Hex()) // 交易hash 0x3bb4e63185ad1d3f2542e7a584c45ab0a4663b50bfa33fffdd91b953c9f99551
	_ = instance
}

// 加载合约并执行方法
func doContract() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/b15360719eb64c42adc56e6b5cc6042d")
	if err != nil {
		log.Fatal(err)
	}
	counterContract, err := counter.NewCounter(common.HexToAddress("0xb30c7063797c9A21bCD9D01b8B9ea1B86F5680f3"), client)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("xxxx")
	if err != nil {
		log.Fatal(err)
	}
	// 初始化交易opt实例
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}
	// 调用合约方法（计数器+1）
	tx, err := counterContract.Increment(opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("交易已发送，哈希: %s\n", tx.Hash().Hex()) // 交易已发送，哈希: 0xc15d09e5b5729cf67d05a1c7f0044304c3a5eac9febd7e5f91c93b2549a562a3
	// 等待交易被执行
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		// 调用合约方法（查询计数器）
		latestCount, err := counterContract.GetCount(nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("当前计数器的值: %d\n", latestCount) // 当前计数器的值: 1
	} else {
		fmt.Println("交易失败")
	}
}
