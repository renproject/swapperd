package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/keystore"
)

type conn struct {
	URL         string
	Network     string
	ChainParams *chaincfg.Params
	BlockchainInfoClient
}

func NewConnWithConfig(conf config.BitcoinNetwork) Conn {
	return NewConn(conf.Network, conf.URL)
}

func NewConn(chain, url string) Conn {
	var chainParams *chaincfg.Params
	switch chain {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}
	return &conn{
		URL:                  url,
		ChainParams:          chainParams,
		Network:              chain,
		BlockchainInfoClient: NewBlockchainInfoClient(url),
	}
}

func (conn *conn) PublishTransaction(stx *wire.MsgTx, postCon func() bool) error {
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

func (conn *conn) SignTransaction(tx *wire.MsgTx, key keystore.BitcoinKey, fee int64) (*wire.MsgTx, bool, error) {
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
		sigScript, err := txscript.SignatureScript(tx, i, txin.SignatureScript, txscript.SigHashAll, key.WIF.PrivKey, key.Compressed)
		if err != nil {
			return nil, false, err
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	value = 0
	for _, j := range tx.TxOut {
		value = value + j.Value
	}
	value = value + fee

	// Verify Transaction
	for _, j := range utxos.Outputs {
		if value <= 0 {
			break
		}

		ScriptPubKey, err := hex.DecodeString(j.ScriptPubKey)
		if err != nil {
			return nil, false, err
		}

		if err := verifyTransaction(ScriptPubKey, tx, 0, j.Amount); err != nil {
			return nil, false, err
		}

		value = value - j.Amount
	}

	return tx, true, nil
}

func (conn *conn) Withdraw(addr string, key keystore.BitcoinKey, value, fee int64) error {
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

	if err := conn.PublishTransaction(stx, func() bool {
		tx := conn.GetRawTransaction(stx.TxHash().String())
		return tx.Version == 2
	}); err != nil {
		return err
	}

	switch key.Network.Name {
	case "mainnet":
		fmt.Printf("The transaction can be viewed at https://www.blockchain.com/btc/tx/%s\n", stx.TxHash().String())
	case "testnet3":
		fmt.Printf("The transaction can be viewed at https://testnet.blockchain.info/tx/%s\n", stx.TxHash().String())
	}
	return nil
}

func (conn *conn) SpendBalance(address string) (*wire.MsgTx, []byte, []int64, error) {
	utxos := UnspentOutputs{}
	for {
		utxos = conn.GetUnspentOutputs(address, 1000, 0)
		if len(utxos.Outputs) > 0 {
			break
		}
		time.Sleep(10 * time.Second)
	}
	tx := wire.NewMsgTx(2)
	values := []int64{}
	pkScript := utxos.Outputs[0].ScriptPubKey
	for _, utxo := range utxos.Outputs {
		// Add Transaction input
		if pkScript != utxo.ScriptPubKey {
			return nil, nil, nil, fmt.Errorf("Invalid transaction")
		}
		hashBytes, err := hex.DecodeString(utxo.TransactionHash)
		if err != nil {
			return nil, nil, nil, err
		}
		txHash, err := chainhash.NewHash(hashBytes)
		if err != nil {
			return nil, nil, nil, err
		}
		values = append(values, utxo.Amount)
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txHash, utxo.TransactionOutputNumber), nil, nil))
	}
	pkScriptBytes, err := hex.DecodeString(pkScript)
	if err != nil {
		return nil, nil, nil, err
	}
	return tx, pkScriptBytes, values, nil
}

func (conn *conn) GetScriptFromSpentP2SH(address string) ([]byte, error) {
	for {
		addrInfo := conn.GetRawAddressInformation(address)
		if addrInfo.Sent > 0 {
			break
		}
	}
	addrInfo := conn.GetRawAddressInformation(address)
	for _, tx := range addrInfo.Transactions {
		for i := range tx.Inputs {
			if tx.Inputs[i].PrevOut.Address == addrInfo.Address {
				return hex.DecodeString(tx.Inputs[i].Script)
			}
		}
	}
	return nil, fmt.Errorf("No spending transactions")
}

func (conn *conn) Net() *chaincfg.Params {
	return conn.ChainParams
}

func (conn *conn) FormatTransactionView(msg, txhash string) string {
	switch conn.ChainParams.Name {
	case "mainnet":
		return fmt.Sprintf("%s, https://www.blockchain.com/btc/tx/%s", msg, txhash)
	case "testnet3":
		return fmt.Sprintf("%s, transaction can be viewed at https://testnet.blockchain.info/tx/%s", msg, txhash)
	default:
		panic(fmt.Sprintf("Unsupported network: %s", conn.ChainParams.Name))
	}
}
