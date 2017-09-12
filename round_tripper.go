package httplogger

import (
	"io"
	"net/http"
	"time"
)

type loggingTransport struct {
	logger httpLogger
	parent http.RoundTripper
}

// NewRoundTripper returns new RoundTripper instance for logging http request and response
func NewRoundTripper(out io.Writer, parent http.RoundTripper) http.RoundTripper {
	return &loggingTransport{
		logger: defaultHTTTPLogger(out),
		parent: parent,
	}
}

// FromLogger creates new logging RoundTripper instance with given log writer
func FromLogger(writer LogWriter, parent http.RoundTripper) http.RoundTripper {
	return &loggingTransport{
		logger: newHTTPLogger(writer),
		parent: parent,
	}
}

func (lt *loggingTransport) parentTransport() http.RoundTripper {
	if lt.parent == nil {
		return http.DefaultTransport
	}
	return lt.parent
}

func (lt *loggingTransport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := lt.parentTransport().(canceler); ok {
		cr.CancelRequest(req)
	}
}

func (lt *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	requestedAt := time.Now()
	lt.logger.LogRequest(&RequestLog{Request: req, RequestedAt: requestedAt})

	resp, err := lt.parentTransport().RoundTrip(req)

	respTime := time.Since(requestedAt)
	lt.logger.LogResponse(&ResponseLog{Response: resp, DurationNano: respTime.Nanoseconds(), Error: err})

	return resp, err
}
