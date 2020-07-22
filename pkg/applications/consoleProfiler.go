package applications

import (
	"fmt"
	"sync"

	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

var consoleProfiler k.Profiler
var cpMutex sync.Mutex

// func init() {
// 	consoleProfiler = k.NewProfiler(1000, listeners.NewConsolLog())
// }

func ConsoleProfiler() k.Profiler {
	if consoleProfiler == nil {
		cpMutex.Lock()
		if consoleProfiler == nil {
			consoleProfiler = k.NewProfiler(1000, listeners.NewConsolLog())
		}
		fmt.Println("ConsoleProfiler is running")
		cpMutex.Unlock()
	} // check - lock - check flow has better performance than lock - check flow
	return consoleProfiler
}
