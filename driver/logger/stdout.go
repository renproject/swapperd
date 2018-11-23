package logger

import (
	"encoding/base64"
	"fmt"

	"github.com/republicprotocol/swapperd/foundation"
)

const white = "\033[m"

type stdOut struct {
}

func NewStdOut() foundation.Logger {
	return &stdOut{}
}

func (logger *stdOut) LogInfo(swapID foundation.SwapID, msg string) {
	clr := pickColor(swapID)
	fmt.Println(fmt.Sprintf("[INF] (%s%s%s) %s", clr, swapID, white, msg))
}

func (logger *stdOut) LogDebug(swapID foundation.SwapID, msg string) {
	clr := pickColor(swapID)
	fmt.Println(fmt.Sprintf("[DEB] (%s%s%s) %s", clr, swapID, white, msg))
}

func (logger *stdOut) LogError(swapID foundation.SwapID, err error) {
	clr := pickColor(swapID)
	fmt.Println(fmt.Sprintf("[ERR] (%s%s%s) %s", clr, swapID, white, err))
}

func (logger *stdOut) GlobalLogInfo(msg string) {
	fmt.Println(fmt.Sprintf("[INF] %s", msg))
}

func (logger *stdOut) GlobalLogDebug(msg string) {
	fmt.Println(fmt.Sprintf("[DEB] %s", msg))
}

func (logger *stdOut) GlobalLogError(err error) {
	fmt.Println(fmt.Sprintf("[ERR] %s", err))
}

func pickColor(swapID foundation.SwapID) string {
	swapIDBytes, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return white
	}
	return fmt.Sprintf("\033[3%dm", int64(swapIDBytes[0])%6+1)
}
