package main

import "github.com/pdelewski/autotel/rtlib"

type element struct {
}

type driver struct {
	e element
}

type i interface {
	foo(p int) int
}

type impl struct {
}

func (i impl) foo(p int) int {
  return 5
}

func foo(p int) int {
  return 1
}

func (d driver) process(a int) {
}

func (e element) get(a int) {
}

func main() {
	rtlib.AutotelEntryPoint__()
	d := driver{}
	d.process(10)
	d.e.get(5)
	var in i
	in = impl{}
	in.foo(10)
}
