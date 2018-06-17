package eth

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/services/swap"
)

type EthereumRequestAtom struct {
	context context.Context
	client  ethclient.Conn
	auth    *bind.TransactOpts
	to      common.Address
	binding *bindings.AtomSwap
	data    EthereumData
}

// NewEthereumRequestAtom returns a new Ethereum RequestAtom instance
func NewEthereumRequestAtom(context context.Context, client ethclient.Conn, auth *bind.TransactOpts, to common.Address, swapID [32]byte) (swap.AtomRequester, error) {
	contract, err := bindings.NewAtomSwap(client.AtomAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return &EthereumRequestAtom{}, err
	}
	return &EthereumRequestAtom{
		context: context,
		client:  client,
		auth:    auth,
		binding: contract,
		to:      to,
		data: EthereumData{
			SwapID: swapID,
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (ethAtom *EthereumRequestAtom) Initiate(hash [32]byte, value *big.Int, expiry int64) error {
	ethAtom.auth.Value = value
	ethAtom.data.HashLock = hash

	tx, err := ethAtom.binding.Initiate(ethAtom.auth, ethAtom.data.SwapID, ethAtom.to, hash, big.NewInt(expiry))
	ethAtom.auth.Value = big.NewInt(0)
	if err != nil {
		return err
	}
	_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	return err
}

// Refund an Atom swap by calling a function on ethereum
func (ethAtom *EthereumRequestAtom) Refund() error {
	tx, err := ethAtom.binding.Refund(ethAtom.auth, ethAtom.data.SwapID)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (ethAtom *EthereumRequestAtom) AuditSecret() ([32]byte, error) {
	return ethAtom.binding.AuditSecret(&bind.CallOpts{}, ethAtom.data.SwapID)
}

// Serialize serializes the atom details into a bytes array
func (ethAtom *EthereumRequestAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(ethAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (ethAtom *EthereumRequestAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &ethAtom.data)
}

// From returns the address of the sender
func (ethAtom *EthereumRequestAtom) From() []byte {
	return ethAtom.auth.From.Bytes()
}

// PriorityCode returns the priority code of the currency.
func (ethAtom *EthereumRequestAtom) PriorityCode() int64 {
	return 1
}
