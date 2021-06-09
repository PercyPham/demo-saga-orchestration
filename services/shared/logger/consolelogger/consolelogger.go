package consolelogger

import (
	"fmt"
	"services.shared/logger"
)

const (
	InfoColor  = "\033[1;34m%s\033[0m"
	TraceColor = "\033[1;36m%s\033[0m"
	WarnColor  = "\033[1;33m%s\033[0m"
	ErrorColor = "\033[1;31m%s\033[0m"
	DebugColor = "\033[0;36m%s\033[0m"
)

type log struct {
	level logger.LogLevel
}

func New() logger.Logger {
	return &log{logger.InfoLevel}
}

func (l *log) Ping() error {
	return nil
}

func (l *log) SetLevel(level logger.LogLevel) {
	l.level = level
}

func (l *log) Trace(args ...interface{}) {
	l.printLog(logger.TraceLevel, TraceColor, args)
}

func (l *log) Debug(args ...interface{}) {
	l.printLog(logger.DebugLevel, DebugColor, args)
}

func (l *log) Info(args ...interface{}) {
	l.printLog(logger.InfoLevel, InfoColor, args)
}

func (l *log) Warn(args ...interface{}) {
	l.printLog(logger.WarnLevel, WarnColor, args)
}

func (l *log) Error(args ...interface{}) {
	l.printLog(logger.ErrorLevel, ErrorColor, args)
}

func (l *log) Fatal(args ...interface{}) {
	l.printLog(logger.FatalLevel, ErrorColor, args)
}

func (l *log) printLog(level logger.LogLevel, color string, args []interface{}) {
	if level >= l.level {
		fmt.Printf(color, "["+logger.LevelText(level)+"] ")
		fmt.Printf(color, args...)
		fmt.Println("")
	}
}
