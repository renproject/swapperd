package order

import (
	"math/big"
)

type Order interface {
	Price() *big.Int
	PriceC() *big.Int
	PriceQ() *big.Int
	Volume() *big.Int
	VolumeC() *big.Int
	VolumeQ() *big.Int
	Type() uint8
	Parity() uint8
	MinimumVolume() *big.Int
	MinimumVolumeC() *big.Int
	MinimumVolumeQ() *big.Int
	Expiry() uint64
	Tokens() uint64
	NonceHash() *big.Int
	SendCurrency() uint32
	RecieveCurrency() uint32
	Trader() []byte
}

type order struct {
	parity         uint8
	orderType      uint8
	expiry         uint64
	tokens         uint64
	priceC         *big.Int
	priceQ         *big.Int
	volumeC        *big.Int
	volumeQ        *big.Int
	minimumVolumeC *big.Int
	minimumVolumeQ *big.Int
	nonceHash      *big.Int
	trader         []byte
}

func NewOrder(parity, orderType uint8, expiry, tokens uint64, priceC, priceQ, volumC, volumeQ, minimumVolumeC, minimumVolumeQ, nonceHash *big.Int, trader []byte) Order {
	return &order{
		parity:         parity,
		orderType:      orderType,
		expiry:         expiry,
		tokens:         tokens,
		priceC:         priceC,
		priceQ:         priceQ,
		volumeC:        volumC,
		volumeQ:        volumeQ,
		minimumVolumeC: minimumVolumeC,
		minimumVolumeQ: minimumVolumeQ,
		nonceHash:      nonceHash,
		trader:         trader,
	}
}

func (order *order) Price() *big.Int {
	return decodePrice(order.priceC, order.priceQ)
}
func (order *order) PriceC() *big.Int {
	return order.priceC
}
func (order *order) PriceQ() *big.Int {
	return order.priceQ
}

func (order *order) Volume() *big.Int {
	return decodeVolume(order.volumeC, order.volumeQ)
}
func (order *order) VolumeC() *big.Int {
	return order.volumeC
}
func (order *order) VolumeQ() *big.Int {
	return order.volumeQ
}

func (order *order) MinimumVolume() *big.Int {
	return decodeVolume(order.minimumVolumeC, order.minimumVolumeQ)
}

func (order *order) MinimumVolumeC() *big.Int {
	return order.minimumVolumeC
}
func (order *order) MinimumVolumeQ() *big.Int {
	return order.minimumVolumeQ
}

func (order *order) NonceHash() *big.Int {
	return order.nonceHash
}

func (order *order) Type() uint8 {
	return order.orderType
}

func (order *order) Parity() uint8 {
	return order.parity
}

func (order *order) Expiry() uint64 {
	return order.expiry
}

func (order *order) Tokens() uint64 {
	return order.tokens
}

func (order *order) SendCurrency() uint32 {
	if order.parity == uint8(0) {
		return priorityToken(order.tokens)
	}
	return nonPriorityToken(order.tokens)
}

func (order *order) RecieveCurrency() uint32 {
	if order.parity == uint8(1) {
		return priorityToken(order.tokens)
	}
	return nonPriorityToken(order.tokens)
}

func (order *order) Trader() []byte {
	return order.trader
}

func decodePrice(priceC, priceQ *big.Int) *big.Int {
	var price *big.Int

	priceQ.Sub(priceQ, big.NewInt(29))
	priceC.Mul(big.NewInt(5), priceC)
	price.Exp(big.NewInt(10), priceQ, nil)
	price.Mul(price, priceC)

	return price
}

func decodeVolume(volumeC, volumeQ *big.Int) *big.Int {
	var volume *big.Int

	volumeQ.Sub(volumeQ, big.NewInt(1))
	volumeC.Mul(big.NewInt(2), volumeC)
	volume.Exp(big.NewInt(10), volumeQ, nil)
	volume.Mul(volume, volumeC)

	return volume
}

func priorityToken(token uint64) uint32 {
	return uint32(token & 0x00000000FFFFFFFF)
}

func nonPriorityToken(token uint64) uint32 {
	return uint32(token >> 32)
}
