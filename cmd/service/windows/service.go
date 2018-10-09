package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/republicprotocol/swapperd/utils"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type RenExSwapper struct {
}

func (m *RenExSwapper) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	go func() {
		cmd := exec.Command("swapper", "http", "--keyphrase", "hello")

		ferr, err := os.OpenFile(utils.GetDefaultSwapperHome()+"\\swapper.err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			// elog.Info("Failed to open std err %q: %v", utils.GetDefaultSwapperHome()+"\\swapper.err", err)
			return
		}
		defer ferr.Close()
		cmd.Stderr = ferr

		fout, err := os.OpenFile(utils.GetDefaultSwapperHome()+"\\swapper.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			// elog.Info("Failed to open std out %q: %v", utils.GetDefaultSwapperHome()+"\\swapper.log", err)
			return
		}
		defer fout.Close()
		cmd.Stdout = fout

		if err := cmd.Run(); err != nil {
			// elog.Info("Error running: %v", err)
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	elog.Info(1, strings.Join(args, "-"))
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &RenExSwapper{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
