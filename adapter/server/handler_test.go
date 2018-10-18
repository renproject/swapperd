package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/foundation"
)

var _ = Describe("Server Adapter", func() {
	Context("when the requestor is being right", func() {
		It("should correctly respond to a get ping request", func() {
			swapCh := make(chan foundation.Swap)
			server := NewServer(swapCh)
			pingResp := server.GetPing()
			Expect(pingResp.Version).Should(Equal("0.1.0"))
			Expect(pingResp.SupportedTokens[0]).Should(Equal(foundation.TokenBTC))
			Expect(pingResp.SupportedTokens[1]).Should(Equal(foundation.TokenETH))
			Expect(pingResp.SupportedTokens[2]).Should(Equal(foundation.TokenWBTC))
		})

		// 	It("should correctly respond to the requestor's get balance request", func() {
		// 		conf := mockConfig()
		// 		ks := mockKeystore()
		// 		watcher := mockWatcher()
		// 		adapter := NewAdapter(conf, ks, watcher)
		// 		balances, err := adapter.GetBalances()
		// 		Expect(err).Should(BeNil())
		// 		ethKey := ks.GetKey(tokens.TokenETH).(keystore.EthereumKey)
		// 		btcKey := ks.GetKey(tokens.TokenBTC).(keystore.BitcoinKey)
		// 		Expect(balances.Ethereum.Address).Should(Equal(ethKey.Address.String()))
		// 		Expect(balances.Bitcoin.Address).Should(Equal(btcKey.Address.String()))
		// 	})

		// 	It("should correctly respond to the requestor's get status request", func() {
		// 		conf := mockConfig()
		// 		ks := mockKeystore()
		// 		watcher := mockWatcher()
		// 		adapter := NewAdapter(conf, ks, watcher)
		// 		orderID := generateRandomOrderID()
		// 		status, err := adapter.GetStatus(orderID)
		// 		Expect(err).Should(BeNil())
		// 		Expect(status.OrderID).Should(Equal(orderID))
		// 		Expect(status.Status).Should(Equal("MOCK"))
		// 	})

		// 	It("should correctly respond to the requestor's post order request", func() {
		// 		conf := mockConfig()
		// 		ks := mockKeystore()
		// 		watcher := mockWatcher()
		// 		ethKey := ks.GetKey(tokens.TokenETH).(keystore.EthereumKey)
		// 		conf.AuthorizedAddresses = append(
		// 			conf.AuthorizedAddresses,
		// 			ethKey.Address.String(),
		// 		)
		// 		adapter := NewAdapter(conf, ks, watcher)
		// 		orderID := generateRandomOrderID()
		// 		id, err := UnmarshalOrderID(orderID)
		// 		Expect(err).Should(BeNil())
		// 		message := append([]byte("Republic Protocol: open: "), id[:]...)
		// 		hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))), message)
		// 		sig, err := ethKey.Sign(hash)
		// 		Expect(err).Should(BeNil())
		// 		order := PostOrder{
		// 			OrderID:   orderID,
		// 			Signature: MarshalSignature(sig),
		// 		}
		// 		result, err := adapter.PostOrder(order)
		// 		Expect(err).Should(BeNil())
		// 		Expect(result.OrderID).Should(Equal(result.OrderID))
		// 		sigBytes, err := hex.DecodeString(result.Signature)
		// 		Expect(err).Should(BeNil())
		// 		orderIDBytes, err := hex.DecodeString(orderID)
		// 		Expect(err).Should(BeNil())
		// 		marshalledPubKey, err := crypto.Ecrecover(orderIDBytes, sigBytes)
		// 		Expect(err).Should(BeNil())
		// 		ecdsaPubKey, err := crypto.UnmarshalPubkey(marshalledPubKey)
		// 		Expect(err).Should(BeNil())
		// 		address := crypto.PubkeyToAddress(*ecdsaPubKey)
		// 		Expect(address.String()).Should(Equal(ethKey.Address.String()))
		// 	})
		// })

		// Context("when the requestor is being dishonest", func() {
		// 	It("should correctly respond to the requestor's get balance request", func() {

		// 	})
	})
})
