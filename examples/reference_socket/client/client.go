package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"

	"go.opentelemetry.io/otel"
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
	_, span := otel.Tracer("main").Start(ctx, "main")
	defer func() {
		span.End()
	}()

	// connect to server
	conn, _ := net.Dial("tcp", "127.0.0.1:8000")
	for {
		// what to send?
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to server
		fmt.Fprintf(conn, text+"\n")
	}
}
