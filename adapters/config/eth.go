package config

// Config are the parameters required to create an ethereum client
type EthereumConfig struct {
	Chain          string `json:"chain"`
	URL            string `json:"url"`
	AtomAddress    string `json:"atom_address"`
	NetworkAddress string `json:"network_address"`
	InfoAddress    string `json:"info_address"`
	WalletAddress  string `json:"wallet_address"`
}

func (config *Config) GetEthereumConfig() EthereumConfig {
	config.mu.RLock()
	defer config.mu.RUnlock()
	return config.Ethereum
}

func (config *Config) SetEthereumConfig(ethereumConfig EthereumConfig) {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.Ethereum = ethereumConfig
	config.Update()
}
