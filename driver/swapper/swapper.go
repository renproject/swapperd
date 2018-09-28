package swapper

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

type Swapper struct {
	config.Config
	keystore.Keystore
}

func NewSwapper(config config.Config, keys keystore.Keystore) Swapper {
	return Swapper{
		Config:   config,
		Keystore: keys,
	}
}
