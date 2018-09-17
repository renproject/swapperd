package config

// Config is the the global config object
type Config struct {
	Version             string          `json:"version"`
	SupportedCurrencies []string        `json:"supportedCurrencies"`
	AuthorizedAddresses []string        `json:"authorizedAddresses"`
	StoreLocation       string          `json:"storeLocation"`
	Ethereum            EthereumNetwork `json:"ethereum"`
	Bitcoin             BitcoinNetwork  `json:"bitcoin"`
	RenEx               RenExNetwork    `json:"renex"`
}

// EthereumNetwork is the ethereum specific config object
type EthereumNetwork struct {
	Network string `json:"network"`
	URL     string `json:"url"`
}

// BitcoinNetwork is the bitcoin specific config object
type BitcoinNetwork struct {
	Network string `json:"network"`
	URL     string `json:"url"`
}

// RenExNetwork is the renex specific config object
type RenExNetwork struct {
	Network    string `json:"network"`
	Watchdog   string `json:"watchdog"`
	Ingress    string `json:"ingress"`
	Settlement string `json:"settlement"`
	Swapper    string `json:"swapper"`
	Orderbook  string `json:"orderbook"`
}
