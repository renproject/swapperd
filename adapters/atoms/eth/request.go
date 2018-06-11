package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/services/atom"
)

// EthereumData
type EthereumData struct {
	SwapID   [32]byte `json:"swap_id"`
	HashLock [32]byte `json:"hash_lock"`
}

type EthereumAtom struct {
	context context.Context
	client  ethclient.Connection
	auth    *bind.TransactOpts
	binding *bindings.Atom
	data    EthereumData
}

// NewEthereumRequestAtom returns a new Ethereum RequestAtom instance
func NewEthereumRequestAtom(context context.Context, client ethclient.Connection, auth *bind.TransactOpts, swapID [32]byte) (atom.RequestAtom, error) {
	contract, err := bindings.NewAtom(client.EthAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return &EthereumAtom{}, err
	}
	return &EthereumAtom{
		context: context,
		client:  client,
		auth:    auth,
		binding: contract,
		data: EthereumData{
			SwapID: swapID,
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Initiate(hash [32]byte, from []byte, to []byte, value *big.Int, expiry int64) error {
	ethAtom.auth.Value = value
	ethAtom.data.HashLock = hash

	if bytes.Compare(ethAtom.auth.From.Bytes(), from) != 0 {
		return errors.New("Refund Address Signing Address Mismatch")
	}

	tx, err := ethAtom.binding.Initiate(ethAtom.auth, ethAtom.data.SwapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
	ethAtom.auth.Value = big.NewInt(0)
	if err != nil {
		return err
	}
	_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	return err
}

// Refund an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Refund() error {
	tx, err := ethAtom.binding.Refund(ethAtom.auth, ethAtom.data.SwapID)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) AuditSecret() ([32]byte, error) {
	return ethAtom.binding.AuditSecret(&bind.CallOpts{}, ethAtom.data.SwapID)
}

// Serialize serializes the atom details into a bytes array
func (ethAtom *EthereumAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(ethAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (ethAtom *EthereumAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &ethAtom.data)
}

// From returns the address of the sender
func (ethAtom *EthereumAtom) From() []byte {
	return ethAtom.auth.From.Bytes()
}

// PriorityCode returns the priority code of the currency.
func (ethAtom *EthereumAtom) PriorityCode() int64 {
	return 1
}
