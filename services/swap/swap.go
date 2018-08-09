package swap

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domains/match"
	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

// Swap is the interface for an atomic swap object
type Swap interface {
	Execute() error
}

type swap struct {
	personalAtom Atom
	foreignAtom  Atom
	order        match.Match
	swapAdapter  SwapAdapter
	state        store.State
}

// NewSwap returns a new Swap instance
func NewSwap(personalAtom Atom, foreignAtom Atom, order match.Match, swapAdapter SwapAdapter, state store.State) Swap {
	return &swap{
		personalAtom: personalAtom,
		foreignAtom:  foreignAtom,
		order:        order,
		swapAdapter:  swapAdapter,
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
	personalOrderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(personalOrderID, "is the requestor")
	if swap.state.Status(personalOrderID) == StatusInfoSubmitted {
		if err := swap.generateDetails(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping generate details")
	}

	if swap.state.Status(personalOrderID) == StatusInitateDetailsAcquired {
		if err := swap.initiate(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping initiate")
	}

	if swap.state.Status(personalOrderID) == StatusInitiated {
		if err := swap.sendDetails(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping send details")
	}

	if swap.state.Status(personalOrderID) == StatusSentSwapDetails {
		if err := swap.receiveDetails(); err != nil {
			swap.swapAdapter.ComplainDelayedResponderInitiation(personalOrderID)
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping receive details")
	}

	if swap.state.Status(personalOrderID) == StatusReceivedSwapDetails {
		if err := swap.requestorAudit(); err != nil {
			swap.swapAdapter.ComplainWrongResponderInitiation(personalOrderID)
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping requestor audit")
	}

	if swap.state.Status(personalOrderID) == StatusAudited {
		if err := swap.redeem(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping redeem")
	}
	return nil
}

func (swap *swap) respond() error {
	personalOrderID := swap.order.PersonalOrderID()

	swap.swapAdapter.LogInfo(personalOrderID, "is the requestor")

	if swap.state.Status(personalOrderID) == StatusInfoSubmitted {
		if err := swap.receiveDetails(); err != nil {
			swap.swapAdapter.ComplainDelayedRequestorInitiation(personalOrderID)
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping generate details")
	}

	if swap.state.Status(personalOrderID) == StatusReceivedSwapDetails {
		if err := swap.responderAudit(); err != nil {
			swap.swapAdapter.ComplainWrongRequestorInitiation(personalOrderID)
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping audit")
	}

	if swap.state.Status(personalOrderID) == StatusAudited {
		if err := swap.initiate(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping initiate")
	}

	if swap.state.Status(personalOrderID) == StatusInitiated {
		if err := swap.sendDetails(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping send details")
	}

	if swap.state.Status(personalOrderID) == StatusSentSwapDetails {
		if err := swap.getRedeemDetails(); err != nil {
			swap.swapAdapter.ComplainDelayedRequestorRedemption(personalOrderID)
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping get redeem details audit")
	}

	if swap.state.Status(personalOrderID) == StatusRedeemDetailsAcquired {
		if err := swap.redeem(); err != nil {
			return err
		}
	} else {
		swap.swapAdapter.LogInfo(personalOrderID, "skipping redeem")
	}
	return nil
}

func (swap *swap) generateDetails() error {
	orderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(orderID, "generating swap details")
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
	swap.swapAdapter.LogInfo(orderID, "generated the swap details")
	return nil
}

func (swap *swap) initiate() error {
	orderID := swap.order.PersonalOrderID()
	expiry, secretHash, err := swap.state.InitiateDetails(orderID)
	if err != nil {
		return err
	}
	swap.swapAdapter.LogInfo(orderID, "initiating the swap")

	foreignAddr, err := swap.swapAdapter.ReceiveOwnerAddress(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	if err = swap.personalAtom.Initiate(foreignAddr, secretHash, swap.order.SendValue(), expiry); err != nil {
		return err
	}

	details, err := swap.personalAtom.Serialize()
	if err != nil {
		return err
	}

	if err := swap.state.PutAtomDetails(swap.order.PersonalOrderID(), details); err != nil {
		return err
	}

	if err := swap.state.PutRedeemable(orderID); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusInitiated); err != nil {
		return err
	}

	swap.swapAdapter.LogInfo(orderID, "initiated the swap")
	return nil
}

func (swap *swap) sendDetails() error {
	orderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(orderID, "sending the swap details")
	personalAtomBytes, err := swap.state.AtomDetails(orderID)
	if err != nil {
		return err
	}
	if err := swap.swapAdapter.SendSwapDetails(orderID, personalAtomBytes); err != nil {
		log.Println("Error Here", err.Error())
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusSentSwapDetails); err != nil {
		return err
	}
	swap.swapAdapter.LogInfo(orderID, "sent the swap details for")
	return nil
}

func (swap *swap) receiveDetails() error {
	personalOrderID := swap.order.PersonalOrderID()
	foreignOrderID := swap.order.ForeignOrderID()
	swap.swapAdapter.LogInfo(personalOrderID, "recieving the swap details")
	foreignAtomBytes, err := swap.swapAdapter.ReceiveSwapDetails(foreignOrderID, true)
	if err != nil {
		return err
	}

	if err := swap.state.PutAtomDetails(foreignOrderID, foreignAtomBytes); err != nil {
		return err
	}

	if err := swap.state.PutStatus(personalOrderID, StatusReceivedSwapDetails); err != nil {
		return err
	}
	swap.swapAdapter.LogInfo(personalOrderID, "recieved the swap details")
	return nil
}

func (swap *swap) redeem() error {
	orderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(orderID, "redeeming the swap details")

	details, err := swap.state.AtomDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Deserialize(details); err != nil {
		return err
	}

	secret, err := swap.state.RedeemDetails(orderID)
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Redeem(secret); err != nil {
		return err
	}

	if err := swap.state.Redeemed(orderID); err != nil {
		return err
	}

	if err := swap.state.PutStatus(orderID, StatusRedeemed); err != nil {
		return err
	}

	swap.swapAdapter.LogInfo(orderID, "redeemed the swap details")
	return nil
}

func (swap *swap) responderAudit() error {
	orderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(orderID, "auditing the swap")

	details, err := swap.state.AtomDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Deserialize(details); err != nil {
		return err
	}
	hashLock, to, value, expiry, err := swap.foreignAtom.Audit()
	if err != nil {
		return err
	}
	newExpiry := expiry - 24*60*60

	personalAddr, err := swap.swapAdapter.ReceiveOwnerAddress(swap.order.PersonalOrderID())
	if err != nil {
		return err
	}

	if bytes.Compare(to, personalAddr) != 0 {
		return errors.New("Receiver Address Mismatch")
	}

	if value.Cmp(swap.order.ReceiveValue()) > 0 {
		return errors.New("Receive value is less than expected")
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

	swap.swapAdapter.LogInfo(orderID, "auditing successful")
	return nil
}

func (swap *swap) requestorAudit() error {
	orderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(orderID, "auditing the swap")

	details, err := swap.state.AtomDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Deserialize(details); err != nil {
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

	personalAddr, err := swap.swapAdapter.ReceiveOwnerAddress(swap.order.PersonalOrderID())
	if err != nil {
		return err
	}

	if bytes.Compare(to, personalAddr) != 0 {
		return errors.New("Receiver Address Mismatch")
	}

	if value.Cmp(swap.order.ReceiveValue()) < 0 {
		return errors.New("Receive value is less than expected")
	}

	if time.Now().Unix() > expiry {
		return errors.New("No time left to do the atomic swap")
	}

	if err := swap.state.PutStatus(orderID, StatusAudited); err != nil {
		return err
	}

	swap.swapAdapter.LogInfo(orderID, "auditing successful")
	return nil
}

func (swap *swap) getRedeemDetails() error {
	orderID := swap.order.PersonalOrderID()
	swap.swapAdapter.LogInfo(orderID, "recieving the redeem details")

	if err := swap.personalAtom.WaitForCounterRedemption(); err != nil {
		return err
	}

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

	swap.swapAdapter.LogInfo(orderID, "recieved the redeem details")
	return nil
}
