package logger

import (
	"encoding/base64"
	"fmt"

	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

const white = "\033[m"

type stdOut struct {
}

func NewStdOut() swapper.Logger {
	return &stdOut{}
}

func (logger *stdOut) LogInfo(swapID foundation.SwapID, msg string) {
	clr := pickColor(swapID)
	fmt.Println(fmt.Sprintf("[INF] (%s%s%s) %s", clr, base64.StdEncoding.EncodeToString(swapID[:]), white, msg))
}

func (logger *stdOut) LogDebug(swapID foundation.SwapID, msg string) {
	clr := pickColor(swapID)
	fmt.Println(fmt.Sprintf("[DEB] (%s%s%s) %s", clr, base64.StdEncoding.EncodeToString(swapID[:]), white, msg))
}

func (logger *stdOut) LogError(swapID foundation.SwapID, err error) {
	clr := pickColor(swapID)
	fmt.Println(fmt.Sprintf("[ERR] (%s%s%s) %s", clr, base64.StdEncoding.EncodeToString(swapID[:]), white, err))
}

func pickColor(orderID [32]byte) string {
	return fmt.Sprintf("\033[3%dm", int64(orderID[0])%6+1)
}
