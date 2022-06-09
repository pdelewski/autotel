package main

import "fmt"

func foo() {}

func main() {

	messages := make(chan string)

	go func() { messages <- "ping" }()

	foo()

	msg := <-messages
	fmt.Println(msg)

}
