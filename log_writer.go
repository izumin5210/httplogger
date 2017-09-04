package httplogger

// LogWriter is interface for writing logs
type LogWriter interface {
	Println(v ...interface{})
}
