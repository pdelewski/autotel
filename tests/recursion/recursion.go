package main

import "github.com/pdelewski/autotel/rtlib"

func recur(n int) {
	if n > 0 {
		recur(n - 1)
	}
}

func main() {
	rtlib.AutotelEntryPoint__()
	recur(5)
}
