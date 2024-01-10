package common

import (
	"github.com/sirupsen/logrus"
)

type LogEntry interface {
	Info(msg string)
	Infof(format string, args ...any)
	Error(msg string)
	Errorf(format string, args ...any)
	WithError(err error) LogEntry
	WithOBUID(obuid int) LogEntry
	WithTraceID(traceID string) LogEntry
}

type Logger interface {
	New() LogEntry
}

type CustomLogger struct {
	logger *logrus.Logger
}

type CustomLogEntry struct {
	entry *logrus.Entry
}

func NewCustomLogger() Logger {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	return &CustomLogger{
		logger: logrus.StandardLogger(),
	}
}

func (e *CustomLogEntry) Info(msg string) {
	e.entry.Info(msg)
}

func (e *CustomLogEntry) Infof(format string, args ...any) {
	e.entry.Infof(format, args...)
}

func (e *CustomLogEntry) Error(msg string) {
	e.entry.Error(msg)
}

func (e *CustomLogEntry) Errorf(format string, args ...any) {
	e.entry.Errorf(format, args...)
}

func (l *CustomLogger) New() LogEntry {
	return &CustomLogEntry{
		entry: logrus.NewEntry(l.logger),
	}
}

func (e *CustomLogEntry) WithError(err error) LogEntry {
	e.entry = e.entry.WithField("error", err)
	return e
}

func (e *CustomLogEntry) WithOBUID(obuid int) LogEntry {
	e.entry = e.entry.WithField("OBUID", obuid)
	return e
}

func (e *CustomLogEntry) WithTraceID(traceID string) LogEntry {
	e.entry = e.entry.WithField("traceID", traceID)
	return e
}
