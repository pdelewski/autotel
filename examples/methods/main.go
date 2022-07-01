package main

import "github.com/pdelewski/autotel/rtlib"

type element struct {
}

type driver struct {
	e element
}

type i interface {
	foo(p int)
}

type impl struct {
}

func (i impl) foo(p int) {
}

func (d driver) process(a int) {
}

func (e element) get(a int) {
}

func main() {
	rtlib.SumoAutoInstrument()
	d := driver{}
	d.process(10)
	d.e.get(5)
	i := impl{}
	i.foo(10)
}
