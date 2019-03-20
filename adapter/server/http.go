package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/renproject/swapperd/adapter/wallet"
	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	Receipts() ([]swap.SwapReceipt, error)
	Transfers() ([]transfer.TransferReceipt, error)
}

type httpServer struct {
	port    string
	handler Handler
	logger  logrus.FieldLogger
}

func NewHttpServer(cap int, port, version string, receiver *Receiver, storage Storage, wallet wallet.Wallet, logger logrus.FieldLogger) Server {
	return &httpServer{port, NewHandler(cap, version, wallet, storage, receiver), logger}
}

// NewHttpListener creates a new http listener
func (server *httpServer) Run(done <-chan struct{}) {
	r := mux.NewRouter()
	r.HandleFunc("/swaps", server.postSwapsHandler(server.handler)).Methods("POST")
	r.HandleFunc("/swaps", server.getSwapsHandler(server.handler)).Methods("GET")
	// r.HandleFunc("/swaps/{id}", server.getSwapHandler(server.handler)).Methods("GET")
	r.HandleFunc("/transfers", server.postTransfersHandler(server.handler)).Methods("POST")
	r.HandleFunc("/transfers", server.getTransfersHandler(server.handler)).Methods("GET")
	r.HandleFunc("/balances", server.getBalancesHandler(server.handler)).Methods("GET")
	r.HandleFunc("/balances/{token}", server.getBalancesHandler(server.handler)).Methods("GET")
	r.HandleFunc("/addresses", server.getAddressesHandler(server.handler)).Methods("GET")
	r.HandleFunc("/addresses/{token}", server.getAddressesHandler(server.handler)).Methods("GET")
	r.HandleFunc("/info", server.getInfoHandler(server.handler)).Methods("GET")
	r.HandleFunc("/id/{type}", server.getIDHandler(server.handler)).Methods("GET")
	r.HandleFunc("/id", server.getIDHandler(server.handler)).Methods("GET")
	r.HandleFunc("/sign/{type}", server.postSignatureHandler(server.handler)).Methods("POST")
	r.Use(recoveryHandler)
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}).Handler(r)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.port))
	if err != nil {
		server.logger.Error(err)
	}
	go func() {
		if err := http.Serve(listener, handler); err != nil {
			server.logger.Error(err)
		}
	}()
	server.logger.Infof("swapperd started listening on http://127.0.0.1:%s", server.port)
	<-done
	server.logger.Infof("swapperd stopped listening on http://127.0.0.1:%s", server.port)
	listener.Close()
	server.handler.Shutdown()
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

// getInfoHandler handles the get info request, it returns the basic information
// of the swapper such as the version, supported tokens addresses.
func (server *httpServer) getInfoHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}
		respBytes, err := json.MarshalIndent(reqHandler.GetInfo(password), "\t", "")
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode info response: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusOK, respBytes)
	}
}

// getSwapHandler handles the get swap request, it returns the status of the
// existing swap with given swap id on the swapper.
func (server *httpServer) getSwapHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		swapID := r.FormValue("id")
		if swapID == "" {
			server.writeError(w, r, http.StatusBadRequest, "requires a swap id")
			return
		}

		resp, err := reqHandler.GetSwap(password, swap.SwapID(swapID))
		if err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot get swap with id (%s): %v", swapID, err))
			return
		}

		respBytes, err := json.MarshalIndent(resp, "\t", "")
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode get swap response: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusOK, respBytes)
	}
}

// getSwapsHandler handles the get swaps request, it returns the status of the
// existing swaps of the person calling it.
func (server *httpServer) getSwapsHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		resp, err := reqHandler.GetSwaps(password)
		if err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot get swaps: %v", err))
			return
		}

		respBytes, err := json.MarshalIndent(resp, "\t", "")
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode swaps response: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusOK, respBytes)
		return
	}
}

// postSwapsHandler handles the post swaps request, it fills incomplete
// information and starts the Atomic Swap.
func (server *httpServer) postSwapsHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		swapReq := PostSwapRequest{}
		if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}
		swapReq.Password = password
		if swapReq.Speed == blockchain.Nil {
			swapReq.Speed = blockchain.Fast
		}

		passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}
		swapReq.PasswordHash = base64.StdEncoding.EncodeToString(passwordHashBytes)

		if swapReq.Delay {
			if err := reqHandler.PostDelayedSwaps(swapReq); err != nil {
				server.writeError(w, r, http.StatusBadRequest, err.Error())
				return
			}
			server.writeResponse(w, r, http.StatusCreated, []byte{})
			return
		}
		patchedSwap, err := reqHandler.PostSwaps(swapReq)
		if err != nil {
			server.writeError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		if patchedSwap.Signature == "" {
			respBytes, err := json.MarshalIndent(PostRedeemSwapResponse{patchedSwap.ID}, "\t", "")
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
				return
			}
			server.writeResponse(w, r, http.StatusCreated, respBytes)
			return
		}

		respBytes, err := json.MarshalIndent(patchedSwap, "\t", "")
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusCreated, respBytes)
		return
	}
}

