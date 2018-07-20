package swap

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
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
	state        store.SwapState
}

// NewSwap returns a new Swap instance
func NewSwap(personalAtom Atom, foreignAtom Atom, info Info, order match.Match, network Network, state store.SwapState) Swap {
	return &swap{
		personalAtom: personalAtom,
		foreignAtom:  foreignAtom,
		order:        order,
		info:         info,
		network:      network,
		state:        state,
	}
}

func (swap *swap) Execute() error {
	if swap.personalAtom.PriorityCode() == swap.foreignAtom.PriorityCode() {
		return fmt.Errorf("Trying to swap between atoms with the same priority code %d and %d", swap.personalAtom.PriorityCode(), swap.foreignAtom.PriorityCode())
	}
	if swap.personalAtom.PriorityCode() < swap.foreignAtom.PriorityCode() {
		return swap.request()
	}
	return swap.respond()
}

func (swap *swap) request() error {
	log.Println("Requestor ", order.ID(swap.order.PersonalOrderID()))

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusInfoSubmitted {
		if err := swap.generateDetails(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping generate details")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusInitateDetailsAcquired {
		if err := swap.initiate(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping initiate")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusInitiated {
		if err := swap.sendDetails(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping send details")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusSentSwapDetails {
		if err := swap.recieveDetails(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping recieve details")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusRecievedSwapDetails {
		if err := swap.requestorAudit(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping requestor audit")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusAudited {
		if err := swap.redeem(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping redeem")
	}
	return nil
}

func (swap *swap) respond() error {
	log.Println("Responder ", order.ID(swap.order.PersonalOrderID()))

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusInfoSubmitted {
		if err := swap.recieveDetails(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping generate details")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusRecievedSwapDetails {
		if err := swap.responderAudit(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping audit")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusAudited {
		if err := swap.initiate(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping initiate")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusInitiated {
		if err := swap.sendDetails(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping send details")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusSentSwapDetails {
		if err := swap.getRedeemDetails(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping get redeem details audit")
	}

	if swap.state.Status(swap.order.PersonalOrderID()) == StatusRedeemDetailsAcquired {
		if err := swap.redeem(); err != nil {
			return err
		}
	} else {
		log.Println("Skipping redeem")
	}

	return nil
}

func (swap *swap) generateDetails() error {
	orderID := swap.order.PersonalOrderID()
	log.Println("Generating the swap details for ", order.ID(orderID))
	expiry := time.Now().Add(48 * time.Hour).Unix()
	secret := make([]byte, 32)
	rand.Read(secret)
	secret32, err := utils.ToBytes32(secret)
	if err != nil {
		return err
	}
	secretHash := sha256.Sum256(secret)

	if err := swap.state.PutInitiateDetails(orderID, expiry, secretHash); err != nil {
		return err
	}

	if err := swap.state.PutRedeemDetails(orderID, secret32); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusInitateDetailsAcquired); err != nil {
		return err
	}
	log.Println("Generated the swap details for ", order.ID(orderID))
	return nil
}

func (swap *swap) initiate() error {
	orderID := swap.order.PersonalOrderID()
	expiry, secretHash, err := swap.state.InitiateDetails(orderID)
	if err != nil {
		return err
	}
	log.Println("Initiating the swap for ", order.ID(orderID))

	foreignAddr, err := swap.info.GetOwnerAddress(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	if err = swap.personalAtom.Initiate(foreignAddr, secretHash, swap.order.SendValue(), expiry); err != nil {
		return err
	}

	if err := swap.personalAtom.Store(swap.state); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusInitiated); err != nil {
		return err
	}

	log.Println("Initiated the swap for", order.ID(orderID))
	return nil
}

func (swap *swap) sendDetails() error {
	orderID := swap.order.PersonalOrderID()
	log.Println("Sending the swap details for ", order.ID(orderID))
	personalAtomBytes, err := swap.state.AtomDetails(orderID)
	if err != nil {
		return err
	}
	if err := swap.network.SendSwapDetails(orderID, personalAtomBytes); err != nil {
		log.Println("Error Here")
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusSentSwapDetails); err != nil {
		return err
	}

	log.Println("Sent the swap details for ", order.ID(orderID))
	return nil
}

func (swap *swap) recieveDetails() error {
	personalOrderID := swap.order.PersonalOrderID()
	foreignOrderID := swap.order.ForeignOrderID()
	log.Println("Recieving the swap details for ", order.ID(personalOrderID))
	foreignAtomBytes, err := swap.network.ReceiveSwapDetails(foreignOrderID)

	if err != nil {
		return err
	}

	if err := swap.state.PutAtomDetails(foreignOrderID, foreignAtomBytes); err != nil {
		return err
	}

	if err := swap.state.PutStatus(personalOrderID, StatusRecievedSwapDetails); err != nil {
		return err
	}

	log.Println("Recieved the swap details for ", order.ID(personalOrderID))
	return nil
}

func (swap *swap) redeem() error {
	orderID := swap.order.PersonalOrderID()
	log.Println("Redeeming the swap for ", order.ID(orderID))

	swap.foreignAtom.Restore(swap.state)

	secret, err := swap.state.RedeemDetails(orderID)
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Redeem(secret); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusRedeemed); err != nil {
		return err
	}

	log.Println("Redeemed the swap for ", order.ID(orderID))
	return nil
}

func (swap *swap) refund() error {
	return nil
}

func (swap *swap) responderAudit() error {
	orderID := swap.order.PersonalOrderID()
	log.Println("Auditing the swap for ", order.ID(orderID))

	if err := swap.foreignAtom.Restore(swap.state); err != nil {
		return err
	}

	hashLock, to, value, expiry, err := swap.foreignAtom.Audit()
	newExpiry := expiry - 24*60*60

	personalAddr, err := swap.info.GetOwnerAddress(swap.order.PersonalOrderID())
	if err != nil {
		return err
	}

	if bytes.Compare(to, personalAddr) != 0 {
		return errors.New("Reciever Address Mismatch")
	}

	if value.Cmp(swap.order.ReceiveValue()) > 0 {
		return errors.New("Recieve value is less than expected")
	}

	if time.Now().Unix() > newExpiry {
		return errors.New("No time left to do the atomic swap")
	}

	if err := swap.state.PutInitiateDetails(orderID, newExpiry, hashLock); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusAudited); err != nil {
		return err
	}

	log.Println("Audit successful for ", order.ID(orderID))

	return nil
}

func (swap *swap) requestorAudit() error {
	orderID := swap.order.PersonalOrderID()
	log.Println("Auditing the swap for ", order.ID(orderID))

	if err := swap.foreignAtom.Restore(swap.state); err != nil {
		return err
	}

	hashLock, to, value, expiry, err := swap.foreignAtom.Audit()
	if err != nil {
		return err
	}
	_, selfHashLock, err := swap.state.InitiateDetails(orderID)
	if err != nil {
		return err
	}

	if hashLock != selfHashLock {
		return fmt.Errorf("Hashlock Mismatch %v %v", hashLock, selfHashLock)
	}

	personalAddr, err := swap.info.GetOwnerAddress(swap.order.PersonalOrderID())
	if err != nil {
		return err
	}

	if bytes.Compare(to, personalAddr) == 0 {
		return errors.New("Reciever Address Mismatch")
	}

	if value.Cmp(swap.order.ReceiveValue()) <= 0 {
		return errors.New("Recieve value is less than expected")
	}

	if time.Now().Unix() > expiry {
		return errors.New("No time left to do the atomic swap")
	}

	if err := swap.state.PutStatus(orderID, StatusAudited); err != nil {
		return err
	}

	log.Println("Audit successful for ", order.ID(orderID))
	return nil
}

func (swap *swap) getRedeemDetails() error {
	orderID := swap.order.PersonalOrderID()
	log.Println("Recieving the redeem details for ", order.ID(orderID))

	if err := swap.personalAtom.WaitForCounterRedemption(); err != nil {
		return err
	}

	log.Println("Counter Party redeemed", order.ID(orderID))

	secret, err := swap.personalAtom.AuditSecret()
	if err != nil {
		return err
	}

	if err := swap.state.PutRedeemDetails(orderID, secret); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusRedeemDetailsAcquired); err != nil {
		return err
	}

	log.Println("Recieved the redeem details for ", order.ID(orderID))
	return nil
}
