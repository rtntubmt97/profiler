package applications

import (
	listeners "github.com/rtntubmt97/profiler/pkg/intervalListeners"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

func NewHttpApisProfiler(port int, name string) k.Profiler {
	httpApi := listeners.NewHttpApi(port, name)
	httpApi.Serve(port, name)
	profiler := k.NewProfiler(httpApi)
	return profiler
}
