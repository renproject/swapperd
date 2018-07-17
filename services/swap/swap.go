package swap

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/utils"
	"github.com/republicprotocol/republic-go/order"
)

// Swap is the interface for an atomic swap object
type Swap interface {
	Execute() error
}

type swap struct {
	personalAtom Atom
	foreignAtom  Atom
	order        match.Match
	network      Network
	info         Info
	stateStr     store.State
}

// NewSwap returns a new Swap instance
func NewSwap(atom1 Atom, atom2 Atom, info Info, order match.Match, network Network, stateStr store.State) Swap {
	personalAtom := atom2
	foreignAtom := atom1

	if atom1.PriorityCode() == order.SendCurrency() {
		personalAtom = atom1
		foreignAtom = atom2
	}

	return &swap{
		personalAtom: personalAtom,
		foreignAtom:  foreignAtom,
		order:        order,
		info:         info,
		network:      network,
		stateStr:     stateStr,
	}
}

func (swap *swap) Execute() error {
	if swap.personalAtom.PriorityCode() < swap.foreignAtom.PriorityCode() {
		return swap.initiate()
	}
	return swap.respond()
}

func (swap *swap) initiate() error {
	type initInfo struct {
		PersonalAddr []byte   `json:"personalAddr"`
		Secret32     [32]byte `json:"secret32"`
		SecretHash   [32]byte `json:"secretHash"`
	}

	var info initInfo

	// TODO: Be defensive when checking statuses. Do not assume that an invalid
	// status cannot happen.
	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "INFO_SUBMITTED" {
		orderID := swap.order.PersonalOrderID()
		log.Println("Initiating the atomic swap for ", order.ID(orderID))

		var err error
		info.PersonalAddr, err = swap.info.GetOwnerAddress(swap.order.PersonalOrderID())
		if err != nil {
			return err
		}
		foreignAddr, err := swap.info.GetOwnerAddress(swap.order.ForeignOrderID())
		if err != nil {
			return err
		}

		expiry := time.Now().Add(48 * time.Hour).Unix()

		secret := make([]byte, 32)
		rand.Read(secret)

		info.Secret32, err = utils.ToBytes32(secret)
		if err != nil {
			return err
		}
		info.SecretHash = sha256.Sum256(secret)

		// TODO: Change all references to "str" to "store" because "str" means
		// string.
		initStr, err := json.Marshal(info)
		if err != nil {
			return err
		}

		if err := swap.stateStr.Write(append([]byte("INIT INFO:"), orderID[:]...), initStr); err != nil {
			return err
		}

		err = swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "INITIATED")
		if err != nil {
			return err
		}

		// TODO: Detect failure to initiate and retry, but ensure that retrying
		// never happens if funds were deposited. Self audit to see whether the
		// previous attempt actually *did* make it to the blockchain so that we
		// always initiate exactly once. (we can use the store to help with
		// this.)
		err = swap.personalAtom.Initiate(foreignAddr, info.SecretHash, swap.order.SendValue(), expiry)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		log.Println("Initiated the atomic swap for ", order.ID(orderID))

	} else {
		orderID := swap.order.PersonalOrderID()
		data, err := swap.stateStr.Read(append([]byte("INIT INFO:"), orderID[:]...))
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &info); err != nil {
			return err
		}
	}

	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "INITIATED" {
		personalSwapDetails, err := swap.personalAtom.Serialize()
		if err != nil {
			return err
		}
		err = swap.network.SendSwapDetails(swap.order.PersonalOrderID(), personalSwapDetails)
		if err != nil {
			return err
		}
		err = swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "WAITING_FOR_SWAP_DETAILS")
		if err != nil {
			return err
		}
	}

	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "WAITING_FOR_SWAP_DETAILS" {
		foreignSwapDetails, err := swap.network.ReceiveSwapDetails(swap.order.ForeignOrderID())
		if err != nil {
			return err
		}
		err = swap.foreignAtom.Deserialize(foreignSwapDetails)
		if err != nil {
			return err
		}

		// TODO: Persistent storage of the swap details is required in case of
		// crashing the system.
		err = swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "RECIEVED_SWAP_DETAILS")
		if err != nil {
			return err
		}
	}

	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "RECIEVED_SWAP_DETAILS" {
		// TODO: What does "Audit" do if we are already refunded? Its probably
		// correct for it to return an error. To prevent that being an issue
		// though, we need to have better state management.
		err := swap.foreignAtom.Audit(info.SecretHash, info.PersonalAddr, swap.order.ReceiveValue(), 60*60)
		if err != nil {
			orderID := swap.order.PersonalOrderID()
			log.Println("Refunding the atomic swap for ", order.ID(orderID))
			// TODO: Using prefix states (e.g. "REFUNDING" not "REFUNDED") will
			// allow a more tolerant, and easier to understand, state machine.
			err2 := swap.personalAtom.Refund()
			if err2 != nil {
				// Should never happen
				return err2
			}
			err2 = swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "REFUNDED")
			if err2 != nil {
				return err2
			}
			log.Println("Refunded the atomic swap for ", order.ID(orderID))
			return err
		}
		if err := swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "AUDITED"); err != nil {
			return err
		}
	}

	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "AUDITED" {

		orderID := swap.order.PersonalOrderID()
		log.Println("Redeeming the atomic swap for ", order.ID(orderID))

		err := swap.foreignAtom.Redeem(info.Secret32)
		if err != nil {
			return err
		}

		err = swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "REDEEMED")
		if err != nil {
			return err
		}

		log.Println("Redeemed the atomic swap for ", order.ID(orderID))
	}

	return nil
}

