package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	Z string
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
	resp, _ := http.Get("https://github.com/rtntubmt97/profiler/blob/master/web/static/summary.html")
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
}
