package persist

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}
