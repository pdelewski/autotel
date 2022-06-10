package main

import "fmt"
import "sumologic.com/autotel/rtlib"

func foo() {}

func main() {
	rtlib.SumoAutoInstrument()
	messages := make(chan string)

	go func() { messages <- "ping" }()

	foo()

	msg := <-messages
	fmt.Println(msg)

}
