package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"github.com/pdelewski/autotel/rtlib"
)

func processRequest(message string) {
	carrier := propagation.MapCarrier{}
	carrier.Set("traceparent", message)

	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	parentCtx := propgator.Extract(context.Background(), carrier)

	_, span := otel.Tracer("processRequest").Start(parentCtx, "processRequest")

	defer func() {
		span.End()
	}()
	fmt.Print("Message Received:", string(message))
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
	_, span := otel.Tracer("main").Start(ctx, "main")
	defer func() {
		span.End()
	}()

	fmt.Println("Start server...")

	// listen on port 8000
	ln, _ := net.Listen("tcp", ":8000")

	// accept connection
	conn, _ := ln.Accept()

	message, _ := bufio.NewReader(conn).ReadString('\n')
	processRequest(message)

}
