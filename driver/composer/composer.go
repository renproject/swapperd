package composer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/renproject/swapperd/driver/notifier"
	"github.com/renproject/swapperd/driver/swapperd"
	"github.com/republicprotocol/co-go"
	"github.com/sirupsen/logrus"
)

type composer struct {
	version    string
	homeDir    string
	executable string
	logger     logrus.FieldLogger
}

func New(version string) (*composer, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	homeDir := filepath.Dir(filepath.Dir(ex))
	logFile, err := os.OpenFile(fmt.Sprintf("%s/swapperd.log", homeDir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetOutput(logFile)
	return &composer{
		version:    version,
		homeDir:    homeDir,
		executable: ex,
		logger:     logger,
	}, nil
}

func Run(version string, done <-chan struct{}) {
	composer, err := New(version)
	if err != nil {
		composer.logger.Error(err)
		os.Exit(1)
	}
	co.ParBegin(
		func() {
			swapperd.New(composer.version, composer.homeDir, "testnet", "17927", composer.logger).Run(done)
		},
		func() {
			swapperd.New(composer.version, composer.homeDir, "mainnet", "7927", composer.logger).Run(done)
		},
		func() {
			notifier.New(composer.logger).Watch(
				done,
				path.Join(composer.homeDir, "config.json"),
				path.Join(composer.homeDir, "mainnet.json"),
				path.Join(composer.homeDir, "testnet.json"),
				composer.executable,
			)
		},
	)
}
