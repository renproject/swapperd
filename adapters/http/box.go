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

	btcClient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/btc"
	ethClient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/general"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/network"
	"github.com/republicprotocol/renex-swapper-go/services/watch"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")

type boxHttpAdapter struct {
	config  config.Config
	network network.Config
	keystr  keystore.Keystore
	watch   watch.Watch
}

func NewBoxHttpAdapter(config config.Config, network network.Config, keystr keystore.Keystore, watcher watch.Watch) BoxHTTPAdapter {
	return &boxHttpAdapter{
		config:  config,
		network: network,
		keystr:  keystr,
		watch:   watcher,
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

	ethKey, err := adapter.keystr.GetKey(1, 0)
	if err != nil {
		return WhoAmI{}, err
	}

	privKey, err := ethKey.GetKey()
	if err != nil {
		return WhoAmI{}, err
	}

	signature, err := ethCrypto.Sign(boxHash, privKey)
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

	go func() {
		if err := adapter.watch.Add(orderID); err != nil {
			return
		}
		adapter.watch.Notify()
	}()

	ethKey, err := adapter.keystr.GetKey(1, 0)
	if err != nil {
		return PostOrder{}, err
	}

	privKey, err := ethKey.GetKey()
	if err != nil {
		return PostOrder{}, err
	}

	signOut, err := ethCrypto.Sign(orderID[:], privKey)

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
	ethKey, err := adapter.keystr.GetKey(1, 0)
	if err != nil {
		return balances, err
	}
	ethBal, err := ethereumBalance(adapter.network, ethKey)
	if err != nil {
		return balances, err
	}
	balances = append(balances, ethBal)

	btcKey, err := adapter.keystr.GetKey(0, 0)
	if err != nil {
		return balances, err
	}
	btcBal, err := bitcoinBalance(adapter.network, btcKey)
	if err != nil {
		return balances, err
	}
	balances = append(balances, btcBal)
	return balances, nil
}

func bitcoinBalance(conf network.Config, key keystore.Key) (Balance, error) {
	conn, err := btcClient.Connect(conf)
	if err != nil {
		return Balance{}, err
	}

	addr, err := key.GetAddress()
	if err != nil {
		return Balance{}, err
	}

	amt, err := conn.Client.GetBalance("Atom")
	if err != nil {
		return Balance{}, err
	}

	return Balance{
		PriorityCode: key.PriorityCode(),
		Address:      string(addr),
		Amount:       uint64(amt.ToUnit(btcutil.AmountSatoshi)),
	}, nil
}

func ethereumBalance(conf network.Config, key keystore.Key) (Balance, error) {
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
