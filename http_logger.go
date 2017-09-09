package httplogger

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type keyRequestedAt struct{}

const (
	defaultPrefix = "[http] "
)

// httpLogger is interface for logging http request
type httpLogger interface {
	LogRequest(req *http.Request) *http.Request
	LogResponse(resp *http.Response)
}

type httpLoggerImpl struct {
	writer LogWriter
}

func defaultHTTTPLogger(out io.Writer) httpLogger {
	return newHTTPLogger(log.New(out, defaultPrefix, log.LstdFlags))
}

func newHTTPLogger(writer LogWriter) httpLogger {
	return &httpLoggerImpl{
		writer: writer,
	}
}

func (l *httpLoggerImpl) LogRequest(req *http.Request) *http.Request {
	dump, _ := httputil.DumpRequest(req, true)
	l.writer.Println(fmt.Sprintf("--> %s", strings.Replace(string(dump), "\r\n", "\n", -1)))
	return setRequestedAt(req)
}

func (l *httpLoggerImpl) LogResponse(resp *http.Response) {
	if resp == nil {
		return
	}
	dump, _ := httputil.DumpResponse(resp, true)
	lines := strings.Split(string(dump), "\r\n")
	lines[0] = fmt.Sprintf("<-- %s (%dms)", lines[0], getRespTimeInMillis(resp))
	l.writer.Println(strings.Join(lines, "\n"))
}

func setRequestedAt(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), keyRequestedAt{}, time.Now())
	return req.WithContext(ctx)
}

func getRequestedAt(resp *http.Response) time.Time {
	return resp.Request.Context().Value(keyRequestedAt{}).(time.Time)
}

func getRespTimeInMillis(resp *http.Response) time.Duration {
	return time.Since(getRequestedAt(resp)) / 1e6
}
