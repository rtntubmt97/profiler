package applications

import (
	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

func NewHttpPageProfiler(port int) k.Profiler {
	httpPage := listeners.NewHttpPage()
	httpPage.Serve(port)
	profiler := k.NewProfiler(1000, httpPage)
	return profiler
}
