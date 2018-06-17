package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Config struct {
	Ethereum EthereumConfig `json:"ethereum"`
	Bitcoin  BitcoinConfig  `json:"bitcoin"`

	mu   *sync.RWMutex
	path string
}

func LoadConfig(path string) (Config, error) {
	var config Config
	config.path = path
	config.mu = new(sync.RWMutex)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}
	json.Unmarshal(raw, &config)
	return config, nil
}

func (config *Config) Update() error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(config.path, data, 700)
}
