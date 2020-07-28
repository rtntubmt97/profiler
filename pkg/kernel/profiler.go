package kernel

import (
	"sync"
	"time"
)

type Profiler interface {
	Record(name string, startNanos int64) int
}

type IntervalListener interface {
	Listen(profiles map[string]*Profile, startTime time.Time)
}

type profiler struct {
	profiles         map[string]*Profile
	startTime        time.Time
	intervalListener IntervalListener
	profilesLock     sync.Mutex
}

func NewProfiler(intervalHandle IntervalListener) *profiler {
	ret := new(profiler)
	ret.intervalListener = intervalHandle

	ret.profiles = make(map[string]*Profile)
	ret.startTime = time.Now()
	ret.Record("LastIntervalUpdate", time.Now().UnixNano())

	go func() {
		for range time.Tick(time.Second) {
			startTime := time.Now().UnixNano()
			for _, profile := range ret.profiles {
				profile.intervalUpdate()
			}
			ret.intervalListener.Listen(ret.profiles, ret.startTime)
			ret.Record("LastIntervalUpdate", startTime)
		}
	}()

	return ret
}

func (pfr *profiler) addProfile(name string) {
	pfr.profilesLock.Lock()
	defer pfr.profilesLock.Unlock()

	if _, ok := pfr.profiles[name]; ok {
		return
	}
	newPf := new(Profile)
	newPf.Name = name
	pfr.profiles[name] = newPf
}

func (pfr *profiler) Record(name string, startNanos int64) int {
	if _, ok := pfr.profiles[name]; !ok {
		pfr.addProfile(name)
	}

	profile, _ := pfr.profiles[name]
	return profile.Record(startNanos)
}
