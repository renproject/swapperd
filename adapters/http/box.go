package http

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	btcClient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethClient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/atom-go/services/watch"
	"github.com/republicprotocol/atom-go/utils"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")

type boxHttpAdapter struct {
	config config.Config
	keystr swap.Keystore
	watch  watch.Watch
}

func NewBoxHttpAdapter(config config.Config, keystr swap.Keystore, watcher watch.Watch) BoxHttpAdapter {
	return &boxHttpAdapter{
		config: config,
		keystr: keystr,
		watch:  watcher,
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

	keys, err := adapter.keystr.LoadKeys()
	if err != nil {
		return WhoAmI{}, err
	}
	key := keys[0].GetKey()

	signature, err := ethCrypto.Sign(boxHash, key)
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

	if err := adapter.watch.Add(orderID); err != nil {
		return PostOrder{}, err
	}
	adapter.watch.Notify()

	keys, err := adapter.keystr.LoadKeys()
	if err != nil {
		return PostOrder{}, err
	}
	key := keys[0].GetKey()

	if err != nil {
		return PostOrder{}, err
	}

	signOut, err := ethCrypto.Sign(orderID[:], key)

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

func (adapter *boxHttpAdapter) GetStatus(orderID string) (Status, error) {
	id, err := UnmarshalOrderID(orderID)
	if err != nil {
		return Status{}, err
	}

	status := adapter.watch.Status(id)

	return Status{
		OrderID: orderID,
		Status:  status,
	}, nil
}

func (adapter *boxHttpAdapter) GetBalances() (Balances, error) {
	balances := Balances{}
	keys, err := adapter.keystr.LoadKeys()
	if err != nil {
		return balances, err
	}
	for _, key := range keys {
		bal, err := getBalance(adapter.config, key)
		if err != nil {
			return balances, err
		}
		balances = append(balances, bal)
	}
	return balances, nil
}

func getBalance(conf config.Config, key swap.Key) (Balance, error) {
	switch key.PriorityCode() {
	case 0:
		return bitcoinBalance(conf, key)
	case 1:
		return ethereumBalance(conf, key)
	default:
		return Balance{}, errors.New("Unknown priority code")
	}
}

func bitcoinBalance(conf config.Config, key swap.Key) (Balance, error) {
	conn, err := btcClient.Connect(conf)
	if err != nil {
		return Balance{}, err
	}

	addr, err := key.GetAddress()
	if err != nil {
		return Balance{}, err
	}

	amt, err := conn.Client.GetBalance("*")
	if err != nil {
		fmt.Println(err)
		return Balance{}, err
	}

	return Balance{
		PriorityCode: key.PriorityCode(),
		Address:      string(addr),
		Amount:       uint64(amt.ToUnit(btcutil.AmountSatoshi)),
	}, nil
}

func ethereumBalance(conf config.Config, key swap.Key) (Balance, error) {
	conn, err := ethClient.Connect(conf)
	if err != nil {
		return Balance{}, err
	}
	addr, err := key.GetAddress()
	if err != nil {
		return Balance{}, err
	}

	address := common.BytesToAddress(addr)
	bal, err := conn.Client().PendingBalanceAt(context.Background(), address)
	if err != nil {
		return Balance{}, err
	}

	return Balance{
		PriorityCode: key.PriorityCode(),
		Address:      address.String(),
		Amount:       bal.Uint64(),
	}, nil
}

func MarshalSignature(signatureIn [65]byte) string {
	return hex.EncodeToString(signatureIn[:])
}

func UnmarshalSignature(signatureIn string) ([65]byte, error) {
	signature := [65]byte{}
	signatureBytes, err := hex.DecodeString(signatureIn)
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
	return hex.EncodeToString(orderIDIn[:])
}

func UnmarshalOrderID(orderIDIn string) ([32]byte, error) {
	orderID := [32]byte{}
	orderIDBytes, err := hex.DecodeString(orderIDIn)
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
	message := append([]byte("Republic Protocol: open: "), id[:]...)
	signatureData := ethCrypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))), message)

	upubKey, err := ethCrypto.Ecrecover(signatureData, signature[:])
	if err != nil {
		return err
	}

	ecdsaPubKey, err := ethCrypto.UnmarshalPubkey(upubKey)
	if err != nil {
		return err
	}
	addr := ethCrypto.PubkeyToAddress(*ecdsaPubKey)

	for _, j := range addresses {
		if j.String() == addr.String() {
			return nil
		}
	}
	return errors.New("Unauthorized Public Key")
}
