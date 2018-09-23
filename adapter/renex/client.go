package renex

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
)

type Conn struct {
	Network         string
	Client          *ethclient.Client
	RenExSettlement common.Address
	RenExOrderbook  common.Address
}

// NewConnWithConfig creates a new ethereum connection with the given config
// file.
func NewConnWithConfig(config config.Config) (Conn, error) {
	return NewConn(config.Ethereum.URL, config.Ethereum.Network, config.RenEx.Settlement, config.RenEx.Orderbook)
}

// NewConn creates a new ethereum connection with the given config parameters.
func NewConn(url, network, settlementAddress, orderbookAddress string) (Conn, error) {
	ethclient, err := ethclient.Dial(url)
	if err != nil {
		return Conn{}, err
	}
	return Conn{
		Client:          ethclient,
		Network:         network,
		RenExSettlement: common.HexToAddress(settlementAddress),
		RenExOrderbook:  common.HexToAddress(orderbookAddress),
	}, nil
}
