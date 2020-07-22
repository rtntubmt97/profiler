package applications

import (
	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

func NewHttpApiProfiler(port int) k.Profiler {
	httpApi := listeners.NewHttpApi()
	httpApi.Serve(port)
	profiler := k.NewProfiler(1000, httpApi)
	return profiler
}
