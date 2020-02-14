package logger

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// log := //logger.With().Logger()
		fn := func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			log := FromContext(ctx) //.WithField("request-id", id)

			// fmt.Println(">>> set logger")
			// ctx := context.WithValue(r.Context(), "logger", log)
			ctx = WithContext(ctx, log)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

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

				// log end request
				log.Info().
					Str("type", "access").Timestamp().
					Fields(map[string]interface{}{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"ua":         r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					}).Msg("request_in")
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// // https://github.com/drone/drone/blob/master/logger/logger.go
// func FromRequest(r *http.Request) *zerolog.Logger {
// 	fmt.Println(">>> get logger From Request")
// 	logger := r.Context().Value("logger").(*zerolog.Logger)
// 	// if !ok {
// 	// 	panic("!!!!!!! no logger")
// 	// }
// 	return logger
// }

type loggerKey struct{}

// L is an alias for the the standard logger.
var L = log.Logger

// WithContext returns a new context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func WithContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext retrieves the current logger from the context. If no
// logger is available, the default logger is returned.
func FromContext(ctx context.Context) *zerolog.Logger {
	logger := ctx.Value(loggerKey{})
	if logger == nil {
		return &L
	}
	return logger.(*zerolog.Logger)
}

// FromRequest retrieves the current logger from the request. If no
// logger is available, the default logger is returned.
func FromRequest(r *http.Request) *zerolog.Logger {
	return FromContext(r.Context())
}
