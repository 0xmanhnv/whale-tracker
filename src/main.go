package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"whale-tracker/src/blockchain"
	"whale-tracker/src/database"
	"whale-tracker/src/handles"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	pwd, _ := os.Getwd()

	// load .env file
	err := godotenv.Load(path.Join(pwd, ".env"))
	fmt.Println(path.Join(pwd, ".env"))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func bootstrap() {
	LoadEnv() // init load env
	database.CreateDBInstance()

	// crons.Init()
}

func InitBSCWhaleTracker() {
	// init client
	client := blockchain.GetClient(56)
	defer client.Close()

	/*
		load transfer logs of tokens
	*/
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	topics := [][]common.Hash{
		{logTransferSigHash},
	}

	// tokens
	addresses := []common.Address{
		common.HexToAddress("0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F"),
		common.HexToAddress("0x3244b3b6030f374bafa5f8f80ec2f06aaf104b64"),
	}

	fromBlock := 16618413
	toBlock := fromBlock + 1000

	logs := handles.LoadLogs(client, addresses, int64(fromBlock), int64(toBlock), topics)
	// end load logs

	// fmt.Println(logs)

	handles.LogHandleToWhale(client, logs, addresses)

	// holder := models.Holder{
	// 	Id:           primitive.NewObjectID(),
	// 	Address:      "0x58c34146316a9a60BFA5dA1d7F451e46BDd51215",
	// 	TokenAddress: "0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F",
	// }

	// services.CreateHolder(holder)

	// services.GetTokenHolderWithCovalenthqAPI("0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F", "56")
}

func main() {
	bootstrap()
	InitBSCWhaleTracker()
}
