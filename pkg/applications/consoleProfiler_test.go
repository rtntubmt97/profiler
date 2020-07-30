package applications

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rtntubmt97/profiler/pkg/kernel"
)

func TestConsoleProfiler(t *testing.T) {
	profiler := ConsoleProfiler()

	api1 := func() {
		mark := kernel.CreateMark()
		defer profiler.Record("api1", mark)
		// fmt.Println("api1")
		sleepTime := rand.Int() % 400
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	api2 := func() {
		mark := kernel.CreateMark()
		defer profiler.Record("api2", mark)
		// fmt.Println("api2")
		sleepTime := rand.Int() % 400
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	go func() {
		for range time.Tick(time.Millisecond * time.Duration(500)) {
			api1()
		}
	}()

	go func() {
		for range time.Tick(time.Millisecond * time.Duration(2500)) {
			api2()
		}
	}()

	time.Sleep(10 * time.Second)
}
