package btc

import (
	"encoding/json"
	"math/big"

	bindings "github.com/republicprotocol/atom-go/adapters/bindings/btc"
	"github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/services/swap"
)

// BitcoinAtomRequester is a struct for Bitcoin AtomRequester
type BitcoinAtomRequester struct {
	personalAddress string
	foreignAddress  string
	connection      btc.Conn
	data            BitcoinData
}

// NewBitcoinAtomRequester returns a new Bitcoin AtomRequester instance
func NewBitcoinAtomRequester(connection btc.Conn, personalAddress, foreignAddress string) swap.AtomRequester {
	return &BitcoinAtomRequester{
		personalAddress: personalAddress,
		foreignAddress:  foreignAddress,
		connection:      connection,
	}
}

// Initiate a new Atom swap by calling Bitcoin
func (btcAtom *BitcoinAtomRequester) Initiate(hash [32]byte, value *big.Int, expiry int64) error {
	result, err := bindings.Initiate(btcAtom.connection, btcAtom.personalAddress, btcAtom.foreignAddress, value.Int64(), hash[:], expiry)
	if err != nil {
		return err
	}
	btcAtom.data.HashLock = hash
	btcAtom.data.Contract = result.Contract
	btcAtom.data.ContractHash = result.ContractHash
	btcAtom.data.ContractTx = result.ContractTx
	btcAtom.data.ContractTxHash = result.ContractTxHash
	btcAtom.data.RefundTx = result.RefundTx
	btcAtom.data.RefundTxHash = result.RefundTxHash
	btcAtom.data.SecretHash = hash
	return nil
}

// Refund an Atom swap by calling Bitcoin
func (btcAtom *BitcoinAtomRequester) Refund() error {
	return bindings.Refund(btcAtom.connection, btcAtom.personalAddress, btcAtom.data.Contract, btcAtom.data.ContractTx)
}

// AuditSecret audits the secret of an Atom swap by calling Bitcoin
func (btcAtom *BitcoinAtomRequester) AuditSecret() ([32]byte, error) {
	result, err := bindings.AuditSecret(btcAtom.connection, btcAtom.data.RedeemTx, btcAtom.data.SecretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

// Serialize serializes the atom details into a bytes array
func (btcAtom *BitcoinAtomRequester) Serialize() ([]byte, error) {
	b, err := json.Marshal(btcAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (btcAtom *BitcoinAtomRequester) Deserialize(b []byte) error {
	return json.Unmarshal(b, &btcAtom.data)
}

// From returns the address of the sender
func (btcAtom *BitcoinAtomRequester) From() []byte {
	return []byte(btcAtom.personalAddress)
}

// PriorityCode returns the priority code of the currency.
func (btcAtom *BitcoinAtomRequester) PriorityCode() int64 {
	return 0
}
