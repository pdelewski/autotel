package main

import (
	"bufio"
	"fmt"
	"net"
)

func processRequest(message string) {
	fmt.Print("Message Received:", string(message))
}

func main() {
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
