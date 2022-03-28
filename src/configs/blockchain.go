package configs

type NetworkConfig struct {
	Name    string `json:"name"`
	ChainId int    `json:"chainId"`
	RPC     string `json:"rpc"`
}

type BlockChainConfig struct {
	Networks []NetworkConfig
}

func GetNetwork(chainId int) NetworkConfig {
	blockchainConfig := BlockChainConfig{
		Networks: []NetworkConfig{
			{
				Name:    "BSC",
				ChainId: 97,
				RPC:     "https://bsc-dataseed.binance.org/",
			},
		},
	}

	for _, network := range blockchainConfig.Networks {
		if network.ChainId == chainId {
			return network
		}
	}

	return blockchainConfig.Networks[1]
}
