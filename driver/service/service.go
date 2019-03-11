package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

func Create(name, binLocation string) error {
	switch runtime.GOOS {
	case "linux":
		homeDir := filepath.Dir(filepath.Dir(binLocation))
		serviceLocation := path.Join(os.Getenv("HOME"), ".config", "systemd", "user")
		if err := exec.Command("mkdir", "-p", serviceLocation).Run(); err != nil {
			return err
		}
		serviceContent := fmt.Sprintf("[Unit]\nDescription=%s daemon\nAssertPathExists=%s\n\n[Service]\nWorkingDirectory=%s\nExecStart=%s\nRestart=on-failure\nPrivateTmp=true\nNoNewPrivileges=true\n\n# Specifies which signal to use when killing a service. Defaults to SIGTERM.\n# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)\nKillSignal=SIGHUP\n\n[Install]\nWantedBy=default.target", name, homeDir, homeDir, binLocation)
		servicePath := path.Join(serviceLocation, fmt.Sprintf("%s.service", name))
		if err := ioutil.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
			return err
		}
		if err := exec.Command("loginctl", "enable-linger", os.Getenv("whoami")).Run(); err != nil {
			return err
		}
		return exec.Command("systemctl", "--user", "enable", fmt.Sprintf("%s.service", name)).Run()
	case "darwin":
		serviceContent := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\t\n<dict>\t\t\n<key>Label</key>\t\t\n<string>%s</string>\t\t\n<key>ProgramArguments</key>\t\t\n<array>\t\t\t\t\n<string>%s</string>\t\t\n</array>\t\t\n<key>KeepAlive</key>\t\t\n<true/>\t\t\n<key>StandardOutPath</key>\t\t\n<string>/dev/null</string>\t\t\n<key>StandardErrorPath</key>\t\t\n<string>/dev/null</string>\t\n</dict>\n</plist>", name, binLocation)
		servicePath := path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", fmt.Sprintf("%s.plist", name))
		return ioutil.WriteFile(servicePath, []byte(serviceContent), 0755)
	default:
		return fmt.Errorf("unsupported Operating System: %s", runtime.GOOS)
	}
}

func Start(name string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("systemctl", "--user", "start", fmt.Sprintf("%s.service", name)).Run()
	case "darwin":
		servicePath := path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", fmt.Sprintf("%s.plist", name))
		return exec.Command("launchctl", "load", "-w", servicePath).Run()
	default:
		return fmt.Errorf("unsupported Operating System: %s", runtime.GOOS)
	}
}

func Stop(name string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("systemctl", "--user", "stop", fmt.Sprintf("%s.service", name)).Run()
	case "darwin":
		servicePath := path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", fmt.Sprintf("%s.plist", name))
		return exec.Command("launchctl", "unload", "-w", servicePath).Run()
	default:
		return fmt.Errorf("unsupported Operating System: %s", runtime.GOOS)
	}
}

func Delete(name string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("rm", path.Join(os.Getenv("HOME"), ".config", "systemd", "user", fmt.Sprintf("%s.service", name))).Run()
	case "darwin":
		servicePath := path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", fmt.Sprintf("%s.plist", name))
		return exec.Command("rm", servicePath).Run()
	default:
		return fmt.Errorf("unsupported Operating System: %s", runtime.GOOS)
	}
}
