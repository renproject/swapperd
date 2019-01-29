package server

import (
	"github.com/republicprotocol/tau"
)

type Server interface {
	Run(doneCh <-chan struct{})
	Receive() (tau.Message, error)
}
