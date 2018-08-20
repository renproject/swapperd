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
	switch int64(orderID[0]) % 7 {
	case 0:
		return "\033[30m"
	case 1:
		return "\033[31m"
	case 2:
		return "\033[32m"
	case 3:
		return "\033[33m"
	case 4:
		return "\033[34m"
	case 5:
		return "\033[35m"
	default:
		return "\033[36m"
	}
}
