package main

import (
	"fmt"
	"context"
	"github.com/pdelewski/autotel/rtlib"
	otel "go.opentelemetry.io/otel"
)

func foo(__tracing_ctx context.Context) {
	__child_tracing_ctx, span := otel.Tracer("foo").Start(__tracing_ctx, "foo")
	_ = __child_tracing_ctx
	defer span.End()
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
	messages := make(chan string)

	go func(__tracing_ctx context.Context) {
		__child_tracing_ctx, span := otel.Tracer("anonymous").Start(__tracing_ctx, "anonymous")
		_ = __child_tracing_ctx
		defer span.End()
		messages <- "ping"
	}(__child_tracing_ctx)

	foo(__child_tracing_ctx)

	msg := <-messages
	fmt.Println(msg)

}
