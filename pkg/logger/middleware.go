package logger

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/go-chi/chi/middleware"
)

func Middleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = ksuid.New().String()
			}
			ctx := r.Context()
			log := logger.With().Str("request-id", id).Logger()
			ctx = WithContext(ctx, &log)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))

			// log end request
			end := time.Now()
			log.Info().
				Str("type", "access").Timestamp().
				Fields(map[string]interface{}{
					"remote_ip":  r.RemoteAddr,
					"url":        r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"ua":         r.Header.Get("User-Agent"),
					"status":     ww.Status(),
					"latency_ms": float64(start.Sub(end).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  ww.BytesWritten(),
				}).Msg("request_in")

		}
		return http.HandlerFunc(fn)
	}
}

type loggerKey struct{}

func WithContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *zerolog.Logger {
	logger := ctx.Value(loggerKey{})
	l := logger.(*zerolog.Logger)
	return l
}

func FromRequest(r *http.Request) *zerolog.Logger {
	return FromContext(r.Context())
}
