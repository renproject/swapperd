package config

type BitcoinConfig struct {
	Chain    string `json:"chain"`
	User     string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

func (config *Config) GetBitcoinConfig() BitcoinConfig {
	config.mu.RLock()
	defer config.mu.RUnlock()
	return config.Bitcoin
}

func (config *Config) SetBitcoinConfig(bitcoinConfig BitcoinConfig) {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.Bitcoin = bitcoinConfig
	config.Update()
}
