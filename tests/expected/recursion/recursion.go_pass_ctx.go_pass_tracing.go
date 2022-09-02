package main

import (
	"github.com/pdelewski/autotel/rtlib"
	otel "go.opentelemetry.io/otel"
	"context"
)

func recur(__tracing_ctx context.Context, n int) {
	__child_tracing_ctx, span := otel.Tracer("recur").Start(__tracing_ctx, "recur")
	_ = __child_tracing_ctx
	defer span.End()
	if n > 0 {
		recur(__child_tracing_ctx, n-1)
	}
}

func main() {
	__child_tracing_ctx := context.TODO()
	_ = __child_tracing_ctx
	ts := rtlib.NewTracingState()
	defer rtlib.Shutdown(ts)
	otel.SetTracerProvider(ts.Tp)
	ctx := context.Background()
	__child_tracing_ctx, span := otel.Tracer("main").Start(ctx, "main")
	defer span.End()
	rtlib.AutotelEntryPoint__()
	recur(__child_tracing_ctx, 5)
}
