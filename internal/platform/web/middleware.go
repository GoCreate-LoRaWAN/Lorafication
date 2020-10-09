package web

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/pborman/uuid"
	"go.uber.org/zap"
)

// requestIDHeader contains the key of the header field that stores a request ID.
const requestIDHeader = "X-Request-ID"

// responseWriter wraps an http.ResponseWriter so we can
// capture the status code.
type responseWriter struct {
	status int
	http.ResponseWriter
}

// WriteHeader captures the statusCode and then writes it the
// wrapped ResponseWriter.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Hijack implements the http.Hijacker interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("ResponseWriter does not implement http.Hijacker")
	}
	return h.Hijack()
}

// RequestMW is a middleware that creates a request id for each request
// and sets it on the header field X-Request-Id. Also logs the start and
// end of each request.
func RequestMW(logger *zap.Logger, next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()

		ww := &responseWriter{
			status:         http.StatusOK,
			ResponseWriter: w,
		}

		// Check if request ID was passed in header, if it wasn't, generate a request ID.
		id := r.Header.Get(requestIDHeader)
		if id == "" {
			id = uuid.New()
		}

		logger.Info("request received",
			zap.String("method", r.Method),
			zap.String("requestID", id),
			zap.String("requestURI", r.RequestURI))

		defer func() {
			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("requestID", id),
				zap.String("requestURI", r.RequestURI),
				zap.Int64("time_ms", time.Since(st).Milliseconds()),
				zap.Int("status", ww.status))
		}()

		ww.Header().Set(requestIDHeader, id)
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(f)
}
