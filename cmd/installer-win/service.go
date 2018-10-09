package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/republicprotocol/swapperd/utils"

	"github.com/kardianos/service"
)

// Config is the runner app config structure.
type Config struct {
	Keyphrase string
	HomeDir   string
}

var logger service.Logger

type program struct {
	exit    chan struct{}
	service service.Service

	*Config

	cmd *exec.Cmd
}

func (p *program) Start(s service.Service) error {
	// Look for exec.
	// Verify home directory.
	keyphrase := ""
	if p.Keyphrase != "" {
		keyphrase = fmt.Sprintf("--keyphrase %s", p.Keyphrase)
	}

	fmt.Println("Path: ", os.Getenv("PATH"))

	p.cmd = exec.Command("C:\\Users\\Administrator\\Desktop\\Swapper\\swapper.exe", "http", keyphrase)
	p.cmd.Dir = p.HomeDir
	p.cmd.Env = os.Environ()

	go p.run()
	return nil
}

func (p *program) run() {
	logger.Info("Starting RenEx Swapper")
	defer func() {
		if service.Interactive() {
			p.Stop(p.service)
		} else {
			p.service.Stop()
		}
	}()

	if p.HomeDir != "" {
		ferr, err := os.OpenFile(p.HomeDir+"\\swapper.err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logger.Warningf("Failed to open std err %q: %v", p.HomeDir+"\\swapper.err", err)
			return
		}
		defer ferr.Close()
		p.cmd.Stderr = ferr

		fout, err := os.OpenFile(p.HomeDir+"\\swapper.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logger.Warningf("Failed to open std out %q: %v", p.HomeDir+"\\swapper.log", err)
			return
		}
		defer fout.Close()
		p.cmd.Stdout = fout
	}

	err := p.cmd.Run()
	if err != nil {
		logger.Warningf("Error running: %v", err)
	}

	return
}
func (p *program) Stop(s service.Service) error {
	close(p.exit)
	logger.Info("Stopping RenEx Swapper")
	if p.cmd.ProcessState.Exited() == false {
		p.cmd.Process.Kill()
	}
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func main() {
	keyphrase := flag.String("keyphrase", "", "keyphrase used to encrypt the keystore")
	homeDir := flag.String("location", utils.GetDefaultSwapperHome(), "Location of Swapper's home directory")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "RenEx Swapper",
		DisplayName: "Swapper",
		Description: "Atomic swapping service built to work with RenEx",
	}

	prg := &program{
		exit: make(chan struct{}),
		Config: &Config{
			Keyphrase: *keyphrase,
			HomeDir:   *homeDir,
		},
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
