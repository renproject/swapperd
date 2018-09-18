package network_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	"github.com/republicprotocol/renex-swapper-go/domains/order"
	. "github.com/republicprotocol/renex-swapper-go/drivers/network"
)

var _ = Describe("Ingress Network Driver", func() {
	buildIngress := func(network string) network.Network {
		return NewIngress(fmt.Sprintf("renex-ingress-%s.herokuapp.com", network))
	}

	randomDetails := func() (order.ID, []byte, error) {
		orderID := [32]byte{}
		address := make([]byte, 20)
		if _, err := rand.Read(orderID[:]); err != nil {
			return order.ID{}, nil, err
		}
		if _, err := rand.Read(address); err != nil {
			return order.ID{}, nil, err
		}
		return order.ID(orderID), address, nil
	}

	Context("when communicating with nightly ingress", func() {
		It("should be able to send and retrieve address information when done right", func() {
			ingress := buildIngress("nightly")
			id, sendAddr, err := randomDetails()
			Expect(err).Should(BeNil())
			err = ingress.SendOwnerAddress(id, sendAddr)
			Expect(err).Should(BeNil())
			recvAddr, err := ingress.ReceiveOwnerAddress(id, time.Now().Unix()+100)
			Expect(err).Should(BeNil())
			Expect(bytes.Compare(sendAddr, recvAddr)).Should(Equal(0))
		})

		It("should be able to send and retrieve swap details when done right", func() {
			ingress := buildIngress("nightly")
			id, sendAddr, err := randomDetails()
			Expect(err).Should(BeNil())
			err = ingress.SendSwapDetails(id, sendAddr)
			Expect(err).Should(BeNil())
			recvAddr, err := ingress.ReceiveSwapDetails(id, time.Now().Unix()+100)
			Expect(err).Should(BeNil())
			Expect(bytes.Compare(sendAddr, recvAddr)).Should(Equal(0))
		})
	})

	// Context("when communicating with falcon ingress", func() {
	// 	It("should be able to send and retrieve address information when done right", func() {
	// 		ingress := buildIngress("falcon")
	// 		id, sendAddr, err := randomDetails()
	// 		Expect(err).Should(BeNil())
	// 		err = ingress.SendOwnerAddress(id, sendAddr)
	// 		Expect(err).Should(BeNil())
	// 		recvAddr, err := ingress.ReceiveOwnerAddress(id, time.Now().Unix()+100)
	// 		Expect(err).Should(BeNil())
	// 		Expect(bytes.Compare(sendAddr, recvAddr)).Should(Equal(0))
	// 	})

	// 	It("should be able to send and retrieve swap details when done right", func() {
	// 		ingress := buildIngress("falcon")
	// 		id, sendAddr, err := randomDetails()
	// 		Expect(err).Should(BeNil())
	// 		err = ingress.SendSwapDetails(id, sendAddr)
	// 		Expect(err).Should(BeNil())
	// 		recvAddr, err := ingress.ReceiveSwapDetails(id, time.Now().Unix()+100)
	// 		Expect(err).Should(BeNil())
	// 		Expect(bytes.Compare(sendAddr, recvAddr)).Should(Equal(0))
	// 	})
	// })

	// Context("when communicating with testnet ingress", func() {
	// 	It("should be able to send and retrieve address information when done right", func() {
	// 		ingress := buildIngress("testnet")
	// 		id, sendAddr, err := randomDetails()
	// 		Expect(err).Should(BeNil())
	// 		err = ingress.SendOwnerAddress(id, sendAddr)
	// 		Expect(err).Should(BeNil())
	// 		recvAddr, err := ingress.ReceiveOwnerAddress(id, time.Now().Unix()+100)
	// 		Expect(err).Should(BeNil())
	// 		Expect(bytes.Compare(sendAddr, recvAddr)).Should(Equal(0))
	// 	})

	// 	It("should be able to send and retrieve swap details when done right", func() {
	// 		ingress := buildIngress("testnet")
	// 		id, sendAddr, err := randomDetails()
	// 		Expect(err).Should(BeNil())
	// 		err = ingress.SendSwapDetails(id, sendAddr)
	// 		Expect(err).Should(BeNil())
	// 		recvAddr, err := ingress.ReceiveSwapDetails(id, time.Now().Unix()+100)
	// 		Expect(err).Should(BeNil())
	// 		Expect(bytes.Compare(sendAddr, recvAddr)).Should(Equal(0))
	// 	})
	// })
})
