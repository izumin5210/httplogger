package httplogger

import (
	"io"
	"net/http"
)

type loggingTransport struct {
	logger logger
	parent http.RoundTripper
}

// New returns new RoundTripper instance for logging http request and response
func New(out io.Writer, parent http.RoundTripper) http.RoundTripper {
	return &loggingTransport{
		logger: newLogger(out),
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
	req = lt.logger.LogRequest(req)
	resp, err := lt.parentTransport().RoundTrip(req)
	lt.logger.LogResponse(resp)
	return resp, err
}
