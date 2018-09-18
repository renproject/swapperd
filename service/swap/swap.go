package swap

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/match"

	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"

	"github.com/republicprotocol/renex-swapper-go/utils"
)

type Swap interface {
	Execute() error
}

type swap struct {
	personalAtom Atom
	foreignAtom  Atom
	order        match.Match
	Adapter
}

func (swap *swap) Execute() error {
	if swap.personalAtom.PriorityCode() == swap.foreignAtom.PriorityCode() {
		swap.LogError(swap.order.PersonalOrderID(), fmt.Sprintf("Trying to swap between atoms with the same priority code %d and %d", swap.personalAtom.PriorityCode(), swap.foreignAtom.PriorityCode()))
		return fmt.Errorf("Trying to swap between atoms with the same priority code %d and %d", swap.personalAtom.PriorityCode(), swap.foreignAtom.PriorityCode())
	}
	if swap.personalAtom.PriorityCode() < swap.foreignAtom.PriorityCode() {
		return swap.request()
	}
	return swap.respond()
}

func (swap *swap) request() error {
	personalOrderID := swap.order.PersonalOrderID()
	swap.LogInfo(personalOrderID, "is the requestor")
	if swap.Status(personalOrderID) == swapDomain.StatusInfoSubmitted {
		if err := swap.generateDetails(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to generate details: %v", err))
			return fmt.Errorf("failed to generate details: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping generate details")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusInitiateDetailsAcquired {
		if err := swap.initiate(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to initiate details: %v", err))
			return fmt.Errorf("failed to initiate details: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping initiate")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusInitiated {
		if err := swap.sendDetails(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to send details: %v", err))
			return fmt.Errorf("failed to send details: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping send details")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusSentSwapDetails {
		if err := swap.receiveDetails(); err != nil {
			if err := swap.ComplainDelayedResponderInitiation(personalOrderID); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to complain to the watchdog: %v", err))
				return fmt.Errorf("failed to complain to the watchdog: %v", err)
			}
			swap.LogError(personalOrderID, fmt.Sprintf("failed to receive details: %v", err))
			if err := swap.PutStatus(personalOrderID, swapDomain.StatusComplained); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to change the status: %v", err))
				return fmt.Errorf("failed to change the status: %v", err)
			}
			return fmt.Errorf("failed to receive details: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping receive details")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusReceivedSwapDetails {
		if err := swap.requestorAudit(); err != nil {
			if err := swap.ComplainWrongResponderInitiation(personalOrderID); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to complain to the watch dog: %v", err))
				return fmt.Errorf("failed to complain to the watch dog: %v", err)
			}
			swap.LogError(personalOrderID, fmt.Sprintf("failed to receive swap details: %v", err))
			if err := swap.PutStatus(personalOrderID, swapDomain.StatusComplained); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to update the status: %v", err))
				return fmt.Errorf("failed to update the status: %v", err)
			}
			return fmt.Errorf("failed to receive swap details: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping requestor audit")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusAudited {
		if err := swap.redeem(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to redeem: %v", err))
			return fmt.Errorf("failed to redeem: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping redeem")
	}

	return nil
}

func (swap *swap) respond() error {
	personalOrderID := swap.order.PersonalOrderID()

	swap.LogInfo(personalOrderID, "is the responder")

	if swap.Status(personalOrderID) == swapDomain.StatusInfoSubmitted {
		if err := swap.receiveDetails(); err != nil {
			if err := swap.ComplainDelayedRequestorInitiation(personalOrderID); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to complain to the watch dog: %v", err))
				return fmt.Errorf("failed to complain to the watch dog: %v", err)
			}
			swap.LogError(personalOrderID, fmt.Sprintf("failed to receive details: %v", err))
			if err := swap.PutStatus(personalOrderID, swapDomain.StatusComplained); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to change status: %v", err))
				return fmt.Errorf("failed to change status: %v", err)
			}
			return fmt.Errorf("failed to receive details: %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping generate details")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusReceivedSwapDetails {
		if err := swap.responderAudit(); err != nil {
			if err := swap.ComplainWrongRequestorInitiation(personalOrderID); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to complain to the watch dog %v", err))
				return fmt.Errorf("failed to complain to the watch dog %v", err)
			}
			swap.LogError(personalOrderID, fmt.Sprintf("audit failed %v", err))
			if err := swap.PutStatus(personalOrderID, swapDomain.StatusComplained); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to update status %v", err))
				return fmt.Errorf("failed to update status %v", err)
			}
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping audit")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusAudited {
		if err := swap.initiate(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to initiate %v", personalOrderID))
			return fmt.Errorf("failed to initiate %v", personalOrderID)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping initiate")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusInitiated {
		if err := swap.sendDetails(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to send details %v", personalOrderID))
			return fmt.Errorf("failed to send details %v", personalOrderID)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping send details")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusSentSwapDetails {
		if err := swap.getRedeemDetails(); err != nil {
			if err := swap.ComplainDelayedRequestorRedemption(personalOrderID); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to complain to the watch dog %v", err))
				return fmt.Errorf("failed to complain to the watch dog %v", err)
			}
			swap.LogError(personalOrderID, fmt.Sprintf("failed to get redeem details %v", err))
			if err := swap.PutStatus(personalOrderID, swapDomain.StatusComplained); err != nil {
				swap.LogError(personalOrderID, fmt.Sprintf("failed to update status %v", err))
				return fmt.Errorf("failed to update status %v", err)
			}
			return fmt.Errorf("failed to get redeem details %v", err)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping get redeem details audit")
	}

	if swap.Status(personalOrderID) == swapDomain.StatusRedeemDetailsAcquired {
		if err := swap.redeem(); err != nil {
			swap.LogError(personalOrderID, fmt.Sprintf("failed to redeem %v", personalOrderID))
			return fmt.Errorf("failed to redeem %v", personalOrderID)
		}
	} else {
		swap.LogInfo(personalOrderID, "skipping redeem")
	}
	return nil
}

func (swap *swap) generateDetails() error {
	orderID := swap.order.PersonalOrderID()
	swap.LogInfo(orderID, "generating swap details")
	expiry := time.Now().Add(48 * time.Hour).Unix()
	secret := make([]byte, 32)
	rand.Read(secret)
	secret32, err := utils.ToBytes32(secret)
	if err != nil {
		return err
	}
	secretHash := sha256.Sum256(secret)

	if err := swap.PutInitiateDetails(orderID, secretHash, expiry); err != nil {
		return err
	}

	if err := swap.PutRedeemDetails(orderID, secret32); err != nil {
		return err
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusInitiateDetailsAcquired); err != nil {
		return err
	}
	swap.LogInfo(orderID, "generated the swap details")
	return nil
}

func (swap *swap) initiate() error {
	orderID := swap.order.PersonalOrderID()
	secretHash, expiry, err := swap.InitiateDetails(orderID)
	if err != nil {
		return err
	}
	swap.LogInfo(orderID, "initiating the swap")

	foreignAddr, err := swap.ReceiveOwnerAddress(swap.order.ForeignOrderID(), time.Now().Add(24*time.Hour).Unix())
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

	if err := swap.PutAtomDetails(swap.order.PersonalOrderID(), details); err != nil {
		return err
	}

	if err := swap.PutRedeemable(orderID); err != nil {
		return err
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusInitiated); err != nil {
		return err
	}

	swap.LogInfo(orderID, "initiated the swap")
	return nil
}

func (swap *swap) sendDetails() error {
	orderID := swap.order.PersonalOrderID()
	swap.LogInfo(orderID, "sending the swap details")
	personalAtomBytes, err := swap.AtomDetails(orderID)
	if err != nil {
		return err
	}
	if err := swap.SendSwapDetails(orderID, personalAtomBytes); err != nil {
		log.Println("Error Here", err)
		return err
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusSentSwapDetails); err != nil {
		return err
	}
	swap.LogInfo(orderID, "sent the swap details for")
	return nil
}

func (swap *swap) receiveDetails() error {
	personalOrderID := swap.order.PersonalOrderID()
	foreignOrderID := swap.order.ForeignOrderID()
	swap.LogInfo(personalOrderID, "receiving the swap details")
	foreignAtomBytes, err := swap.ReceiveSwapDetails(foreignOrderID, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		return err
	}

	if err := swap.PutAtomDetails(foreignOrderID, foreignAtomBytes); err != nil {
		return err
	}

	if err := swap.PutStatus(personalOrderID, swapDomain.StatusReceivedSwapDetails); err != nil {
		return err
	}
	swap.LogInfo(personalOrderID, "received the swap details")
	return nil
}

func (swap *swap) redeem() error {
	orderID := swap.order.PersonalOrderID()
	swap.LogInfo(orderID, "redeeming the swap details")

	details, err := swap.AtomDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Deserialize(details); err != nil {
		return err
	}

	secret, err := swap.RedeemDetails(orderID)
	if err != nil {
		return err
	}

	if err := swap.foreignAtom.Redeem(secret); err != nil {
		return err
	}

	if err := swap.Redeemed(orderID); err != nil {
		return err
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusRedeemed); err != nil {
		return err
	}

	swap.LogInfo(orderID, "redeemed the swap details")
	return nil
}

func (swap *swap) responderAudit() error {
	orderID := swap.order.PersonalOrderID()
	swap.LogInfo(orderID, "auditing the swap")

	details, err := swap.AtomDetails(swap.order.ForeignOrderID())
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

	personalAddr, err := swap.ReceiveOwnerAddress(swap.order.PersonalOrderID(), 0)
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

	if err := swap.PutInitiateDetails(orderID, hashLock, expiry); err != nil {
		return err
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusAudited); err != nil {
		return err
	}

	swap.LogInfo(orderID, "auditing successful")
	return nil
}

func (swap *swap) requestorAudit() error {
	orderID := swap.order.PersonalOrderID()
	swap.LogInfo(orderID, "auditing the swap")

	details, err := swap.AtomDetails(swap.order.ForeignOrderID())
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

	selfHashLock, _, err := swap.InitiateDetails(orderID)
	if err != nil {
		return err
	}

	if hashLock != selfHashLock {
		return fmt.Errorf("Hashlock Mismatch %v %v", hashLock, selfHashLock)
	}

	personalAddr, err := swap.ReceiveOwnerAddress(swap.order.PersonalOrderID(), 0)
	if err != nil {
		return err
	}

	if bytes.Compare(to, personalAddr) != 0 {
		swap.LogError(orderID, fmt.Sprintf("Receiver Address Mismatch %s != %s", string(to), string(personalAddr)))
	}

	if value.Cmp(swap.order.ReceiveValue()) < 0 {
		return fmt.Errorf("Receive value: %v is less than expected: %v", value, swap.order.ReceiveValue())
	}

	if time.Now().Unix() > expiry {
		return errors.New("No time left to do the atomic swap")
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusAudited); err != nil {
		return err
	}

	swap.LogInfo(orderID, "auditing successful")
	return nil
}

func (swap *swap) getRedeemDetails() error {
	orderID := swap.order.PersonalOrderID()
	swap.LogInfo(orderID, "receiving the redeem details")

	if err := swap.personalAtom.WaitForCounterRedemption(); err != nil {
		return err
	}

	secret, err := swap.personalAtom.AuditSecret()
	if err != nil {
		return err
	}

	if err := swap.PutRedeemDetails(orderID, secret); err != nil {
		return err
	}

	if err := swap.PutStatus(orderID, swapDomain.StatusRedeemDetailsAcquired); err != nil {
		return err
	}

	swap.LogInfo(orderID, "received the redeem details")
	return nil
}
