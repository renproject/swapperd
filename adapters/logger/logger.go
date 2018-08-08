package logger

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domains/order"
	"github.com/republicprotocol/renex-swapper-go/services/logger"
)

type stdOutLogger struct {
}

func NewStdOutLogger() logger.Logger {
	return &stdOutLogger{}
}

func (logger *stdOutLogger) LogInfo(orderID [32]byte, msg string) {
	log.Println(fmt.Sprintf("(%s)%s", order.Fmt(orderID), msg))
}

func (logger *stdOutLogger) LogDebug(orderID [32]byte, msg string) {
	log.Println(fmt.Sprintf("(%s)%s", order.Fmt(orderID), msg))
}

func (logger *stdOutLogger) LogError(orderID [32]byte, msg string) {
	log.Println(fmt.Sprintf("(%s)%s", order.Fmt(orderID), msg))
}
