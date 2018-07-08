package config

// EthereumConfig are the parameters required to create an ethereum client
type EthereumConfig struct {
	Chain                   string `json:"chain"`
	URL                     string `json:"url"`
	AtomAddress             string `json:"atom_address"`
	InfoAddress             string `json:"info_address"`
	WalletAddress           string `json:"wallet_address"`
	RenExTokens             string `json:"renex_tokens"`
	RenExBalances           string `json:"renex_balances"`
	RewardVault             string `json:"reward_vault"`
	DarkNodeRegistryAddress string `json:"dnr_address"`
	RepublicTokenAddress    string `json:"ren_address"`
	OrderBookAddress        string `json:"ob_address"`
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
