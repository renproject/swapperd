package swap

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/utils"
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
	swapStr      SwapStore
}

// NewSwap returns a new Swap instance
func NewSwap(atom1 Atom, atom2 Atom, info Info, order match.Match, network Network, swapStr SwapStore) Swap {
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
		swapStr:      swapStr,
	}
}

func (swap *swap) Execute() error {
	err := swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "MATCHED")
	if err != nil {
		return err
	}

	if swap.personalAtom.PriorityCode() < swap.foreignAtom.PriorityCode() {
		return swap.initiate()
	}
	return swap.respond()
}

func (swap *swap) initiate() error {
	personalAddr, err := swap.info.GetOwnerAddress(swap.order.PersonalOrderID())
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

	secret32, err := utils.ToBytes32(secret)
	if err != nil {
		return err
	}
	secretHash := sha256.Sum256(secret)

	err = swap.personalAtom.Initiate(foreignAddr, secretHash, swap.order.SendValue(), expiry)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "INITIATED")
	if err != nil {
		return err
	}

	personalSwapDetails, err := swap.personalAtom.Serialize()
	if err != nil {
		return err
	}
	err = swap.network.SendSwapDetails(swap.order.PersonalOrderID(), personalSwapDetails)
	if err != nil {
		return err
	}

	foreignSwapDetails, err := swap.network.ReceiveSwapDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}
	err = swap.foreignAtom.Deserialize(foreignSwapDetails)
	if err != nil {
		return err
	}

	err = swap.foreignAtom.Audit(secretHash, personalAddr, swap.order.ReceiveValue(), 60*60)
	if err != nil {
		err2 := swap.personalAtom.Refund()
		if err2 != nil {
			// Should never happen
			return err2
		}
		err2 = swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "REFUNDED")
		if err2 != nil {
			return err2
		}
		return err
	}

	err = swap.foreignAtom.Redeem(secret32)
	if err != nil {
		return err
	}

	err = swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "REDEEMED")
	if err != nil {
		return err
	}

	return nil
}

func (swap *swap) respond() error {
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

	err = swap.foreignAtom.Audit(hash, personalAddr, swap.order.ReceiveValue(), expiry)
	if err != nil {
		return err
	}

	err = swap.personalAtom.Initiate(foreignAddr, hash, swap.order.SendValue(), expiry)
	if err != nil {
		return err
	}

	err = swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "INITIATED")
	if err != nil {
		return err
	}

	personalSwapDetails, err := swap.personalAtom.Serialize()
	if err != nil {
		return err
	}

	fmt.Println("Sending swap details")
	err = swap.network.SendSwapDetails(swap.order.PersonalOrderID(), personalSwapDetails)
	if err != nil {
		return err
	}

	secret, err := swap.personalAtom.AuditSecret()
	if err != nil {
		err1 := swap.personalAtom.Refund()
		if err1 != nil {
			// SHOULD NEVER HAPPEN
			return err
		}
		err1 = swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "REFUNDED")
		if err1 != nil {
			return err1
		}
		return err
	}

	err = swap.foreignAtom.Redeem(secret)
	if err != nil {
		return err
	}

	err = swap.swapStr.UpdateStatus(swap.order.PersonalOrderID(), "REDEEMED")
	if err != nil {
		return err
	}

	return nil
}
