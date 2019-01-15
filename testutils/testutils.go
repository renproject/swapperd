package testutils

import (
	"math/rand"
	"testing/quick"
	"time"
)

var DefaultQuickCheckConfig = &quick.Config{
	MaxCount: 256,
	Rand:     rand.New(rand.NewSource(time.Now().Unix())),
}
