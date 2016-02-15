package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

// A custom type that extends http.ResponseWriter interface
// to capture and provide an easy access to http status code
type LogResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Easy way to retrieve the status code
func (w *LogResponseWriter) Status() int {
	return w.status
}

func (w *LogResponseWriter) Size() int {
	return w.size
}

// Returns the header to satisty the http.ResponseWriter interface
func (w *LogResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Capture the size of the data written and satisfy the http.ResponseWriter interface
func (w *LogResponseWriter) Write(data []byte) (int, error) {
	written, err := w.ResponseWriter.Write(data)
	w.size += written
	return written, err
}

// Capture the status code and satisfies the http.ResponseWriter interface
func (w *LogResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logrus *logrus.Logger
	// Name is the name of the application as recorded in latency metrics
	Name string

	logStarting bool

	clock timer
}

// NewLogger returns a new *Logger
func NewLogger() *Logger {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	log.Out = os.Stderr
	return &Logger{
		Logrus:      log,
		Name:        "web",
		logStarting: true,
		clock:       &realClock{},
	}
}

// NewFileLogger writes to a file
func NewFileLogger(w io.Writer) *Logger {
	logger := NewLogger()
	logger.Logrus.Out = w
	return logger
}

// NewCustomMiddleware builds a *Logger with the given level and formatter
func NewCustomMiddleware(level logrus.Level, formatter logrus.Formatter, name string) *Logger {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter
	log.Out = os.Stderr
	return &Logger{
		Logrus:      log,
		Name:        name,
		logStarting: true,
		clock:       &realClock{},
	}
}

// NewMiddlewareFromLogger returns a new *Logger which writes to a given logrus logger.
func NewMiddlewareFromLogger(logger *logrus.Logger, name string) *Logger {
	return &Logger{Logrus: logger, Name: name, logStarting: true, clock: &realClock{}}
}

// SetLogStarting accepts a bool to control the logging of "started handling
// request" prior to passing to the next middleware
func (l *Logger) SetLogStarting(v bool) {
	l.logStarting = v
}

func (l *Logger) LoggerMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := l.clock.Now()

		// Try to get the real IP
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteAddr = realIP
		}

		entry := l.Logrus.WithFields(logrus.Fields{
			"request": r.RequestURI,
			"method":  r.Method,
			"remote":  remoteAddr,
		})

		if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
			entry = entry.WithField("request_id", reqID)
		}

		if l.logStarting {
			entry.Info("started handling request")
		}
		res := &LogResponseWriter{ResponseWriter: w}
		h.ServeHTTP(res, r)

		latency := l.clock.Since(start)
		entry.WithFields(logrus.Fields{
			"status":      res.Status(),
			"text_status": http.StatusText(res.Status()),
			"took":        latency,
			fmt.Sprintf("measure#%s.latency", l.Name): latency.Nanoseconds(),
		}).Info("completed handling request")
	}
	return http.HandlerFunc(fn)
}
