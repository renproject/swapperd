package bytesutils

import "errors"

// ToBytes32 converts dynamic byte array of size 32 to static byte array of size 32
func ToBytes32(b []byte) ([32]byte, error) {
	bytes32 := [32]byte{}
	if len(b) != 32 {
		return bytes32, errors.New("Length mismatch")
	}
	for i := range b {
		bytes32[i] = b[i]
	}
	return bytes32, nil
}
