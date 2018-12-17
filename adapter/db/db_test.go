package db_test

import (
	"encoding/base64"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/adapter/db"

	"github.com/republicprotocol/swapperd/foundation"
	"github.com/syndtr/goleveldb/leveldb"
)

var _ bool = Describe("DB", func() {
	Context("when interacting with the DB", func() {
		It("should be able to read/write from/to the db", func() {
			ldb, err := leveldb.OpenFile("./db-test", nil)
			Expect(err).ShouldNot(HaveOccurred())
			db := New(ldb)
			defer ldb.Close()

			test := func(swap foundation.SwapRequest) bool {
				swap.ID = foundation.SwapID(base64.StdEncoding.EncodeToString([]byte(swap.ID)))
				swap.DelayInfo = []byte("null")

				Expect(db.InsertSwap(swap)).ShouldNot(HaveOccurred())
				stored, err := db.PendingSwap(swap.ID)
				Expect(err).ShouldNot(HaveOccurred())

				return reflect.DeepEqual(swap, stored)
			}

			Expect(quick.Check(test, &quick.Config{})).ShouldNot(HaveOccurred())
		})

		It("should be able to read/write from/to the db", func() {
			ldb, err := leveldb.OpenFile("./db-test", nil)
			Expect(err).ShouldNot(HaveOccurred())
			db := New(ldb)

			value, ok := quick.Value(reflect.TypeOf(foundation.SwapRequest{}), rand.New(rand.NewSource(time.Now().Unix())))
			Expect(ok).Should(BeTrue())
			swap := value.Interface().(foundation.SwapRequest)
			swap.ID = foundation.SwapID(base64.StdEncoding.EncodeToString([]byte(swap.ID)))
			swap.DelayInfo = []byte("null")

			Expect(db.InsertSwap(swap)).ShouldNot(HaveOccurred())
			stored, err := db.PendingSwap(swap.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(reflect.DeepEqual(swap, stored)).Should(BeTrue())
		})
	})
})
