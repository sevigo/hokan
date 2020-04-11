package logger

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi/middleware"
)

func Middleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = ksuid.New().String()
			}
			ctx := r.Context()
			log := logger.WithField("request-id", id)
			ctx = WithContext(ctx, log)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.WithFields(logrus.Fields{
						"recover_info": rec,
						"debug_stack":  debug.Stack(),
					}).Error("log system error")

					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))

			// log end request
			end := time.Now()
			log.WithFields(logrus.Fields{
				"headers":    r.Header,
				"remote_ip":  r.RemoteAddr,
				"url":        r.URL.Path,
				"proto":      r.Proto,
				"method":     r.Method,
				"ua":         r.Header.Get("User-Agent"),
				"status":     ww.Status(),
				"latency_ms": float64(start.Sub(end).Nanoseconds()) / 1000000.0,
				"bytes_in":   r.Header.Get("Content-Length"),
				"bytes_out":  ww.BytesWritten(),
			}).Info("request_in")
		}
		return http.HandlerFunc(fn)
	}
}

type loggerKey struct{}

func DefaultLoggert() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}

func WithContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})
	l := logger
	if l == nil {
		l = DefaultLoggert()
	}
	return l.(*logrus.Entry)
}

func FromRequest(r *http.Request) *logrus.Entry {
	return FromContext(r.Context())
}
