package network

type BitcoinNetwork struct {
	Network  string `json:"network"`
	User     string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

func (network *Config) GetBitcoinNetwork() BitcoinNetwork {
	network.mu.RLock()
	defer network.mu.RUnlock()
	return network.Bitcoin
}

func (network *Config) SetBitcoinNetwork(bitcoinNetwork BitcoinNetwork) {
	network.mu.Lock()
	defer network.mu.Unlock()
	network.Bitcoin = bitcoinNetwork
	network.Update()
}
