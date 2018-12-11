package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/core/bootload"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	InsertSwap(swap swap.SwapBlob) error
	bootload.Storage
}
type httpServer struct {
	bootloader   bootload.Bootloader
	wallet       wallet.Wallet
	storage      Storage
	logger       logrus.FieldLogger
	passwordHash []byte
	port         string
	loggedin     bool
}

func NewHttpServer(wallet wallet.Wallet, storage Storage, logger logrus.FieldLogger, passwordHash []byte, port string) Server {
	return &httpServer{bootload.New(storage, logger), wallet, storage, logger, passwordHash, port, false}
}

// NewHttpListener creates a new http listener
func (listener *httpServer) Run(doneCh <-chan struct{}, swaps chan<- swap.SwapBlob, receipts chan<- swap.SwapReceipt, statusQueries chan<- swap.ReceiptQuery, balanceQueries chan<- balance.BalanceQuery) {
	reqHandler := NewHandler(listener.passwordHash, listener.wallet, listener.storage, listener.bootloader)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", postSwapsHandler(reqHandler, receipts, swaps)).Methods("POST")
	r.HandleFunc("/swaps", getSwapsHandler(reqHandler, statusQueries)).Methods("GET")
	r.HandleFunc("/transfers", postTransfersHandler(reqHandler)).Methods("POST")
	r.HandleFunc("/balances", getBalancesHandler(reqHandler, balanceQueries)).Methods("GET")
	r.HandleFunc("/bootload", postBootloadHandler(reqHandler, receipts, swaps)).Methods("POST")
	r.HandleFunc("/info", getInfoHandler(reqHandler)).Methods("GET")
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
func writeError(w http.ResponseWriter, statusCode int, err string) {
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
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("%v", r))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// postBootloadHandler handles the post login request, it loads pending swaps and
// historical swap receipts into memory.
func postBootloadHandler(reqHandler Handler, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		if !reqHandler.VerifyPassword(password) {
			writeError(w, http.StatusUnauthorized, "incorrect password")
			return
		}
		reqHandler.PostBootload(password, receipts, swaps)
		w.WriteHeader(http.StatusOK)
	}
}

// getInfoHandler handles the get info request, it returns the basic information
// of the swapper such as the version, supported tokens addresses.
func getInfoHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(reqHandler.GetInfo()); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode info response: %v", err))
			return
		}
	}
}

// getSwapsHandler handles the get swaps request, it returns the status of all
// the existing swaps on the swapper.
func getSwapsHandler(reqHandler Handler, statusQueries chan<- swap.ReceiptQuery) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := reqHandler.GetSwaps(statusQueries)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot get swaps: %v", err))
			return
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swaps response: %v", err))
			return
		}
	}
}

// postSwapsHandler handles the post swaps request, it fills incomplete
// information and starts the Atomic Swap.
func postSwapsHandler(reqHandler Handler, receipts chan<- swap.SwapReceipt, swapRequests chan<- swap.SwapBlob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		if !reqHandler.VerifyPassword(password) {
			writeError(w, http.StatusUnauthorized, "incorrect password")
			return
		}

		swapReq := PostSwapRequest{}
		if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}
		swapReq.Password = password

		patchedSwap, err := reqHandler.PostSwaps(swapReq, receipts, swapRequests)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(patchedSwap); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
			return
		}
	}
}

// postTransferHandler handles the post withdrawal
func postTransfersHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		transferReq := PostTransfersRequest{}
		if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode transfers request: %v", err))
			return
		}
		transferReq.Password = password

		transferResp, err := reqHandler.PostTransfers(transferReq)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode transfers request: %v", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(transferResp); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode transfers response: %v", err))
			return
		}
	}
}

// getBalancesHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func getBalancesHandler(reqHandler Handler, balancesQuery chan<- balance.BalanceQuery) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(reqHandler.GetBalances(balancesQuery)); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}
