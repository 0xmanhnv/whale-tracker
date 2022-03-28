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
	"github.com/ethereum/go-ethereum/crypto"
)

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
	TxHash common.Hash
}

func GetCurrentBlock(chainId int) uint64 {
	client := blockchain.GetClient(chainId)
	currentBlock, _ := client.BlockNumber(context.Background())

	return currentBlock
}

/*
Load logs from token address
@param: chainId - network chain id
@param: tokenAddress - token address on evm or bsc
@param: fromBlock - from block
@param: toBlock - to block

return []types.Log
*/
func LoadLogs(chainId int, tokenAddress string, fromBlock int64, toBlock int64) []types.Log {
	client := blockchain.GetClient(chainId)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: []common.Address{
			common.HexToAddress(tokenAddress),
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	return logs
}

func LogHandle(logs []types.Log) {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	var transferHash []common.Hash

	for _, vLog := range logs {

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			transferEvent.TxHash = vLog.TxHash

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("TxHash: %s\n", transferEvent.TxHash)

			transferHash = append(transferHash, transferEvent.TxHash)

		case logApprovalSigHash.Hex():
			//TODO: handle fmt.Printf("Log Name: Approval\n")
		}

		fmt.Printf("\n\n")

		// fmt.Printf("\nProcessed events from block %v", log.BlockNumber)

		// if vlog.Address == common.HexToAddress(address) {
		// 	// fmt.Println(common.HexToAddress(address))
		// 	fmt.Println(vlog.Address)
		// }
	}

	fmt.Println("LogHandle Done")
}
