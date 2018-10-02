package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

type Conn struct {
	URL         string
	Network     string
	ChainParams *chaincfg.Params
	BlockchainInfoClient
}

func NewConnWithConfig(conf config.BitcoinNetwork) *Conn {
	return NewConn(conf.Network, conf.URL)
}

func NewConn(chain, url string) *Conn {
	var chainParams *chaincfg.Params
	switch chain {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}
	return &Conn{
		URL:                  url,
		ChainParams:          chainParams,
		Network:              chain,
		BlockchainInfoClient: NewBlockchainInfoClient(url),
	}
}

func (conn *Conn) PublishTransaction(stx *wire.MsgTx, postCon func() bool) error {
	var stxBuffer bytes.Buffer
	stxBuffer.Grow(stx.SerializeSize())
	if err := stx.Serialize(&stxBuffer); err != nil {
		return err
	}
	for {
		if err := conn.BlockchainInfoClient.PublishTransaction(stxBuffer.Bytes()); err != nil {
			return err
		}
		for i := 0; i < 20; i++ {
			if postCon() {
				return nil
			}
			time.Sleep(15 * time.Second)
		}
	}
}

func (conn *Conn) SignTransaction(tx *wire.MsgTx, key keystore.BitcoinKey, fee int64) (*wire.MsgTx, bool, error) {
	var value int64
	for _, j := range tx.TxOut {
		value = value + j.Value
	}
	value = value + fee
	unspentValue := conn.Balance(key.AddressString, 0)
	if value > unspentValue {
		return nil, false, fmt.Errorf("Not enough balance in %s "+
			"required:%d current:%d", key.AddressString, value, unspentValue)
	}
	utxos := conn.GetUnspentOutputs(key.AddressString, 1000, 0)
	for _, j := range utxos.Outputs {
		if value <= 0 {
			break
		}
		hashBytes, err := hex.DecodeString(j.TransactionHash)
		if err != nil {
			return nil, false, err
		}
		hash, err := chainhash.NewHash(hashBytes)
		if err != nil {
			return nil, false, err
		}
		ScriptPubKey, err := hex.DecodeString(j.ScriptPubKey)
		if err != nil {
			return nil, false, err
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(hash, j.TransactionOutputNumber), ScriptPubKey, [][]byte{}))
		value = value - j.Amount
	}
	if value < 0 {
		P2PKHScript, err := txscript.PayToAddrScript(key.Address)
		if err != nil {
			return nil, false, err
		}
		tx.AddTxOut(wire.NewTxOut(int64(-value), P2PKHScript))
	}
	for i, txin := range tx.TxIn {
		sigScript, err := txscript.SignatureScript(tx, i, txin.SignatureScript, txscript.SigHashAll, key.WIF.PrivKey, true)
		if err != nil {
			return nil, false, err
		}
		tx.TxIn[i].SignatureScript = sigScript
	}
	return tx, true, nil
}

func (conn *Conn) Withdraw(addr string, key keystore.BitcoinKey, value, fee int64) error {
	tx := wire.NewMsgTx(2)
	address, err := btcutil.DecodeAddress(addr, conn.ChainParams)
	P2PKHScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return err
	}
	balance := conn.Balance(key.AddressString, 0)
	if value == 0 || value > balance-fee {
		value = balance - fee
	}
	tx.AddTxOut(wire.NewTxOut(value, P2PKHScript))
	stx, complete, err := conn.SignTransaction(tx, key, fee)
	if err != nil || !complete {
		return fmt.Errorf("Failed to sign the transaction complete: %v error: %v", complete, err)
	}
	return conn.PublishTransaction(stx, func() bool {
		tx := conn.GetRawTransaction(stx.TxHash().String())
		return tx.Version == 2
	})
}
