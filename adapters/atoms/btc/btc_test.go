package btc_test

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/btcsuite/btcutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/renex-swapper-go/adapters/atoms/btc"
	btcclient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/btc"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/domains/order"
	"github.com/republicprotocol/renex-swapper-go/drivers/btc/regtest"
	"github.com/republicprotocol/renex-swapper-go/services/swap"
)

var _ = Describe("bitcoin", func() {

	var connection btcclient.Conn
	// var cmd *exec.Cmd
	var aliceAddr, bobAddr string // btcutil.Address
	var aliceAddrBytes, bobAddrBytes []byte
	var _aliceAddr, _bobAddr btcutil.Address
	var orderID, failedOrderID [32]byte

	var value *big.Int
	var validity int64
	var secret, secretHash [32]byte
	var err error
	var reqAtom, reqAtomFailed swap.Atom
	var resAtom swap.Atom
	var data []byte
	adapter := NewMockAdapter()

	BeforeSuite(func() {
		connection, err = btcclient.ConnectWithParams("regtest", "localhost:18443", "testuser", "testpassword")
		Expect(err).ShouldNot(HaveOccurred())

		rand.Read(orderID[:])
		rand.Read(failedOrderID[:])

		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		// time.Sleep(5 * time.Second)

		alicePrivKey, err := keystore.RandomBitcoinKeyString("regtest")
		Expect(err).ShouldNot(HaveOccurred())

		bobPrivKey, err := keystore.RandomBitcoinKeyString("regtest")
		Expect(err).ShouldNot(HaveOccurred())

		aliceKey, err := keystore.NewKey(alicePrivKey, 0, "regtest")
		Expect(err).ShouldNot(HaveOccurred())

		bobKey, err := keystore.NewKey(bobPrivKey, 0, "regtest")
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddrBytes, err = aliceKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		bobAddrBytes, err = bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddr = string(aliceAddrBytes)
		bobAddr = string(bobAddrBytes)

		_aliceAddr, err = btcutil.DecodeAddress(aliceAddr, connection.ChainParams)
		Expect(err).ShouldNot(HaveOccurred())
		_bobAddr, err = btcutil.DecodeAddress(bobAddr, connection.ChainParams)
		Expect(err).ShouldNot(HaveOccurred())

		btcvalue, err := btcutil.NewAmount(0.5)
		Expect(err).ShouldNot(HaveOccurred())

		connection.Client.SendToAddress(_aliceAddr, btcvalue)
		connection.Client.SendToAddress(_bobAddr, btcvalue)

		aliceWIF, err := btcutil.DecodeWIF(alicePrivKey)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKey(aliceWIF)
		Expect(err).ShouldNot(HaveOccurred())

		bobWIF, err := btcutil.DecodeWIF(bobPrivKey)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKey(bobWIF)
		Expect(err).ShouldNot(HaveOccurred())

		reqAtom = NewBitcoinAtom(&adapter, connection, aliceKey, orderID)
		reqAtomFailed = NewBitcoinAtom(&adapter, connection, aliceKey, failedOrderID)
		resAtom = NewBitcoinAtom(&adapter, connection, bobKey, orderID)

		value = big.NewInt(1000000)
		validity = time.Now().Unix() + 48*60*60
	})

	It("can initiate a btc atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtom.Initiate([]byte(bobAddr), secretHash, value, validity)
		Expect(err).ShouldNot(HaveOccurred())
		data, err = reqAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
		adapter.SendSwapDetails(order.ID(orderID), data)
	})

	It("can audit a btc atomic swap", func() {
		Expect(err).ShouldNot(HaveOccurred())
		_, _, _, _, err = resAtom.Audit()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can redeem a btc atomic swap", func() {
		before, err := connection.Client.GetReceivedByAddress(_bobAddr)
		Expect(err).ShouldNot(HaveOccurred())
		err = resAtom.Redeem(secret)
		Expect(err).ShouldNot(HaveOccurred())
		after, err := connection.Client.GetReceivedByAddress(_bobAddr)
		Expect(err).ShouldNot(HaveOccurred())
		data, err = resAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(after - before).Should(Equal(btcutil.Amount(990000)))
	})

	It("can audit secret after a btc atomic swap", func() {
		err = reqAtom.Deserialize(data)
		Expect(err).ShouldNot(HaveOccurred())
		_secret, err := reqAtom.AuditSecret()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(_secret).Should(Equal(secret))
	})

	It("can refund a btc atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtomFailed.Initiate([]byte(aliceAddr), secretHash, value, 0)
		Expect(err).ShouldNot(HaveOccurred())
		// before, err := connection.Client.GetReceivedByAddress(_aliceAddr)
		// Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Refund()
		Expect(err).ShouldNot(HaveOccurred())
		// after, err := connection.Client.GetReceivedByAddress(_aliceAddr)
		// Expect(err).ShouldNot(HaveOccurred())
		// Expect(after - before).Should(Equal(btcutil.Amount(990000)))
	})
})

type mockAdapter struct {
	swaps map[order.ID][]byte
}

func NewMockAdapter() mockAdapter {
	return mockAdapter{
		swaps: map[order.ID][]byte{},
	}
}

func (adapter *mockAdapter) ReceiveSwapDetails(orderID order.ID, waitTill int64) ([]byte, error) {
	return adapter.swaps[orderID], nil
}

func (adapter *mockAdapter) SendSwapDetails(orderID order.ID, details []byte) error {
	adapter.swaps[orderID] = details
	return nil
}
