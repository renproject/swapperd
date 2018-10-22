package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetHome() string {
	unix := os.Getenv("HOME")
	if unix != "" {
		return unix
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return windows
	}
	panic("unknown Operating System")
}

func GetDefaultSwapperHome() string {
	unix := os.Getenv("HOME")
	if unix != "" {
		return unix + "/.swapper"
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return strings.Join(strings.Split(windows, "\\"), "\\\\") + "\\swapper"
	}
	panic("unknown Operating System")
}

func CreateDir(loc string) error {
	unix := os.Getenv("HOME")
	if unix != "" {
		cmd := exec.Command("mkdir", "-p", loc)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}

	windows := os.Getenv("userprofile")
	if windows != "" {
		cmd := exec.Command("cmd", "/C", "md", loc)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
		return nil
	}

	return errors.New("unknown Operating System")
}

func BuildDBPath(homeDir string) string {
	unix := os.Getenv("HOME")
	fmt.Print(unix)
	if unix != "" {
		return homeDir + "/db"
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return homeDir + "\\db"
	}
	panic("unknown Operating System")
}

func BuildKeystorePath(homeDir, token, renExNet string, unsafe bool) string {
	if unsafe {
		renExNet = renExNet + "-unsafe"
	}
	unix := os.Getenv("HOME")
	fmt.Print(unix)
	if unix != "" {
		return fmt.Sprintf("%s/%s-%s.json", homeDir, token, renExNet)
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return fmt.Sprintf("%s\\%s-%s.json", homeDir, token, renExNet)
	}
	panic("unknown Operating System")
}
