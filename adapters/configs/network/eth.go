package network

// EthereumNetwork are the parameters required to create an ethereum client
type EthereumNetwork struct {
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

func (network *Config) GetEthereumNetwork() EthereumNetwork {
	network.mu.RLock()
	defer network.mu.RUnlock()
	return network.Ethereum
}

func (config *Config) SetEthereumNetwork(ethereumConfig EthereumNetwork) {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.Ethereum = ethereumConfig
	config.Update()
}
