package server

import (
	"context"
	"errors"
	"time"

	"github.com/republicprotocol/tau"
)

var ErrReceiverIsShuttingDown = errors.New("receiver is shutting down")

type Receiver struct {
	done chan struct{}
	buf  chan tau.Message
}

func NewReceiver(cap int) *Receiver {
	return &Receiver{
		done: make(chan struct{}),
		buf:  make(chan tau.Message, cap),
	}
}

// Shutdown shutsdown the receiver, by closing the done channel.
func (receiver *Receiver) Shutdown() {
	close(receiver.done)
}

// Receive blocks until a message can be read from the Receiver buffer.
func (receiver *Receiver) Receive() (tau.Message, error) {
	ticker := time.NewTicker(30 * time.Second)
	select {
	case <-receiver.done:
		return nil, ErrReceiverIsShuttingDown
	case message := <-receiver.buf:
		return message, nil
	case <-ticker.C:
		return tau.NewTick(time.Now()), nil
	}
}

// Write tries to write to the buffer if the buffer is full, it waits until the
// context expires before returning context deadline exceeded error.
func (receiver *Receiver) Write(ctx context.Context, msg tau.Message) error {
	select {
	case <-receiver.done:
		return ErrReceiverIsShuttingDown
	case <-ctx.Done():
		return ctx.Err()
	case receiver.buf <- msg:
		return nil
	}
}
