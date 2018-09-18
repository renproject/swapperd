package keystore_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/renex-swapper-go/drivers/keystore"
)

var _ = Describe("Keystore Adapter", func() {

	Context("when generating encrypted keystore files", func() {
		It("should generate a random ethereum keystore", func() {
			Expect(StoreKeyToFile("./keystore-ethereum.json", "ethereum", "kovan", "really secure passphrase")).To(BeNil())
		})

		It("should generate a random bitcoin keystore", func() {
			Expect(StoreKeyToFile("./keystore-bitcoin.json", "bitcoin", "testnet", "really secure passphrase")).To(BeNil())
		})
	})

	Context("when generating unencrypted keystore files", func() {
		It("should generate a random ethereum keystore", func() {
			Expect(StoreKeyToFile("./keystore-ethereum-unsafe.json", "ethereum", "kovan", "")).To(BeNil())
		})

		It("should generate a random bitcoin keystore", func() {
			Expect(StoreKeyToFile("./keystore-bitcoin-unsafe.json", "bitcoin", "testnet", "")).To(BeNil())
		})
	})

	Context("when decoding encrypted keystore files", func() {
		It("should decode an ethereum keystore file", func() {
			_, err := LoadKeyFromFile("./keystore-ethereum.json", "ethereum", "kovan", "really secure passphrase")
			Expect(err).To(BeNil())
		})

		It("should decode a bitcoin keystore file", func() {
			_, err := LoadKeyFromFile("./keystore-bitcoin.json", "bitcoin", "testnet", "really secure passphrase")
			Expect(err).To(BeNil())
		})
	})

	Context("when decoding unencrypted keystore files", func() {
		It("should decode an ethereum keystore file", func() {
			_, err := LoadKeyFromFile("./keystore-ethereum-unsafe.json", "ethereum", "kovan", "")
			Expect(err).To(BeNil())
		})

		It("should decode a bitcoin keystore file", func() {
			_, err := LoadKeyFromFile("./keystore-bitcoin-unsafe.json", "bitcoin", "testnet", "")
			Expect(err).To(BeNil())
		})
	})

	// Cleanup
	AfterSuite(func() {
		err := os.Remove("./keystore-ethereum.json")
		Expect(err).To(BeNil())
		err = os.Remove("./keystore-bitcoin.json")
		Expect(err).To(BeNil())
		err = os.Remove("./keystore-ethereum-unsafe.json")
		Expect(err).To(BeNil())
		err = os.Remove("./keystore-bitcoin-unsafe.json")
		Expect(err).To(BeNil())
	})
})
