package autoupdater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/renproject/swapperd/driver/notifier"
	"github.com/republicprotocol/co-go"
	"github.com/sirupsen/logrus"
)

func Run(done <-chan struct{}) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	homeDir := filepath.Dir(binPath)

	logFile, err := os.OpenFile(fmt.Sprintf("%s/swapperd-updater.log", homeDir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetOutput(logFile)

	frequency := time.Hour

	if configData, err := ioutil.ReadFile(fmt.Sprintf("%s/config.json", homeDir)); err == nil {
		config := struct {
			Frequency time.Duration `json:"frequency"`
		}{}
		if err := json.Unmarshal(configData, &config); err == nil {
			frequency = config.Frequency * time.Second
		}
	}

	updaterPath := filepath.Join(homeDir, fmt.Sprintf("updater%s", path.Ext(ex)))
	ticker := time.NewTicker(frequency)
	co.ParBegin(
		func() {
			for {
				select {
				case <-done:
				case <-ticker.C:
					if err := update(updaterPath); err != nil {
						logger.Error(err)
					}
				}
			}
		},
		func() {
			notifier.New(logger).Watch(done, filepath.Join(homeDir, "config.json"))
		},
	)
}

func update(updaterPath string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command(updaterPath).Run()
	case "linux":
		return exec.Command(updaterPath).Run()
	case "windows":
		return exec.Command("cmd", "/C", updaterPath).Run()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
