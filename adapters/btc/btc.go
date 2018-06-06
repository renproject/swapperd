package btc

import (
	"encoding/json"
	"math/big"
)

type BitcoinData struct {
	ContractHash   string   `json:"contract_hash"`
	Contract       []byte   `json:"contract"`
	ContractTxHash []byte   `json:"contract_txhash"`
	ContractTx     []byte   `json:"contract_tx"`
	RefundTxHash   [32]byte `json:"refund_txhash"`
	RefundTx       []byte   `json:"refund_tx"`
	RedeemTxHash   [32]byte `json:"redeem_txhash"`
	RedeemTx       []byte   `json:"redeem_tx"`
	SecretHash     [32]byte `json:"secret_hash"`
}

type BitcoinAtom struct {
	connection Connection
	ledgerData BitcoinData
}

// NewBitcoinAtom returns an atom object
func NewBitcoinAtom(connection Connection) *BitcoinAtom {
	return &BitcoinAtom{
		connection: connection,
	}
}

func (atom *BitcoinAtom) Initiate(hash [32]byte, from, to []byte, value *big.Int, expiry int64) (err error) {
	result, err := initiate(atom.connection, string(to), value.Int64(), hash[:], expiry)
	if err != nil {
		return err
	}
	atom.ledgerData = result
	atom.ledgerData.SecretHash = hash
	return nil
}

func (atom *BitcoinAtom) Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error) {
	result, err := read(atom.connection, atom.ledgerData.Contract, atom.ledgerData.ContractTx)
	if err != nil {
		return [32]byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.refundAddress, result.recipientAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (atom *BitcoinAtom) Redeem(secret [32]byte) error {
	result, err := redeem(atom.connection, atom.ledgerData.Contract, atom.ledgerData.ContractTx, secret)
	if err != nil {
		return err
	}
	atom.ledgerData.RedeemTx = result.redeemTx
	atom.ledgerData.RedeemTxHash = result.redeemTxHash
	return nil
}

func (atom *BitcoinAtom) AuditSecret() (secret [32]byte, err error) {
	result, err := readSecret(atom.connection, atom.ledgerData.RedeemTx, atom.ledgerData.SecretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

func (atom *BitcoinAtom) Refund() error {
	return refund(atom.connection, atom.ledgerData.Contract, atom.ledgerData.ContractTx)
}

func (atom *BitcoinAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(atom.ledgerData)
	return b, err
}

func (atom *BitcoinAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &atom.ledgerData)
}
