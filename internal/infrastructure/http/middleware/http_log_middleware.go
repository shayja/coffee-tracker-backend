// file: internal/infrastructure/http/middleware/http_log_middleware.go
package middleware

import (
	http_utils "coffee-tracker-backend/internal/infrastructure/http"
	"log"
	"net/http"
	"time"
)

// RequestLogger logs HTTP request metadata (method, path, IP, status, duration, bytes written)
// It is panic-safe and production-ready.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			duration := time.Since(start)
			clientIP := http_utils.GetUserIpAddress(r)

			// Recover from panic (prevents server crash and still logs)
			if rec := recover(); rec != nil {
				log.Printf("[PANIC] %s %s %s recovered: %v", r.Method, clientIP, r.URL.Path, rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

			log.Printf("[%s] %s %s %d %dB %dms",
				r.Method,
				clientIP,
				r.URL.Path,
				lrw.statusCode,
				lrw.bytesWritten,
				duration.Milliseconds(),
			)
		}()

		next.ServeHTTP(lrw, r)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += n
	return n, err
}
