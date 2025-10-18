package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var LogMode = os.Getenv("APP_ENV") // "dev" ou "prod"

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, &bytes.Buffer{}}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := lrw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacker not supported")
}

func (lrw *loggingResponseWriter) Flush() {
	if f, ok := lrw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (lrw *loggingResponseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := lrw.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		if isWebSocketRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		latency := time.Since(start)

		logData := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      lrw.statusCode,
			"latency_ms":  latency.Milliseconds(),
			"req_body":    string(reqBody),
			"res_body":    lrw.body.String(),
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
		}

		if LogMode == "prod" {
			// JSON
			j, _ := json.Marshal(logData)
			fmt.Println(string(j))
		} else {
			fullPath := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.Query().Encode())
			fmt.Printf(
				"%s %s %s %d (%dms)\nReq: %s\nRes: %s\n\n",
				time.Now().Format("2006-01-02 15:04:05"),
				r.Method,
				fullPath,
				lrw.statusCode,
				latency.Milliseconds(),
				string(reqBody),
				lrw.body.String(),
			)
		}
	})
}

func isWebSocketRequest(r *http.Request) bool {
	conn := r.Header.Get("Connection")
	upg := r.Header.Get("Upgrade")
	if conn == "upgrade" || conn == "Upgrade" {
		return true
	}
	if upg == "websocket" || upg == "Websocket" || upg == "WebSocket" {
		return true
	}
	return false
}
