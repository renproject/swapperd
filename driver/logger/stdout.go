package logger

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
)

const white = "\033[m"

type stdOut struct {
}

func NewStdOut() logger.Logger {
	return &stdOut{}
}

func (logger *stdOut) LogInfo(orderID [32]byte, msg string) {
	clr := pickColor(orderID)
	log.Println(fmt.Sprintf("[INF] (%s%s%s) %s", clr, order.Fmt(orderID), white, msg))
}

func (logger *stdOut) LogDebug(orderID [32]byte, msg string) {
	clr := pickColor(orderID)
	log.Println(fmt.Sprintf("[DEB] (%s%s%s) %s", clr, order.Fmt(orderID), white, msg))
}

func (logger *stdOut) LogError(orderID [32]byte, msg string) {
	clr := pickColor(orderID)
	log.Println(fmt.Sprintf("[ERR] (%s%s%s) %s", clr, order.Fmt(orderID), white, msg))
}

func pickColor(orderID [32]byte) string {
	return fmt.Sprintf("\033[3%dm", int64(orderID[0])%7)
}
