package eth

import (
	"bytes"
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

// EthereumData
type EthereumData struct {
	SwapID   [32]byte `json:"swap_id"`
	HashLock [32]byte `json:"hash_lock"`
}

type EthereumAtom struct {
	context context.Context
	client  ethclient.Conn
	key     swap.Key
	binding *bindings.AtomSwap
	data    EthereumData
}

// NewEthereumAtom returns a new Ethereum RequestAtom instance
func NewEthereumAtom(client ethclient.Conn, key swap.Key) (swap.Atom, error) {

	contract, err := bindings.NewAtomSwap(client.AtomAddress(), bind.ContractBackend(client.Client()))
	if err != nil {
		return &EthereumAtom{}, err
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])

	return &EthereumAtom{
		context: context.Background(),
		client:  client,
		key:     key,
		binding: contract,
		data: EthereumData{
			SwapID: swapID,
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	auth := bind.NewKeyedTransactor(ethAtom.key.GetKey())
	auth.Value = value
	auth.GasLimit = 3000000
	ethAtom.data.HashLock = hash
	tx, err := ethAtom.binding.Initiate(auth, ethAtom.data.SwapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
	if err != nil {
		return err
	}
	_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	return err
}

// Refund an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Refund() error {
	auth := bind.NewKeyedTransactor(ethAtom.key.GetKey())
	auth.GasLimit = 3000000
	tx, err := ethAtom.binding.Refund(auth, ethAtom.data.SwapID)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// Redeem an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Redeem(secret [32]byte) error {
	auth := bind.NewKeyedTransactor(ethAtom.key.GetKey())
	auth.GasLimit = 3000000
	tx, err := ethAtom.binding.Redeem(auth, ethAtom.data.SwapID, secret)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Audit(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
	for start := time.Now(); time.Since(start) < 24*time.Hour; {

		auditReport, err := ethAtom.binding.Audit(&bind.CallOpts{}, ethAtom.data.SwapID)
		if auditReport.SecretLock == [32]byte{} {
			time.Sleep(2 * time.Second)
			continue
		}

		if err != nil {
			return err
		}

		if hash != auditReport.SecretLock {
			ethAtom.data.HashLock = auditReport.SecretLock
			// return fmt.Errorf("HashLock mismatch %v %v", hash, auditReport.SecretLock)
		}

		if bytes.Compare(to, auditReport.To.Bytes()) != 0 {
			return fmt.Errorf("Eth: To Address mismatch %v %v", to, auditReport.To.Bytes())
		}

		// if value.Cmp(auditReport.Value) > 0 {
		// 	return errors.New("Value mismatch")
		// }

		// if expiry > (auditReport.Timelock.Int64() - time.Now().Unix()) {
		// 	return errors.New("Expiry mismatch")
		// }
		return nil
	}
	return errors.New("Audit failed")
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) AuditSecret() ([32]byte, error) {
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
func (ethAtom *EthereumAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(ethAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (ethAtom *EthereumAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &ethAtom.data)
}

// From returns the address of the sender
func (ethAtom *EthereumAtom) From() ([]byte, error) {
	return ethAtom.key.GetAddress()
}

// PriorityCode returns the priority code of the currency.
func (ethAtom *EthereumAtom) PriorityCode() uint32 {
	return ethAtom.key.PriorityCode()
}

// GetSecretHash returns the Secret Hash of the atom.
func (eth *EthereumAtom) GetSecretHash() [32]byte {
	return eth.data.HashLock
}

// GetKey returns the key of the atom.
func (eth *EthereumAtom) GetKey() swap.Key {
	return eth.key
}
