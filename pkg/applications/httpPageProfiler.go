package applications

import (
	"sync"

	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

var httpPageProfiler k.Profiler
var hppMutex sync.Mutex

func HttpPageProfiler() k.Profiler {
	if httpPageProfiler == nil {
		hppMutex.Lock()
		if httpPageProfiler == nil {
			httpPageProfiler = NewHttpPageProfiler(9081, "default")
		}
		hppMutex.Unlock()
	}
	return httpPageProfiler
}

func NewHttpPageProfiler(port int, name string) k.Profiler {
	httpPage := listeners.NewHttpPage(port, name)
	httpPage.Serve()
	handlers := make(listeners.IntervalListeners, 2)
	handlers[0] = httpPage
	handlers[1] = &listeners.HistoryLimiter{LimitLength: 10000000}
	profiler := k.NewProfiler(httpPage)
	return profiler
}
