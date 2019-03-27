package wallet

// Testnet is the Swapperd's testnet config object
var Testnet = Config{
	Bitcoin: BlockchainConfig{
		Network: Network{
			Name: "testnet",
		},
	},
	ZCash: BlockchainConfig{
		Network: Network{
			Name: "testnet",
		},
	},
	Ethereum: BlockchainConfig{
		Network: Network{
			Name: "kovan",
			URL:  "https://kovan.infura.io",
		},
	},
}

// Mainnet is the Swapperd's mainnet config object
var Mainnet = Config{
	Bitcoin: BlockchainConfig{
		Network: Network{
			Name: "mainnet",
		},
	},
	ZCash: BlockchainConfig{
		Network: Network{
			Name: "mainnet",
		},
	},
	Ethereum: BlockchainConfig{
		Network: Network{
			Name: "mainnet",
			URL:  "https://mainnet.infura.io",
		},
	},
}
