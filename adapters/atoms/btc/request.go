package btc

import (
	"encoding/json"
	"math/big"

	bindings "github.com/republicprotocol/atom-go/adapters/bindings/btc"
	"github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/services/atom"
)

type BitcoinRequestAtom struct {
	personalAddress string
	foreignAddress  string
	connection      btc.Connection
	data            BitcoinData
}

// NewBitcoinRequestAtom returns a new Bitcoin RequestAtom instance
func NewBitcoinRequestAtom(connection btc.Connection, personalAddress, foreignAddress string) atom.RequestAtom {
	return &BitcoinRequestAtom{
		personalAddress: personalAddress,
		foreignAddress:  foreignAddress,
		connection:      connection,
	}
}

// Initiate a new Atom swap by calling Bitcoin
func (btcAtom *BitcoinRequestAtom) Initiate(hash [32]byte, value *big.Int, expiry int64) error {
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
func (btcAtom *BitcoinRequestAtom) Refund() error {
	return bindings.Refund(btcAtom.connection, btcAtom.personalAddress, btcAtom.data.Contract, btcAtom.data.ContractTx)
}

// AuditSecret audits the secret of an Atom swap by calling Bitcoin
func (btcAtom *BitcoinRequestAtom) AuditSecret() ([32]byte, error) {
	result, err := bindings.AuditSecret(btcAtom.connection, btcAtom.data.RedeemTx, btcAtom.data.SecretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

// Serialize serializes the atom details into a bytes array
func (btcAtom *BitcoinRequestAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(btcAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (btcAtom *BitcoinRequestAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &btcAtom.data)
}

// From returns the address of the sender
func (btcAtom *BitcoinRequestAtom) From() []byte {
	return []byte(btcAtom.personalAddress)
}

// PriorityCode returns the priority code of the currency.
func (btcAtom *BitcoinRequestAtom) PriorityCode() int64 {
	return 0
}
