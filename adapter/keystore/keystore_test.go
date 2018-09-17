package keystore_test

import (
	"github.com/btcsuite/btcd/btcec"

	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domains/tokens"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keystore Adapter", func() {

	randomBitcoinKey := func() (BitcoinKey, error) {
		priv, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			return BitcoinKey{}, err
		}
		return NewBitcoinKey(priv, "testnet")
	}

	randomEthereumKey := func() (EthereumKey, error) {
		priv, err := crypto.GenerateKey()
		if err != nil {
			return EthereumKey{}, err
		}
		return NewEthereumKey(priv, "kovan")
	}

	buildKeyMap := func() (KeyMap, error) {
		ethKey, err := randomEthereumKey()
		if err != nil {
			return nil, err
		}

		btcKey, err := randomBitcoinKey()
		if err != nil {
			return nil, err
		}
		keyMap := KeyMap{}
		keyMap[tokens.TokenBTC] = btcKey
		keyMap[tokens.TokenETH] = ethKey
		return keyMap, nil
	}

	buildKeystore := func() (Keystore, error) {
		keyMap, err := buildKeyMap()
		if err != nil {
			return nil, err
		}
		return New(keyMap), nil
	}

	Context("when creating random keys", func() {
		It("should generate a random bitcoin key", func() {
			_, err := randomBitcoinKey()
			Expect(err).To(BeNil())
		})

		It("should generate a random ethereum key", func() {
			_, err := randomEthereumKey()
			Expect(err).To(BeNil())
		})
	})

	Context("when retrieving keys from the keystore", func() {
		It("should not panic when type casted properly", func() {
			ks, err := buildKeystore()
			Expect(err).To(BeNil())
			Expect(func() { _ = ks.GetKey(tokens.TokenBTC).(BitcoinKey) }).ShouldNot(Panic())
			Expect(func() { _ = ks.GetKey(tokens.TokenETH).(EthereumKey) }).ShouldNot(Panic())
		})

		It("should panic when type casted improperly", func() {
			ks, err := buildKeystore()
			Expect(err).To(BeNil())
			Expect(func() { _ = ks.GetKey(tokens.TokenBTC).(EthereumKey) }).Should(Panic())
			Expect(func() { _ = ks.GetKey(tokens.TokenETH).(BitcoinKey) }).Should(Panic())
		})
	})
})
