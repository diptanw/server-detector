// Package logger implements a simple library for output logging
package logger

import (
	"fmt"
	"io"
)

// Available logging levels
const (
	Error = iota
	Warn
	Info
	Debug
)

// Logger is simple logger
type Logger struct {
	writer io.Writer
	level  int
}

// New return a new Logger instance configured with a given
// logging level and writer
func New(writer io.Writer, level int) Logger {
	return Logger{
		writer: writer,
		level:  level,
	}
}

// Infof writes to the output with a Info logging level
func (l Logger) Infof(format string, vals ...interface{}) {
	l.write(Info, format, vals)
}

// Errorf writes to the output with a Error logging level
func (l Logger) Errorf(format string, vals ...interface{}) {
	l.write(Error, format, vals)
}

// Warnf writes to the output with a Warn logging level
func (l Logger) Warnf(format string, vals ...interface{}) {
	l.write(Warn, format, vals)
}

// Debugf writes to the output with a Debug logging level
func (l Logger) Debugf(format string, vals ...interface{}) {
	l.write(Debug, format, vals)
}

func (l Logger) prefix(level int) (string, bool) {
	var prefix string

	switch level {
	case Error:
		prefix = "ERR: "
	case Warn:
		prefix = "WRN: "
	case Info:
		prefix = "INF: "
	case Debug:
		prefix = "DBG: "
	}

	return prefix, level <= l.level
}

func (l Logger) write(level int, format string, a []interface{}) {
	if prefix, ok := l.prefix(level); ok {
		fmt.Fprintln(l.writer, fmt.Sprintf(prefix+format, a...)) //nolint[errcheck]
	}
}
