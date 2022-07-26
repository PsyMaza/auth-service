package middlewares

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

func Logger(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			defer func() {
				logger.Info().
					Str("request_id", GetReqID(r.Context())).
					Str("request_path", r.URL.Path).
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Str("method", r.Method).
					Str("query", r.URL.RawQuery).
					Str("ip", r.RemoteAddr).
					Str("trace.id", trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()).
					Str("user-agent", r.UserAgent()).
					Dur("latency", time.Since(start)).
					Msg("request completed")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
