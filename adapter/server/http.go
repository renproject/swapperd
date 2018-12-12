package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	PutSwap(blob swap.SwapBlob) error
	PendingSwaps() ([]swap.SwapBlob, error)
}
type httpServer struct {
	wallet       wallet.Wallet
	storage      Storage
	logger       logrus.FieldLogger
	passwordHash []byte
	port         string
	loggedin     bool
}

func NewHttpServer(wallet wallet.Wallet, storage Storage, logger logrus.FieldLogger, passwordHash []byte, port string) Server {
	return &httpServer{wallet, storage, logger, passwordHash, port, false}
}

// NewHttpListener creates a new http listener
func (listener *httpServer) Run(doneCh <-chan struct{}, swaps, delayedSwaps chan<- swap.SwapBlob, receipts chan<- swap.SwapReceipt, statusQueries chan<- status.ReceiptQuery, balanceQueries chan<- balance.BalanceQuery) {
	reqHandler := NewHandler(listener.passwordHash, listener.wallet, listener.storage, listener.logger)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", postSwapsHandler(reqHandler, receipts, swaps, delayedSwaps, listener.logger)).Methods("POST")
	r.HandleFunc("/swaps", getSwapsHandler(reqHandler, statusQueries, listener.logger)).Methods("GET")
	r.HandleFunc("/transfers", postTransfersHandler(reqHandler, listener.logger)).Methods("POST")
	r.HandleFunc("/balances", getBalancesHandler(reqHandler, balanceQueries, listener.logger)).Methods("GET")
	r.HandleFunc("/balances/{token}", getBalancesHandler(reqHandler, balanceQueries, listener.logger)).Methods("GET")
	r.HandleFunc("/addresses", getAddressesHandler(reqHandler, listener.logger)).Methods("GET")
	r.HandleFunc("/addresses/{token}", getAddressesHandler(reqHandler, listener.logger)).Methods("GET")
	r.HandleFunc("/bootload", postBootloadHandler(reqHandler, swaps, delayedSwaps, listener.logger)).Methods("POST")
	r.HandleFunc("/info", getInfoHandler(reqHandler, listener.logger)).Methods("GET")
	r.Use(recoveryHandler)
	httpHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}).Handler(r)

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%s", listener.port))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := http.Serve(httpListener, httpHandler); err != nil {
			panic(err)
		}
	}()
	listener.logger.Info(fmt.Sprintf("listening for swaps on http://127.0.0.1:%s", listener.port))
	<-doneCh
	httpListener.Close()
}

// writeError response.
func writeError(w http.ResponseWriter, logger logrus.FieldLogger, statusCode int, err string) {
	logger.Error(err)
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}

// recoveryHandler handles errors while processing the requests and populates
// the errors in the response.
func recoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				writeError(w, logrus.StandardLogger(), http.StatusInternalServerError, fmt.Sprintf("%v", r))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// postBootloadHandler handles the post login request, it loads pending swaps and
// historical swap receipts into memory.
func postBootloadHandler(reqHandler Handler, swaps, delayedSwaps chan<- swap.SwapBlob, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, logger, http.StatusUnauthorized, "authentication required")
			return
		}

		if !reqHandler.VerifyPassword(password) {
			writeError(w, logger, http.StatusUnauthorized, "incorrect password")
			return
		}

		if err := reqHandler.PostBootload(password, swaps, delayedSwaps); err != nil {
			writeError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}
}

// getInfoHandler handles the get info request, it returns the basic information
// of the swapper such as the version, supported tokens addresses.
func getInfoHandler(reqHandler Handler, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(reqHandler.GetInfo()); err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode info response: %v", err))
			return
		}
	}
}

// getSwapsHandler handles the get swaps request, it returns the status of all
// the existing swaps on the swapper.
func getSwapsHandler(reqHandler Handler, statusQueries chan<- status.ReceiptQuery, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := reqHandler.GetSwaps(statusQueries)
		if err != nil {
			writeError(w, logger, http.StatusBadRequest, fmt.Sprintf("cannot get swaps: %v", err))
			return
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode swaps response: %v", err))
			return
		}
	}
}

// postSwapsHandler handles the post swaps request, it fills incomplete
// information and starts the Atomic Swap.
func postSwapsHandler(reqHandler Handler, receipts chan<- swap.SwapReceipt, swaps, delayedSwaps chan<- swap.SwapBlob, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, logger, http.StatusUnauthorized, "authentication required")
			return
		}

		if !reqHandler.VerifyPassword(password) {
			writeError(w, logger, http.StatusUnauthorized, "incorrect password")
			return
		}

		swapReq := PostSwapRequest{}
		if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
			writeError(w, logger, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}
		swapReq.Password = password

		if swapReq.Delay {
			if err := reqHandler.PostDelayedSwaps(swapReq, receipts, delayedSwaps); err != nil {
				writeError(w, logger, http.StatusBadRequest, err.Error())
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}

		patchedSwap, err := reqHandler.PostSwaps(swapReq, receipts, swaps)
		if err != nil {
			writeError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(patchedSwap); err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
			return
		}
	}
}

// postTransferHandler handles the post withdrawal
func postTransfersHandler(reqHandler Handler, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, logger, http.StatusUnauthorized, "authentication required")
			return
		}

		transferReq := PostTransfersRequest{}
		if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
			writeError(w, logger, http.StatusBadRequest, fmt.Sprintf("cannot decode transfers request: %v", err))
			return
		}
		transferReq.Password = password

		transferResp, err := reqHandler.PostTransfers(transferReq)
		if err != nil {
			writeError(w, logger, http.StatusBadRequest, fmt.Sprintf("cannot decode transfers request: %v", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(transferResp); err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode transfers response: %v", err))
			return
		}
	}
}

// getBalancesHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func getBalancesHandler(reqHandler Handler, balancesQuery chan<- balance.BalanceQuery, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := mux.Vars(r)
		tokenName := opts["token"]
		balances := reqHandler.GetBalances(balancesQuery)

		if tokenName == "" {
			if err := json.NewEncoder(w).Encode(balances); err != nil {
				writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
				return
			}
			return
		}

		token, err := blockchain.PatchToken(tokenName)
		if err != nil {
			writeError(w, logger, http.StatusBadRequest, fmt.Sprintf("invalid token name: %s", tokenName))
		}

		if err := json.NewEncoder(w).Encode(balances[token.Name]); err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}

// getAddressesHandler handles the get addresses request, and returns the addresses
// of the accounts held by the swapper.
func getAddressesHandler(reqHandler Handler, logger logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := mux.Vars(r)
		tokenName := opts["token"]
		addresses, err := reqHandler.GetAddresses()
		if err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot get addresses: %v", err))
			return
		}

		if tokenName == "" {
			if err := json.NewEncoder(w).Encode(addresses); err != nil {
				writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
				return
			}
			return
		}

		token, err := blockchain.PatchToken(tokenName)
		if err != nil {
			writeError(w, logger, http.StatusBadRequest, fmt.Sprintf("invalid token name: %s", tokenName))
		}

		if err := json.NewEncoder(w).Encode(addresses[token.Name]); err != nil {
			writeError(w, logger, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}
