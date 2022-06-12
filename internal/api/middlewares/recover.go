package middlewares

import (
	"fmt"
	"github.com/go-stack/stack"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func Recover(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%v", p)
					}

					var stackTrace stack.CallStack
					traces := stack.Trace().TrimRuntime()

					for i := 0; i < len(traces); i++ {
						t := traces[i]
						tFunc := t.Frame().Function

						if tFunc == "runtime.gopanic" || tFunc == "go.opentelemetry.io/otel/sdk/trace.(*span).End" {
							continue
						}

						if tFunc == "net/http.HandlerFunc.ServeHTTP" {
							break
						}
						stackTrace = append(stackTrace, t)
					}

					logger.WithLevel(zerolog.PanicLevel).
						Err(err).
						Str("trace.id", trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()).
						Str("request-id", GetReqID(r.Context())).
						Str("stack", fmt.Sprintf("%+v", stackTrace)).
						Msg("panic")

					http.Error(rw, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}
