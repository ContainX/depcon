// Logging configuration and management
package logger

// Loosely wraps go-logging to offer a possible facade in the future for other logging frameworks

import (
	"github.com/op/go-logging"
)

type LogLevel int

const (
	CRITICAL LogLevel = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelTypes = []logging.Level{
	logging.CRITICAL,
	logging.ERROR,
	logging.WARNING,
	logging.NOTICE,
	logging.INFO,
	logging.DEBUG,
}

var dlog *logging.Logger

func InitWithDefaultLogger(module string) {
	dlog = logging.MustGetLogger(module)
}

func (p LogLevel) unWrap() logging.Level {
	return levelTypes[p]
}

func SetLevel(level LogLevel, module string) {
	logging.SetLevel(level.unWrap(), module)
}

func GetLogger(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}

func Logger() *logging.Logger {
	return dlog
}
