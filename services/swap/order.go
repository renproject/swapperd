package swap

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type OrderMatch interface {
	PersonalOrderID() [32]byte
	ForeignOrderID() [32]byte
	SendValue() *big.Int
	RecieveValue() *big.Int
	SendCurrency() uint32
	RecieveCurrency() uint32
}

type orderMatch struct {
	personalOrderID [32]byte
	foreignOrderID  [32]byte
	sendValue       *big.Int
	recieveValue    *big.Int
	sendCurrency    uint32
	recieveCurrency uint32
}

type order struct {
	Parity         uint8
	OrderType      uint8
	Expiry         uint64
	Tokens         uint64
	PriceC         *big.Int
	PriceQ         *big.Int
	VolumeC        *big.Int
	VolumeQ        *big.Int
	MinimumVolumeC *big.Int
	MinimumVolumeQ *big.Int
	NonceHash      *big.Int
	Trader         common.Address
}

func NewOrderMatch(personalOrder, foreignOrder order, personalOrderID, foreignOrderID [32]byte) OrderMatch {
	var price, volume *big.Int
	var sendCurrency, recieveCurrency uint32
	var sendValue, recieveValue *big.Int

	personalOrderPrice, personalOrderVolume := decode(personalOrder.PriceC, personalOrder.VolumeC, personalOrder.PriceQ, personalOrder.VolumeQ)
	foreignOrderPrice, foreignOrderVolume := decode(foreignOrder.PriceC, foreignOrder.VolumeC, foreignOrder.PriceQ, foreignOrder.VolumeQ)

	price.Add(personalOrderPrice, personalOrderPrice)
	price.Div(price, big.NewInt(2))

	if foreignOrderVolume.Cmp(personalOrderVolume) > 0 {
		volume = personalOrderVolume
	} else {
		volume = foreignOrderVolume
	}

	priorityCode := personalOrder.Tokens & 0x00000000FFFFFFFF
	nonPriorityCode := foreignOrder.Tokens >> 32

	if personalOrder.Parity == 0 {
		sendCurrency = nonPriorityCode
		recieveCurrency = priorityCode
		sendValue, recieveValue = calculateValues(volume, price)
	} else {
		sendCurrency = priorityCode
		recieveCurrency = nonPriorityCode
		recieveValue, sendValue = calculateValues(volume, price)
	}

	return &orderMatch{
		personalOrderID: personalOrderID,
		foreignOrderID:  foreignOrderID,
		sendValue:       sendValue,
		recieveValue:    recieveValue,
		sendCurrency:    sendCurrency,
		recieveCurrency: recieveCurrency,
	}
}

func (order *orderMatch) PersonalOrderID() [32]byte {
	return order.personalOrderID
}
func (order *orderMatch) ForeignOrderID() [32]byte {
	return order.foreignOrderID
}
func (order *orderMatch) SendValue() *big.Int {
	return order.sendValue
}
func (order *orderMatch) RecieveValue() *big.Int {
	return order.recieveValue
}
func (order *orderMatch) SendCurrency() uint32 {
	return order.sendCurrency
}
func (order *orderMatch) RecieveCurrency() uint32 {
	return order.recieveCurrency
}

func decode(priceC, volumeC, priceQ, volumeQ *big.Int) (*big.Int, *big.Int) {
	var volume, price *big.Int

	priceQ.Sub(priceQ, big.NewInt(29))
	priceC.Mul(big.NewInt(5), priceC)
	price.Exp(big.NewInt(10), priceQ, nil)
	price.Mul(price, priceC)

	volumeQ.Sub(volumeQ, big.NewInt(1))
	volumeC.Mul(big.NewInt(2), volumeC)
	volume.Exp(big.NewInt(10), volumeQ, nil)
	volume.Mul(volume, volumeC)

	return price, volume
}

func calculateValues(price, volume *big.Int) (*big.Int, *big.Int) {
	return price, volume
}

func priorityToken(token uint64) uint32 {
	return uint32(token & 0x00000000FFFFFFFF)
}

func nonPriorityToken(token uint64) uint32 {
	return uint32(token >> 32)
}
