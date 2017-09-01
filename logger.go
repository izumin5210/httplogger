package httplogger

// Logger is interface that has the same methods as log.Logger type
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
