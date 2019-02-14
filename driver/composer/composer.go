package composer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/renproject/swapperd/driver/notifier"
	"github.com/renproject/swapperd/driver/swapperd"
	"github.com/renproject/swapperd/driver/updater"
	"github.com/republicprotocol/co-go"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Version   string        `json:"version"`
	Frequency time.Duration `json:"frequency"`
}

func Run(done <-chan struct{}) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	homeDir := filepath.Dir(filepath.Dir(ex))
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	configData, err := ioutil.ReadFile(fmt.Sprintf("%s/config.json", homeDir))
	if err != nil {
		panic(err)
	}
	config := Config{}
	if err := json.Unmarshal(configData, &config); err != nil {
		panic(err)
	}

	co.ParBegin(
		func() {
			updater.New(config.Version, homeDir, config.Frequency*time.Second, logger).Run(done)
		},
		func() {
			swapperd.New(config.Version, homeDir, "testnet", "17927", logger).Run(done)
		},
		func() {
			swapperd.New(config.Version, homeDir, "mainnet", "7927", logger).Run(done)
		},
		func() {
			notifier.New(homeDir, logger).Watch(done, "config.json", "mainnet.json", "testnet.json")
		},
	)
}
