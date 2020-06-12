package logging

type ZapLogger struct {

	logger Logger
}

// NewZapLogger 创建封装了zap的对象，该对象是对LoggerV2接口的实现
func NewZapLogger(logger *Logger) *ZapLogger {
	return &ZapLogger{
		logger: *logger,
	}
}

// Info returns
func (zl *ZapLogger) Info(args ...interface{}) {
	zl.logger.Info(args)
}

// Infoln returns
func (zl *ZapLogger) Infoln(args ...interface{}) {
	zl.logger.Info(args...)
}

// Infof returns
func (zl *ZapLogger) Infof(format string, args ...interface{}) {
	zl.logger.Infof(format, args...)
}

// Warning returns
func (zl *ZapLogger) Warning(args ...interface{}) {
	zl.logger.Warn(args...)
}

// Warningln returns
func (zl *ZapLogger) Warningln(args ...interface{}) {
	zl.logger.Warn(args...)
}

// Warningf returns
func (zl *ZapLogger) Warningf(format string, args ...interface{}) {
	zl.logger.Warnf(format, args...)
}

// Error returns
func (zl *ZapLogger) Error(args ...interface{}) {
	zl.logger.Error(args...)
}

// Errorln returns
func (zl *ZapLogger) Errorln(args ...interface{}) {
	zl.logger.Error(args...)
}

// Errorf returns
func (zl *ZapLogger) Errorf(format string, args ...interface{}) {
	zl.logger.Errorf(format, args...)
}

// Fatal returns
func (zl *ZapLogger) Fatal(args ...interface{}) {
	zl.logger.Fatal(args...)
}

// Fatalln returns
func (zl *ZapLogger) Fatalln(args ...interface{}) {
	zl.logger.Fatal(args...)
}

// Fatalf logs to fatal level
func (zl *ZapLogger) Fatalf(format string, args ...interface{}) {
	zl.logger.Fatalf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (zl *ZapLogger) V(v int) bool {
	return false
}