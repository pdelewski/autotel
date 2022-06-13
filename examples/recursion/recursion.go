package main

func recur(n int) {
	if n > 0 {
		recur(n - 1)
	}
}

func main() {
	rtlib.SumoAutoInstrument()
	recur(5)
}
