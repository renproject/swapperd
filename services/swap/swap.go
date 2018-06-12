package swap

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/republicprotocol/atom-go/bytesutils"
	"github.com/republicprotocol/atom-go/services/atom"
	"github.com/republicprotocol/atom-go/services/axc"
	"github.com/republicprotocol/atom-go/services/network"
	"github.com/republicprotocol/atom-go/services/order"
)

// Swap is the interface for an atomic swap object
type Swap interface {
	Execute() error
}

type swap struct {
	personalAtom atom.RequestAtom
	foreignAtom  atom.ResponseAtom
	order        order.Order
	network      network.Network
	axc          axc.AXC
}

// NewSwap returns a new Swap instance
func NewSwap(personalAtom atom.RequestAtom, foreignAtom atom.ResponseAtom, axc axc.AXC, order order.Order, network network.Network) Swap {
	return &swap{
		personalAtom: personalAtom,
		foreignAtom:  foreignAtom,
		order:        order,
		axc:          axc,
		network:      network,
	}
}

func (swap *swap) Execute() error {
	if swap.personalAtom.PriorityCode() < swap.foreignAtom.PriorityCode() {
		return swap.initiate()
	}
	return swap.respond()
}

func (swap *swap) initiate() error {
	personalAddr, err := swap.axc.GetOwnerAddress(swap.order.PersonalOrderID())

	expiry := time.Now().Add(48 * time.Hour).Unix()

	secret := make([]byte, 32)
	rand.Read(secret)

	secret32, err := bytesutils.ToBytes32(secret)
	if err != nil {
		return err
	}

	secretHash := sha256.Sum256(secret)

	fmt.Println("Initiating the atomic Swap")
	err = swap.personalAtom.Initiate(secretHash, swap.order.SendValue(), expiry)
	if err != nil {
		return err
	}
	fmt.Println("Initiated the atomic Swap")

	personalSwapDetails, err := swap.personalAtom.Serialize()
	if err != nil {
		return err
	}

	fmt.Println("Sending swap details")
	err = swap.network.SendSwapDetails(swap.order.PersonalOrderID(), personalSwapDetails)
	if err != nil {
		return err
	}
	fmt.Println("Sent swap details")

	foreignSwapDetails, err := swap.network.RecieveSwapDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	fmt.Println("deserializing swap details")

	err = swap.foreignAtom.Deserialize(foreignSwapDetails)
	if err != nil {
		return err
	}

	fmt.Println("auditing swap details")

	err = swap.foreignAtom.Audit(secretHash, personalAddr, swap.order.RecieveValue(), 60*60)

	if err != nil {
		fmt.Println("Initiating a refund")
		err2 := swap.personalAtom.Refund()
		if err2 != nil {
			// Should never happen
			return err2
		}
		return err
	}

	fmt.Println("redeeming swap details")

	return swap.foreignAtom.Redeem(secret32)
}

func (swap *swap) respond() error {
	personalAddr, err := swap.axc.GetOwnerAddress(swap.order.PersonalOrderID())

	fmt.Println("Trying to retrieve swap details")

	foreignSwapDetails, err := swap.network.RecieveSwapDetails(swap.order.ForeignOrderID())
	if err != nil {
		return err
	}

	fmt.Println("Retrieved swap details")

	err = swap.foreignAtom.Deserialize(foreignSwapDetails)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(24 * time.Hour).Unix()

	hash := swap.foreignAtom.GetSecretHash()

	err = swap.foreignAtom.Audit(hash, personalAddr, swap.order.RecieveValue(), expiry)
	if err != nil {
		return err
	}

	err = swap.personalAtom.Initiate(hash, swap.order.SendValue(), expiry)
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

	time.Sleep(5 * time.Second)
	secret, err := swap.personalAtom.AuditSecret()
	if err != nil {
		// Should never happen.
		fmt.Println("Audit secret failed trying to refund:", err)
		return swap.personalAtom.Refund()
	}

	return swap.foreignAtom.Redeem(secret)
}
