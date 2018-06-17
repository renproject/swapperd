package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/services/swap"
)

type EthereumResponseAtom struct {
	context context.Context
	client  ethclient.Conn
	auth    *bind.TransactOpts
	binding *bindings.AtomSwap
	data    EthereumData
}

// NewEthereumResponseAtom returns a new Ethereum ResponseAtom instance
func NewEthereumResponseAtom(context context.Context, client ethclient.Conn, auth *bind.TransactOpts) (swap.AtomResponder, error) {
	contract, err := bindings.NewAtomSwap(client.AtomAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return &EthereumResponseAtom{}, err
	}
	return &EthereumResponseAtom{
		context: context,
		client:  client,
		auth:    auth,
		binding: contract,
	}, nil
}

// Redeem an Atom swap by calling a function on ethereum
func (eth *EthereumResponseAtom) Redeem(secret [32]byte) error {
	tx, err := eth.binding.Redeem(eth.auth, eth.data.SwapID, secret)
	if err == nil {
		_, err = eth.client.PatchedWaitMined(eth.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (eth *EthereumResponseAtom) Audit(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
	auditReport, err := eth.binding.Audit(&bind.CallOpts{}, eth.data.SwapID)
	if err != nil {
		return err
	}

	if hash != auditReport.SecretLock {
		return errors.New("HashLock mismatch")
	}

	if bytes.Compare(to, auditReport.To.Bytes()) != 0 {
		return errors.New("To Address mismatch")
	}

	// if value.Cmp(auditReport.Value) > 0 {
	// 	return errors.New("Value mismatch")
	// }

	// if expiry > (auditReport.Timelock.Int64() - time.Now().Unix()) {
	// 	return errors.New("Expiry mismatch")
	// }

	return nil
}

// Serialize serializes the atom details into a bytes array
func (eth *EthereumResponseAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(eth.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (eth *EthereumResponseAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &eth.data)
}

// From returns the address of the sender
func (eth *EthereumResponseAtom) From() []byte {
	return eth.auth.From.Bytes()
}

// PriorityCode returns the priority code of the currency.
func (eth *EthereumResponseAtom) PriorityCode() int64 {
	return 1
}

// GetSecretHash returns the Secret Hash of the atom.
func (eth *EthereumResponseAtom) GetSecretHash() [32]byte {
	return eth.data.HashLock
}
