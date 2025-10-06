package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Configuração para habilitar formato dev ou prod
var LogMode = os.Getenv("APP_ENV") // "dev" ou "prod"

// responseWriter customizado para capturar body e status
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
	lrw.body.Write(b) // salva a resposta
	return lrw.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Lendo o request body
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
		}
		// restaurando o body para que o próximo handler use
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		// Response Writer customizado
		lrw := newLoggingResponseWriter(w)

		// Executa a rota
		next.ServeHTTP(lrw, r)

		// Tempo de execução
		latency := time.Since(start)

		// Dados para log
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

		// Impressão formatada
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
