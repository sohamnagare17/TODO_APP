package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer() func() {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		//otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),

		// ✅ IMPORTANT: Add service name
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("todo-app"),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
