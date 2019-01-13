package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"

	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/core/transfer"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
)

var _ = Describe("Server Adapter", func() {
	done := make(chan struct{})

	startServer := func() {
		Expect(keystore.Generate("../../secrets", "testnet", "weird")).Should(BeNil())
		blockchain, err := keystore.Wallet("../../secrets", "testnet")
		Expect(err).Should(BeNil())
		ldb, err := leveldb.NewStore("../../secrets", "testnet")
		Expect(err).Should(BeNil())
		storage := db.New(ldb)
		logger := logger.NewStdOut()
		swapperdTask := swapper.New(128, storage, binder.NewBuilder(blockchain, logger), callback.New())
		walletTask := transfer.New(128, blockchain, storage, logger)
		go func() {
			httpServer := NewHttpServer(blockchain, logger, swapperdTask, walletTask, "27927")
			httpServer.Run(done)
		}()
	}

	buildSwap := func(password string) swap.SwapBlob {
		wallet, err := keystore.Wallet("../../secrets", "testnet")
		Expect(err).Should(BeNil())
		ethAddr, err := wallet.GetAddress(password, blockchain.Ethereum)
		Expect(err).Should(BeNil())
		btcAddr, err := wallet.GetAddress(password, blockchain.Bitcoin)
		Expect(err).Should(BeNil())

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

	fullSwap := func(initiatorPassword, responderPassword string) swap.SwapID {
		aliceSwap := buildSwap(responderPassword)
		data, err := json.MarshalIndent(aliceSwap, "", "  ")
		Expect(err).Should(BeNil())
		req1, err := http.NewRequest("POST", "http://localhost:27927/swaps", bytes.NewBuffer(data))
		Expect(err).Should(BeNil())
		req1.SetBasicAuth("", initiatorPassword)
		resp1, err := http.DefaultClient.Do(req1)
		Expect(err).Should(BeNil())
		respBytes, err := ioutil.ReadAll(resp1.Body)
		fmt.Println(string(respBytes))
		Expect(resp1.StatusCode).Should(Equal(http.StatusCreated))
		swapResp := PostSwapResponse{}
		Expect(json.Unmarshal(respBytes, &swapResp)).Should(BeNil())
		bobSwapData, err := json.MarshalIndent(swapResp.Swap, "", "  ")
		Expect(err).Should(BeNil())
		req2, err := http.NewRequest("POST", "http://localhost:27927/swaps", bytes.NewBuffer(bobSwapData))
		Expect(err).Should(BeNil())
		req2.SetBasicAuth("", responderPassword)
		resp2, err := http.DefaultClient.Do(req2)
		Expect(err).Should(BeNil())
		Expect(resp2.StatusCode).Should(Equal(http.StatusCreated))
		respBytes2, err := ioutil.ReadAll(resp2.Body)
		Expect(err).Should(BeNil())
		swapResp2 := PostSwapResponse{}
		Expect(json.Unmarshal(respBytes2, &swapResp2)).Should(BeNil())
		return swapResp2.ID
	}

	waitForSwap := func(responderPassword string, id swap.SwapID) {
		for {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:27927/swaps?id=%s", url.QueryEscape(string(id))), nil)
			if err != nil {
				panic(err)
			}
			req.SetBasicAuth("", responderPassword)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			swapResp := GetSwapResponse{}
			respBytes, err := ioutil.ReadAll(resp.Body)
			json.Unmarshal(respBytes, &swapResp)

			swap := swap.SwapReceipt(swapResp)
			if swap.Status == 5 {
				break
			}

			time.Sleep(10 * time.Second)
		}
	}

	bootload := func() {
		req1, err := http.NewRequest("POST", "http://localhost:27927/bootload", nil)
		if err != nil {
			panic(err)
		}
		req1.SetBasicAuth("", "Alice")
		resp1, err := http.DefaultClient.Do(req1)
		if err != nil {
			panic(err)
		}
		Expect(resp1.StatusCode).Should(Equal(http.StatusOK))

		req2, err := http.NewRequest("POST", "http://localhost:27927/bootload", nil)
		if err != nil {
			panic(err)
		}
		req2.SetBasicAuth("", "Bob")
		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			panic(err)
		}
		Expect(resp2.StatusCode).Should(Equal(http.StatusOK))
	}

	BeforeSuite(func() {
		startServer()
		bootload()
	})

	Context("basic requests", func() {
		It("when getting swapperd info", func() {
			req, err := http.NewRequest("GET", "http://localhost:27927/info", nil)
			Expect(err).Should(BeNil())
			req.SetBasicAuth("", "Alice")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			fmt.Println(string(respBytes))
		})

		It("when doing an atomic swap", func() {
			bobSwapID := fullSwap("Alice", "Bob")
			waitForSwap("Bob", bobSwapID)
		})
	})

	AfterSuite(func() {
		close(done)
	})
})
