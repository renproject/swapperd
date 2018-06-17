package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/orders", postOrders).Methods("POST")
	http.ListenAndServe(":18516", r)
}

func postOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
