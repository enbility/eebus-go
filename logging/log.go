package logging

// Logging needs to be implemented, if the internal logs should be printed
type Logging interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// NoLogging is an empty implementation of Logging which does nothing.
type NoLogging struct{}

func (l *NoLogging) Trace(args ...interface{})                 {}
func (l *NoLogging) Tracef(format string, args ...interface{}) {}
func (l *NoLogging) Debug(args ...interface{})                 {}
func (l *NoLogging) Debugf(format string, args ...interface{}) {}
func (l *NoLogging) Info(args ...interface{})                  {}
func (l *NoLogging) Infof(format string, args ...interface{})  {}
func (l *NoLogging) Error(args ...interface{})                 {}
func (l *NoLogging) Errorf(format string, args ...interface{}) {}
