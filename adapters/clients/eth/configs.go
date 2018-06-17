package ethclient

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Config are the parameters required to create an ethereum client
type Config struct {
	Chain          string `json:"chain"`
	URL            string `json:"url"`
	AtomAddress    string `json:"atom_address"`
	NetworkAddress string `json:"network_address"`
	InfoAddress    string `json:"info_address"`
	WalletAddress  string `json:"wallet_address"`
}

// ConfigFile has the configs of different ethereum chains
type ConfigFile struct {
	Ganache Config `json:"ganache"`
	Ropsten Config `json:"ropsten"`
	Kovan   Config `json:"kovan"`
	Mainnet Config `json:"mainnet"`
}

var configPath = path.Join(os.Getenv("GOPATH"), "/src/github.com/republicprotocol/atom-go/adapters/clients/secrets/ethConfig.json")

func readConfig(chain string) (Config, error) {
	var config Config
	var configFile ConfigFile

	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	json.Unmarshal(raw, &configFile)

	switch chain {
	case "ganache":
		config = configFile.Ganache
	case "ropsten":
		config = configFile.Ropsten
	case "kovan":
		config = configFile.Kovan
	case "mainnet":
		config = configFile.Mainnet
	}
	return config, nil
}

func writeConfig(config Config) error {
	var configFile ConfigFile

	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	json.Unmarshal(raw, &configFile)

	switch config.Chain {
	case "ganache":
		configFile.Ganache = config
	case "ropsten":
		configFile.Ropsten = config
	case "kovan":
		configFile.Kovan = config
	case "mainnet":
		configFile.Mainnet = config
	}

	data, err := json.Marshal(configFile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, data, 700)
}
