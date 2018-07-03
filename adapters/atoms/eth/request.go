package eth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

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
	binding *bindings.AtomSwap
	data    EthereumData
}

// NewEthereumRequestAtom returns a new Ethereum RequestAtom instance
func NewEthereumRequestAtom(client ethclient.Conn, auth *bind.TransactOpts) (swap.AtomRequester, error) {
	contract, err := bindings.NewAtomSwap(client.AtomAddress(), bind.ContractBackend(client.Client()))
	if err != nil {
		return &EthereumRequestAtom{}, err
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])

	return &EthereumRequestAtom{
		context: context.Background(),
		client:  client,
		auth:    auth,
		binding: contract,
		data: EthereumData{
			SwapID: swapID,
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (ethAtom *EthereumRequestAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	ethAtom.auth.Value = value
	ethAtom.data.HashLock = hash

	tx, err := ethAtom.binding.Initiate(ethAtom.auth, ethAtom.data.SwapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
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
	for start := time.Now(); time.Since(start) < 24*time.Hour; {
		secret, err := ethAtom.binding.AuditSecret(&bind.CallOpts{}, ethAtom.data.SwapID)
		fmt.Println("Audit Secret Tried", secret)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if secret == [32]byte{} {
			time.Sleep(2 * time.Second)
			continue
		}
		return secret, nil
	}
	return [32]byte{}, errors.New("Failed to audit the secret")
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
