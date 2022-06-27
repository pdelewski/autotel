package main

import (
	"context"
	"fmt"

	otel "go.opentelemetry.io/otel"
	"github.com/pdelewski/autotel/rtlib"
)

func main() {

	ts := rtlib.NewTracingState()
	defer func() {
		if err := ts.Tp.Shutdown(context.Background()); err != nil {
			ts.Logger.Fatal(err)
		}
	}()

	otel.SetTracerProvider(ts.Tp)
	ctx := context.Background()
	newCtx, span := otel.Tracer("main").Start(ctx, "main")
	defer func() {
		span.End()
	}()

	messages := make(chan string)

	go func(ctx context.Context) {
		_, span := otel.Tracer("fib").Start(ctx, "anonymous")
		defer func() {
			span.End()
		}()
		messages <- "ping"
	}(newCtx)

	msg := <-messages
	fmt.Println(msg)
}
