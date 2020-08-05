package main

import (
	"fmt"
	"net/http"
)

type Mark int64

func proc(mark Mark) {
	fmt.Println(mark)
}

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":9080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("X-Forwarded-For"))
	fmt.Println(r.Header.Get("x-forwarded-for"))
	fmt.Println(r.Header.Get("X-FORWARDED-FOR"))
	fmt.Println(r.Header.Get("Remote_Addr"))
	fmt.Println(r.RemoteAddr)
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
