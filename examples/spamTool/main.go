package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	totalReqCount := 0
	loop := 1000
	for i := 0; i < loop; i++ {
		// be careful of reqCount value, running spam tool and webserver on the same host may cause the request rate too high
		// golang httpServer will drop the request if the request queue reach its limit
		reqCount := rand.Int() % 100
		totalReqCount += reqCount
		fmt.Printf("ReqCount: %d\n", reqCount)

		for j := 0; j < reqCount; j++ {
			resp, _ := http.Get("http://127.0.0.1:9080/foo")
			resp.Body.Close()
		}
		time.Sleep(time.Second)
	}
	fmt.Printf("TotalReqCount: %d\n", totalReqCount)
}
