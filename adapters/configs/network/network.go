package network

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

type Config struct {
	Ethereum EthereumNetwork `json:"ethereum"`
	Bitcoin  BitcoinNetwork  `json:"bitcoin"`

	mu   *sync.RWMutex
	path string
}

var ErrUnSupportedPriorityCode = errors.New("Un Supported Priority Code")

func LoadNetwork(path string) (Config, error) {
	var network Config
	network.path = path
	network.mu = new(sync.RWMutex)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return network, err
	}
	json.Unmarshal(raw, &network)
	return network, nil
}

func (network *Config) Update() error {
	data, err := json.Marshal(network)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(network.path, data, 700)
}
