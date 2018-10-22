package utils

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/beth-go"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	bip39 "github.com/tyler-smith/go-bip39"

	"github.com/republicprotocol/libbtc-go"

	"github.com/onsi/ginkgo"
	"github.com/republicprotocol/swapperd/adapter/account"
)

// LocalContext allows you to mark a ginkgo context as being local-only.
// It won't run if the CI environment variable is true.
func LocalContext(description string, f func()) {
	var local bool

	ciEnv := os.Getenv("CI")
	ci, err := strconv.ParseBool(ciEnv)
	if err != nil {
		ci = false
	}

	// Assume tests are running locally if CI environment variable is not defined
	local = !ci

	if local {
		ginkgo.Context(description, f)
	} else {
		ginkgo.PContext(description, func() {
			ginkgo.It("SKIPPING LOCAL TESTS", func() {})
		})
	}
}

type TestKeystore struct {
	Mnemonic   string `json:"mnemonic"`
	Passphrase string `json:"passphrase"`
}

func LoadAccounts(loc string) (account.Accounts, error) {
	keystores := TestKeystore{}
	data, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &keystores); err != nil {
		return nil, err
	}
	seed := bip39.NewSeed(keystores.Mnemonic, keystores.Passphrase)
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}

	btcKey, err := loadKey(masterKey, 44, 1, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	ethKey, err := loadKey(masterKey, 44, 60, 0, 0)
	if err != nil {
		return nil, err
	}

	ethAccount, err := beth.NewAccount("https://kovan.infura.io", ethKey)
	if err != nil {
		return nil, err
	}

	if err := ethAccount.WriteAddress("ERC20:WBTC", common.HexToAddress("0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0")); err != nil {
		return nil, err
	}

	if err := ethAccount.WriteAddress("SWAPPER:WBTC", common.HexToAddress("0x2218fa20c33765e7e01671ee6aaca75fbaf3a974")); err != nil {
		return nil, err
	}

	btcAccount := libbtc.NewAccount(libbtc.NewBlockchainInfoClient("testnet"), btcKey)

	return account.New(
		btcAccount,
		ethAccount,
	), nil
}

func loadKey(key *hdkeychain.ExtendedKey, path ...uint32) (*ecdsa.PrivateKey, error) {
	var err error
	for _, val := range path {
		key, err = key.Child(val)
		if err != nil {
			return nil, err
		}
	}
	privKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privKey.ToECDSA(), nil
}
