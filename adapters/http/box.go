package http

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	btcClient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethClient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"
	"github.com/republicprotocol/atom-go/adapters/store/leveldb"
	wal "github.com/republicprotocol/atom-go/adapters/wallet/eth"
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

func NewBoxHttpAdapter(config config.Config, keystr swap.Keystore) (BoxHttpAdapter, error) {
	watcher, err := buildWatcher(config, keystr)
	if err != nil {
		return nil, err
	}

	return &boxHttpAdapter{
		config: config,
		keystr: keystr,
		watch:  watcher,
	}, nil
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

	go func() {
		err := adapter.watch.Run(orderID)
		if err != nil {
			panic(err)
		}
	}()

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

func buildWatcher(config config.Config, kstr swap.Keystore) (watch.Watch, error) {
	ethConn, err := ethClient.Connect(config)
	if err != nil {
		return nil, err
	}

	btcConn, err := btcClient.Connect(config)
	if err != nil {
		return nil, err
	}

	keys, err := kstr.LoadKeys()
	if err != nil {
		return nil, err
	}

	ethKey := keys[0]
	btcKey := keys[1]

	_WIF, err := btcKey.GetKeyString()
	if err != nil {
		return nil, err
	}

	WIF, err := btcutil.DecodeWIF(_WIF)
	if err != nil {
		return nil, err
	}

	err = btcConn.Client.ImportPrivKey(WIF)
	if err != nil {
		return nil, err
	}

	owner := bind.NewKeyedTransactor(ethKey.GetKey())
	owner.GasLimit = 3000000

	ethNet, err := net.NewEthereumNetwork(ethConn, owner)
	if err != nil {
		return nil, err
	}

	ethInfo, err := ax.NewEtereumAtomInfo(ethConn, owner)
	if err != nil {
		return nil, err
	}

	ethWallet, err := wal.NewEthereumWallet(ethConn, *owner)
	if err != nil {
		return nil, err
	}

	ethAtom, err := eth.NewEthereumAtom(ethConn, ethKey)
	if err != nil {
		return nil, err
	}

	btcAtom := btc.NewBitcoinAtom(btcConn, btcKey)

	loc := config.StoreLocation()
	str := swap.NewSwapStore(leveldb.NewLDBStore(loc))

	watcher := watch.NewWatch(ethNet, ethInfo, ethWallet, ethAtom, btcAtom, str)
	return watcher, nil
}

func (adapter *boxHttpAdapter) GetStatus(orderID string) (Status, error) {
	id, err := UnmarshalOrderID(orderID)
	if err != nil {
		return Status{}, err
	}

	status, err := adapter.watch.Status(id)
	if err != nil {
		return Status{}, err
	}

	return Status{
		OrderID: orderID,
		Status:  status,
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
