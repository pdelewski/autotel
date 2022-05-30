package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func extract(tp *sdktrace.TracerProvider, traceid string) {
	carrier := propagation.MapCarrier{}
	carrier.Set("traceparent", traceid)
	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	parentCtx := propgator.Extract(context.Background(), carrier)

	_, childSpan := tp.Tracer("propagator").Start(parentCtx, "extract")
	childSpan.End()
}

func main() {

	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	bsp := sdktrace.NewSimpleSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)
	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	ctx, span := tp.Tracer("propagator").Start(context.Background(), "main")
	defer span.End()

	// Serialize the context into carrier
	carrier := propagation.MapCarrier{}
	propgator.Inject(ctx, carrier)
	// This carrier is sent accros the process
	fmt.Println("------")
	fmt.Println(carrier)
	fmt.Println("------")
	// Extract the context and start new span as child
	// In your receiving function
	extract(tp, carrier.Get("traceparent"))
}
