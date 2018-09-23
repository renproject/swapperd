package network

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/renex-ingress-go/httpadapter"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

type ingress struct {
	hostAddress string
	ethKey      keystore.EthereumKey
}

func NewIngress(hostAddress string, ethKey keystore.EthereumKey) renex.Network {
	return &ingress{
		hostAddress: hostAddress,
		ethKey:      ethKey,
	}
}

func (ingress *ingress) SendSwapDetails(orderID [32]byte, swapDetails renex.SwapDetails) error {
	swapBytes, err := json.Marshal(swapDetails)
	if err != nil {
		return err
	}

	info := httpadapter.PostSwapInfo{
		OrderID: base64.StdEncoding.EncodeToString(orderID[:]),
		Swap:    hex.EncodeToString(swapBytes),
	}

	infoBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	hash := crypto.Keccak256(infoBytes)

	signature, err := crypto.Sign(hash, ingress.ethKey.PrivateKey)
	if err != nil {
		return err
	}
	sig65, err := utils.ToBytes65(signature)
	if err != nil {
		return err
	}

	req := httpadapter.PostSwapRequest{
		Info:      info,
		Signature: httpadapter.MarshalSignature(sig65),
	}

	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)

	resp, err := http.Post(fmt.Sprintf("https://"+ingress.hostAddress+"/swap"), "application/json", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		return nil
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("Unexpected error %d: %s", resp.StatusCode, respBytes)
}

func (ingress *ingress) ReceiveSwapDetails(orderID [32]byte, waitTill int64) (renex.SwapDetails, error) {
	swapDetails := renex.SwapDetails{}
	for {
		resp, err := http.Get(fmt.Sprintf("https://" + ingress.hostAddress + "/swap/" + hex.EncodeToString(orderID[:])))
		if err != nil {
			return swapDetails, err
		}
		if resp.StatusCode == 200 {
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return swapDetails, err
			}
			swapBytes, err := hex.DecodeString(string(respBytes))
			if err != nil {
				return swapDetails, err
			}
			if err := json.Unmarshal(swapBytes, &swapDetails); err != nil {
				return swapDetails, err
			}
			resp.Body.Close()
			return swapDetails, nil
		}
		resp.Body.Close()
		if time.Now().Unix() > waitTill {
			return swapDetails, fmt.Errorf("Receive Swap Details Timedout")
		}
		time.Sleep(10 * time.Second)
	}
}
