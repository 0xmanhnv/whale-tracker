package handles

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"whale-tracker/src/blockchain"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func LoadLogs(chainId int, lastQueryBlock int64) []types.Log {
	client := blockchain.GetClient(chainId)
	currentBlock, _ := client.BlockNumber(context.Background())

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(lastQueryBlock),
		ToBlock:   big.NewInt(int64(currentBlock)),
	}

	// fmt.Println(query)

	logs, err := client.FilterLogs(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	return logs
}

func LogHandle(logs []types.Log, address string) {
	for _, log := range logs {
		// fmt.Printf("\nProcessed events from block %v", log.BlockNumber)

		if log.Address == common.HexToAddress(address) {
			// fmt.Println(common.HexToAddress(address))
			fmt.Println(log.Address)
		}
	}
}
