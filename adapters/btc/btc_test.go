package btc_test

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/btc"
	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	"github.com/republicprotocol/atom-go/services/atom"
)

const CHAIN = "regtest"
const RPC_USERNAME = "testuser"
const RPC_PASSWORD = "testpassword"

var err error

func randomBytes32() [32]byte {
	randString := [32]byte{}
	_, err := rand.Read(randString[:])
	if err != nil {
		panic(err)
	}
	return randString
}

var _ = Describe("Bitcoin", func() {

	// Don't run on CI
	atom.LocalContext("atom swap", func() {

		var connection Connection
		// var cmd *exec.Cmd
		var aliceAddr, bobAddr string // btcutil.Address
		var _aliceAddr, _bobAddr btcutil.Address

		BeforeSuite(func() {
			var err error

			// cmd = regtest.Start()
			time.Sleep(5 * time.Second)

			connection, err = Connect("regtest", RPC_USERNAME, RPC_PASSWORD)
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				err = regtest.Mine(connection)
				Expect(err).ShouldNot(HaveOccurred())
			}()

			time.Sleep(5 * time.Second)

			_aliceAddr, err = regtest.GetAddressForAccount(connection, "alice")
			Expect(err).ShouldNot(HaveOccurred())
			aliceAddr = _aliceAddr.EncodeAddress()

			_bobAddr, err = regtest.GetAddressForAccount(connection, "bob")
			Expect(err).ShouldNot(HaveOccurred())
			bobAddr = _bobAddr.EncodeAddress()

			Ω(err).Should(BeNil())

			fmt.Println("Alice")
			aliceBalance, err := connection.Client.GetReceivedByAddress(_aliceAddr)
			Ω(err).Should(BeNil())
			fmt.Println(aliceAddr, aliceBalance)

			fmt.Println("Bob")
			bobBalance, err := connection.Client.GetReceivedByAddress(_bobAddr)
			Ω(err).Should(BeNil())
			fmt.Println(bobAddr, bobBalance)
		})

		AfterSuite(func() {
			connection.Shutdown()
			// regtest.Stop(cmd)
		})

		It("can initiate a bitcoin atomic swap", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			err := BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(300000000), time.Now().Unix()+10000)
			Ω(err).Should(BeNil())
		})

		It("can redeem a bitcoin atomic swap with correct secret", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			err = BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(300000000), time.Now().Unix()+10000)
			Ω(err).Should(BeNil())
			err = BTCAtom.Redeem(secret)
			Ω(err).Should(BeNil())
		})

		It("cannot redeem a bitcoin atomic swap with a wrong secret", func() {
			secret := randomBytes32()
			wrongSecret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			err = BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(3000000), time.Now().Unix()+10000)
			Ω(err).Should(BeNil())
			err = BTCAtom.Redeem(wrongSecret)
			Ω(err).Should(Not(BeNil()))
		})

		It("can read a bitcoin atomic swap", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			to := []byte(aliceAddr)
			from := []byte(bobAddr)
			value := big.NewInt(300000000)
			expiry := time.Now().Unix() + 10000
			err := BTCAtom.Initiate(hashLock, from, to, value, expiry)
			Ω(err).Should(BeNil())
			readHashLock, readFrom, readTo, readValue, readExpiry, readErr := BTCAtom.Audit()
			Ω(readErr).Should(BeNil())
			Ω(readHashLock).Should(Equal(hashLock))
			Ω(readTo).Should(Equal(to))
			Ω(readFrom).Should(Equal(from))
			Ω(readValue).Should(Equal(value))
			Ω(readExpiry).Should(Equal(expiry))
		})

		It("can read the correct secret from a bitcoin atomic swap", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			err := BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(300000000), time.Now().Unix()+10000)
			Ω(err).Should(BeNil())
			err = BTCAtom.Redeem(secret)
			Ω(err).Should(BeNil())
			readSecret, err := BTCAtom.AuditSecret()
			Ω(err).Should(BeNil())
			Ω(readSecret).Should(Equal(secret))
		})

		It("cannot refund a bitcoin atomic swap before expiry", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			err := BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(300000000), time.Now().Unix()+10000)
			Ω(err).Should(BeNil())
			err = BTCAtom.Refund()
			Ω(err).Should(Not(BeNil()))
		})

		It("can refund a bitcoin atomic swap", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret[:])
			BTCAtom := NewBitcoinAtom(connection)
			err := BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(300000000), time.Now().Unix()+10)
			Ω(err).Should(BeNil())
			time.Sleep(30 * time.Second)
			err = BTCAtom.Refund()
			Ω(err).Should(BeNil())
		})
	})
})
