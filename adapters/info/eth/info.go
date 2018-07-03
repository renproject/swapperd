package eth

import (
	"bytes"
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	client "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/services/swap"
)

type ethereumAtomInfo struct {
	conn client.Conn
	auth *bind.TransactOpts
	info *bindings.AtomInfo
	ctx  context.Context
}

func NewEtereumAtomInfo(conn client.Conn, auth *bind.TransactOpts) (swap.Info, error) {
	info, err := bindings.NewAtomInfo(conn.InfoAddress(), bind.ContractBackend(conn.Client()))
	if err != nil {
		return &ethereumAtomInfo{}, err
	}
	return &ethereumAtomInfo{
		conn: conn,
		auth: auth,
		info: info,
		ctx:  context.Background(),
	}, nil
}

func (info *ethereumAtomInfo) SetOwnerAddress(orderID [32]byte, address []byte) error {
	tx, err := info.info.SetOwnerAddress(info.auth, orderID, address)
	if err != nil {
		return err
	}
	_, err = info.conn.PatchedWaitMined(info.ctx, tx)
	return err
}

func (info *ethereumAtomInfo) GetOwnerAddress(orderID [32]byte) ([]byte, error) {

	for {
		owner, err := info.info.GetOwnerAddress(&bind.CallOpts{}, orderID)
		if err != nil {
			break
		}
		if bytes.Compare(owner, []byte{}) != 0 {
			return owner, nil
		}
	}
	return []byte{}, errors.New("Failed to get owner details")
}
