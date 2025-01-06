package logger

import (
	"io"
	"log"
	"os"
)

var LoggerInstance *Logger

func init() {
	LoggerInstance = NewLogger(os.Stdout)
}

func NewLogger(StdOut io.Writer) *Logger {
	var logger = new(Logger)

	logger.STDOUT = os.Stdout
	logger.STDERR = os.Stderr
	logger.debug = false
	logger.showTime = true
	logger.log = log.New(StdOut, "", 0)

	return logger
}

func Info(args ...any) {
	LoggerInstance.Info(args...)
}

func Good(args ...any) {
	LoggerInstance.Good(args...)
}

func Debug(args ...any) {
	LoggerInstance.Debug(args...)
}

func DebugError(args ...any) {
	LoggerInstance.DebugError(args...)
}

func Warn(args ...any) {
	LoggerInstance.Warn(args...)
}

func Error(args ...any) {
	LoggerInstance.Error(args...)
}

func Fatal(args ...any) {
	LoggerInstance.Fatal(args...)
}

func Panic(args ...any) {
	LoggerInstance.Panic(args...)
}

func SetDebug(enable bool) {
	LoggerInstance.SetDebug(enable)
}

func ShowTime(time bool) {
	LoggerInstance.ShowTime(time)
}

func SetStdOut(w io.Writer) {
	LoggerInstance.log.SetOutput(w)
}
