package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	log "svclog"
)

// responseWriter wraps the std http.ResponseWriter
// to keep track of the status code.
type responseWriter struct {
	w    http.ResponseWriter
	code int
}

func (m responseWriter) Header() http.Header {
	return m.w.Header()
}

func (m *responseWriter) Write(bytes []byte) (int, error) {
	m.code = http.StatusOK
	return m.w.Write(bytes)
}

func (m *responseWriter) WriteHeader(statusCode int) {
	m.code = statusCode
	m.w.WriteHeader(statusCode)
}

// generateRequestID returns a crypto/rand 16 byte string.
func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// withRequestID configures the current request with a unique id.
func withRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate a random id for the current request.
		reqID := generateRequestID()
		ctx := context.WithValue(r.Context(), "requestID", reqID)

		// Create a custom response writer to track
		// response status code.
		rw := &responseWriter{w: w, code: 0}
		next(rw, r.WithContext(ctx))
	}
}

// withLogger configures a logging context for incoming request.
func withLogger(logger log.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value("requestID").(string)
		reqLogger := logger.With("requestID", reqID, "requestURI", r.RequestURI, "method", r.Method)
		reqLogger.Print("Request.")

		t := time.Now()
		ctx := context.WithValue(r.Context(), "reqLogger", reqLogger)
		next(w, r.WithContext(ctx))

		rw := w.(*responseWriter)
		reqLogger.Print("Response.", "code", rw.code,
			"duration", fmt.Sprintf("%dms", time.Now().Sub(t).Milliseconds()))
	}
}

// LoggerFromContext extracts the request logger from the context.
// Returns a Nil logger if it doesn't exist in context.
func LoggerFromContext(ctx context.Context) log.Logger {
	l, ok := ctx.Value("requestLogger").(log.Logger)
	if !ok {
		return log.NewNilLogger()
	}
	return l
}

// protectedRoute enriches the request context with protected data.
func protectedRoute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update the request logger with protected data.
		requestLogger := LoggerFromContext(r.Context())
		requestLogger = requestLogger.With("protected", true, "access", "admin/*")

		// Update the request logger in context.
		ctx := context.WithValue(r.Context(), "requestLogger", requestLogger)
		next(w, r.WithContext(ctx))
	}
}

// stubHandler is a simple http request handler.
func stubHandler(w http.ResponseWriter, r *http.Request) {
	logger := LoggerFromContext(r.Context())
	logger.Print("stubHandler")

	w.WriteHeader(http.StatusOK)
}

func main() {
	svcLogger := log.NewJSONLogger()
	if os.Getenv("ENV") == "DEV" {
		svcLogger = log.NewKeyvalLogger(log.ColorYellow)
	}

	// Replace the global default logger with a configured logger.
	svcLogger = svcLogger.With("service", "log-example", "version", 1.0)
	log.SetLogger(svcLogger)

	http.HandleFunc("/", withRequestID(withLogger(svcLogger, stubHandler)))
	http.HandleFunc("/admin", withRequestID(withLogger(svcLogger, protectedRoute(stubHandler))))

	// Use the global logger to log.
	log.With("port", "8080").Print("listening")
	_ = http.ListenAndServe(":8080", nil)
}
