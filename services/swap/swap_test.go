package swap_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	"github.com/republicprotocol/atom-go/drivers/eth/ganache"
	"github.com/republicprotocol/atom-go/services/axc"
	"github.com/republicprotocol/atom-go/services/network"
	"github.com/republicprotocol/atom-go/services/order"
	. "github.com/republicprotocol/atom-go/services/swap"

	btc "github.com/republicprotocol/atom-go/adapters/atoms/btc"
	eth "github.com/republicprotocol/atom-go/adapters/atoms/eth"
	ax "github.com/republicprotocol/atom-go/adapters/axc/mock"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	net "github.com/republicprotocol/atom-go/adapters/network/mock"
	ord "github.com/republicprotocol/atom-go/adapters/orders/mock"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap", func() {

	const CHAIN = "regtest"
	const RPC_USERNAME = "testuser"
	const RPC_PASSWORD = "testpassword"

	var aliceSwap, bobSwap Swap

	BeforeSuite(func() {

		var mockAXC axc.AXC
		var mockNetwork network.Network
		var aliceOrder, bobOrder order.Order
		var aliceOrderID, bobOrderID [32]byte
		var aliceSendValue, bobSendValue *big.Int
		var aliceRecieveValue, bobRecieveValue *big.Int
		var aliceCurrency, bobCurrency string
		var alice, bob *bind.TransactOpts
		var aliceBitcoinAddress, bobBitcoinAddress string
		var bobEthereumAddress common.Address
		var swapID [32]byte
		aliceOrderID[0] = 0x12
		bobOrderID[0] = 0x13

		swapID[0] = 0x14

		aliceCurrency = "ETHEREUM"
		bobCurrency = "BITCOIN"

		conn, err := ganache.Connect("http://localhost:8545")
		Expect(err).ShouldNot(HaveOccurred())

		alice, _, err = ganache.NewAccount(conn, big.NewInt(1000000000000000000))
		Expect(err).ShouldNot(HaveOccurred())
		alice.GasLimit = 3000000

		bob, bobEthereumAddress, err = ganache.NewAccount(conn, big.NewInt(1000000000000000000))
		Expect(err).ShouldNot(HaveOccurred())
		bob.GasLimit = 3000000

		time.Sleep(5 * time.Second)
		connection, err := btcclient.Connect("regtest", RPC_USERNAME, RPC_PASSWORD)
		Expect(err).ShouldNot(HaveOccurred())

		aliceSendValue = big.NewInt(10000000)
		bobSendValue = big.NewInt(10000000)

		aliceRecieveValue = big.NewInt(99990000)
		bobRecieveValue = big.NewInt(8000000)

		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		time.Sleep(5 * time.Second)

		aliceAddr, err := regtest.GetAddressForAccount(connection, "alice")
		Expect(err).ShouldNot(HaveOccurred())
		aliceBitcoinAddress = aliceAddr.EncodeAddress()

		bobAddr, err := regtest.GetAddressForAccount(connection, "bob")
		Expect(err).ShouldNot(HaveOccurred())
		bobBitcoinAddress = bobAddr.EncodeAddress()
		Expect(err).Should(BeNil())

		mockNetwork = net.NewMockNetwork()
		mockAXC = ax.NewMockAXC()
		aliceOrder = ord.NewMockOrder(aliceOrderID, bobOrderID, aliceSendValue, aliceRecieveValue, aliceCurrency, bobCurrency)
		bobOrder = ord.NewMockOrder(bobOrderID, aliceOrderID, bobSendValue, bobRecieveValue, bobCurrency, aliceCurrency)

		mockAXC.SetOwnerAddress(aliceOrderID, []byte(aliceBitcoinAddress))
		mockAXC.SetOwnerAddress(bobOrderID, bob.From.Bytes())

		reqAlice, err := eth.NewEthereumRequestAtom(context.Background(), conn, alice, bobEthereumAddress, swapID)
		Expect(err).Should(BeNil())

		reqBob := btc.NewBitcoinRequestAtom(connection, bobBitcoinAddress, aliceBitcoinAddress)
		resAlice := btc.NewBitcoinResponseAtom(connection, aliceBitcoinAddress, bobBitcoinAddress)

		resBob, err := eth.NewEthereumResponseAtom(context.Background(), conn, bob)
		Expect(err).Should(BeNil())

		aliceSwap = NewSwap(reqAlice, resAlice, mockAXC, aliceOrder, mockNetwork)
		bobSwap = NewSwap(reqBob, resBob, mockAXC, bobOrder, mockNetwork)

	})

	It("can do an eth - btc atomic swap", func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := aliceSwap.Execute()
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 1")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := bobSwap.Execute()
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 2")
		}()

		wg.Wait()
	})
})
