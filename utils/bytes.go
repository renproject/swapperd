package utils

import (
	"errors"
)

// ToBytes32 converts dynamic byte array of size 32 to static byte array of size 32
func ToBytes32(b []byte) ([32]byte, error) {
	bytes32 := [32]byte{}
	for i := range b {
		bytes32[i] = b[i]
	}
	return bytes32, nil
}

// ToBytes65 converts dynamic byte array of size 65 to static byte array of size 65
func ToBytes65(b []byte) ([65]byte, error) {
	bytes65 := [65]byte{}
	if len(b) != 65 {
		return bytes65, errors.New("Length mismatch")
	}
	for i := range b {
		bytes65[i] = b[i]
	}
	return bytes65, nil
}
