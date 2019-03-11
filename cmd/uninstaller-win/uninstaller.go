package uninstaller

import "github.com/renproject/swapperd/driver/winexec"

func main() {
	winexec.StopService("swapperd")
	winexec.StopService("swapperd-win")
	winexec.DeleteService("swapperd-win")
	winexec.DeleteService("swapperd")
}
