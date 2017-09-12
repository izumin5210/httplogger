package httplogger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type httpTestContext struct {
	mux    *http.ServeMux
	server *httptest.Server
}

func newHTTPTestContext() *httpTestContext {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &httpTestContext{
		mux:    mux,
		server: server,
	}
}

func Test_LoggerRoundTripper(t *testing.T) {
	cases := []struct {
		prefix             string
		createRoundTripper func(out io.Writer, prefix string) http.RoundTripper
	}{
		{
			prefix: defaultPrefix,
			createRoundTripper: func(out io.Writer, prefix string) http.RoundTripper {
				return NewRoundTripper(out, nil)
			},
		},
		{
			prefix: "[httplogger] ",
			createRoundTripper: func(out io.Writer, prefix string) http.RoundTripper {
				return FromSimpleLogger(log.New(out, prefix, log.LstdFlags), nil)
			},
		},
	}

	var (
		path     = "/ping"
		query    = "baz=qux"
		reqBody  = "{\"baz\": \"qux\"}"
		respBody = "{\"message\": \"pong\"}"
		status   = 201
	)

	ctx := newHTTPTestContext()
	defer ctx.server.Close()

	ctx.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, respBody)
	})

	for _, c := range cases {
		buf := bytes.NewBufferString("")
		client := &http.Client{Transport: c.createRoundTripper(buf, c.prefix)}

		_, err := client.Post(fmt.Sprintf("%s%s?%s", ctx.server.URL, path, query), "application/json", strings.NewReader(reqBody))

		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if got, want := buf.String(), c.prefix; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain prefix %q", got, want)
		}

		if got, want := buf.String(), path; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain path %q", got, want)
		}

		if got, want := buf.String(), query; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain query %q", got, want)
		}

		if got, want := buf.String(), reqBody; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain request body %q", got, want)
		}

		if got, want := buf.String(), fmt.Sprint(status); !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain status code %q", got, want)
		}

		if got, want := buf.String(), respBody; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain response body %q", got, want)
		}
	}
}
