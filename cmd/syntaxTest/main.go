package main

import (
	"fmt"
)

type Foo struct {
	bar   []int
	bfunc func()
}

type Out struct {
	In
	In2
	fooo func(int) int
}

type In struct {
	i int
}

type In2 struct {
	i int
}

func (in *In) foo() {
	fmt.Println("foo")
}

type A struct {
	F []byte
}

func (a *A) foo() {
	fmt.Println("A")
}

type B struct {
	A
}

func (b *B) foo() {
	fmt.Println("B")
	b.A.foo()
}

func main() {
	var fooMap map[string]bool
	fooMap = nil
	if fooMap == nil {
		fmt.Println("nil")
	}
	t := fooMap["d"]
	fmt.Println(t)
}
