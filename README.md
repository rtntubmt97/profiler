example

```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	app "github.com/rtntubmt97/profiler/pkg/applications"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

func main() {
	var profiler k.Profiler = app.HttpPageProfiler()
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		mark := k.CreateMark()
		defer profiler.Record("GetProfile", mark)
		fmt.Fprint(w, "Welcome to foo!")
	})
	http.ListenAndServe(":9000", nil)
}

```