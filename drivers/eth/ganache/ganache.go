package ganache

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// Start a local Ganache instance.
func Start() *exec.Cmd {
	cmd := exec.Command("ganache-cli", fmt.Sprintf("--account=0x2aba04ee8a322b8648af2a784144181a0c793f1a2e80519418f3d20bbfb22249,1000000000000000000000"))
	cmd.Start()
	time.Sleep(10 * time.Second)
	return cmd
}

// Stop will kill the local Ganache instance
func Stop(cmd *exec.Cmd) {
	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Wait()
}
