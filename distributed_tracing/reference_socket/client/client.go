package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"github.com/pdelewski/autotel/rtlib"
)

func sendRequest(ctx context.Context, conn net.Conn) {
	thisCtx, span := otel.Tracer("sendRequest").Start(ctx, "sendRequest")
	defer func() {
		span.End()
	}()

	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	carrier := propagation.MapCarrier{}
	propgator.Inject(thisCtx, carrier)

	traceid := carrier["traceparent"]
	fmt.Println(traceid)

	fmt.Fprintf(conn, traceid)
	time.Sleep(2 * time.Second)

}

func main() {
	ts := rtlib.NewTracingState()
	defer func() {
		if err := ts.Tp.Shutdown(context.Background()); err != nil {
			ts.Logger.Fatal(err)
		}
	}()

	otel.SetTracerProvider(ts.Tp)
	ctx := context.Background()
	parentCtx, span := otel.Tracer("main").Start(ctx, "main")
	defer func() {
		span.End()
	}()

	// connect to server
	conn, _ := net.Dial("tcp", "127.0.0.1:8000")

	sendRequest(parentCtx, conn)

}
