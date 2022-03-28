package main

import (
	"fmt"
	"whale-tracker/src/handles"

	"github.com/ethereum/go-ethereum/core/types"
)

func GetTokenHolder(tokenAddress string, logs []types.Log) {
	fmt.Println(tokenAddress)
}

func main() {
	chainId := 97
	fromBlock := 16446905
	// maxRange := 5000
	tokenAddress := "0xE3233fdb23F1c27aB37Bd66A19a1f1762fCf5f3F"

	for {

		currentBlock := handles.GetCurrentBlock(chainId)
		fmt.Println(currentBlock)

		logs := handles.LoadLogs(chainId, int64(fromBlock), int64(currentBlock))

		fmt.Println(logs)
		GetTokenHolder(tokenAddress, logs)

		// for i := fromBlock; i < int(currentBlock)-fromBlock; i++ {
		// 	fmt.Println(i)
		// 	logs := handles.LoadLogs(chainId, int64(fromBlock), int64(currentBlock))

		// 	fmt.Println(logs)

		// 	GetTokenHolder(tokenAddress, logs)
		// }
		break

		// fromBlock = fromBlock + maxRange

		// if uint64(fromBlock) > currentBlock {
		// 	break
		// }
	}
	// logs := handles.LoadLogs(chainId, int64(fromBlock), int64(currentBlock))

	// fmt.Println(logs)
}
