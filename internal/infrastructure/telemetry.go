package infrastructure

import (
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func NewTelemetry(cfg *config.Config, logger *zerolog.Logger) {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Jaeger.Host), jaeger.WithAgentPort(cfg.Jaeger.Port)),
	)

	if err != nil {
		logger.Fatal().Err(err).Msg("failed connecting to apm exporter")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Jaeger.Service),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
