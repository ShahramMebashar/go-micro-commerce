package logger

import (
	"fmt"
	"io"
	"log"
	"maps"
	"os"
)

type LogLevel int

const (
	// Debug level for detailed troubleshooting
	Debug LogLevel = iota
	// Info level for general operational information
	Info
	// Warn level for potentially harmful situations
	Warn
	// Error level for error events
	Error
	// Fatal level for very severe error events that will lead the application to abort
	Fatal
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
)

func (l LogLevel) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// getColorForLevel returns the ANSI color code for the given log level
func getColorForLevel(level LogLevel) string {
	switch level {
	case Debug:
		return Cyan
	case Info:
		return Green
	case Warn:
		return Yellow
	case Error:
		return Red
	case Fatal:
		return Bold + Red
	default:
		return Reset
	}
}

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	Fatal(format string, args ...any)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}

// SimpleLogger is a basic implementation of the Logger interface
type SimpleLogger struct {
	level     LogLevel
	logger    *log.Logger
	fields    map[string]any
	useColors bool
}

func NewLogger(level LogLevel, out io.Writer, useColors bool) *SimpleLogger {
	if out == nil {
		out = os.Stdout
	}
	return &SimpleLogger{
		level:     level,
		logger:    log.New(out, "", log.LstdFlags),
		fields:    make(map[string]any, 0),
		useColors: useColors,
	}
}

// log logs a message at the specified level
func (l *SimpleLogger) log(level LogLevel, format string, args ...any) {
	if level < l.level {
		return
	}
	// Format the message
	msg := fmt.Sprintf(format, args...)

	// Add fields if any
	if len(l.fields) > 0 {
		fields := ""
		for k, v := range l.fields {
			fields += fmt.Sprintf(" %s=%v", k, v)
		}
		msg = msg + fields
	}

	levelStr := level.String()

	if l.useColors {
		color := getColorForLevel(level)
		levelStr = color + levelStr + Reset
	}
	// Log with level prefix
	l.logger.Printf("[%s] %s", levelStr, msg)

	// Exit if fatal
	if level == Fatal {
		os.Exit(1)
	}
}

func (l *SimpleLogger) Debug(format string, args ...any) {
	l.log(Debug, format, args...)
}

func (l *SimpleLogger) Info(format string, args ...any) {
	l.log(Info, format, args...)
}

func (l *SimpleLogger) Warn(format string, args ...any) {
	l.log(Warn, format, args...)
}

func (l *SimpleLogger) Error(format string, args ...any) {
	l.log(Error, format, args...)
}

func (l *SimpleLogger) Fatal(format string, args ...any) {
	l.log(Fatal, format, args...)
}

func (l *SimpleLogger) WithField(key string, value any) Logger {
	newLogger := &SimpleLogger{
		level:     l.level,
		logger:    l.logger,
		fields:    make(map[string]any),
		useColors: l.useColors, // Preserve color setting
	}

	maps.Copy(newLogger.fields, l.fields)
	newLogger.fields[key] = value
	return newLogger
}

func (l *SimpleLogger) WithFields(fields map[string]any) Logger {
	newLogger := &SimpleLogger{
		level:     l.level,
		logger:    l.logger,
		fields:    make(map[string]any),
		useColors: l.useColors, // Preserve color setting
	}

	maps.Copy(newLogger.fields, l.fields)
	maps.Copy(newLogger.fields, fields)

	return newLogger
}

// EnableColors enables or disables color output
func (l *SimpleLogger) EnableColors(enable bool) {
	l.useColors = enable
}

// Default logger instance
var defaultLogger = NewLogger(Info, os.Stdout, true)

// SetDefaultLogger sets the default logger
func SetDefaultLogger(logger *SimpleLogger) {
	defaultLogger = logger
}

// GetDefaultLogger returns the default logger
func GetDefaultLogger() *SimpleLogger {
	return defaultLogger
}

// Debug logs a debug message using the default logger
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message using the default logger
func Infof(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message using the default logger
func Errorf(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal logs a fatal message and exits using the default logger
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// WithField returns a new logger with the field added using the default logger
func WithField(key string, value interface{}) Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields returns a new logger with the fields added using the default logger
func WithFields(fields map[string]interface{}) Logger {
	return defaultLogger.WithFields(fields)
}

// EnableColors enables or disables color output for the default logger
func EnableColors(enable bool) {
	defaultLogger.EnableColors(enable)
}
