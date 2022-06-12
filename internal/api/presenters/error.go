package presenters

import (
	"fmt"
	"gitlab.com/g6834/team17/auth-service/internal/api/middlewares"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func (p *presenters) Error(w http.ResponseWriter, r *http.Request, err error) {
	span := trace.SpanFromContext(r.Context())
	span.RecordError(err)

	switch e := err.(type) {
	case models.StatusError:
		p.logger.Error().
			Err(err).
			Str("caller", err.(models.StatusError).Caller).
			Str("request-id", middlewares.GetReqID(r.Context())).
			Str("trace.id", span.SpanContext().TraceID().String()).
			Msg("error.go")

		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), e.Code)
		return
	default:
		p.logger.Error().
			Err(err).
			Str("request-id", middlewares.GetReqID(r.Context())).
			Str("trace.id", span.SpanContext().TraceID().String()).
			Msg("unhandled error.go")

		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusInternalServerError)
		return
	}
}
