package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/btcjson"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/network"
)

type Conn struct {
	Client      *rpc.Client
	ChainParams *chaincfg.Params
	Network     string
}

func Connect(networkConfig network.Config) (Conn, error) {
	connParams := networkConfig.GetBitcoinNetwork()
	return ConnectWithParams(connParams.Chain, connParams.URL, connParams.User, connParams.Password)
}

func ConnectWithParams(chain, url, user, password string) (Conn, error) {
	var chainParams *chaincfg.Params
	var connect string
	var err error

	switch chain {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	if url == "" {
		connect, err = normalizeAddress("localhost", walletPort(chainParams))
		if err != nil {
			return Conn{}, fmt.Errorf("wallet server address: %v", err)
		}
	} else {
		connect = url
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         user,
		Pass:         password,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	rpcClient, err := rpc.New(connConfig, nil)
	if err != nil {
		return Conn{}, fmt.Errorf("rpc connect: %v", err)
	}

	// Should call the following after this function:
	/*
		defer func() {
			rpcClient.Shutdown()
			pcClient.WaitForShutdown()
		}()
	*/

	return Conn{
		Client:      rpcClient,
		ChainParams: chainParams,
		Network:     chain,
	}, nil
}

func (conn *Conn) FundRawTransaction(tx *wire.MsgTx) (fundedTx *wire.MsgTx, err error) {
	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	tx.Serialize(&buf)
	param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	params := []json.RawMessage{param0}
	rawResp, err := conn.Client.RawRequest("fundrawtransaction", params)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Hex       string  `json:"hex"`
		Fee       float64 `json:"fee"`
		ChangePos float64 `json:"changepos"`
	}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	fundedTxBytes, err := hex.DecodeString(resp.Hex)
	if err != nil {
		return nil, err
	}
	fundedTx = &wire.MsgTx{}
	err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
	if err != nil {
		return nil, err
	}
	return fundedTx, nil
}

func (conn *Conn) PromptPublishTx(tx *wire.MsgTx, name string) (*chainhash.Hash, error) {
	// FIXME: Transaction fees are set to high, change it before deploying to mainnet. By changing the booleon to false.
	txHash, err := conn.Client.SendRawTransaction(tx, true)
	if err != nil {
		return nil, fmt.Errorf("sendrawtransaction: %v", err)
	}
	return txHash, nil
}

func (conn *Conn) WaitForConfirmations(txHash *chainhash.Hash, requiredConfirmations int64) error {
	confirmations := int64(0)
	for confirmations < requiredConfirmations {
		txDetails, err := conn.Client.GetTransaction(txHash)
		if err != nil {
			return err
		}
		confirmations = txDetails.Confirmations

		// TODO: Base delay on chain config
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (conn *Conn) Shutdown() {
	conn.Client.Shutdown()
	conn.Client.WaitForShutdown()
}

func normalizeAddress(addr string, defaultPort string) (hostport string, err error) {
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err = net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}

func walletPort(params *chaincfg.Params) string {
	switch params {
	case &chaincfg.MainNetParams:
		return "8332"
	case &chaincfg.TestNet3Params:
		return "18332"
	case &chaincfg.RegressionNetParams:
		return "18443"
	default:
		return ""
	}
}

func (conn *Conn) FundTransaction(tx *wire.MsgTx, addresses []btcutil.Address) (fundedTx *wire.MsgTx, inputs []btcjson.RawTxInput, err error) {
	var value, unspentValue float64
	for _, j := range tx.TxOut {
		value = value + float64(j.Value)
	}
	Unspents, err := conn.Client.ListUnspentMinMaxAddresses(0, 99999, addresses)
	if err != nil {
		return nil, nil, err
	}

	value = value / 100000000
	for _, j := range Unspents {
		unspentValue = unspentValue + j.Amount
	}
	if value > unspentValue {
		return nil, nil, fmt.Errorf("Not enough balance required:%f current:%f", value, unspentValue)
	}
	selectedTxIns := []btcjson.RawTxInput{}

	for _, j := range Unspents {
		if value <= 0 {
			break
		}
		hashBytes, err := hex.DecodeString(j.TxID)
		if err != nil {
			return nil, nil, err
		}
		hash, err := chainhash.NewHash(reverse(hashBytes))
		if err != nil {
			return nil, nil, err
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(hash, j.Vout), []byte{}, [][]byte{}))
		selectedTxIns = append(selectedTxIns, btcjson.RawTxInput{
			Txid:         j.TxID,
			Vout:         j.Vout,
			ScriptPubKey: j.ScriptPubKey,
			RedeemScript: j.RedeemScript,
		})
		value = value - j.Amount
	}

	P2PKHscript, err := txscript.PayToAddrScript(addresses[0])
	if err != nil {
		return nil, nil, err
	}

	if value <= 0 {
		tx.AddTxOut(wire.NewTxOut(int64(-value*100000000)-10000, P2PKHscript))
	}

	return tx, selectedTxIns, nil
}

// TODO: Implement Sign Transaction logic here so that we do not have to import
// the privatekey onto the bitcoin node and people can submit signed
// transactions to arbitrary nodes.
func (conn *Conn) SignTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	buf := bytes.NewBuffer([]byte{})

	// for _, txin := range tx.TxIn {
	// 	fmt.Println(*txin)
	// }

	// for _, txout := range tx.TxOut {
	// 	fmt.Println(*txout)
	// }

	if err := tx.Serialize(buf); err != nil {
		panic(err)
	}
	// hash1 := sha256.Sum256(buf.Bytes())
	// hash2 := sha256.Sum256(hash1[:])
	// fmt.Println(hash2)

	stx, complete, err := conn.Client.SignRawTransaction(tx)
	if err != nil {
		return nil, false, err
	}

	// for _, txin := range stx.TxIn {
	// 	fmt.Println(*txin)
	// }

	// for _, txout := range stx.TxOut {
	// 	fmt.Println(*txout)
	// }

	return stx, complete, nil
}

func reverse(arr []byte) []byte {
	for i, j := 0, len(arr)-1; i < len(arr)/2; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
