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

type BitcoinData struct {
	HashLock       [32]byte `json:"hash_lock"`
	ContractHash   string   `json:"contract_hash"`
	Contract       []byte   `json:"contract"`
	ContractTxHash []byte   `json:"contract_tx_hash"`
	ContractTx     []byte   `json:"contract_tx"`
	RefundTxHash   [32]byte `json:"refund_tx_hash"`
	RefundTx       []byte   `json:"refund_tx"`
	RedeemTxHash   [32]byte `json:"redeem_tx_hash"`
	RedeemTx       []byte   `json:"redeem_tx"`
	SecretHash     [32]byte `json:"secret_hash"`
}

// BitcoinAtom is a struct for Bitcoin Atom
type BitcoinAtom struct {
	key        swap.Key
	connection btc.Conn
	data       BitcoinData
}

// NewBitcoinAtom returns a new Bitcoin Atom instance
func NewBitcoinAtom(connection btc.Conn, key swap.Key) swap.Atom {
	return &BitcoinAtom{
		key:        key,
		connection: connection,
	}
}

// Initiate a new Atom swap by calling Bitcoin
func (btcAtom *BitcoinAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	from, err := btcAtom.From()
	if err != nil {
		return err
	}
	result, err := bindings.Initiate(btcAtom.connection, string(from), string(to), value.Int64(), hash[:], expiry)
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

// Redeem an Atom swap by calling a function on Bitcoin
func (btcAtom *BitcoinAtom) Redeem(secret [32]byte) error {
	from, err := btcAtom.From()
	if err != nil {
		return err
	}
	result, err := bindings.Redeem(btcAtom.connection, string(from), btcAtom.data.Contract, btcAtom.data.ContractTx, secret)
	if err != nil {
		return err
	}
	btcAtom.data.RedeemTx = result.RedeemTx
	btcAtom.data.RedeemTxHash = result.RedeemTxHash
	return nil
}

// Refund an Atom swap by calling Bitcoin
func (btcAtom *BitcoinAtom) Refund() error {
	from, err := btcAtom.From()
	if err != nil {
		return err
	}
	return bindings.Refund(btcAtom.connection, string(from), btcAtom.data.Contract, btcAtom.data.ContractTx)
}

// Audit an Atom swap by calling a function on Bitcoin
func (btcAtom *BitcoinAtom) Audit(hashLock [32]byte, to []byte, value *big.Int, expiry int64) error {

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

// AuditSecret audits the secret of an Atom swap by calling Bitcoin
func (btcAtom *BitcoinAtom) AuditSecret() ([32]byte, error) {
	result, err := bindings.AuditSecret(btcAtom.connection, btcAtom.data.RedeemTx, btcAtom.data.SecretHash[:])
	if err != nil {
		return [32]byte{}, errors.New("Cannot read the secret")
	}
	if result != [32]byte{} {
		return result, nil
	}
	return [32]byte{}, errors.New("Cannot read the secret")
}

// Serialize serializes the atom details into a bytes array
func (btcAtom *BitcoinAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(btcAtom.data)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (btcAtom *BitcoinAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &btcAtom.data)
}

// PriorityCode returns the priority code of the currency.
func (btcAtom *BitcoinAtom) PriorityCode() uint32 {
	return btcAtom.key.PriorityCode()
}

// GetSecretHash returns the Secret Hash of the atom.
func (btcAtom *BitcoinAtom) GetSecretHash() [32]byte {
	return btcAtom.data.HashLock
}

// From returns the address of the sender
func (btcAtom *BitcoinAtom) From() ([]byte, error) {
	return btcAtom.key.GetAddress()
}

// GetKey returns the key of the atom.
func (btcAtom *BitcoinAtom) GetKey() swap.Key {
	return btcAtom.key
}
