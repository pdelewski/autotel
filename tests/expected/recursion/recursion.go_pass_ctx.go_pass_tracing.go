package main

import (
	"github.com/pdelewski/autotel/rtlib"
	otel "go.opentelemetry.io/otel"
	"context"
)

func recur(n int, __tracing_ctx context.Context) {
	__child_tracing_ctx, span := otel.Tracer("recur").Start(__tracing_ctx, "recur")
	_ = __child_tracing_ctx
	defer span.End()
	if n > 0 {
		recur(n-1, __child_tracing_ctx)
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
	recur(5, __child_tracing_ctx)
}
