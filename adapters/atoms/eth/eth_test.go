package eth_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/renex-swapper-go/adapters/atoms/eth"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/eth"
	config "github.com/republicprotocol/renex-swapper-go/adapters/configs/general"
	"github.com/republicprotocol/renex-swapper-go/adapters/key/eth"
	"github.com/republicprotocol/renex-swapper-go/adapters/owner"
	"github.com/republicprotocol/renex-swapper-go/services/swap"
)

var _ = Describe("ether", func() {
	var conn ethclient.Conn
	var aliceOrderID, bobOrderID [32]byte
	var aliceKey, bobKey swap.Key
	var value *big.Int
	var validity int64
	var secret, secretHash [32]byte
	var err error
	var reqAtom, reqAtomFailed swap.Atom
	var resAtom swap.Atom
	var data []byte
	var confPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/renex-swapper-go/secrets/configLocal.json"
	var ownPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/renex-swapper-go/secrets/owner.json"

	BeforeSuite(func() {
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		conn, err = ethclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		own, err := owner.LoadOwner(ownPath)
		Expect(err).ShouldNot(HaveOccurred())

		pk, err := crypto.HexToECDSA(own.Ganache)
		Expect(err).ShouldNot(HaveOccurred())
		owner := bind.NewKeyedTransactor(pk)

		alice, err := crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		aliceKey, err = eth.NewEthereumKey(hex.EncodeToString(crypto.FromECDSA(alice)), "ganache")
		Expect(err).ShouldNot(HaveOccurred())

		bob, err := crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		bobKey, err = eth.NewEthereumKey(hex.EncodeToString(crypto.FromECDSA(bob)), "ganache")
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddrBytes, err := aliceKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		bobAddrBytes, err := bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())

		err = conn.Transfer(common.BytesToAddress(aliceAddrBytes), owner, 1000000000000000000)
		Expect(err).ShouldNot(HaveOccurred())

		err = conn.Transfer(common.BytesToAddress(bobAddrBytes), owner, 1000000000000000000)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(err).ShouldNot(HaveOccurred())
		value = big.NewInt(10)
		validity = int64(time.Hour * 24)
		aliceOrderID[0] = 0x33
		bobOrderID[0] = 0x4a
		reqAtom, err = NewEthereumAtom(conn, aliceKey)
		Expect(err).ShouldNot(HaveOccurred())
		reqAtomFailed, err = NewEthereumAtom(conn, aliceKey)
		Expect(err).ShouldNot(HaveOccurred())
		resAtom, err = NewEthereumAtom(conn, bobKey)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can initiate an eth atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		bobAddr, err := bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		err = reqAtom.Initiate(bobAddr, secretHash, value, validity)
		Expect(err).ShouldNot(HaveOccurred())
		data, err = reqAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit an eth atomic swap", func() {
		err = resAtom.Deserialize(data)
		Expect(err).ShouldNot(HaveOccurred())
		bobAddr, err := bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		err = resAtom.Audit(secretHash, bobAddr, value, 60*60)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can redeem an eth atomic swap", func() {
		err = resAtom.Redeem(secret)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit secret after an eth atomic swap", func() {
		_secret, err := reqAtom.AuditSecret()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(_secret).Should(Equal(secret))
	})

	It("can refund an eth atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		bobAddr, err := bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Initiate(bobAddr, secretHash, value, 0)
		Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Refund()
		Expect(err).ShouldNot(HaveOccurred())
	})
})
