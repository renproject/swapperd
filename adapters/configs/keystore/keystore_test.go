package keystore_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/configs/keystore"
)

var _ = Describe("Keystore", func() {
	var keystoreObj Keystore
	var err error

	BeforeSuite(func() {
		_, err := NewKeystore([]uint32{0, 1}, []string{"regtest", "ganache"}, "./local_keystore.json")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		err := os.Remove("./local_keystore.json")
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can load from a file", func() {
		keystoreObj, err = Load("./local_keystore.json")
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can get the ethereum key string", func() {
		ethKey, err := keystoreObj.GetKey(1, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_ = ethKey.GetKeyString()
	})

	It("can get the bitcoin key string", func() {
		btcKey, err := keystoreObj.GetKey(0, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_ = btcKey.GetKeyString()
	})

	It("can get the ethereum address", func() {
		ethKey, err := keystoreObj.GetKey(1, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = ethKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can get the bitcoin address", func() {
		btcKey, err := keystoreObj.GetKey(0, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = btcKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can get the ethereum ecdsa private key", func() {
		ethKey, err := keystoreObj.GetKey(1, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = ethKey.GetKey()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can get the bitcoin ecdsa private key", func() {
		btcKey, err := keystoreObj.GetKey(0, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = btcKey.GetKey()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can get the ethereum priority code", func() {
		ethKey, err := keystoreObj.GetKey(1, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_ = ethKey.PriorityCode()
	})

	It("can get the bitcoin priority code", func() {
		btcKey, err := keystoreObj.GetKey(0, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_ = btcKey.PriorityCode()
	})

	It("can get the ethereum priority code", func() {
		ethKey, err := keystoreObj.GetKey(1, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_ = ethKey.Chain()
	})

	It("can get the bitcoin priority code", func() {
		btcKey, err := keystoreObj.GetKey(0, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_ = btcKey.Chain()
	})
})
