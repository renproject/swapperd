package ethnetwork

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	client "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/services/swap"
)

type ethereumNetwork struct {
	conn client.Conn
	auth *bind.TransactOpts
	net  *bindings.AtomicInfo
	ctx  context.Context
}

func NewEthereumNetwork(conn client.Conn, auth *bind.TransactOpts) (swap.Network, error) {
	net, err := bindings.NewAtomicInfo(conn.InfoAddress(), bind.ContractBackend(conn.Client()))
	if err != nil {
		return &ethereumNetwork{}, err
	}
	return &ethereumNetwork{
		conn: conn,
		auth: auth,
		net:  net,
		ctx:  context.Background(),
	}, nil
}

func (net *ethereumNetwork) SendSwapDetails(orderID [32]byte, swapDetails []byte) error {
	net.auth.GasLimit = 3000000
	tx, err := net.net.SubmitDetails(net.auth, orderID, swapDetails)
	if err != nil {
		return err
	}

	_, err = net.conn.PatchedWaitMined(net.ctx, tx)
	return err
}

func (net *ethereumNetwork) ReceiveSwapDetails(orderID [32]byte) ([]byte, error) {
	for {
		swap, err := net.net.SwapDetails(&bind.CallOpts{}, orderID)
		if err != nil {
			return []byte{}, fmt.Errorf("Failed to get swap details: %v", err)
		}
		if bytes.Compare(swap, []byte{}) != 0 {
			return net.net.SwapDetails(&bind.CallOpts{}, orderID)
		}
	}
}
