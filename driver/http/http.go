package http

import (
	"encoding/json"
	"fmt"
	netHttp "net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/renex-swapper-go/adapter/http"
	"github.com/rs/cors"
)

// NewServer creates a new http handler
func NewServer(adapter http.Adapter) netHttp.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/orders", PostOrdersHandler(adapter)).Methods("POST")
	r.HandleFunc("/status/{orderID}", GetStatusHandler(adapter)).Methods("GET")
	r.HandleFunc("/whoami/{challenge}", WhoAmIHandler(adapter)).Methods("GET")
	r.HandleFunc("/balances", GetBalancesHandler(adapter)).Methods("GET")
	r.Use(RecoveryHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}).Handler(r)
	return handler
}

// RecoveryHandler handles errors while processing the requests and populates
// the errors in the response
func RecoveryHandler(h netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		defer func() {
			if r := recover(); r != nil {
				writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("%v", r))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// PostOrdersHandler handles post orders request, it gets the signed order id,
// checks whether the signer is authorized, if the signer is authorized this
// function adds the order id to the queue.
func PostOrdersHandler(adapter http.Adapter) netHttp.HandlerFunc {
	return func(w netHttp.ResponseWriter, r *netHttp.Request) {
		postOrder := http.PostOrder{}
		if err := json.NewDecoder(r.Body).Decode(&postOrder); err != nil {
			writeError(w, netHttp.StatusBadRequest, fmt.Sprintf("cannot decode json into post order format: %v", err))
			return
		}

		processedOrder, err := adapter.PostOrder(postOrder)
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot process the order: %v", err))
			return
		}

		orderJSON, err := json.Marshal(processedOrder)
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot marshal the processed order: %v", err))
			return
		}

		w.WriteHeader(netHttp.StatusCreated)
		w.Write(orderJSON)
	}
}

// WhoAmIHandler handles the get whoami request,it gets a challenge from the
// caller signs it and sends back the signed challenge with it's version
// information.
func WhoAmIHandler(adapter http.Adapter) netHttp.HandlerFunc {
	return func(w netHttp.ResponseWriter, r *netHttp.Request) {

		params := mux.Vars(r)
		whoami, err := adapter.WhoAmI(params["challenge"])
		if err != nil {

			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot get the whoami information: %v", err))
			return
		}
		whoamiJSON, err := json.Marshal(whoami)
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot marshal whoami information: %v", err))
			return
		}

		w.WriteHeader(netHttp.StatusOK)
		w.Write(whoamiJSON)
	}
}

// GetStatusHandler handles the get status request, it gets an order ID and
// returns the atomic swap status of the given order ID.
func GetStatusHandler(adapter http.Adapter) netHttp.HandlerFunc {
	return func(w netHttp.ResponseWriter, r *netHttp.Request) {
		params := mux.Vars(r)
		status, err := adapter.GetStatus(params["orderID"])
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot get the status information: %v", err))
			return
		}

		statusJSON, err := json.Marshal(status)
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot marshal status information: %v", err))
			return
		}

		w.WriteHeader(netHttp.StatusOK)
		w.Write(statusJSON)
	}
}

// GetBalancesHandler handles the get balance request, it returns the balances
// of the addresses in the atomic swapper.
func GetBalancesHandler(adapter http.Adapter) netHttp.HandlerFunc {
	return func(w netHttp.ResponseWriter, r *netHttp.Request) {
		balances, err := adapter.GetBalances()
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot get the balances: %v", err))
			return
		}

		balancesJSON, err := json.Marshal(balances)
		if err != nil {
			writeError(w, netHttp.StatusInternalServerError, fmt.Sprintf("cannot marshal the balance information: %v", err))
			return
		}

		w.WriteHeader(netHttp.StatusOK)
		w.Write(balancesJSON)
	}
}

// func RequestWhoAmI(addr string, challenge string) error {
// 	resp, err := netHttp.Get(fmt.Sprintf("https://" + addr + "/whoami/" + challenge))
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode == 200 {
// 		return nil
// 	}
// 	return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
// }

func writeError(w netHttp.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}
