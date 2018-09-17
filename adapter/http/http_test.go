package http_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/renex-swapper-go/adapter/http"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domains/tokens"
	configDriver "github.com/republicprotocol/renex-swapper-go/drivers/config"
	keystoreDriver "github.com/republicprotocol/renex-swapper-go/drivers/keystore"
	"github.com/republicprotocol/renex-swapper-go/services/watch"
)

type mockWatcher struct {
}

func (watcher *mockWatcher) Start() <-chan error {
	return nil
}

func (watcher *mockWatcher) Add([32]byte) error {
	return nil
}

func (watcher *mockWatcher) Status([32]byte) string {
	return "MOCK"
}

func (watcher *mockWatcher) Notify() {
}

func (watcher *mockWatcher) Stop() {
}

var _ = Describe("HTTP Adapter", func() {
	mockConfig := func() config.Config {
		return configDriver.New("nightly")
	}
	mockKeystore := func() keystore.Keystore {
		return keystoreDriver.GenerateRandom("nightly")
	}

	mockWatcher := func() watch.Watch {
		return &mockWatcher{}
	}

	generateRandomChallenge := func() string {
		challenge := make([]byte, 20)
		_, err := rand.Read(challenge)
		if err != nil {
			panic(err)
		}
		return base64.StdEncoding.EncodeToString(challenge)
	}

	generateRandomOrderID := func() string {
		orderID := [32]byte{}
		_, err := rand.Read(orderID[:])
		if err != nil {
			panic(err)
		}
		return MarshalOrderID(orderID)
	}

	Context("when the requestor is being honest", func() {
		It("should correctly respond to the requestor's challenge", func() {
			conf := mockConfig()
			ks := mockKeystore()
			watcher := mockWatcher()
			ethKey := ks.GetKey(tokens.TokenETH).(keystore.EthereumKey)
			addr := ethKey.Address.String()
			adapter := NewAdapter(conf, ks, watcher)
			challenge := generateRandomChallenge()
			signedResponse, err := adapter.WhoAmI(challenge)
			Expect(err).Should(BeNil())
			Expect(signedResponse.WhoAmI.AuthorizedAddresses).Should(Equal(conf.AuthorizedAddresses))
			Expect(signedResponse.WhoAmI.SupportedCurrencies).Should(Equal(conf.SupportedCurrencies))
			Expect(signedResponse.WhoAmI.Version).Should(Equal(conf.Version))
			Expect(signedResponse.WhoAmI.Challenge).Should(Equal(challenge))
			whoAmIBytes, err := MarshalWhoAmI(signedResponse.WhoAmI)
			Expect(err).Should(BeNil())
			hash := crypto.Keccak256(whoAmIBytes)
			sig, err := hex.DecodeString(signedResponse.Signature)
			Expect(err).Should(BeNil())
			marshalledPubKey, err := crypto.Ecrecover(hash, sig)
			Expect(err).Should(BeNil())
			ecdsaPubKey, err := crypto.UnmarshalPubkey(marshalledPubKey)
			Expect(err).Should(BeNil())
			address := crypto.PubkeyToAddress(*ecdsaPubKey)
			Expect(address.String()).Should(Equal(addr))
		})

		It("should correctly respond to the requestor's get balance request", func() {
			conf := mockConfig()
			ks := mockKeystore()
			watcher := mockWatcher()
			adapter := NewAdapter(conf, ks, watcher)
			balances, err := adapter.GetBalances()
			Expect(err).Should(BeNil())
			ethKey := ks.GetKey(tokens.TokenETH).(keystore.EthereumKey)
			btcKey := ks.GetKey(tokens.TokenBTC).(keystore.BitcoinKey)
			Expect(balances.Ethereum.Address).Should(Equal(ethKey.Address.String()))
			Expect(balances.Bitcoin.Address).Should(Equal(btcKey.Address.String()))
		})

		It("should correctly respond to the requestor's get status request", func() {
			conf := mockConfig()
			ks := mockKeystore()
			watcher := mockWatcher()
			adapter := NewAdapter(conf, ks, watcher)
			orderID := generateRandomOrderID()
			status, err := adapter.GetStatus(orderID)
			Expect(err).Should(BeNil())
			Expect(status.OrderID).Should(Equal(orderID))
			Expect(status.Status).Should(Equal("MOCK"))
		})

		It("should correctly respond to the requestor's post order request", func() {
			conf := mockConfig()
			ks := mockKeystore()
			watcher := mockWatcher()
			ethKey := ks.GetKey(tokens.TokenETH).(keystore.EthereumKey)
			conf.AuthorizedAddresses = append(
				conf.AuthorizedAddresses,
				ethKey.Address.String(),
			)
			adapter := NewAdapter(conf, ks, watcher)
			orderID := generateRandomOrderID()
			id, err := UnmarshalOrderID(orderID)
			Expect(err).Should(BeNil())
			message := append([]byte("Republic Protocol: open: "), id[:]...)
			hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))), message)
			sig, err := ethKey.Sign(hash)
			Expect(err).Should(BeNil())
			order := PostOrder{
				OrderID:   orderID,
				Signature: MarshalSignature(sig),
			}
			result, err := adapter.PostOrder(order)
			Expect(err).Should(BeNil())
			Expect(result.OrderID).Should(Equal(result.OrderID))
			sigBytes, err := hex.DecodeString(result.Signature)
			Expect(err).Should(BeNil())
			orderIDBytes, err := hex.DecodeString(orderID)
			Expect(err).Should(BeNil())
			marshalledPubKey, err := crypto.Ecrecover(orderIDBytes, sigBytes)
			Expect(err).Should(BeNil())
			ecdsaPubKey, err := crypto.UnmarshalPubkey(marshalledPubKey)
			Expect(err).Should(BeNil())
			address := crypto.PubkeyToAddress(*ecdsaPubKey)
			Expect(address.String()).Should(Equal(ethKey.Address.String()))
		})
	})

	Context("when the requestor is being dishonest", func() {
		It("should correctly respond to the requestor's get balance request", func() {

		})
	})
})
