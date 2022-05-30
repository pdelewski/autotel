package main

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/otel"
	"sumologic.com/autotel/rtlib"
)

func processRequest(message string) {
	fmt.Print("Message Received:", string(message))
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

	// run loop forever (or until ctrl-c)
	for {
		// get message, output
		message, _ := bufio.NewReader(conn).ReadString('\n')
		processRequest(message)
	}
}