// postTransferHandler handles the post withdrawal
func (server *httpServer) postTransfersHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		transferReq := PostTransfersRequest{}
		if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot decode transfers request: %v", err))
			return
		}
		transferReq.Password = password

		if err := reqHandler.PostTransfers(transferReq); err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot decode transfers request: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusCreated, []byte{})
	}
}

// getBalancesHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func (server *httpServer) getBalancesHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := mux.Vars(r)
		tokenName := opts["token"]

		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		if tokenName == "" {
			balances, err := reqHandler.GetBalances(password)
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve balances: %v", err))
				return
			}

			respBytes, err := json.MarshalIndent(balances, "\t", "")
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
				return
			}
			server.writeResponse(w, r, http.StatusOK, respBytes)
			return
		}

		token := tokens.ParseToken(tokenName)
		if token == tokens.InvalidToken {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid token name: %s", tokenName))
		}

		balance, err := reqHandler.GetBalance(password, token)
		if err != nil {
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("unable to retrieve balance for token: %s", tokenName))
			return
		}

		respBytes, err := json.MarshalIndent(balance, "", "\t")
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusOK, respBytes)
	}
}

// getTransfersHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func (server *httpServer) getTransfersHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transfers GetTransfersResponse
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		transfers, err := reqHandler.GetTransfers(password)
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot get transfers: %v", err))
			return
		}

		respBytes, err := json.MarshalIndent(transfers, "\t", "")
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
		server.writeResponse(w, r, http.StatusOK, respBytes)
	}
}

// getAddressesHandler handles the get addresses request, and returns the addresses
// of the accounts held by the swapper.
func (server *httpServer) getAddressesHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := mux.Vars(r)
		tokenName := opts["token"]

		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		if tokenName == "" {
			addresses, err := reqHandler.GetAddresses(password)
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot get addresses: %v", err))
				return
			}
			respBytes, err := json.MarshalIndent(addresses, "\t", "")
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
				return
			}
			server.writeResponse(w, r, http.StatusOK, respBytes)
		} else {
			token := tokens.ParseToken(tokenName)
			if token == tokens.InvalidToken {
				server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid token name: %s", tokenName))
				return
			}
			address, err := reqHandler.GetAddress(password, token)
			if err != nil {
				server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("unable to retrieve address for token: %s", tokenName))
				return
			}

			server.writeResponse(w, r, http.StatusOK, []byte(address))
		}
	}
}

func (server *httpServer) getIDHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		vars := mux.Vars(r)
		idType := vars["type"]
		resp, err := reqHandler.GetID(password, idType)
		if err != nil {
			server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot get id: %v", err))
			return
		}

		server.writeResponse(w, r, http.StatusOK, []byte(resp))
	}
}

func (server *httpServer) postSignatureHandler(reqHandler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok {
			server.writeError(w, r, http.StatusUnauthorized, "authentication required")
			return
		}

		msg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			server.writeError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		switch vars["type"] {
		case "json":
			resp, err := reqHandler.GetJSONSignature(password, msg)
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("failed to sign the message: %v", err))
				return
			}
			respBytes, err := json.MarshalIndent(resp, "\t", "")
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode signature response: %v", err))
				return
			}
			server.writeResponse(w, r, http.StatusCreated, respBytes)
		case "base64":
			resp, err := reqHandler.GetBase64Signature(password, string(msg))
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("failed to sign the message: %v", err))
				return
			}
			respBytes, err := json.MarshalIndent(resp, "\t", "")
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode signature response: %v", err))
				return
			}
			server.writeResponse(w, r, http.StatusCreated, respBytes)
		case "hex":
			resp, err := reqHandler.GetHexSignature(password, string(msg))
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("failed to sign the message: %v", err))
				return
			}
			respBytes, err := json.MarshalIndent(resp, "\t", "")
			if err != nil {
				server.writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot encode signature response: %v", err))
				return
			}
			server.writeResponse(w, r, http.StatusCreated, respBytes)
		default:
			server.writeError(w, r, http.StatusBadRequest, fmt.Sprintf("unknown message type: %s", vars["type"]))
			return
		}
	}
}

func (server *httpServer) writeResponse(w http.ResponseWriter, r *http.Request, statusCode int, resp []byte) {
	logger := server.logger
	logger = logger.WithField("method", r.Method)
	logger = logger.WithField("url", r.URL.String())
	logger = logger.WithField("port", server.port)
	logger = logger.WithField("status", statusCode)
	logger.Info("successfully responded to a http request")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

func (server *httpServer) writeError(w http.ResponseWriter, r *http.Request, statusCode int, err string) {
	logger := server.logger
	logger = logger.WithField("method", r.Method)
	logger = logger.WithField("url", r.URL)
	logger = logger.WithField("port", server.port)
	logger = logger.WithField("status", statusCode)
	logger = logger.WithError(fmt.Errorf(err))
	logger.Warnf("failed to respond to a http request")
	writeError(w, statusCode, err)
}
