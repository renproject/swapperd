package main

import (
	"syscall"
)

var (
	beepFunc = syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
)

func beep() {
	beepFunc.Call(0xffffffff)
}
