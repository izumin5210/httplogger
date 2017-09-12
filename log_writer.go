package httplogger

import (
	"net/http"
	"time"
)

// RequestLog contains http(s) request information
type RequestLog struct {
	*http.Request
	RequestedAt time.Time
}

// RequestLog contains http(s) response information or errors
type ResponseLog struct {
	*http.Response
	DurationNano int64
	Error        error
}

// LogWriter is interface for writing logs
type LogWriter interface {
	Println(v ...interface{})
}
