package server

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/rs/cors"
	"golang.org/x/crypto/sha3"
)

// NewHandler creates a new http handler
func NewHandler(username, passwordHash string, swaps chan<- swapper.Query, statusQueries chan<- status.Query) http.Handler {
	s := NewServer(swaps, statusQueries)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", postSwapsHandler(s, username, passwordHash)).Methods("POST")
	r.HandleFunc("/swaps", getSwapsHandler(s)).Methods("GET")
	r.HandleFunc("/balances", getBalancesHandler(s)).Methods("GET")
	r.HandleFunc("/ping", getPingHandler(s)).Methods("GET")
	r.Use(recoveryHandler)
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}).Handler(r)
	return handler
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

// getPingHandler handles the get ping request, it returns the basic information
// of the swapper such as the version and supported tokens.
func getPingHandler(server *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(server.GetPing()); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode ping response: %v", err))
			return
		}
	}
}

// getSwapsHandler handles the get swaps request, it returns the status of all
// the existing swaps on the swapper.
func getSwapsHandler(server *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(server.GetSwaps()); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swaps response: %v", err))
			return
		}
	}
}

// postSwapsHandler handles the post orders request, it fills incomplete
// information and starts the Atomic Swap.
func postSwapsHandler(server *server, username, passwordHash string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		passHash := sha3.Sum256([]byte(pass))
		passwordHashBytes, err := base64.StdEncoding.DecodeString(passwordHash)

		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode password hash: %v", err))
			return
		}

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare(passHash[:], passwordHashBytes) != 1 {
			writeError(w, http.StatusUnauthorized, fmt.Sprintf("incorrect username or password"))
			return
		}

		swapReq := PostSwapRequestResponse{}
		if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}

		swapRes, err := server.PostSwaps(swapReq, pass)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot execute swap: %v", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(swapRes); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
			return
		}
	}
}

// getBalancesHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func getBalancesHandler(server *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balancesRes, err := server.GetBalances()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get balances: %v", err))
			return
		}

		if err := json.NewEncoder(w).Encode(balancesRes); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode balances response: %v", err))
			return
		}
	}
}