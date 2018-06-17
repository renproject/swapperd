package http

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/republicprotocol/atom-go/utils"
	repCrypto "github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")

type boxHttpAdapter struct {
	signer repCrypto.Signer
}

func NewBoxHttpAdapter() BoxHttpAdapter {
	return &boxHttpAdapter{}
}

func (adapter *boxHttpAdapter) WhoAmI(challenge string) (WhoAmI, error) {

	version := getVersion()
	suppCurrencies := getSupportedCurrencies()

	boxInfo := BoxInfo{
		challenge:           challenge,
		version:             version,
		supportedCurrencies: suppCurrencies,
	}

	boxBytes, err := MarshalBoxInfo(boxInfo)
	if err != nil {
		return WhoAmI{}, err
	}
	boxHash := repCrypto.Keccak256(boxBytes)

	signature, err := adapter.signer.Sign(boxHash)
	if err != nil {
		return WhoAmI{}, err
	}

	sig65, err := utils.ToBytes65(signature)
	if err != nil {
		return WhoAmI{}, err
	}

	return WhoAmI{
		signature: MarshalSignature(sig65),
		boxInfo:   boxInfo,
	}, nil
}

func (adapter *boxHttpAdapter) PostOrder(order PostOrder) (PostOrder, error) {
	orderID, err := UnmarshalOrderID(order.OrderID)
	if err != nil {
		return PostOrder{}, err
	}
	sigIn, err := UnmarshalSignature(order.Signature)
	if err != nil {
		return PostOrder{}, err
	}

	err = validate(orderID, sigIn)
	if err != nil {
		return PostOrder{}, err
	}

	signOut, err := sign(orderID)
	if err != nil {
		return PostOrder{}, err
	}

	return PostOrder{
		order.OrderID,
		MarshalSignature(signOut),
	}, nil
}

func MarshalSignature(signatureIn [65]byte) string {
	return base64.StdEncoding.EncodeToString(signatureIn[:])
}

func UnmarshalSignature(signatureIn string) ([65]byte, error) {
	signature := [65]byte{}
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureIn)
	if err != nil {
		return signature, fmt.Errorf("cannot decode signature %v: %v", signatureIn, err)
	}
	if len(signatureBytes) != 65 {
		return signature, ErrInvalidSignatureLength
	}
	copy(signature[:], signatureBytes)
	return signature, nil
}

func MarshalOrderID(orderIDIn order.ID) string {
	return base64.StdEncoding.EncodeToString(orderIDIn[:])
}

func UnmarshalOrderID(orderIDIn string) (order.ID, error) {
	orderID := order.ID{}
	orderIDBytes, err := base64.StdEncoding.DecodeString(orderIDIn)
	if err != nil {
		return orderID, fmt.Errorf("cannot decode order id %v: %v", orderIDIn, err)
	}
	if len(orderIDBytes) != 32 {
		return orderID, ErrInvalidOrderIDLength
	}
	copy(orderID[:], orderIDBytes)
	return orderID, nil
}

func MarshalBoxInfo(boxInfo BoxInfo) ([]byte, error) {
	return json.Marshal(boxInfo)
}

func UnmarshalBoxInfo(boxInfo []byte) (BoxInfo, error) {
	var box BoxInfo

	err := json.Unmarshal(boxInfo, box)

	if err != nil {
		return BoxInfo{}, err
	}

	return box, nil
}

func getVersion() string {
	return ""
}
func getSupportedCurrencies() []string {
	return []string{}
}

func validate(id order.ID, signature [65]byte) error {
	return nil
}

func sign(id order.ID) ([65]byte, error) {
	return [65]byte{}, nil
}
