package main

import (
	"encoding/json"
	"net/http"
)

type AtomicSwap struct {
	MyOrderID       string `json:"myorderid"`
	MatchingOrderID string `json:"matchingorderid"`
	Network         string `json:"network"`
	PrivateKey      string `json:"privatekey"`
	Value           string `json:"value"`
}

func PostAtomicSwap(w http.ResponseWriter, r *http.Request) {
	var swap AtomicSwap
	err := json.NewDecoder(r.Body).Decode(&swap)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(201)

	InitiateAtomicSwap(swap)
}
