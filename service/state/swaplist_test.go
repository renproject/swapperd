package state_test

import (
	"crypto/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/renex-swapper-go/adapter/state"
	"github.com/republicprotocol/renex-swapper-go/driver/logger"
	"github.com/republicprotocol/renex-swapper-go/driver/store"
	. "github.com/republicprotocol/renex-swapper-go/service/state"
)

var _ = Describe("State", func() {
	buildSwapList := func() SwapList {
		return SwapList{
			List: [][32]byte{},
		}
	}

	buildState := func() State {
		ldbStore, err := store.NewLevelDB("../../temp/db")
		Expect(err).Should(BeNil())
		stdLogger := logger.NewStdOut()
		return NewState(state.New(ldbStore, stdLogger))
	}

	randomBytes32 := func() [32]byte {
		bytes32 := [32]byte{}
		rand.Read(bytes32[:])
		return bytes32
	}

	It("should be able to add and delete swaps from local memory", func() {
		list := buildSwapList()

		a := randomBytes32()
		b := randomBytes32()
		c := randomBytes32()
		d := randomBytes32()
		e := randomBytes32()
		f := randomBytes32()

		list.Add(a)
		list.Delete(a)
		list.Add(b)
		list.Add(c)
		list.Delete(c)
		list.Add(d)
		list.Add(e)
		list.Delete(d)
		list.Delete(b)
		list.Add(f)
		list.Delete(e)
		list.Delete(f)
	})

	It("should be able to add and delete swaps from persistent storage", func() {
		state := buildState()

		a := randomBytes32()
		b := randomBytes32()
		c := randomBytes32()
		d := randomBytes32()
		e := randomBytes32()
		f := randomBytes32()

		Expect(state.AddSwap(a)).Should(BeNil())
		Expect(state.DeleteSwap(a)).Should(BeNil())
		Expect(state.AddSwap(b)).Should(BeNil())
		Expect(state.AddSwap(c)).Should(BeNil())
		Expect(state.DeleteSwap(c)).Should(BeNil())
		Expect(state.AddSwap(d)).Should(BeNil())
		Expect(state.AddSwap(e)).Should(BeNil())
		Expect(state.DeleteSwap(d)).Should(BeNil())
		Expect(state.DeleteSwap(b)).Should(BeNil())
		Expect(state.AddSwap(f)).Should(BeNil())
		Expect(state.DeleteSwap(e)).Should(BeNil())
		Expect(state.DeleteSwap(f)).Should(BeNil())
	})
})
