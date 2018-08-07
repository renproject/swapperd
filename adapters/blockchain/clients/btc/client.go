package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/network"
)

type Conn struct {
	Client      *rpc.Client
	ChainParams *chaincfg.Params
	Network     string
}

func Connect(networkConfig network.Config) (Conn, error) {
	var chainParams *chaincfg.Params
	var connect string
	var err error

	connParams := networkConfig.GetBitcoinNetwork()

	switch connParams.Chain {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	if connParams.URL == "" {
		connect, err = normalizeAddress("localhost", walletPort(chainParams))
		if err != nil {
			return Conn{}, fmt.Errorf("wallet server address: %v", err)
		}
	} else {
		connect = connParams.URL
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         connParams.User,
		Pass:         connParams.Password,
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
		Network:     networkConfig.Bitcoin.Chain,
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

// func (conn *Conn) FundTransaction(tx *wire.MsgTx, addresses []btcutil.Address) (fundedTx *wire.MsgTx, err error) {

// 	// FIXME: update the output selection policy
// 	// policy := wallet.OutputSelectionPolicy{
// 	// 	Account:               req.Account,
// 	// 	RequiredConfirmations: req.RequiredConfirmations,
// 	// }
// 	// unspentOutputs, err := s.wallet.UnspentOutputs(policy)
// 	// if err != nil {
// 	// 	return nil, translateError(err)
// 	// }

// 	unspentOutputs, err := conn.Client.ListUnspentMinMaxAddresses(1, 999999999, addresses);

// 	selectedOutputs := make([]*pb.FundTransactionResponse_PreviousOutput, 0, len(unspentOutputs))
// 	var totalAmount btcutil.Amount
// 	for _, output := range unspentOutputs {
// 		selectedOutputs = append(selectedOutputs, &pb.FundTransactionResponse_PreviousOutput{

// 	TxID          string  `json:"txid"`
// 	Vout          uint32  `json:"vout"`
// 	Address       string  `json:"address"`
// 	Account       string  `json:"account"`
// 	ScriptPubKey  string  `json:"scriptPubKey"`
// 	RedeemScript  string  `json:"redeemScript,omitempty"`
// 	Amount        float64 `json:"amount"`
// 	Confirmations int64   `json:"confirmations"`
// 	Spendable     bool    `json:"spendable"`

// 			TransactionHash: output.OutPoint.Hash[:],
// 			OutputIndex:     output.OutPoint.Index,
// 			Amount:          output.Output.Value,
// 			PkScript:        output.Output.PkScript,
// 			ReceiveTime:     output.ReceiveTime.Unix(),
// 			FromCoinbase:    output.OutputKind == wallet.OutputKindCoinbase,
// 		})
// 		totalAmount += btcutil.Amount(output.Output.Value)

// 		if req.TargetAmount != 0 && totalAmount > btcutil.Amount(req.TargetAmount) {
// 			break
// 		}
// 	}

// 	var changeScript []byte
// 	if req.IncludeChangeScript && totalAmount > btcutil.Amount(req.TargetAmount) {
// 		changeAddr, err := s.wallet.NewChangeAddress(req.Account, waddrmgr.KeyScopeBIP0044)
// 		if err != nil {
// 			return nil, translateError(err)
// 		}
// 		changeScript, err = txscript.PayToAddrScript(changeAddr)
// 		if err != nil {
// 			return nil, translateError(err)
// 		}
// 	}

// 	return &pb.FundTransactionResponse{
// 		SelectedOutputs: selectedOutputs,
// 		TotalAmount:     int64(totalAmount),
// 		ChangePkScript:  changeScript,
// 	}, nil
// }
