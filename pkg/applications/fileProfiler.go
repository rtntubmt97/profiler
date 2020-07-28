package applications

import (
	"fmt"
	"sync"

	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

var fileProfiler k.Profiler
var fMutex sync.Mutex

func init() {
	// log, _ := listeners.NewFileLog("/tmp/profiler_data/out.txt")
	// fileProfiler = k.NewProfiler(1000, log)
}

func FileProfiler() k.Profiler {
	if fileProfiler == nil {
		fMutex.Lock()
		if fileProfiler == nil {
			log, _ := listeners.NewFileLog("/tmp/profiler_data/out.txt")
			fileProfiler = k.NewProfiler(log)
		}
		fmt.Println("FileProfiler is running")
		fMutex.Unlock()
	} // check - lock - check flow has better performance than lock - check flow
	return fileProfiler
}

func NewFileProfiler(path string) k.Profiler {
	log, _ := listeners.NewFileLog(path)
	ret := k.NewProfiler(log)
	fmt.Printf("FileProfiler is writting at %s\n", path)
	return ret
}
