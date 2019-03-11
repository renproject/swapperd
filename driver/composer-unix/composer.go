package composer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/renproject/swapperd/driver/notifier"
	"github.com/renproject/swapperd/driver/swapperd"
	"github.com/renproject/swapperd/driver/updater"
	"github.com/republicprotocol/co-go"
	"github.com/sirupsen/logrus"
)

type Composer struct {
	version    string
	homeDir    string
	executable string
	frequency  time.Duration
	logger     logrus.FieldLogger
	updater    *updater.Updater
}

func New() (*Composer, error) {
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
	configData, err := ioutil.ReadFile(fmt.Sprintf("%s/config.json", homeDir))
	if err != nil {
		return nil, err
	}
	config := struct {
		Version   string        `json:"version"`
		Frequency time.Duration `json:"frequency"`
	}{}
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	updater, err := updater.New(nil, nil)
	if err != nil {
		return nil, err
	}
	return &Composer{
		frequency:  config.Frequency * time.Second,
		version:    config.Version,
		homeDir:    homeDir,
		executable: ex,
		updater:    updater,
		logger:     logger,
	}, nil
}

func Run(done <-chan struct{}) {
	composer, err := New()
	if err != nil {
		composer.logger.Error(err)
		os.Exit(1)
	}
	co.ParBegin(
		func() {
			ticker := time.NewTicker(composer.frequency)
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if err := composer.updater.Update(); err != nil {
						composer.logger.Error(err)
					}
				}
			}
		},
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
