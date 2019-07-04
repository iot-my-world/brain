package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func Debug(args ...interface{}) {
	logrus.Debug(args)
}

func Info(args ...interface{}) {
	logrus.Info(args)
}

func Warn(args ...interface{}) {
	logrus.Warn(args, "\n", string(debug.Stack()))
}

func Error(args ...interface{}) {
	logrus.Error(args, "\n", string(debug.Stack()))
}

func Panic(args ...interface{}) {
	logrus.Panic(args)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args, "\n", string(debug.Stack()))
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args)
}
