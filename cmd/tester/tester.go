package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/renproject/tokens"

	"github.com/renproject/swapperd/adapter/server"
	"github.com/renproject/swapperd/driver/keystore"
	"github.com/renproject/swapperd/foundation/swap"
)

func main() {
	for i := 0; i < 5; i++ {
		doSwap("", "")
	}
}

func buildSwap(initiatorPassword string) swap.SwapBlob {
	wallet, err := keystore.Wallet(os.Getenv("HOME")+"/.swapperd", "mainnet")
	if err != nil {
		panic(err)
	}

	ethAddr, err := wallet.GetAddress(initiatorPassword, tokens.ETHEREUM)
	if err != nil {
		panic(err)
	}

	btcAddr, err := wallet.GetAddress(initiatorPassword, tokens.BITCOIN)
	if err != nil {
		panic(err)
	}

	return swap.SwapBlob{
		SendToken:           "BTC",
		ReceiveToken:        "ETH",
		SendAmount:          "20000",
		ReceiveAmount:       "2000000000000",
		SendTo:              btcAddr,
		ReceiveFrom:         ethAddr,
		ShouldInitiateFirst: true,
	}
}

func doSwap(initiatorPassword, responderPassword string) {
	aliceSwap := buildSwap(initiatorPassword)
	data, err := json.MarshalIndent(aliceSwap, "", "  ")
	if err != nil {
		panic(err)
	}
	req1, err := http.NewRequest("POST", "http://localhost:7927/swaps", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req1.SetBasicAuth("", initiatorPassword)
	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		panic(err)
	}
	respBytes, err := ioutil.ReadAll(resp1.Body)
	fmt.Println(string(respBytes))
	if resp1.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("unexpected status code: %d", resp1.StatusCode))
	}
	swapResp := server.PostSwapResponse{}
	if err := json.Unmarshal(respBytes, &swapResp); err != nil {
		panic(err)
	}
	bobSwapData, err := json.MarshalIndent(swapResp.Swap, "", "  ")
	if err != nil {
		panic(err)
	}
	req2, err := http.NewRequest("POST", "http://localhost:7927/swaps", bytes.NewBuffer(bobSwapData))
	if err != nil {
		panic(err)
	}
	req2.SetBasicAuth("", responderPassword)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		panic(err)
	}
	if resp2.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("unexpected status code: %d", resp2.StatusCode))
	}
	respBytes2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		panic(err)
	}
	swapResp2 := server.PostSwapResponse{}
	if err := json.Unmarshal(respBytes2, &swapResp2); err != nil {
		panic(err)
	}
}
