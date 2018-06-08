package btcatom

import (
	"encoding/json"
	"math/big"

	"github.com/republicprotocol/atom-go/adapters/btcclient"
)

type BitcoinData struct {
	Owner          []byte   `json:"owner"`
	ContractHash   string   `json:"contract_hash"`
	Contract       []byte   `json:"contract"`
	ContractTxHash []byte   `json:"contract_txhash"`
	ContractTx     []byte   `json:"contract_tx"`
	RedeemTxHash   []byte   `json:"redeem_txhash"`
	RedeemTx       []byte   `json:"redeem_tx"`
	SecretHash     [32]byte `json:"secret_hash"`
}

type BitcoinAtom struct {
	from         []byte
	connection   btcclient.Connection
	personalData BitcoinData
	foreignData  BitcoinData
}

// NewBitcoinAtom returns an atom object
func NewBitcoinAtom(from []byte, connection btcclient.Connection) *BitcoinAtom {
	return &BitcoinAtom{
		from:       from,
		connection: connection,
		personalData: BitcoinData{
			Owner: from,
		},
	}
}

func (atom *BitcoinAtom) Initiate(hash [32]byte, value *big.Int, expiry int64) (err error) {
	result, err := initiate(atom.connection, string(atom.foreignData.Owner), value.Int64(), hash[:], expiry)
	if err != nil {
		return err
	}
	atom.personalData = result
	atom.personalData.SecretHash = hash
	return nil
}

func (atom *BitcoinAtom) Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error) {
	result, err := read(atom.connection, atom.foreignData.Contract, atom.foreignData.ContractTx)
	if err != nil {
		return [32]byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.refundAddress, result.recipientAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (atom *BitcoinAtom) Redeem(secret [32]byte) error {
	result, err := redeem(atom.connection, atom.foreignData.Contract, atom.foreignData.ContractTx, secret)
	if err != nil {
		return err
	}
	atom.personalData.RedeemTx = result.redeemTx
	atom.personalData.RedeemTxHash = result.redeemTxHash[:]
	return nil
}

func (atom *BitcoinAtom) AuditSecret() (secret [32]byte, err error) {
	result, err := readSecret(atom.connection, atom.foreignData.RedeemTx, atom.foreignData.SecretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

func (atom *BitcoinAtom) Refund() error {
	return refund(atom.connection, atom.personalData.Contract, atom.personalData.ContractTx)
}

func (atom *BitcoinAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(atom.personalData)
	return b, err
}

func (atom *BitcoinAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &atom.foreignData)
}

// From returns the spenders bitcoin address
func (atom *BitcoinAtom) From() []byte {
	return atom.from
}

// PriorityCode returns the priority code of the currency.
func (atom *BitcoinAtom) PriorityCode() int64 {
	return 0
}
