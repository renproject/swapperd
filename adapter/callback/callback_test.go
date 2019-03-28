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
	. "github.com/renproject/swapperd/adapter/callback"
	"github.com/renproject/tokens"

	"github.com/gorilla/mux"
	"github.com/renproject/swapperd/foundation/swap"
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
		sendToken, receiveToken tokens.Name
	}{
		{tokens.NameBTC, tokens.NameETH},
		{tokens.NameBTC, tokens.NameWBTC},
		{tokens.NameETH, tokens.NameBTC},
		{tokens.NameETH, tokens.NameWBTC},
		{tokens.NameWBTC, tokens.NameBTC},
		{tokens.NameWBTC, tokens.NameETH},
	}

	amountOptions := []struct {
		sendAmount, receiveAmount, minReceiveAmount string
	}{
		{"1000000", "100000", "0"},
		{"1000000", "100000", "10000"},
		{"1000000", "100000", "100000"},

		{"100000", "1000000", "0"},
		{"100000", "1000000", "100000"},
		{"100000", "1000000", "1000000"},
	}

	dprOptions := []int64{
		0, // 100, 300,
	}

	type UpdateAmountOptions struct {
		sendAmount, receiveAmount string
	}

	honestUpdateAmountOptions := map[int64][]UpdateAmountOptions{
		0: []UpdateAmountOptions{
			{"800000", "85000"},
			{"100000", "10000"},
			{"80000", "100000"},

			{"80000", "950000"},
			{"10000", "100000"},
			{"8000", "1200000"},
		},
	}

	maliciousUpdateAmountOptions := map[int64][]UpdateAmountOptions{
		0: []UpdateAmountOptions{
			{"1200000", "100000"},
			{"1000000", "80000"},
			{"900000", "90000"},

			{"100000", "800000"},
			{"120000", "1000000"},
			{"10000", "100000"},
		},
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
			for _, dprOption := range dprOptions {
				for _, initiationOption := range initiationOptions {
					delayInfo, _ := json.Marshal(TestDelayInfo{i})
					swap := swap.SwapBlob{
						ID:                   swap.SwapID(randomString()),
						SendToken:            tokenPairOption.sendToken,
						ReceiveToken:         tokenPairOption.receiveToken,
						SendAmount:           amountOption.sendAmount,
						ReceiveAmount:        amountOption.receiveAmount,
						MinimumReceiveAmount: amountOption.minReceiveAmount,
						ShouldInitiateFirst:  initiationOption,
						DelayPriceRange:      dprOption,
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

				swap.SendAmount = honestUpdateAmountOptions[swap.DelayPriceRange][testDelayInfo.Index].sendAmount
				swap.ReceiveAmount = honestUpdateAmountOptions[swap.DelayPriceRange][testDelayInfo.Index].receiveAmount

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

				swap.SendAmount = maliciousUpdateAmountOptions[swap.DelayPriceRange][testDelayInfo.Index].sendAmount
				swap.ReceiveAmount = maliciousUpdateAmountOptions[swap.DelayPriceRange][testDelayInfo.Index].receiveAmount

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
