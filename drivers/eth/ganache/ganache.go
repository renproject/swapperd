package ganache

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/atom-go/adapters/eth"
)

var genesisPrivateKey, genesisTransactor = genesis()

// GenesisPrivateKey used by Ganache.
func GenesisPrivateKey() *ecdsa.PrivateKey {
	return genesisPrivateKey
}

// GenesisTransactor used by Ganache.
func GenesisTransactor() *bind.TransactOpts {
	return genesisTransactor
}

// Start a local Ganache instance.
func Start() *exec.Cmd {
	cmd := exec.Command("ganache-cli", fmt.Sprintf("--account=0x%x,1000000000000000000000", crypto.FromECDSA(genesisPrivateKey)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd
}

// Connect to a local Ganache instance.
func Connect(ganacheRPC string) (eth.Connection, error) {
	ethclient, err := ethclient.Dial(ganacheRPC)
	if err != nil {
		return eth.Connection{}, err
	}

	_, ethAddress, err := deployEthSwap(context.Background(), eth.Connection{
		Client:  ethclient,
		Network: eth.NetworkGanache,
	}, genesisTransactor)

	if err != nil {
		return eth.Connection{}, err
	}

	return eth.Connection{
		Client:     ethclient,
		EthAddress: ethAddress,
		Network:    eth.NetworkGanache,
	}, nil
}

// DeployContracts to Ganache deploys REN and DNR contracts using the genesis private key
func DeployContracts(conn eth.Connection) error {
	return deployContracts(conn, genesisTransactor)
}

// DistributeEth transfers ETH to each of the addresses
func DistributeEth(conn eth.Connection, amount *big.Int, addresses ...common.Address) error {

	for _, address := range addresses {
		err := conn.TransferEth(context.Background(), genesisTransactor, address, amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewAccount(conn eth.Connection, amount *big.Int) (*bind.TransactOpts, common.Address, error) {
	ethereumPair, err := crypto.GenerateKey()
	if err != nil {
		return nil, common.Address{}, err
	}
	addr := crypto.PubkeyToAddress(ethereumPair.PublicKey)
	account := bind.NewKeyedTransactor(ethereumPair)
	if amount.Cmp(big.NewInt(0)) > 0 {
		if err := DistributeEth(conn, amount, addr); err != nil {
			return nil, common.Address{}, err
		}
	}

	return account, addr, nil
}

func genesis() (*ecdsa.PrivateKey, *bind.TransactOpts) {
	deployerKey, err := crypto.HexToECDSA("2aba04ee8a322b8648af2a784144181a0c793f1a2e80519418f3d20bbfb22249")
	if err != nil {
		log.Fatalf("cannot read genesis key: %v", err)
		return nil, nil
	}
	deployerAuth := bind.NewKeyedTransactor(deployerKey)
	return deployerKey, deployerAuth
}

func deployContracts(conn eth.Connection, transactor *bind.TransactOpts) error {
	_, _, err := deployEthSwap(context.Background(), conn, transactor)
	if err != nil {
		return err
	}
	return nil
}

func deployEthSwap(ctx context.Context, conn eth.Connection, auth *bind.TransactOpts) (*eth.Atom, common.Address, error) {
	address, tx, ethAtom, err := eth.DeployAtom(auth, conn.Client)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy Eth Atom contract: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return ethAtom, address, nil
}
