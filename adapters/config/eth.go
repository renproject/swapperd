package config

// Config are the parameters required to create an ethereum client
type EthereumNetworkConfig struct {
	Chain          string `json:"chain"`
	URL            string `json:"url"`
	AtomAddress    string `json:"atom_address"`
	NetworkAddress string `json:"network_address"`
	InfoAddress    string `json:"info_address"`
	WalletAddress  string `json:"wallet_address"`
}

type EthereumConfig struct {
	Ganache EthereumNetworkConfig `json:"ganache"`
	Ropsten EthereumNetworkConfig `json:"ropsten"`
	Kovan   EthereumNetworkConfig `json:"kovan"`
	Mainnet EthereumNetworkConfig `json:"mainnet"`
}

func (config *Config) GetEthereumConfig(chain string) (EthereumNetworkConfig, error) {
	config.mu.RLock()
	defer config.mu.RUnlock()
	switch chain {
	case "ganache":
		return config.Ethereum.Ganache, nil
	case "ropsten":
		return config.Ethereum.Ropsten, nil
	case "kovan":
		return config.Ethereum.Kovan, nil
	case "mainnet":
		return config.Ethereum.Mainnet, nil
	default:
		return EthereumNetworkConfig{}, ErrUnknownChain
	}
}

func (config *Config) SetEthereumConfig(chain string, ethereumConfig EthereumNetworkConfig) error {
	config.mu.Lock()
	defer config.mu.Unlock()
	switch chain {
	case "ganache":
		config.Ethereum.Ganache = ethereumConfig
	case "ropsten":
		config.Ethereum.Ropsten = ethereumConfig
	case "kovan":
		config.Ethereum.Kovan = ethereumConfig
	case "mainnet":
		config.Ethereum.Mainnet = ethereumConfig
	default:
		return ErrUnknownChain
	}
	config.Update()
	return nil
}
