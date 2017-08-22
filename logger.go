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

type logger interface {
	LogRequest(req *http.Request) *http.Request
	LogResponse(resp *http.Response)
}

type loggerImpl struct {
	logger *log.Logger
}

func newLogger(out io.Writer) logger {
	return &loggerImpl{
		logger: log.New(out, "", log.LstdFlags),
	}
}

func (l *loggerImpl) LogRequest(req *http.Request) *http.Request {
	l.logger.SetPrefix("[http] --> ")
	dump, _ := httputil.DumpRequest(req, true)
	l.logger.Printf("%s", string(dump))
	return setRequestedAt(req)
}

func (l *loggerImpl) LogResponse(resp *http.Response) {
	dump, _ := httputil.DumpResponse(resp, true)
	lines := strings.Split(string(dump), "\n")
	lines[0] = fmt.Sprintf("%s (%dms)", lines[0], getRespTimeInMillis(resp))
	l.logger.SetPrefix("[http] <-- ")
	l.logger.Print(strings.Join(lines, "\n"))
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
