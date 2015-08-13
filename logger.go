package hitch

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[hitch] ", 0)}
}

func (l *Logger) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l.Printf("Started %s %s", r.Method, r.URL.Path)

		next(w, r)

		res := w.(ResponseWriter)
		l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
	})
}
