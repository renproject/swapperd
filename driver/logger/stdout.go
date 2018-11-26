package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const white = "\033[m"

type stdOut struct {
}

func NewStdOut() logrus.FieldLogger {
	fieldLogger := logrus.New()
	fieldLogger.SetOutput(os.Stdout)
	return fieldLogger
}
