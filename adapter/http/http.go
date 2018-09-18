package http

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")

type Adapter interface {
	WhoAmI(challenge string) (WhoAmISigned, error)
	PostOrder(order PostOrder) (PostOrder, error)
	GetStatus(orderID string) (Status, error)
	GetBalances() (Balances, error)
}

type adapter struct {
	config config.Config
	keystr keystore.Keystore
	renex  renex.RenEx
}

func NewAdapter(config config.Config, keystr keystore.Keystore, renExSwapper renex.RenEx) Adapter {
	return &adapter{
		config: config,
		keystr: keystr,
		renex:  renExSwapper,
	}
}

func (adapter *adapter) WhoAmI(challenge string) (WhoAmISigned, error) {
	whoAmI := NewWhoAmI(challenge, adapter.config)
	infoBytes, err := MarshalWhoAmI(whoAmI)
	if err != nil {
		return WhoAmISigned{}, err
	}
	infoHash := crypto.Keccak256(infoBytes)
	ethKey := adapter.keystr.GetKey(token.ETH).(keystore.EthereumKey)
	sig, err := ethKey.Sign(infoHash)
	return WhoAmISigned{
		Signature: MarshalSignature(sig),
		WhoAmI:    whoAmI,
	}, nil
}

func (adapter *adapter) PostOrder(order PostOrder) (PostOrder, error) {
	orderID, err := UnmarshalOrderID(order.OrderID)
	if err != nil {
		return PostOrder{}, err
	}
	if err := validate(orderID, order.Signature, adapter.config.AuthorizedAddresses); err != nil {
		return PostOrder{}, err
	}
	go func() {
		if err := adapter.renex.Add(orderID); err != nil {
			return
		}
		adapter.renex.Notify()
	}()
	key := adapter.keystr.GetKey(token.ETH).(keystore.EthereumKey)
	sig, err := key.Sign(orderID[:])
	return PostOrder{
		order.OrderID,
		MarshalSignature(sig),
	}, nil
}

func (adapter *adapter) GetStatus(orderID string) (Status, error) {
	id, err := UnmarshalOrderID(orderID)
	if err != nil {
		return Status{}, err
	}
	status := adapter.renex.Status(id)
	return Status{
		OrderID: orderID,
		Status:  string(status),
	}, nil
}

func (adapter *adapter) GetBalances() (Balances, error) {
	ethBal, err := ethereumBalance(
		adapter.config,
		adapter.keystr.GetKey(token.ETH).(keystore.EthereumKey),
	)
	if err != nil {
		return Balances{}, err
	}
	btcBal, err := bitcoinBalance(
		adapter.config,
		adapter.keystr.GetKey(token.BTC).(keystore.BitcoinKey),
	)
	if err != nil {
		return Balances{}, err
	}
	return Balances{
		Ethereum: ethBal,
		Bitcoin:  btcBal,
	}, nil
}

func bitcoinBalance(conf config.Config, key keystore.BitcoinKey) (Balance, error) {
	conn, err := btc.NewConnWithConfig(conf.Bitcoin)
	if err != nil {
		return Balance{}, err
	}
	balance, err := conn.Balance(key.AddressString)
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: key.AddressString,
		Amount:  uint64(balance),
	}, nil
}

func ethereumBalance(conf config.Config, key keystore.EthereumKey) (Balance, error) {
	conn, err := eth.Connect(conf.Ethereum)
	if err != nil {
		return Balance{}, err
	}
	bal, err := conn.Balance(key.Address)
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: key.Address.String(),
		Amount:  bal.Uint64(),
	}, nil
}

func validate(id [32]byte, signature string, addresses []string) error {
	sig, err := UnmarshalSignature(signature)
	if err != nil {
		return err
	}

	message := append([]byte("Republic Protocol: open: "), id[:]...)
	signatureData := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))), message)

	marshalledPubKey, err := crypto.Ecrecover(signatureData, sig)
	if err != nil {
		return err
	}

	ecdsaPubKey, err := crypto.UnmarshalPubkey(marshalledPubKey)
	if err != nil {
		return err
	}
	addr := crypto.PubkeyToAddress(*ecdsaPubKey)

	for _, address := range addresses {
		if address == addr.String() {
			return nil
		}
	}
	return errors.New("Unauthorized Public Key")
}
