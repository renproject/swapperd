package ethclient

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
)

// Conn is the ethereum client connection object
type Conn struct {
	Network        string
	Client         *ethclient.Client
	AtomAddress    common.Address
	NetworkAddress common.Address
	WalletAddress  common.Address
	InfoAddress    common.Address
}

// Connect to an ethereum network.
func Connect(chain string) (Conn, error) {
	var config Config

	config, err := readConfig(chain)
	ethclient, err := ethclient.Dial(config.URL)
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		Network:        config.Chain,
		Client:         ethclient,
		AtomAddress:    common.HexToAddress(config.AtomAddress),
		NetworkAddress: common.HexToAddress(config.NetworkAddress),
		InfoAddress:    common.HexToAddress(config.InfoAddress),
		WalletAddress:  common.HexToAddress(config.WalletAddress),
	}, nil
}

// NewConn Deploys all the contracts to the given ethereum network and creates a connection.
func NewConn(chain string) (Conn, error) {

	var conn Conn

	keyStore, err := readKeyStore(chain)
	if err != nil {
		return conn, err
	}

	ethclient, err := ethclient.Dial(keyStore.URL)
	if err != nil {
		return conn, err
	}

	conn = Conn{
		Network: chain,
		Client:  ethclient,
	}

	ownerECDSA, err := crypto.HexToECDSA(keyStore.PrivateKey)
	if err != nil {
		return conn, err
	}
	owner := bind.NewKeyedTransactor(ownerECDSA)

	// Deploy Atom contract
	AtomAddress, tx, _, err := bindings.DeployAtomSwap(owner, ethclient)
	if err != nil {
		return conn, err
	}

	_, err = conn.PatchedWaitDeployed(context.Background(), tx)
	if err != nil {
		return conn, err
	}

	// Deploy Network contract
	NetworkAddress, tx, _, err := bindings.DeployAtomNetwork(owner, ethclient)
	if err != nil {
		return conn, err
	}

	_, err = conn.PatchedWaitDeployed(context.Background(), tx)
	if err != nil {
		return conn, err
	}

	// Deploy Info contract
	InfoAddress, tx, _, err := bindings.DeployAtomInfo(owner, ethclient)
	if err != nil {
		return conn, err
	}

	_, err = conn.PatchedWaitDeployed(context.Background(), tx)
	if err != nil {
		return conn, err
	}

	// Deploy Wallet contract
	WalletAddress, tx, _, err := bindings.DeployAtomWallet(owner, ethclient)
	if err != nil {
		return conn, err
	}

	_, err = conn.PatchedWaitDeployed(context.Background(), tx)
	if err != nil {
		return conn, err
	}

	config := Config{
		Chain:          keyStore.Chain,
		URL:            keyStore.URL,
		AtomAddress:    AtomAddress.Hex(),
		NetworkAddress: NetworkAddress.Hex(),
		InfoAddress:    InfoAddress.Hex(),
		WalletAddress:  WalletAddress.Hex(),
	}

	err = writeConfig(config)

	return Conn{
		Network:        keyStore.Chain,
		Client:         ethclient,
		AtomAddress:    AtomAddress,
		NetworkAddress: NetworkAddress,
		WalletAddress:  WalletAddress,
		InfoAddress:    InfoAddress,
	}, err
}

// NewAccount creates a new account and funds it wit ether
func (b *Conn) NewAccount(value int64) (common.Address, *bind.TransactOpts, error) {
	account, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, &bind.TransactOpts{}, err
	}

	accountAddress := crypto.PubkeyToAddress(account.PublicKey)
	accountAuth := bind.NewKeyedTransactor(account)

	return accountAddress, accountAuth, b.Transfer(accountAddress, value)
}

// Transfer is a helper function for sending ETH to an address
func (b *Conn) Transfer(to common.Address, value int64) error {
	fromKeyStore, err := readKeyStore(b.Network)
	if err != nil {
		return err
	}

	fromECDSA, err := crypto.HexToECDSA(fromKeyStore.PrivateKey)
	if err != nil {
		return err
	}
	from := bind.NewKeyedTransactor(fromECDSA)

	transactor := &bind.TransactOpts{
		From:     from.From,
		Nonce:    from.Nonce,
		Signer:   from.Signer,
		Value:    big.NewInt(value),
		GasPrice: from.GasPrice,
		GasLimit: 30000,
		Context:  from.Context,
	}

	// Why is there no ethclient.Transfer?
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.Client, nil)
	tx, err := bound.Transfer(transactor)
	if err != nil {
		return err
	}
	_, err = b.PatchedWaitMined(context.Background(), tx)
	return err
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Conn) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	switch b.Network {
	case "ganache":
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	default:
		return bind.WaitMined(ctx, b.Client, tx)
	}
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Conn) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	switch b.Network {
	case "ganache":
		time.Sleep(100 * time.Millisecond)
		return common.Address{}, nil
	default:
		return bind.WaitDeployed(ctx, b.Client, tx)
	}
}
