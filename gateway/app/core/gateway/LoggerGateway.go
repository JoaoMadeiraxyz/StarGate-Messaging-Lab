package gateway

type LoggerGateway interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	With(fields ...interface{}) LoggerGateway
}

var globalLogger LoggerGateway

func GetLogger() LoggerGateway {
	return globalLogger
}

func SetLogger(logger LoggerGateway) {
	globalLogger = logger
}
