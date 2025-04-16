package registry

type LoggerInterface interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Error(format string, args ...any)
}

type defaultLogger struct{}

func (l *defaultLogger) Debug(format string, args ...any) {
	// Implement debug logging
}
func (l *defaultLogger) Info(format string, args ...any) {
	// Implement info logging
}
func (l *defaultLogger) Error(format string, args ...any) {
	// Implement error logging
}
