package main

import (
	"bytes"
	"errors"
)

func goFirst(myOrderID, matchingOrderID []byte) bool {
	status := bytes.Compare(myOrderID, matchingOrderID)
	if status == -1 {
		return true
	}
	return false
}

func toByte32(b []byte) ([32]byte, error) {
	bytes32 := [32]byte{}
	if len(b) != 32 {
		return bytes32, errors.New("Deserialization failed due to malformed input")
	}
	for i := range b {
		bytes32[i] = b[i]
	}
	return bytes32, nil
}
