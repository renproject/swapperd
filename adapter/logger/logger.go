package logger

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domains/order"
	"github.com/republicprotocol/renex-swapper-go/services/logger"
)

const white = "\033[m"

type stdOutLogger struct {
}

func NewStdOutLogger() logger.Logger {
	return &stdOutLogger{}
}

func (logger *stdOutLogger) LogInfo(orderID [32]byte, msg string) {
	clr := pickColor(orderID)
	log.Println(fmt.Sprintf("[INF] (%s%s%s) %s", clr, order.Fmt(orderID), white, msg))
}

func (logger *stdOutLogger) LogDebug(orderID [32]byte, msg string) {
	clr := pickColor(orderID)
	log.Println(fmt.Sprintf("[DEB] (%s%s%s) %s", clr, order.Fmt(orderID), white, msg))
}

func (logger *stdOutLogger) LogError(orderID [32]byte, msg string) {
	clr := pickColor(orderID)
	log.Println(fmt.Sprintf("[ERR] (%s%s%s) %s", clr, order.Fmt(orderID), white, msg))
}

func pickColor(orderID [32]byte) string {
	return fmt.Sprintf("\033[3%dm", int64(orderID[0])%7)
}
