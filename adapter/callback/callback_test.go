package callback_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/adapter/callback"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/rs/cors"
)

var _ = Describe("Server Adapter", func() {
	writeError := func(w http.ResponseWriter, statusCode int, err string) {
		w.WriteHeader(statusCode)
		w.Write([]byte(err))
		return
	}

	recoveryHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					writeError(w, http.StatusInternalServerError, fmt.Sprintf("%v", r))
				}
			}()
			h.ServeHTTP(w, r)
		})
	}

	startTestServer := func(postSwapsHandler func() http.HandlerFunc, doneCh <-chan struct{}, port int64) {
		r := mux.NewRouter()
		r.HandleFunc("/swaps", postSwapsHandler()).Methods("POST")
		r.Use(recoveryHandler)
		handler := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST"},
		}).Handler(r)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			panic(err)
		}

		if err := http.Serve(listener, handler); err != nil {
			panic(err)
		}
	}

	randomString := func() string {
		id := [32]byte{}
		rand.Read(id[:])
		return base64.StdEncoding.EncodeToString(id[:])
	}

	tokenPairOptions := []struct {
		sendToken, receiveToken string
	}{
		{blockchain.TokenBTC.String(), blockchain.TokenETH.String()},
		{blockchain.TokenBTC.String(), blockchain.TokenWBTC.String()},
		{blockchain.TokenETH.String(), blockchain.TokenBTC.String()},
		{blockchain.TokenETH.String(), blockchain.TokenWBTC.String()},
		{blockchain.TokenWBTC.String(), blockchain.TokenBTC.String()},
		{blockchain.TokenWBTC.String(), blockchain.TokenETH.String()},
	}

	amountOptions := []struct {
		sendAmount, receiveAmount, minReceiveAmount string
	}{
		{"1000", "100", "0"},
		{"1000", "100", "10"},
		{"1000", "100", "100"},

		{"100", "1000", "0"},
		{"100", "1000", "100"},
		{"100", "1000", "1000"},
	}

	honestUpdateAmountOptions := []struct {
		sendAmount, receiveAmount string
	}{
		{"800", "60"},
		{"100", "10"},
		{"80", "100"},

		{"80", "950"},
		{"10", "100"},
		{"8", "1200"},
	}

	maliciousUpdateAmountOptions := []struct {
		sendAmount, receiveAmount string
	}{
		{"1200", "100"},
		{"1000", "120"},
		{"100", "10"},

		{"100", "800"},
		{"120", "1000"},
		{"10", "100"},
	}

	initiationOptions := []bool{
		true,
		false,
	}

	type TestDelayInfo struct {
		Index int `json:"index"`
	}

	partialSwaps := []swap.SwapBlob{}
	for _, tokenPairOption := range tokenPairOptions {
		for i, amountOption := range amountOptions {
			for _, initiationOption := range initiationOptions {
				delayInfo, _ := json.Marshal(TestDelayInfo{i})
				swap := swap.SwapBlob{
					ID:                   swap.RandomID(),
					SendToken:            tokenPairOption.sendToken,
					ReceiveToken:         tokenPairOption.receiveToken,
					SendAmount:           amountOption.sendAmount,
					ReceiveAmount:        amountOption.receiveAmount,
					MinimumReceiveAmount: amountOption.minReceiveAmount,
					ShouldInitiateFirst:  initiationOption,
					Delay:                true,
					DelayInfo:            delayInfo,
				}

				if initiationOption {
					swap.SecretHash = randomString()
					swap.TimeLock = time.Now().Unix()
				}

				partialSwaps = append(partialSwaps, swap)
			}
		}
	}

	Context("when the broker is not changing amounts being swapped", func() {
		doneCh := make(chan struct{})
		go startTestServer(func() http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				swap := swap.SwapBlob{}
				if err := json.NewDecoder(r.Body).Decode(&swap); err != nil {
					writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
					return
				}

				if !swap.ShouldInitiateFirst {
					swap.SecretHash = randomString()
					swap.TimeLock = time.Now().Unix()
				}
				swap.SendTo = fmt.Sprintf("Address:%s", swap.SendToken)
				swap.ReceiveFrom = fmt.Sprintf("Address:%s", swap.ReceiveToken)

				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(swap); err != nil {
					writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
					return
				}
			}
		}, doneCh, 17777)

		for _, pendingSwap := range partialSwaps {
			It(fmt.Sprintf("verification should succeed, and delay should be set to false"), func() {
				pendingSwap.DelayCallbackURL = "http://127.0.0.1:17777/swaps"
				swapFiller := New()
				filledSwap, err := swapFiller.DelayCallback(pendingSwap)
				Expect(err).Should(BeNil())
				Expect(filledSwap.Delay).Should(BeFalse())
			})
		}

		close(doneCh)
	})

	Context("when the broker is changing amounts being swapped honestly", func() {
		doneCh := make(chan struct{})
		go startTestServer(func() http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				swap := swap.SwapBlob{}
				if err := json.NewDecoder(r.Body).Decode(&swap); err != nil {
					writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
					return
				}

				testDelayInfo := TestDelayInfo{}
				if err := json.Unmarshal(swap.DelayInfo, &testDelayInfo); err != nil {
					writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse delay info: %v", err))
					return
				}

				swap.SendAmount = honestUpdateAmountOptions[testDelayInfo.Index].sendAmount
				swap.ReceiveAmount = honestUpdateAmountOptions[testDelayInfo.Index].receiveAmount

				if !swap.ShouldInitiateFirst {
					swap.SecretHash = randomString()
					swap.TimeLock = time.Now().Unix()
				}
				swap.SendTo = fmt.Sprintf("Address:%s", swap.SendToken)
				swap.ReceiveFrom = fmt.Sprintf("Address:%s", swap.ReceiveToken)
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(swap); err != nil {
					writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
					return
				}
			}
		}, doneCh, 17778)

		for _, pendingSwap := range partialSwaps {
			It(fmt.Sprintf("verification should succeed, and delay should be set to false %v", pendingSwap), func() {
				pendingSwap.DelayCallbackURL = "http://127.0.0.1:17778/swaps"
				swapFiller := New()
				filledSwap, err := swapFiller.DelayCallback(pendingSwap)
				Expect(err).Should(BeNil())
				Expect(filledSwap.Delay).Should(BeFalse())
			})
		}
		close(doneCh)
	})

	Context("when the broker is changing amounts being swapped maliciously", func() {
		doneCh := make(chan struct{})
		go startTestServer(func() http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				swap := swap.SwapBlob{}
				if err := json.NewDecoder(r.Body).Decode(&swap); err != nil {
					writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
					return
				}

				testDelayInfo := TestDelayInfo{}
				if err := json.Unmarshal(swap.DelayInfo, &testDelayInfo); err != nil {
					writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse delay info: %v", err))
					return
				}

				swap.SendAmount = maliciousUpdateAmountOptions[testDelayInfo.Index].sendAmount
				swap.ReceiveAmount = maliciousUpdateAmountOptions[testDelayInfo.Index].receiveAmount

				if !swap.ShouldInitiateFirst {
					swap.SecretHash = randomString()
					swap.TimeLock = time.Now().Unix()
				}
				swap.SendTo = fmt.Sprintf("Address:%s", swap.SendToken)
				swap.ReceiveFrom = fmt.Sprintf("Address:%s", swap.ReceiveToken)
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(swap); err != nil {
					writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
					return
				}
			}
		}, doneCh, 17779)

		for _, pendingSwap := range partialSwaps {
			It(fmt.Sprintf("verification should fail %v", pendingSwap), func() {
				pendingSwap.DelayCallbackURL = "http://127.0.0.1:17779/swaps"
				swapFiller := New()
				_, err := swapFiller.DelayCallback(pendingSwap)
				Expect(err).ShouldNot(BeNil())
			})
		}
		close(doneCh)
	})
})
