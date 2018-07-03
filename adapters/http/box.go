package http

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/services/watch"
	"github.com/republicprotocol/atom-go/utils"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")

type boxHttpAdapter struct {
	config config.Config
	key    *ecdsa.PrivateKey
}

func NewBoxHttpAdapter(config config.Config, key *ecdsa.PrivateKey) BoxHttpAdapter {
	return &boxHttpAdapter{
		config: config,
		key:    key,
	}
}

func (adapter *boxHttpAdapter) WhoAmI(challenge string) (WhoAmI, error) {

	version := adapter.config.GetVersion()

	suppCurrencies := adapter.config.GetSupportedCurrencies()

	authorizedAddresses := adapter.config.AuthorizedAddresses

	boxInfo := BoxInfo{
		Challenge:           challenge,
		Version:             version,
		SupportedCurrencies: suppCurrencies,
		AuthorizedAddresses: authorizedAddresses,
	}

	boxBytes, err := MarshalBoxInfo(boxInfo)
	if err != nil {
		return WhoAmI{}, err
	}
	boxHash := ethCrypto.Keccak256(boxBytes)

	signature, err := ethCrypto.Sign(boxHash, adapter.key)
	if err != nil {
		return WhoAmI{}, err
	}

	sig65, err := utils.ToBytes65(signature)
	if err != nil {
		return WhoAmI{}, err
	}

	return WhoAmI{
		Signature: MarshalSignature(sig65),
		BoxInfo:   boxInfo,
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

	addrs := adapter.config.GetAuthorizedAddresses()

	err = validate(orderID, sigIn, addrs)
	if err != nil {
		return PostOrder{}, err
	}

	signOut, err := ethCrypto.Sign(orderID[:], adapter.key)

	if err != nil {
		return PostOrder{}, err
	}

	sig65, err := utils.ToBytes65(signOut)
	if err != nil {
		return PostOrder{}, err
	}

	return PostOrder{
		order.OrderID,
		MarshalSignature(sig65),
	}, nil
}

func (adapter *boxHttpAdapter) BuildWatcher() (watch.Watch, error) {
	ethConn, err := ethclient.Connect(adapter.config)
	if err != nil {
		return watch.Watch{}, err
	}

	btcConn, err := btcclient.Connect(config)
	if err != nil {
		return watch.Watch{}, err
	}

	ownerECDSA, err := keystore.LoadKeypair("ethereum")
	if err != nil {
		return watch.Watch{}, err
	}
	owner := bind.NewKeyedTransactor(ownerECDSA)
	owner.GasLimit = 3000000

	ethNet, err = net.NewEthereumNetwork(ethConn, owner)
	if err != nil {
		return watch.Watch{}, err
	}

	ethInfo, err = ax.NewEtereumAtomInfo(ethConn, owner)
	if err != nil {
		return watch.Watch{}, err
	}

	ethWallet, err := wal.NewEthereumWallet(ethConn, *owner)
	if err != nil {
		return watch.Watch{}, err
	}

	reqAtom, err := eth.NewEthereumRequestAtom(ethConn, owner)
	if err != nil {
		return watch.Watch{}, err
	}
	resAtom := btc.NewBitcoinAtomResponder(btcConn, bobBitcoinAddress)

	watcher = NewWatch(ethNet, ethInfo, ethWallet, reqAtom, resAtom)

	return watcher, nil
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

func MarshalOrderID(orderIDIn [32]byte) string {
	return base64.StdEncoding.EncodeToString(orderIDIn[:])
}

func UnmarshalOrderID(orderIDIn string) ([32]byte, error) {
	orderID := [32]byte{}
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

func validate(id [32]byte, signature [65]byte, addresses []common.Address) error {
	upubKey, err := ethCrypto.Ecrecover(id[:], signature[:])
	if err != nil {
		return err
	}
	addr := ethCrypto.PubkeyToAddress(*ethCrypto.ToECDSAPub(upubKey))

	for _, j := range addresses {
		if j == addr {
			return nil
		}
	}

	return errors.New("Unauthorized Public Key")
}
