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
	"github.com/republicprotocol/renex-swapper-go/service/swap"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

type ingress struct {
	hostAddress string
	ethKey      keystore.EthereumKey
}

func NewIngress(hostAddress string, ethKey keystore.EthereumKey) swap.Network {
	return &ingress{
		hostAddress: hostAddress,
		ethKey:      ethKey,
	}
}

func (ingress *ingress) SendOwnerAddress(orderID order.ID, address []byte) error {
	info := httpadapter.PostAddressInfo{
		OrderID: base64.StdEncoding.EncodeToString(orderID[:]),
		Address: hex.EncodeToString(address),
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

	req := httpadapter.PostAddressRequest{
		Info:      info,
		Signature: httpadapter.MarshalSignature(sig65),
	}

	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)

	resp, err := http.Post(fmt.Sprintf("https://"+ingress.hostAddress+"/address"), "application/json", buf)
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

	return fmt.Errorf("Unexpected status error %d: %s", resp.StatusCode, respBytes)
}

func (ingress *ingress) SendSwapDetails(orderID order.ID, swapDetails []byte) error {
	info := httpadapter.PostSwapInfo{
		OrderID: base64.StdEncoding.EncodeToString(orderID[:]),
		Swap:    hex.EncodeToString(swapDetails),
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

func (ingress *ingress) ReceiveOwnerAddress(orderID order.ID, waitTill int64) ([]byte, error) {
	for {
		resp, err := http.Get(fmt.Sprintf("https://" + ingress.hostAddress + "/address/" + hex.EncodeToString(orderID[:])))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 200 {
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			address, err := hex.DecodeString(string(respBytes))
			resp.Body.Close()
			return address, nil
		}
		resp.Body.Close()
		if time.Now().Unix() > waitTill {
			return nil, fmt.Errorf("Receive Address Timedout")
		}
		time.Sleep(10 * time.Second)
	}
}

func (ingress *ingress) ReceiveSwapDetails(orderID order.ID, waitTill int64) ([]byte, error) {
	for {
		resp, err := http.Get(fmt.Sprintf("https://" + ingress.hostAddress + "/swap/" + hex.EncodeToString(orderID[:])))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 200 {
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			swapDetails, err := hex.DecodeString(string(respBytes))
			resp.Body.Close()
			return swapDetails, nil
		}
		resp.Body.Close()
		if time.Now().Unix() > waitTill {
			return nil, fmt.Errorf("Receive Swap Details Timedout")
		}
		time.Sleep(10 * time.Second)
	}
}
