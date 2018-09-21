package btc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

type Conn struct {
	URL         string
	ChainParams *chaincfg.Params
	Network     string
}

func NewConnWithConfig(conf config.BitcoinNetwork) (Conn, error) {
	return NewConn(conf.Network, conf.URL)
}

func NewConn(chain, url string) (Conn, error) {
	var chainParams *chaincfg.Params

	switch chain {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	return Conn{
		URL:         url,
		ChainParams: chainParams,
		Network:     chain,
	}, nil
}

func (conn Conn) SignTransaction(tx *wire.MsgTx, key keystore.BitcoinKey, fee int64) (*wire.MsgTx, bool, error) {
	var value int64
	for _, j := range tx.TxOut {
		value = value + j.Value
	}
	value = value + fee

	unspentValue, err := conn.Balance(key.AddressString)
	if err != nil {
		return nil, false, err
	}

	if value > unspentValue {
		return nil, false, fmt.Errorf("Not enough balance in %s "+
			"required:%d current:%d", key.AddressString, value, unspentValue)
	}

	utxos, err := conn.GetUnspentOutputs(key.AddressString)
	if err != nil {
		return nil, false, err
	}

	for _, j := range utxos.Outputs {
		if value <= 0 {
			break
		}
		hashBytes, err := hex.DecodeString(j.TxID)
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
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(hash, j.Vout), ScriptPubKey, [][]byte{}))
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

	return tx, true, nil
}

func (conn Conn) Balance(address string) (int64, error) {
	utxos, err := conn.GetUnspentOutputs(address)
	if err != nil {
		return -1, err
	}
	var balance int64
	for _, utxo := range utxos.Outputs {
		balance = balance + utxo.Amount
	}
	return balance, nil
}

// WaitTillMined waits for the transactions to be mined, and gets the given
// number of confirmations.
func (conn Conn) WaitTillMined(txHash *chainhash.Hash, confirmations int64) error {
	for {
		mined, err := conn.Mined(txHash.String(), confirmations)
		if err != nil {
			return err
		}

		if mined {
			return nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (conn Conn) GetUnspentOutputs(address string) (UnspentOutputs, error) {
	resp, err := http.Get(fmt.Sprintf(conn.URL + "/unspent?active=" + address + "&confirmations=0"))
	if err != nil {
		return UnspentOutputs{}, err
	}
	defer resp.Body.Close()
	utxoBytes, err := ioutil.ReadAll(resp.Body)
	utxos := UnspentOutputs{}
	json.Unmarshal(utxoBytes, &utxos)
	return utxos, nil
}

func (conn Conn) PublishTransaction(signedTransaction []byte) error {
	data := url.Values{}
	data.Set("tx", hex.EncodeToString(signedTransaction))
	client := &http.Client{}
	r, err := http.NewRequest("POST", conn.URL+"/pushtx", strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if _, err = client.Do(r); err != nil {
		return err
	}
	return nil
}

func (conn Conn) GetScriptCreationTransaction(address, scriptHash string) (*wire.MsgTx, error) {
	return nil, nil
}

func (conn Conn) Net() *chaincfg.Params {
	return conn.ChainParams
}

type RawTransaction struct {
	BlockHeight      int64    `json:"block_height"`
	VinSize          uint32   `json:"vin_sz"`
	VoutSize         uint32   `json:"vout_sz"`
	Version          uint8    `json:"ver"`
	TransactionHash  string   `json:"hash"`
	TransactionIndex uint64   `json:"tx_index"`
	Inputs           []Input  `json:"inputs"`
	Outputs          []Output `json:"out"`
}

type Input struct {
	PrevOut PreviousOut `json:"prev_out"`
	Script  string      `json:"script"`
}

type PreviousOut struct {
	TransactionHash  string `json:"hash"`
	Value            uint64 `json:"value"`
	TransactionIndex uint64 `json:"tx_index"`
	VoutNumber       uint8  `json:"n"`
}

type Output struct {
	TransactionHash string `json:"hash"`
	Value           uint64 `json:"value"`
	Script          string `json:"script"`
}

func (conn Conn) GetRawTransaction(txhash string) (RawTransaction, error) {
	resp, err := http.Get(fmt.Sprintf(conn.URL + "/rawtx/" + txhash))
	if err != nil {
		return RawTransaction{}, err
	}
	defer resp.Body.Close()
	txBytes, err := ioutil.ReadAll(resp.Body)
	transaction := RawTransaction{}
	if err := json.Unmarshal(txBytes, &transaction); err != nil {
		return RawTransaction{}, err
	}
	return transaction, nil
}

type LatestBlock struct {
	BlockHash          string  `json:"hash"`
	Time               int64   `json:"time"`
	BlockIndex         int64   `json:"block_index"`
	Height             int64   `json:"height"`
	TransactionIndexes []int64 `json:"txIndexes"`
}

func (conn Conn) Mined(txhash string, confirmations int64) (bool, error) {
	if confirmations <= 0 {
		return true, nil
	}

	confirmations = confirmations - 1
	latestBlock := LatestBlock{}

	resp, err := http.Get(fmt.Sprintf(conn.URL + "/latestblock"))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	blockBytes, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(blockBytes, &latestBlock); err != nil {
		return false, err
	}

	tx, err := conn.GetRawTransaction(txhash)
	if err != nil {
		return false, err
	}

	if tx.BlockHeight != 0 {
		return true, nil
	}

	return false, nil
}

type UnspentOutput struct {
	TxID         string `json:"tx_hash"`
	Vout         uint32 `json:"tx_output_n"`
	ScriptPubKey string `json:"script"`
	Amount       int64  `json:"value"`
}

type UnspentOutputs struct {
	Outputs []UnspentOutput `json:"unspent_outputs"`
}
