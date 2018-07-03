package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// opening of an order.
func NewServer(adapter BoxHttpAdapter) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/orders", PostOrdersHandler(adapter)).Methods("POST")
	r.HandleFunc("/whoami/{challenge}", WhoAmIHandler(adapter)).Methods("GET")
	r.Use(RecoveryHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}).Handler(r)
	return handler
}

// RecoveryHandler handles errors while processing the requests and populates the errors in the response
func RecoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("%v", r))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func PostOrdersHandler(boxHttpAdapter BoxHttpAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postOrder := PostOrder{}
		if err := json.NewDecoder(r.Body).Decode(&postOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into post order format: %v", err))
			return
		}

		processedOrder, err := boxHttpAdapter.PostOrder(postOrder)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot process the order: %v", err))
			return
		}

		processedOrderID, err := UnmarshalOrderID(processedOrder.OrderID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot process the order: %v", err))
			return
		}

		atomWatcher, err := boxHttpAdapter.BuildWatcher()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start the atomic swap process: %v", err))
			return
		}
		go atomWatcher.Run(processedOrderID)

		orderJSON, err := json.Marshal(processedOrder)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal the processed order: %v", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(orderJSON)
	}
}

func WhoAmIHandler(adapter BoxHttpAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		whoami, err := adapter.WhoAmI(params["challenge"])
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get the woami information: %v", err))
			return
		}

		whoamiJSON, err := json.Marshal(whoami)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal who am i information: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(whoamiJSON)
	}
}

func writeError(w http.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}
