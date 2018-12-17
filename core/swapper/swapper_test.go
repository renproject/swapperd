package swapper_test

import (
	"math/rand"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	. "github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/testutils"
)

var Random *rand.Rand
var Logger *logrus.Logger

func init() {
	Random = rand.New(rand.NewSource(time.Now().Unix()))
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
}

var _ = Describe("Swapper", func() {

	init := func() Swapper {
		callback := testutils.NewMockCallback()
		builder := testutils.MockContractBuilder{}
		storage := testutils.NewMockStorage()
		swapper := New(callback, builder, storage, Logger)

		return swapper
	}

	Context("when running the swapper", func() {
		It("should initiate the swap", func() {

		})
	})

	Expect(true).Should(BeTrue())
})
