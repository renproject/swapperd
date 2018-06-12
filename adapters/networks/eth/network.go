package ethnetwork

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	client "github.com/republicprotocol/atom-go/adapters/ethclient"
	"github.com/republicprotocol/atom-go/services/network"
)

type EthereumNetwork struct {
	conn client.Connection
	auth *bind.TransactOpts
	net  *Network
	ctx  context.Context
}

func NewEthereumNetwork(context context.Context, conn client.Connection, auth *bind.TransactOpts) (network.Network, error) {
	net, err := NewNetwork(conn.NetworkAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		return &EthereumNetwork{}, err
	}
	return &EthereumNetwork{
		conn: conn,
		auth: auth,
		net:  net,
		ctx:  context,
	}, nil
}

func (net EthereumNetwork) Send(orderID [32]byte, swapDetails []byte) error {
	tx, err := net.net.SubmitDetails(net.auth, orderID, swapDetails)
	if err != nil {
		return err
	}
	_, err = net.conn.PatchedWaitMined(net.ctx, tx)
	return err
}

func (net EthereumNetwork) Recieve(orderID [32]byte) ([]byte, error) {
	return net.net.SwapDetails(&bind.CallOpts{}, orderID)
}
