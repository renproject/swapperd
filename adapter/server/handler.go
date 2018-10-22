package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/rs/cors"
)

// NewHandler creates a new http handler
func NewHandler(swapQueries chan<- swapper.Query, statusQueries chan<- status.Query) http.Handler {
	s := NewServer(swapQueries, statusQueries)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", postSwapsHandler(s)).Methods("POST")
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
func postSwapsHandler(server *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swapReq := PostSwapRequestResponse{}
		if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode swap request: %v", err))
			return
		}

		swapRes, err := server.PostSwaps(swapReq)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot execute swap: %v", err))
			return
		}

		if err := json.NewEncoder(w).Encode(swapRes); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot encode swap response: %v", err))
			return
		}
		w.WriteHeader(http.StatusCreated)
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
