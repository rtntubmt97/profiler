example

```golang
package main

import (
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"time"

	app "github.com/rtntubmt97/profiler/pkg/applications"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

// this default profiler will run on localhost:9081
var profiler k.Profiler = app.HttpPageProfiler()

func main() {
	http.HandleFunc("/foo", FooHandler)
	http.ListenAndServe(":9080", nil)
}

func FooHandler(rsp http.ResponseWriter, request *http.Request) {
	mark := k.CreateMark()
	defer profiler.Record("FooHandler", mark)

	// doing something in handler
	processTimeMillis := rand.Int() % 100
	time.Sleep(time.Microsecond * time.Duration(processTimeMillis))
	fmt.Fprintf(rsp, "Hello, %q", html.EscapeString(request.URL.Path))

	// process something in db
	n := rand.Int() % 5
	for i := 0; i < n; i++ {
		dbProcess()
	}
}

func dbProcess() {
	mark := k.CreateMark()
	defer profiler.Record("dbProcess", mark)

	processTimeMillis := rand.Int() % 100
	time.Sleep(time.Microsecond * time.Duration(processTimeMillis))
}

```
