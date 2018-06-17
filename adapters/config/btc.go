package config

import "errors"

var (
	ErrUnknownChain = errors.New("unknown chain")
)

type BitcoinNetworkConfig struct {
	Chain    string `json:"chain"`
	User     string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

type BitcoinConfig struct {
	Regtest BitcoinNetworkConfig `json:"regtest"`
	Testnet BitcoinNetworkConfig `json:"testnet"`
	Mainnet BitcoinNetworkConfig `json:"mainnet"`
}

func (config *Config) GetBitcoinConfig(chain string) (BitcoinNetworkConfig, error) {
	config.mu.RLock()
	defer config.mu.RUnlock()
	switch chain {
	case "regtest":
		return config.Bitcoin.Regtest, nil
	case "testnet":
		return config.Bitcoin.Testnet, nil
	case "mainnet":
		return config.Bitcoin.Mainnet, nil
	default:
		return BitcoinNetworkConfig{}, ErrUnknownChain
	}
}

func (config *Config) SetBitcoinConfig(chain string, bitcoinConfig BitcoinNetworkConfig) error {
	config.mu.Lock()
	defer config.mu.Unlock()
	switch chain {
	case "regtest":
		config.Bitcoin.Regtest = bitcoinConfig
	case "testnet":
		config.Bitcoin.Testnet = bitcoinConfig
	case "mainnet":
		config.Bitcoin.Mainnet = bitcoinConfig
	default:
		return ErrUnknownChain
	}
	config.Update()
	return nil
}
