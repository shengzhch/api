package log

var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

// Default returns a default logger instance.
func Default() *Logger {
	return defaultLogger
}

// SetDefault sets the default logger instance.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Debug prints a debug-level log by default logger instance.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info prints a info-level log by default logger instance.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn prints a warn-level log by default logger instance.
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error prints a error-level log by default logger instance.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Fatal prints a fatal-level log by default logger instance.
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Panic prints a panic-level log by default logger instance.
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Debugf prints a debug-level log with format by default logger instance.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof prints a info-level log with format by default logger instance.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warnf prints a warn-level log with format by default logger instance.
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Errorf prints a error-level log with format by default logger instance.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatalf prints a fatal-level log with format by default logger instance.
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Panicf prints a panic-level log with format by default logger instance.
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}
