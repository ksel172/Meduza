package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

var LoggerInstance *Logger

func init() {
	LoggerInstance = NewLoggerWithEnv(os.Stdout)
}

// NewLoggerWithEnv creates a logger and configures it from environment variables.
func NewLoggerWithEnv(StdOut io.Writer) *Logger {
	logger := NewLogger(StdOut)

	// Advanced: support levels: debug, info, warn, error, fatal, panic
	// Only one env var controls the level: LOG_LEVEL (default: info)
	level := os.Getenv("LOG_LEVEL")
	level = strings.ToLower(level)

	// Set debug mode only if level is debug
	logger.SetDebug(level == "debug")

	// ShowTime: LOG_SHOW_TIME (default: true)
	showTime := true
	if v, ok := os.LookupEnv("LOG_SHOW_TIME"); ok {
		if v == "false" || v == "0" {
			showTime = false
		}
	}
	logger.ShowTime(showTime)

	// Optionally: allow log output to file
	if path := os.Getenv("LOG_FILE"); path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			logger.log.SetOutput(f)
		}
	}

	return logger
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
