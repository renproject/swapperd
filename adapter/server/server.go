package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/rs/cors"
)

// NewHandler creates a new http handler
func NewHandler(swapCh chan<- foundation.Swap) http.Handler {
	server := NewServer(swapCh)
	r := mux.NewRouter()
	r.HandleFunc("/swaps", server.postSwapsHandler()).Methods("POST")
	r.HandleFunc("/swaps", server.getSwapsHandler()).Methods("GET")
	r.HandleFunc("/balances", server.getBalancesHandler()).Methods("GET")
	r.HandleFunc("/ping", server.getPingHandler()).Methods("GET")
	r.Use(recoveryHandler)
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}).Handler(r)
	return handler
}

// recoveryHandler handles errors while processing the requests and populates
// the errors in the response
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
func (server *server) getPingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pingJSONBytes, err := json.Marshal(server.GetPing())
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal swapperd's configuration: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(pingJSONBytes)
	}
}

// getSwapsHandler handles the get swaps request, it returns the status of all
// the existing swaps on the swapper.
func (server *server) getSwapsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swapsJSON, err := server.GetSwaps()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get the statuses of existing swaps: %v", err))
			return
		}
		swapsJSONBytes, err := json.Marshal(swapsJSON)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal swapperd's configuration: %v", err))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(swapsJSONBytes)
	}
}

// postSwapsHandler handles the post orders request, it fills incomplete
// information and starts the Atomic Swap.
func (server *server) postSwapsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postSwap := PostSwapMessage{}
		if err := json.NewDecoder(r.Body).Decode(&postSwap); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into post swap format: %v", err))
			return
		}

		completedSwap, err := server.PostSwaps(postSwap)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot execute the swap: %v", err))
			return
		}

		swapJSONBBytes, err := json.Marshal(completedSwap)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal the processed swap: %v", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(swapJSONBBytes)
	}
}

// getBalancesHandler handles the get balances request, and returns the balances
// of the accounts held by the swapper.
func (server *server) getBalancesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balancesJSON, err := server.GetBalances()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get balances: %v", err))
			return
		}

		balancesJSONBytes, err := json.Marshal(balancesJSON)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal balances: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(balancesJSONBytes)
	}
}

func writeError(w http.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}
