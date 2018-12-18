package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/adapter/server"

	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
)

var _ = Describe("Server Adapter", func() {

	buildHandler := func() Handler {
		homeDir := ""
		network := ""

		wallet, err := keystore.Wallet(homeDir, network)
		Expect(err).Should(BeNil())

		ldb, err := leveldb.NewStore(homeDir, network)
		Expect(err).Should(BeNil())

		storage := db.New(ldb)

		passwordHash, err := keystore.LoadPasswordHash(homeDir, network)
		Expect(err).Should(BeNil())

		logger := logger.NewStdOut()
		return NewHandler(passwordHash, wallet, storage, logger)
	}

	Context("building swap response", func() {
		_ = buildHandler()
	})
})
