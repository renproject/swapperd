package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/core/transfer"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	swapper.Storage
	transfer.Storage
}
type httpServer struct {
	wallet      wallet.Wallet
	logger      logrus.FieldLogger
	port        string
	loggedin    bool
	swapperTask tau.Task
	walletTask  tau.Task
}

func NewHttpServer(wallet wallet.Wallet, logger logrus.FieldLogger, swapperTask, walletTask tau.Task, port string) Server {
	return &httpServer{wallet, logger, port, false, swapperTask, walletTask}
}

// NewHttpListener creates a new http listener
func (listener *httpServer) Run(done <-chan struct{}) {
	go listener.swapperTask.Run(done)
	go listener.walletTask.Run(done)

	reqHandler := NewHandler(listener.swapperTask, listener.walletTask, listener.wallet, listener.logger)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", postSwapsHandler(reqHandler)).Methods("POST")
	r.HandleFunc("/swaps", getSwapsHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/swap", getSwapsHandler(reqHandler)).Queries("id", "{id}").Methods("GET")
	r.HandleFunc("/transfers", postTransfersHandler(reqHandler)).Methods("POST")
	r.HandleFunc("/transfers", getTransfersHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/balances", getBalancesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/balances/{token}", getBalancesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/addresses", getAddressesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/addresses/{token}", getAddressesHandler(reqHandler)).Methods("GET")
	r.HandleFunc("/bootload", postBootloadHandler(reqHandler)).Methods("POST")
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
	go listener.ticker(done)
	listener.logger.Info(fmt.Sprintf("listening for swaps on http://127.0.0.1:%s", listener.port))
	<-done
	httpListener.Close()
}

func (listener *httpServer) ticker(done <-chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			listener.walletTask.IO().InputWriter() <- tau.NewTick(time.Now())
			listener.swapperTask.IO().InputWriter() <- tau.NewTick(time.Now())
		}
	}
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
func postBootloadHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		if err := reqHandler.PostBootload(password); err != nil {
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
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		if err := json.NewEncoder(w).Encode(reqHandler.GetInfo(password)); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode info response: %v", err))
			return
		}
	}
}

// // getSwapsHandler handles the get swaps request, it returns the status of all
// // the existing swaps on the swapper.
// func getSwapsHandler(reqHandler Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		_, password, ok := r.BasicAuth()
// 		if !ok {
// 			writeError(w, http.StatusUnauthorized, "authentication required")
// 			return
// 		}

// 		resp, err := reqHandler.GetSwaps(password)
// 		if err != nil {
// 			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot get swaps: %v", err))
// 			return
// 		}

// 		if err := json.NewEncoder(w).Encode(resp); err != nil {
// 			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swaps response: %v", err))
// 			return
// 		}
// 	}
// }

// getSwapHandler handles the get swaps request, it returns the status of the
// existing swap with given a swap id on the swapper.
func getSwapsHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		swapID := r.FormValue("id")
		if swapID == "" {
			resp, err := reqHandler.GetSwaps(password)
			if err != nil {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot get swaps: %v", err))
				return
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swaps response: %v", err))
				return
			}
			return
		}

		resp, err := reqHandler.GetSwap(password, swap.SwapID(swapID))
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot get swap with id (%s): %v", swapID, err))
			return
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode get swap response: %v", err))
			return
		}
	}
}

// postSwapsHandler handles the post swaps request, it fills incomplete
// information and starts the Atomic Swap.
func postSwapsHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		swapReq := PostSwapRequest{}
		if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}
		swapReq.Password = password

		passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}
		swapReq.PasswordHash = base64.StdEncoding.EncodeToString(passwordHashBytes)

		if swapReq.Delay {
			if err := reqHandler.PostDelayedSwaps(swapReq); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		patchedSwap, err := reqHandler.PostSwaps(swapReq)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		if patchedSwap.Signature == "" {
			if err := json.NewEncoder(w).Encode(PostRedeemSwapResponse{patchedSwap.ID}); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
				return
			}
			return
		}

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
		transferResp.PasswordHash = ""

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

		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		balances, err := reqHandler.GetBalances(password)
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
		var transfers GetTransfersResponse
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		transfers, err := reqHandler.GetTransfers(password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get transfers: %v", err))
			return
		}

		if err := json.NewEncoder(w).Encode(transfers); err != nil {
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

		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		addresses, err := reqHandler.GetAddresses(password)
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

		if _, err := w.Write([]byte(addresses[token.Name])); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}

func getIDHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp GetIDResponse
		_, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		resp, err := reqHandler.GetID(password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get id: %v", err))
			return
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
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
