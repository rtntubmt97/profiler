package main

import "fmt"

type Mark int64

func proc(mark Mark) {
	fmt.Println(mark)
}

func main() {
	var mark Mark
	mark = 64
	mark2 := Mark(2)
	proc(mark)
	proc(mark2)
}
