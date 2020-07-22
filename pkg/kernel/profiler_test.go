package kernel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestProfilerRecord(t *testing.T) {
	profiler1 := profiler{}
	profiler1.profiles = make(map[string]*Profile)
	var wg sync.WaitGroup
	loopRecord := func(name string, sample []int) {
		defer wg.Done()
		for _, v := range sample {
			profiler1.Record(name, int64(v))
		}
	}

	sample1 := randomIntSlice(5)
	sample2 := randomIntSlice(10)
	wg.Add(2)
	go loopRecord("foo", sample1)
	go loopRecord("foo", sample2)
	wg.Wait()

	//just a simple check
	if len(sample1)+len(sample2) != len(profiler1.profiles["foo"].CurrentProcTimes) {
		t.Errorf("lens mismatch\n")
	}
}

func TestNewProfiler(t *testing.T) {
	profiler1 := NewProfiler(1000, new(counterListener))
	profiler1.Record("foo", 1)

	//just manual check
	time.Sleep(time.Second * 5)
}

type counterListener struct {
	count int
}

func (handler *counterListener) Listen(profiles map[string]*Profile) {
	fmt.Printf("interval %d finished\n", handler.count)
	handler.count++
}
