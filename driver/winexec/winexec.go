package winexec

import (
	"os/exec"
)

func CreateService(name, binLocation string) error {
	if err := Command("sc", "create", name, "start=", "auto", "binpath=", binLocation).Run(); err != nil {
		return err
	}
	return Command("sc", "failure", name, "reset=", "0", "actions=", "restart/0/restart/0/restart/0").Run()
}

func StartService(name string) error {
	return Command("sc", "start", name).Run()
}

func StopService(name string) error {
	return Command("sc", "stop", name).Run()
}

func DeleteService(name string) error {
	return Command("sc", "delete", name).Run()
}

func Run(binLocation string) error {
	return Command(binLocation).Run()
}

func Command(args ...string) *exec.Cmd {
	cmd := exec.Command("cmd", append([]string{"/C"}, args...)...)
	// TODO: Update to not show the window
	return cmd
}
