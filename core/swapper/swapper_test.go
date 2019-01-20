package swapper_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	"github.com/republicprotocol/swapperd/core/swapper/delayed"
	"github.com/republicprotocol/swapperd/core/swapper/immediate"

	"github.com/republicprotocol/swapperd/core/swapper/status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation/swap"

	"github.com/republicprotocol/tau"
)

var _ = Describe("Swapper", func() {

	newTask := func(reduce tau.ReduceFunc) tau.Task {
		if reduce == nil {
			return tau.New(tau.NewIO(2048), tau.ReduceFunc(func(message tau.Message) tau.Message {
				fmt.Printf("%T\n", message)
				Expect(message).ShouldNot(HaveOccurred())
				return nil
			}))
		}
		return tau.New(tau.NewIO(2048), reduce)
	}

	quickCheckConfig := func() *quick.Config {
		return &quick.Config{
			Rand:     rand.New(rand.NewSource(time.Now().Unix())),
			MaxCount: 512,
		}
	}

	init := func(storage Storage, delayed, immediate, status tau.Task) tau.Task {
		swapper := New(2048, storage, delayed, immediate, status)
		return swapper
	}

	Context("when receiving bootload message", func() {
		It("when there are no swaps in the database", func() {
			delayed := newTask(nil)
			immediate := newTask(nil)
			status := newTask(nil)

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayed, immediate, status)
			go swapper.Run(done)

			test := func() bool {
				swapper.IO().InputWriter() <- Bootload{}
				msg := <-swapper.IO().OutputReader()
				msgBatch, ok := msg.(tau.MessageBatch)
				Expect(ok).Should(BeTrue())
				Expect(len(msgBatch)).Should(Equal(2))
				Expect(msgBatch[0]).Should(BeNil())
				Expect(msgBatch[1]).Should(BeNil())
				return true
			}

			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})

		It("when there are swaps in the database", func() {
			delayed := newTask(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case delayed.DelayedSwapRequest:
				default:
					Expect(msg).ShouldNot(HaveOccurred())
				}
				return nil
			})

			immediate := newTask(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case immediate.SwapRequest:
				default:
					Expect(msg).ShouldNot(HaveOccurred())
				}
				return nil
			})

			status := newTask(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case status.Receipt, status.ReceiptUpdate:
				default:
					Expect(msg).ShouldNot(HaveOccurred())
				}
				return nil
			})

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayed, immediate, status)
			go swapper.Run(done)

			test := func(blob swap.SwapBlob) bool {
				storage.PutSwap(blob)
				swapper.IO().InputWriter() <- Bootload{}
				msg := <-swapper.IO().OutputReader()
				msgBatch, ok := msg.(tau.MessageBatch)
				Expect(ok).Should(BeTrue())
				Expect(len(msgBatch)).Should(Equal(2))
				Expect(msgBatch[0]).Should(BeNil())
				Expect(msgBatch[1]).Should(BeNil())
				return true
			}

			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving error message", func() {
		It("should return the error message", func() {
			delayed := newTask(nil)
			immediate := newTask(nil)
			status := newTask(nil)

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayed, immediate, status)
			go swapper.Run(done)

			test := func(errStr string) bool {
				err := tau.NewError(fmt.Errorf(errStr))
				swapper.IO().InputWriter() <- err
				msg := <-swapper.IO().OutputReader()
				errMsg, ok := msg.(tau.Error)
				Expect(ok).Should(BeTrue())
				return reflect.DeepEqual(err, errMsg)
			}

			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving tick message", func() {
		It("should send tick message to all children tasks", func() {
			reducer := func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case tau.Tick:
				default:
					Expect(msg).ShouldNot(HaveOccurred())
				}
				return nil
			}

			delayed := newTask(reducer)
			immediate := newTask(reducer)
			status := newTask(reducer)

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayed, immediate, status)
			go swapper.Run(done)

			test := func(errStr string) bool {
				swapper.IO().InputWriter() <- tau.Tick{}
				return true
			}

			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving receipt query message", func() {
		It("should send receipt query message to status task", func() {
			reducer := func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case status.ReceiptQuery:
				default:
					Expect(msg).ShouldNot(HaveOccurred())
				}
				return nil
			}

			delayedTask := newTask(nil)
			immediateTask := newTask(nil)
			statusTask := newTask(reducer)

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayedTask, immediateTask, statusTask)
			go swapper.Run(done)

			test := func() bool {
				swapper.IO().InputWriter() <- status.ReceiptQuery{}
				return true
			}
			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving receipt update message", func() {
		It("should send receipt update to status task", func() {
			delayedTask := newTask(nil)
			immediateTask := newTask(nil)
			statusTask := newTask(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case status.ReceiptUpdate:
				default:
					Expect(msg).ShouldNot(HaveOccurred())
				}
				return nil
			})

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayedTask, immediateTask, statusTask)
			go swapper.Run(done)

			test := func() bool {
				swapper.IO().InputWriter() <- tau.RandomMessage{}
				err := <-swapper.IO().OutputReader()
				_, ok := err.(tau.Error)
				return ok
			}
			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving an unknown message type", func() {
		It("should return an error", func() {
			delayedTask := newTask(nil)
			immediateTask := newTask(nil)
			statusTask := newTask(nil)

			storage := NewMockStorage()
			done := make(chan struct{})
			defer close(done)
			swapper := init(storage, delayedTask, immediateTask, statusTask)
			go swapper.Run(done)

			test := func() bool {
				swapper.IO().InputWriter() <- tau.RandomMessage{}
				err := <-swapper.IO().OutputReader()
				_, ok := err.(tau.Error)
				return ok
			}
			Expect(quick.Check(test, quickCheckConfig())).ShouldNot(HaveOccurred())
		})
	})

})
