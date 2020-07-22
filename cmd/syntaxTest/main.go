package main

import (
	"fmt"
	"log"
	"os"
	"time"

	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	httpPage := listeners.NewHttpPage()
	httpPage.Serve(9081)
	time.Sleep(time.Hour)
}
