package autoupdater

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

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
	updaterPath := filepath.Join(homeDir, fmt.Sprintf("updater%s", path.Ext(ex)))
	ticker := time.NewTicker(frequency)
	for {
		select {
		case <-done:
		case <-ticker.C:
			if err := update(updaterPath); err != nil {
				logger.Error(err)
			}
		}
	}
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
