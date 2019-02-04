package server

import (
	"fmt"

	"github.com/republicprotocol/tau"
)

type service struct {
	receiver *Receiver
}

func NewService(cap int, receiver *Receiver) tau.Task {
	return tau.New(tau.NewIO(cap), &service{receiver: receiver})
}

func (service *service) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case AcceptRequest:
		request, err := service.receiver.Receive()
		if err != nil {
			return tau.NewError(err)
		}
		return NewAcceptedRequest(request)
	default:
		return tau.NewError(fmt.Errorf("unknown message in server: %T", msg))
	}
}

type AcceptRequest struct {
}

func (req AcceptRequest) IsMessage() {
}

type AcceptedRequest struct {
	Message tau.Message
}

func (req AcceptedRequest) IsMessage() {
}

func NewAcceptedRequest(msg tau.Message) AcceptedRequest {
	return AcceptedRequest{
		Message: msg,
	}
}
