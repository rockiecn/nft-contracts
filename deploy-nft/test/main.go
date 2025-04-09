package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/data-market/deploy-nft/go-contracts/nft"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// dev chain
	devChainURL  = "https://devchain.metamemo.one:8501"         // 开发链地址
	contractAddr = "0xE570a8460D26F35Ec287F4f92B659f8604838851" // 合约地址

	adminSk   = "0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4" // 测试用私钥（不要提交到代码库！）
	adminAddr = "0x1c111472F298E4119150850c198C657DA1F8a368"
	toAddr    = "0x4398e134036b85d22E04693eC20641A4a4d18016"
)

func main() {
	// 1. 连接开发链
	client, err := ethclient.Dial(devChainURL)
	if err != nil {
		log.Fatalf("Failed to connect to dev chain: %v", err)
	}
	defer client.Close()

	// 2. 加载合约
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := nft.NewNft(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to load contract: %v", err)
	}

	// 3. 准备认证信息
	privKey, err := crypto.HexToECDSA(adminSk)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create auth: %v", err)
	}

	// 4. 测试 mint 功能
	fmt.Println("Testing mint...")
	//tokenId := big.NewInt(1) // 要铸造的NFT ID
	// mint 3 nfts
	// first
	mintTx, err := instance.Mint(auth, fromAddress)
	if err != nil {
		log.Fatalf("Mint failed: %v", err)
	}
	fmt.Printf("Mint transaction sent: %s\n", mintTx.Hash().Hex())
	// second
	mintTx, err = instance.Mint(auth, fromAddress)
	if err != nil {
		log.Fatalf("Mint failed: %v", err)
	}
	fmt.Printf("Mint transaction sent: %s\n", mintTx.Hash().Hex())
	// third
	mintTx, err = instance.Mint(auth, fromAddress)
	if err != nil {
		log.Fatalf("Mint failed: %v", err)
	}
	fmt.Printf("Mint transaction sent: %s\n", mintTx.Hash().Hex())

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	// 5. 测试 listNft 功能
	fmt.Println("\nTesting get tokens of owner...")
	tokenIDs, err := instance.TokensOfOwner(&bind.CallOpts{}, common.HexToAddress(adminAddr))
	if err != nil {
		log.Fatalf("get tokens failed: %v", err)
	}
	fmt.Println("token list: ", tokenIDs)

	// 6. 测试 approve 功能
	fmt.Println("\nTesting approve...")
	operator := common.HexToAddress(toAddr) // 替换为实际地址
	tokenId := big.NewInt(1)
	approveTx, err := instance.Approve(auth, operator, tokenId)
	if err != nil {
		log.Fatalf("Approve failed: %v", err)
	}
	fmt.Printf("Approve transaction sent: %s\n", approveTx.Hash().Hex())

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	// // 7. 测试 shareNft 功能
	// fmt.Println("\nTesting shareNft...")
	// recipient := common.HexToAddress("0xRecipientAddress") // 替换为实际地址
	// shareTx, err := instance.ShareNft(auth, tokenId, recipient)
	// if err != nil {
	// 	log.Fatalf("ShareNft failed: %v", err)
	// }
	// fmt.Printf("ShareNft transaction sent: %s\n", shareTx.Hash().Hex())

	// 8. 验证状态变化
	fmt.Println("\nVerifying state changes...")

	// 验证所有者
	owner, err := instance.OwnerOf(nil, tokenId)
	if err != nil {
		log.Fatalf("Failed to get owner: %v", err)
	}
	fmt.Printf("Token %d owner: %s\n", tokenId, owner.Hex())

	// 验证批准
	approved, err := instance.GetApproved(nil, tokenId)
	if err != nil {
		log.Fatalf("Failed to get approved address: %v", err)
	}
	fmt.Printf("Token %d approved: %s\n", tokenId, approved.Hex())
}
