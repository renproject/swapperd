package btc

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	bindings "github.com/republicprotocol/atom-go/adapters/bindings/btc"
	"github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/services/swap"
)

// BitcoinAtomResponder is a struct for Bitcoin AtomResponder
type BitcoinAtomResponder struct {
	personalAddress string
	connection      btc.Conn
	data            BitcoinData
}

// NewBitcoinAtomResponder returns a new Bitcoin AtomResponder instance
func NewBitcoinAtomResponder(connection btc.Conn, address string) swap.AtomResponder {
	return &BitcoinAtomResponder{
		personalAddress: address,
		connection:      connection,
	}
}

// Redeem an Atom swap by calling a function on Bitcoin
func (btcAtom *BitcoinAtomResponder) Redeem(secret [32]byte) error {

	result, err := bindings.Redeem(btcAtom.connection, btcAtom.personalAddress, btcAtom.data.Contract, btcAtom.data.ContractTx, secret)
	if err != nil {
		return err
	}
	btcAtom.data.RedeemTx = result.RedeemTx
	btcAtom.data.RedeemTxHash = result.RedeemTxHash
	return nil
}

// Audit an Atom swap by calling a function on Bitcoin
func (btcAtom *BitcoinAtomResponder) Audit(hashLock [32]byte, to []byte, value *big.Int, expiry int64) error {

	result, err := bindings.Audit(btcAtom.connection, btcAtom.data.Contract, btcAtom.data.ContractTx)
	if err != nil {
		return err
	}

	btcAtom.data.HashLock = result.SecretHash

	if bytes.Compare(to, result.RecipientAddress) != 0 {
		return errors.New("Btc: To Address mismatch")
	}

	// if value.Cmp(big.NewInt(result.Amount)) > 0 {
	// 	return errors.New("Value mismatch")
	// }

	// if expiry > (result.LockTime - time.Now().Unix()) {
	// 	return errors.New("Expiry mismatch")
	// }
	println("Audit Successful")
	return nil
}

// Serialize serializes the atom details into a bytes array
func (btcAtom *BitcoinAtomResponder) Serialize() ([]byte, error) {
	b, err := json.Marshal(btcAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (btcAtom *BitcoinAtomResponder) Deserialize(b []byte) error {
	return json.Unmarshal(b, &btcAtom.data)
}

// PriorityCode returns the priority code of the currency.
func (btcAtom *BitcoinAtomResponder) PriorityCode() int64 {
	return 0
}

// GetSecretHash returns the Secret Hash of the atom.
func (btcAtom *BitcoinAtomResponder) GetSecretHash() [32]byte {
	return btcAtom.data.HashLock
}
