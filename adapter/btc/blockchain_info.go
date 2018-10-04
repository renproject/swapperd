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
)

type PreviousOut struct {
	TransactionHash  string `json:"hash"`
	Value            uint64 `json:"value"`
	TransactionIndex uint64 `json:"tx_index"`
	VoutNumber       uint8  `json:"n"`
	Address          string `json:"addr"`
}

type Input struct {
	PrevOut PreviousOut `json:"prev_out"`
	Script  string      `json:"script"`
}

type Output struct {
	Value           uint64 `json:"value"`
	TransactionHash string `json:"hash"`
	Script          string `json:"script"`
}

type Transaction struct {
	TransactionHash  string   `json:"hash"`
	Version          uint8    `json:"ver"`
	VinSize          uint32   `json:"vin_sz"`
	VoutSize         uint32   `json:"vout_sz"`
	Size             int64    `json:"size"`
	RelayedBy        string   `json:"relayed_by"`
	BlockHeight      int64    `json:"block_height"`
	TransactionIndex uint64   `json:"tx_index"`
	Inputs           []Input  `json:"inputs"`
	Outputs          []Output `json:"out"`
}

type Block struct {
	BlockHash         string        `json:"hash"`
	Version           uint8         `json:"ver"`
	PreviousBlockHash string        `json:"prev_block"`
	MerkleRoot        string        `json:"mrkl_root"`
	Time              int64         `json:"time"`
	Bits              int64         `json:"bits"`
	Nonce             int64         `json:"nonce"`
	TransactionCount  int           `json:"n_tx"`
	Size              int64         `json:"size"`
	BlockIndex        uint64        `json:"block_index"`
	MainChain         bool          `json:"main_chain"`
	Height            int64         `json:"height"`
	ReceivedTime      int64         `json:"received_time"`
	RelayedBy         string        `json:"relayed_by"`
	Transactions      []Transaction `json:"tx"`
}

type Blocks struct {
	Blocks []Block `json:"block"`
}

type SingleAddress struct {
	PublicKeyHash              string        `json:"hash160"`
	Address                    string        `json:"address"`
	TransactionCount           int64         `json:"n_tx"`
	UnredeemedTransactionCount int64         `json:"n_unredeemed"`
	Received                   int64         `json:"total_received"`
	Sent                       int64         `json:"total_sent"`
	Balance                    int64         `json:"final_balance"`
	Transactions               []Transaction `json:"txs"`
}

type Address struct {
	PublicKeyHash    string `json:"hash160"`
	Address          string `json:"address"`
	TransactionCount int64  `json:"n_tx"`
	Received         int64  `json:"total_received"`
	Sent             int64  `json:"total_sent"`
	Balance          int64  `json:"final_balance"`
}

type MultiAddress struct {
	Addresses    []Address     `json:"addresses"`
	Transactions []Transaction `json:"txs"`
}

type UnspentOutput struct {
	TransactionAge          string `json:"tx_age"`
	TransactionHash         string `json:"tx_hash"`
	TransactionIndex        uint32 `json:"tx_index"`
	TransactionOutputNumber uint32 `json:"tx_output_n"`
	ScriptPubKey            string `json:"script"`
	Amount                  int64  `json:"value"`
}

type UnspentOutputs struct {
	Outputs []UnspentOutput `json:"unspent_outputs"`
}

type BlockchainInfoClient struct {
	URL string
}

func NewBlockchainInfoClient(url string) BlockchainInfoClient {
	return BlockchainInfoClient{
		URL: url,
	}
}

func (client BlockchainInfoClient) GetUnspentOutputs(address string, limit, confitmations int64) UnspentOutputs {
	if limit == 0 {
		limit = 250
	}
	for {
		resp, err := http.Get(fmt.Sprintf("%s/unspent?active=%s&confirmations=%d&limit=%d", client.URL, address, confitmations, limit))
		if err != nil {
			fmt.Println(err, " will try again in 10 sec")
			time.Sleep(10 * time.Second)
			continue
		}
		defer resp.Body.Close()
		utxoBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err, " will try again in 10 sec")
			time.Sleep(10 * time.Second)
			continue
		}
		if string(utxoBytes) == "No free outputs to spend" {
			return UnspentOutputs{
				Outputs: []UnspentOutput{},
			}
		}
		utxos := UnspentOutputs{}
		if err := json.Unmarshal(utxoBytes, &utxos); err != nil {
			fmt.Println(err, " will try again in 10 sec")
			time.Sleep(10 * time.Second)
			continue
		}
		return utxos
	}
}

func (client BlockchainInfoClient) Balance(address string, confirmations int64) int64 {
	utxos := client.GetUnspentOutputs(address, 1000, confirmations)
	var balance int64
	for _, utxo := range utxos.Outputs {
		balance = balance + utxo.Amount
	}
	return balance
}

func (client BlockchainInfoClient) GetRawTransaction(txhash string) Transaction {
	for {
		resp, err := http.Get(fmt.Sprintf("%s/rawtx/%s", client.URL, txhash))
		if err != nil {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		defer resp.Body.Close()
		txBytes, err := ioutil.ReadAll(resp.Body)
		transaction := Transaction{}
		if err := json.Unmarshal(txBytes, &transaction); err != nil {
			time.Sleep(10 * time.Second)
			fmt.Println(err)
			continue
		}
		return transaction
	}
}

func (client BlockchainInfoClient) GetRawAddressInformation(addr string) SingleAddress {
	for {
		resp, err := http.Get(fmt.Sprintf("%s/rawaddr/%s", client.URL, addr))
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		defer resp.Body.Close()
		addrBytes, err := ioutil.ReadAll(resp.Body)
		addressInfo := SingleAddress{}
		if err := json.Unmarshal(addrBytes, &addressInfo); err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		return addressInfo
	}
}

func (client BlockchainInfoClient) ScriptSpent(address string) bool {
	rawAddress := client.GetRawAddressInformation(address)
	return rawAddress.Sent > 0
}

func (client BlockchainInfoClient) ScriptFunded(address string, value int64) (bool, int64) {
	rawAddress := client.GetRawAddressInformation(address)
	return rawAddress.Received >= value, rawAddress.Received
}

func (client BlockchainInfoClient) PublishTransaction(signedTransaction []byte) error {
	data := url.Values{}
	data.Set("tx", hex.EncodeToString(signedTransaction))
	httpClient := &http.Client{}
	r, err := http.NewRequest("POST", fmt.Sprintf("%s/pushtx", client.URL), strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// TODO: Handle response and return an error if the transaction fails
	resp, err := httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
