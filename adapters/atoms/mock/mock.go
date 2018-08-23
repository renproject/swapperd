package mock

// import (
// 	"encoding/json"
// 	"math/big"

// 	"github.com/republicprotocol/renex-swapper-go/services/swap"
// )

// const mockPriorityCode = 4294967295

// type mockData struct {
// }

// type mockAtom struct {
// 	data mockData
// }

// // NewMockAtom returns a new Mock Atom instance
// func NewMockAtom(key swap.Key) swap.Atom {
// 	return &mockAtom{}
// }

// // Initiate a new Atom swap by calling a function on ethereum
// func (atom *mockAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
// }

// // Refund an Atom swap by calling a function on ethereum
// func (atom *mockAtom) Refund() error {
// }

// // Redeem an Atom swap by calling a function on ethereum
// func (atom *mockAtom) Redeem(secret [32]byte) error {
// }

// // Audit an Atom swap by calling a function on ethereum
// func (atom *mockAtom) Audit(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
// }

// // AuditSecret audits the secret of an Atom swap by calling a function on ethereum
// func (atom *mockAtom) AuditSecret() ([32]byte, error) {
// }

// // Serialize serializes the atom details into a bytes array
// func (atom *mockAtom) Serialize() ([]byte, error) {
// 	b, err := json.Marshal(atom.data)
// 	return b, err
// }

// // Deserialize deserializes the atom details from a bytes array
// func (atom *mockAtom) Deserialize(b []byte) error {
// 	return json.Unmarshal(b, &atom.data)
// }

// // From returns the address of the caller
// func (atom *mockAtom) From() ([]byte, error) {
// 	return atom.key.GetAddress()
// }

// // PriorityCode returns the priority code of the currency.
// func (atom *mockAtom) PriorityCode() uint32 {
// 	return uint32(-1)
// }

// // GetKey returns the key of the atom.
// func (eth *mockAtom) GetKey() swap.Key {
// 	return eth.key
// }
