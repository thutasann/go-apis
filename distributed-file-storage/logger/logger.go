package logger

import (
	"fmt"
	"time"
)

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

func Success(format string, a ...interface{}) {
	logf("SUCCESS", green, format, a...)
}

func Warn(format string, a ...interface{}) {
	logf("WARNING", yellow, format, a...)
}

func Error(format string, a ...interface{}) {
	logf("ERROR", red, format, a...)
}

func logf(level string, color string, format string, a ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s[%s] [%s] %s%s\n", color, timestamp, level, msg, reset)
}
