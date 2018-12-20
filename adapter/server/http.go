package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/status"
	walletTask "github.com/republicprotocol/swapperd/core/wallet"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
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
func (listener *httpServer) Run(doneCh <-chan struct{}, swaps, delayedSwaps chan<- swap.SwapBlob, receipts chan<- swap.SwapReceipt, statusQueries chan<- status.ReceiptQuery, walletIO tau.IO) {
	walletIO.InputWriter() <- walletTask.Bootload{}
	reqHandler := NewHandler(listener.passwordHash, walletIO, listener.wallet, listener.storage, listener.logger)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", postSwapsHandler(reqHandler, receipts, swaps, delayedSwaps)).Methods("POST")
	r.HandleFunc("/swaps", getSwapsHandler(reqHandler, statusQueries)).Methods("GET")
	r.HandleFunc("/transfers", postTransfersHandler(reqHandler)).Methods("POST")
	r.HandleFunc("/transfers", getTransfersHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/balances", getBalancesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/balances/{token}", getBalancesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/addresses", getAddressesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/addresses/{token}", getAddressesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/bootload", postBootloadHandler(reqHandler, swaps, delayedSwaps)).Methods("POST")
	r.HandleFunc("/info", getInfoHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/id", getIDHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/sign/{type}", postSignatureHandler(reqHandler)).Methods("POST")
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
	w.Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", err)))
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
func postBootloadHandler(reqHandler Handler, swaps, delayedSwaps chan<- swap.SwapBlob) http.HandlerFunc {
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

		if err := reqHandler.PostBootload(password, swaps, delayedSwaps); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
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
func getSwapsHandler(reqHandler Handler, statusQueries chan<- status.ReceiptQuery) http.HandlerFunc {
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
func postSwapsHandler(reqHandler Handler, receipts chan<- swap.SwapReceipt, swaps, delayedSwaps chan<- swap.SwapBlob) http.HandlerFunc {
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

		if swapReq.Delay {
			if err := reqHandler.PostDelayedSwaps(swapReq, receipts, delayedSwaps); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}

		patchedSwap, err := reqHandler.PostSwaps(swapReq, receipts, swaps)
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
func getBalancesHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := mux.Vars(r)
		tokenName := opts["token"]
		balances := reqHandler.GetBalances()

		if tokenName == "" {
			if err := json.NewEncoder(w).Encode(balances); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
				return
			}
			return
		}

		token, err := blockchain.PatchToken(tokenName)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid token name: %s", tokenName))
		}

		if err := json.NewEncoder(w).Encode(balances[token.Name]); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}

// getTransfersHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func getTransfersHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(reqHandler.GetTransfers()); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}

// getAddressesHandler handles the get addresses request, and returns the addresses
// of the accounts held by the swapper.
func getAddressesHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := mux.Vars(r)
		tokenName := opts["token"]
		addresses, err := reqHandler.GetAddresses()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get addresses: %v", err))
			return
		}

		if tokenName == "" {
			if err := json.NewEncoder(w).Encode(addresses); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
				return
			}
			return
		}

		token, err := blockchain.PatchToken(tokenName)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid token name: %s", tokenName))
			return
		}

		if err := json.NewEncoder(w).Encode(addresses[token.Name]); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}

func getIDHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(reqHandler.GetID()); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode get id response: %v", err))
			return
		}
	}
}

func postSignatureHandler(reqHandler Handler) http.HandlerFunc {
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

		msg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		switch vars["type"] {
		case "json":
			resp, err := reqHandler.GetJSONSignature(password, msg)
			if err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign the message: %v", err))
				return
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode signature response: %v", err))
				return
			}
		case "base64":
			resp, err := reqHandler.GetBase64Signature(password, string(msg))
			if err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign the message: %v", err))
				return
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode signature response: %v", err))
				return
			}
		case "hex":
			resp, err := reqHandler.GetHexSignature(password, string(msg))
			if err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign the message: %v", err))
				return
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode signature response: %v", err))
				return
			}
		default:
			writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown message type: %s", vars["type"]))
			return
		}
	}
}
