package cronjob

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"whale-tracker/src/blockchain"
	"whale-tracker/src/handles"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func WhaleTrackerCron() {
	client := blockchain.GetClient(56)
	defer client.Close()

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	lastBlock := os.Getenv("LAST_BLOCK")
	fmt.Println("Lastblock: ----------" + lastBlock)
	fromBlock, err := strconv.Atoi(lastBlock)

	toBlock := header.Number.Int64()
	fmt.Println("Toblcok: ----------" + header.Number.String())

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
		common.HexToAddress("0x0eb3a705fc54725037cc9e008bdede697f62f335"),
		common.HexToAddress("0x3244b3b6030f374bafa5f8f80ec2f06aaf104b64"),
	}

	// fromBlock := 16618405
	// toBlock := fromBlock + 1000

	logs := handles.LoadLogs(client, addresses, int64(fromBlock), toBlock, topics)
	// end load logs

	// fmt.Println(logs)

	handles.LogHandleToWhale(client, logs, addresses)
	fmt.Println("Reset lastblock")
	os.Setenv("LAST_BLOCK", header.Number.String())
}
