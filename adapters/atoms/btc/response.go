package btc

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	bindings "github.com/republicprotocol/atom-go/adapters/bindings/btc"
	"github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/services/atom"
)

type BitcoinResponseAtom struct {
	personalAddress string
	foreignAddress  string
	connection      btc.Connection
	data            BitcoinData
}

// NewBitcoinResponseAtom returns a new Bitcoin ResponseAtom instance
func NewBitcoinResponseAtom(connection btc.Connection, personalAddress, foreignAddress string) atom.ResponseAtom {
	return &BitcoinResponseAtom{
		personalAddress: personalAddress,
		foreignAddress:  foreignAddress,
		connection:      connection,
	}
}

// Redeem an Atom swap by calling a function on Bitcoin
func (btcAtom *BitcoinResponseAtom) Redeem(secret [32]byte) error {

	result, err := bindings.Redeem(btcAtom.connection, btcAtom.personalAddress, btcAtom.data.Contract, btcAtom.data.ContractTx, secret)
	if err != nil {
		return err
	}
	btcAtom.data.RedeemTx = result.RedeemTx
	btcAtom.data.RedeemTxHash = result.RedeemTxHash
	return nil
}

// Audit an Atom swap by calling a function on Bitcoin
func (btcAtom *BitcoinResponseAtom) Audit(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
	result, err := bindings.Audit(btcAtom.connection, btcAtom.data.Contract, btcAtom.data.ContractTx)
	if err != nil {
		return err
	}

	if hash != result.SecretHash {
		return errors.New("HashLock mismatch")
	}

	if bytes.Compare(to, result.RecipientAddress) != 0 {
		return errors.New("To Address mismatch")
	}

	// if value.Cmp(big.NewInt(result.Amount)) > 0 {
	// 	return errors.New("Value mismatch")
	// }

	// if expiry > (result.LockTime - time.Now().Unix()) {
	// 	return errors.New("Expiry mismatch")
	// }

	return nil
}

// Serialize serializes the atom details into a bytes array
func (btcAtom *BitcoinResponseAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(btcAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (btcAtom *BitcoinResponseAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &btcAtom.data)
}

// PriorityCode returns the priority code of the currency.
func (btcAtom *BitcoinResponseAtom) PriorityCode() int64 {
	return 0
}

// GetSecretHash returns the Secret Hash of the atom.
func (btcAtom *BitcoinResponseAtom) GetSecretHash() [32]byte {
	return btcAtom.data.HashLock
}
