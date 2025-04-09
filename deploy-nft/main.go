package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/data-market/deploy-nft/go-contracts/nft"
	"github.com/ethereum/go-ethereum/ethclient"

	com "github.com/memoio/contractsv2/common"
)

var (
	allGas = make([]uint64, 0)
)

// go run main.go -eth=dev -sk=0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4
func main() {
	// env select
	env := flag.String("eth", "dev", "eth api Address;") //dev test or product
	// sk for send tx
	sk := flag.String("sk", "", "signature for sending transaction")

	flag.Parse()

	chain := *env

	// get endpoint
	_, endPoint := com.GetInsEndPointByChain(chain)
	fmt.Println()
	fmt.Println("endPoint:", endPoint)
	adminSk := *sk

	fmt.Println()

	// connect chain
	client, err := ethclient.DialContext(context.TODO(), endPoint)
	if err != nil {
		log.Fatal(err)
	}

	// make auth to send transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		chainID = big.NewInt(666)
	}

	// make auth for admin
	adminAuth, err := com.MakeAuth(chainID, adminSk)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("deploying nft")

	// deploy nft
	nftAddr, tx, nftIns, err := nft.DeployNft(adminAuth, client, "DMNFT", "DMNFT")
	if err != nil {
		log.Fatal("deploy nft err:", err)
	}

	fmt.Println("nftAddr: ", nftAddr.Hex())
	go com.PrintGasUsed(endPoint, tx.Hash(), "deploy nft gas:", &allGas)
	_ = nftIns

}
