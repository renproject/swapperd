package http_test

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	. "github.com/republicprotocol/atom-go/adapters/http"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/drivers/btc/regtest"

	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	wal "github.com/republicprotocol/atom-go/adapters/wallet/eth"
	"github.com/republicprotocol/atom-go/domains/match"
)

var _ = Describe("HTTP", func() {

	var aliceOrderID, bobOrderID [32]byte

	rand.Read(aliceOrderID[:])
	rand.Read(bobOrderID[:])

	BeforeSuite(func() {
		var confPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/config.json"
		var ksPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/keystore.json"
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		keystore := keystore.NewKeystore(ksPath)
		key, err := keystore.LoadKeypair("ethereum")
		Expect(err).ShouldNot(HaveOccurred())

		ganache, err := ethclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		connection, err := btcclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())
		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		time.Sleep(5 * time.Second)

		box, err := NewBoxHttpAdapter(config, key)

		aliceInfo, err := ax.NewEtereumAtomInfo(ganache, alice)
		Expect(err).Should(BeNil())

		bobInfo, err := ax.NewEtereumAtomInfo(ganache, bob)
		Expect(err).Should(BeNil())

		atomMatch := match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, bobSendValue, aliceCurrency, bobCurrency)
		mockWallet, err := wal.NewEthereumWallet(ganache, *owner)
		Expect(err).Should(BeNil())

		err = mockWallet.SetMatch(atomMatch)
		Expect(err).Should(BeNil())

		aliceInfo.SetOwnerAddress(aliceOrderID, []byte(aliceBitcoinAddress))
		bobInfo.SetOwnerAddress(bobOrderID, bob.From.Bytes())

	})

	It("can do an eth - btc atomic swap", func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := aliceWatch.Run(aliceOrderID)
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 1")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := bobWatch.Run(bobOrderID)
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 2")
		}()

		wg.Wait()
	})
})
