package db_test

import (
	"reflect"
	"testing/quick"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/swapperd/adapter/db"

	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/swapperd/testutils"
	"github.com/syndtr/goleveldb/leveldb"
)

var _ bool = Describe("DB", func() {
	Context("when interacting with the DB", func() {
		It("should be able to read/write from/to the db", func() {
			ldb, err := leveldb.OpenFile("./db-test", nil)
			Expect(err).ShouldNot(HaveOccurred())
			db := New(ldb)
			defer ldb.Close()

			test := func(swap swap.SwapBlob) bool {
				swap.DelayInfo = []byte("null")
				Expect(db.PutSwap(swap)).ShouldNot(HaveOccurred())
				stored, err := db.PendingSwap(swap.ID)
				swap.Password = ""
				Expect(err).ShouldNot(HaveOccurred())
				return reflect.DeepEqual(swap, stored)
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})
})
