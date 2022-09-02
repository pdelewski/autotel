// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
