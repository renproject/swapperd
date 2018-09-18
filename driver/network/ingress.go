package network

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/republicprotocol/renex-ingress-go/httpadapter"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
)

type ingress struct {
	hostAddress string
}

func NewIngress(hostAddress string) network.Network {
	return &ingress{
		hostAddress: hostAddress,
	}
}

func (ingress *ingress) SendOwnerAddress(orderID order.ID, address []byte) error {
	req := httpadapter.PostAddressRequest{
		OrderID: hex.EncodeToString(orderID[:]),
		Address: hex.EncodeToString(address),
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
	return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
}

func (ingress *ingress) SendSwapDetails(orderID order.ID, swapDetails []byte) error {
	req := httpadapter.PostSwapRequest{
		OrderID: hex.EncodeToString(orderID[:]),
		Swap:    hex.EncodeToString(swapDetails),
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
	return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
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
