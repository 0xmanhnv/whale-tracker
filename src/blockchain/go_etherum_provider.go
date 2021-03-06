package blockchain

import (
	"log"
	"whale-tracker/src/configs"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetClient(chainId int) *ethclient.Client {
	network := configs.GetNetwork(chainId)

	client, err := ethclient.Dial(network.RPC)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	return client
}
