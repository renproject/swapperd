package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func main() {
	switch runtime.GOOS {
	case "linux":
		if err := exec.Command("systemctl", "--user", "stop", "swapperd.service").Run(); err != nil {
			panic(err)
		}
		if err := exec.Command("rm", path.Join(os.Getenv("HOME"), ".config", "systemd", "user", "swapperd.service")).Run(); err != nil {
			panic(err)
		}
	case "darwin":
		servicePath := path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "ren.swapperd.plist")
		if err := exec.Command("launchctl", "unload", "-w", servicePath).Run(); err != nil {
			panic(err)
		}
		if err := exec.Command("rm", servicePath).Run(); err != nil {
			panic(err)
		}
	case "windows":
		if err := exec.Command("cmd", "/C", "sc", "stop", "swapperd").Run(); err != nil {
			panic(err)
		}
		if err := exec.Command("cmd", "/C", "sc", "delete", "swapperd").Run(); err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unsupported OS: %s", runtime.GOOS))
	}
}