func (swap *swap) respond() error {
	var personalSwapDetails []byte
	orderID := swap.order.PersonalOrderID()

	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "INFO_SUBMITTED" {
		log.Println("Waiting for other trader to initiate for ", order.ID(orderID))

		personalAddr, err := swap.info.GetOwnerAddress(swap.order.PersonalOrderID())
		foreignAddr, err := swap.info.GetOwnerAddress(swap.order.ForeignOrderID())

		foreignSwapDetails, err := swap.network.ReceiveSwapDetails(swap.order.ForeignOrderID())
		if err != nil {
			return err
		}

		err = swap.foreignAtom.Deserialize(foreignSwapDetails)
		if err != nil {
			return err
		}

		expiry := time.Now().Add(24 * time.Hour).Unix()
		hash := swap.foreignAtom.GetSecretHash()

		log.Println("Auditing the swap for ", order.ID(orderID))

		err = swap.foreignAtom.Audit(hash, personalAddr, swap.order.ReceiveValue(), expiry)
		if err != nil {
			return err
		}

		log.Println("Initiating the swap for ", order.ID(orderID))

		if err := swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "INITIATED"); err != nil {
			return err
		}

		if err = swap.personalAtom.Initiate(foreignAddr, hash, swap.order.SendValue(), expiry); err != nil {
			return err
		}

		log.Println("Atomic Swap Initiated for", order.ID(orderID))

		personalSwapDetails, err = swap.personalAtom.Serialize()
		if err != nil {
			return err
		}

		orderID := swap.order.PersonalOrderID()
		swap.stateStr.Write(append([]byte("RESPOND INITIATED:"), orderID[:]...), personalSwapDetails)
	} else {
		var err error
		orderID := swap.order.PersonalOrderID()
		personalSwapDetails, err = swap.stateStr.Read(append([]byte("RESPOND INITIATED:"), orderID[:]...))
		if err != nil {
			return err
		}
	}

	var secret [32]byte
	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "INITIATED" {

		log.Println("Waiting for the other trader to redeem", order.ID(orderID))
		var err error
		if err := swap.network.SendSwapDetails(swap.order.PersonalOrderID(), personalSwapDetails); err != nil {
			return err
		}

		secret, err = swap.personalAtom.AuditSecret()
		if err != nil {
			log.Println("Refunding the atomic swap", order.ID(orderID))
			err1 := swap.personalAtom.Refund()
			if err1 != nil {
				// SHOULD NEVER HAPPEN
				return err
			}
			err1 = swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "REFUNDED")
			if err1 != nil {
				return err1
			}
			log.Println("Refunded the atomic swap", order.ID(orderID))

			return err
		}

		orderID := swap.order.PersonalOrderID()
		err = swap.stateStr.Write(append([]byte("secret:"), orderID[:]...), secret[:])
		if err != nil {
			return err
		}

		if err := swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "AUDITED"); err != nil {
			return err
		}
	} else {
		orderID := swap.order.PersonalOrderID()
		secretBytes, err := swap.stateStr.Read(append([]byte("secret:"), orderID[:]...))
		if err != nil {
			return err
		}
		secret, err = utils.ToBytes32(secretBytes)
		if err != nil {
			return err
		}
	}

	log.Println("Redeeming the atomic swap", order.ID(orderID))

	if swap.stateStr.ReadStatus(swap.order.PersonalOrderID()) == "AUDITED" {
		if err := swap.foreignAtom.Redeem(secret); err != nil {
			return err
		}

		if err := swap.stateStr.UpdateStatus(swap.order.PersonalOrderID(), "REDEEMED"); err != nil {
			return err
		}
	}

	log.Println("Redeemed the atomic swap", order.ID(orderID))

	return nil
}
