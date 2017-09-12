package httplogger

import (
	"fmt"
	"io"
	"log"
	"net/http/httputil"
	"strings"
)

const (
	defaultPrefix = "[http] "
)

// httpLogger is interface for logging http request
type httpLogger interface {
	LogRequest(reqLog *RequestLog)
	LogResponse(respLog *ResponseLog)
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

func (l *httpLoggerImpl) LogRequest(reqLog *RequestLog) {
	dump, _ := httputil.DumpRequest(reqLog.Request, true)
	l.writer.Println(fmt.Sprintf("--> %s", strings.Replace(string(dump), "\r\n", "\n", -1)))
}

func (l *httpLoggerImpl) LogResponse(respLog *ResponseLog) {
	if respLog.Response == nil {
		return
	}
	dump, _ := httputil.DumpResponse(respLog.Response, true)
	lines := strings.Split(string(dump), "\r\n")
	lines[0] = fmt.Sprintf("<-- %s (%dms)", lines[0], respLog.DurationNano/1e6)
	l.writer.Println(strings.Join(lines, "\n"))
}
