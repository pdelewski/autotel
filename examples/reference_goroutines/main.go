package main

import (
	"context"
	"fmt"

	otel "go.opentelemetry.io/otel"
	"sumologic.com/autotel/rtlib"
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

	go func() {
		_, span := otel.Tracer("fib").Start(newCtx, "anonymous")
		defer func() {
			span.End()
		}()
		messages <- "ping"
	}()

	msg := <-messages
	fmt.Println(msg)
}
