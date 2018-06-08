package swap

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/republicprotocol/atom-go/bytesutils"
	"github.com/republicprotocol/atom-go/services/atom"
	"github.com/republicprotocol/atom-go/services/network"
	"github.com/republicprotocol/atom-go/services/order"
)

// Swap is the interface for an atomic swap object
type Swap interface {
	Execute() error
	initiate() error
	respond() error
	store() error
	retrieve() error
}

type swap struct {
	myAtom      atom.Atom
	tradingAtom atom.Atom
	order       order.Order
	network     network.Network
	expiry      int64
}

// NewSwap returns a new Swap instance
func NewSwap(myAtom atom.Atom, tradingAtom atom.Atom, order order.Order, network network.Network) (Swap, error) {
	return &swap{
		myAtom:      myAtom,
		tradingAtom: tradingAtom,
		order:       order,
	}, nil
}

func (swap *swap) Execute() error {
	if swap.myAtom.PriorityCode() > swap.tradingAtom.PriorityCode() {
		return swap.initiate()
	}
	return swap.respond()
}

func (swap *swap) initiate() error {
	swap.store()
	swap.retrieve()

	expiry := time.Now().Add(48 * time.Hour).Unix()

	secret := make([]byte, 32)
	rand.Read(secret)

	secret32, err := bytesutils.ToBytes32(secret)
	if err != nil {
		return err
	}

	secretHash := sha256.Sum256(secret)

	err = swap.myAtom.Initiate(secretHash, swap.order.SendValue(), expiry)
	if err != nil {
		return err
	}
	err = swap.wait1()

	if err != nil {
		err2 := swap.myAtom.Refund()
		if err2 != nil {
			// Should never happen
			return err2
		}
		return err
	}

	return swap.tradingAtom.Redeem(secret32)
}

func (swap *swap) respond() error {
	swap.store()
	swap.retrieve()
	expiry := time.Now().Add(24 * time.Hour).Unix()
	err := swap.wait2(time.Now().Unix() + 2*60*60)
	if err != nil {
		return err
	}

	hash, _, _, _, _, err := swap.tradingAtom.Audit()
	if err != nil {
		return err
	}

	err = swap.myAtom.Initiate(hash, swap.order.SendValue(), expiry)
	if err != nil {
		return err
	}

	err = swap.wait3()
	if err != nil {
		return swap.myAtom.Refund()
	}

	secret, err := swap.myAtom.AuditSecret()
	if err != nil {
		// Should never happen.
		return err
	}

	return swap.tradingAtom.Redeem(secret)
}

func (swap *swap) store() error {
	myDetails, err := swap.myAtom.Serialize()
	if err != nil {
		return err
	}
	tradingDetails, err := swap.tradingAtom.Serialize()
	if err != nil {
		return err
	}
	mySwapDetails := append(myDetails, tradingDetails...)
	return swap.network.Send(swap.order.MyOrderID(), mySwapDetails)
}

func (swap *swap) retrieve() error {
	tradingSwapDetails, err := swap.network.Recieve(swap.order.TradingOrderID())
	if err != nil {
		return err
	}
	err = swap.tradingAtom.Deserialize(tradingSwapDetails[:52])
	if err != nil {
		return err
	}
	return swap.myAtom.Deserialize(tradingSwapDetails[52:])
}

func (swap *swap) wait1() error {
	for time.Unix(swap.expiry, 0).Sub(time.Now()).Seconds() > 0 {
		time.Sleep(10 * time.Second)
		_, _to, _, _value, _expiry, err := swap.tradingAtom.Audit()
		if err != nil {
			continue
		}
		if bytes.Compare(swap.myAtom.From(), _to) != 0 || swap.order.RecieveValue().Cmp(_value) != 0 || time.Unix(_expiry, 0).Sub(time.Now()).Seconds() < 60*60 {
			continue
		}
		return nil
	}
	return errors.New("Timeout")
}

func (swap *swap) wait2(waitTill int64) error {
	for time.Unix(waitTill, 0).Sub(time.Now()).Seconds() > 0 {
		time.Sleep(10 * time.Second)
		_, _to, _, _value, _expiry, err := swap.tradingAtom.Audit()
		if err != nil {
			continue
		}
		if bytes.Compare(swap.myAtom.From(), _to) != 0 || swap.order.RecieveValue().Cmp(_value) != 1 || _expiry-time.Now().Unix() < 24*60*60 {
			continue
		}
		return nil
	}
	return errors.New("Timeout")
}

func (swap *swap) wait3() error {
	for time.Unix(swap.expiry, 0).Sub(time.Now()).Seconds() > 0 {
		time.Sleep(10 * time.Second)
		_secret, err := swap.myAtom.AuditSecret()
		if err != nil {
			continue
		}
		if _secret == [32]byte{} {
			continue
		}
		return nil
	}
	return errors.New("Timeout")
}
