package main

import "sumologic.com/autotel/rtlib"

type driver struct {
}

func (d driver) process(a int) {
}

func main() {
	rtlib.SumoAutoInstrument()
	d := driver{}
	d.process(10)
}
