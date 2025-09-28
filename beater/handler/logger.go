package handler

import (
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/siddontang/go-log/loggers"
	"go.uber.org/zap"
)

type LogpProxyLogger struct {
	loggers.Advanced
	logger *logp.Logger
}

func NewLogpProxyLogger(logger *logp.Logger) *LogpProxyLogger {
	newLogger := logger.
		WithOptions(zap.AddCallerSkip(1)).
		Named("canal")
	return &LogpProxyLogger{
		logger: newLogger,
	}
}

func (l *LogpProxyLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogpProxyLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogpProxyLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogpProxyLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogpProxyLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *LogpProxyLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *LogpProxyLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *LogpProxyLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogpProxyLogger) Debugln(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogpProxyLogger) Infoln(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogpProxyLogger) Warnln(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogpProxyLogger) Errorln(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogpProxyLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *LogpProxyLogger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}
func (l *LogpProxyLogger) Fatalln(args ...any) {
	l.logger.Fatal(args...)
}
func (l *LogpProxyLogger) Print(args ...any) {
	l.logger.Info(args...)
}
func (l *LogpProxyLogger) Printf(format string, args ...any) {
	l.logger.Infof(format, args...)
}
func (l *LogpProxyLogger) Println(args ...any) {
	l.logger.Info(args...)
}
