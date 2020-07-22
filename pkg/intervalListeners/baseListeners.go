package intervalListeners

import (
	"fmt"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type IntervalListeners []k.IntervalListener

func (listeners IntervalListeners) Listen(profiles map[string]*k.Profile) {
	for i := range listeners {
		listeners[i].Listen(profiles)
	}
}

type ShallowCopy struct {
	profiles map[string]*k.Profile
}

func (listener *ShallowCopy) Listen(profiles map[string]*k.Profile) {
	listener.profiles = profiles
}

type HistoryLimiter struct {
	LimitLength int
}

func (historyLimiter *HistoryLimiter) Listen(profiles map[string]*k.Profile) {
	if len(profiles) == 0 {
		return
	}

	for _, profile := range profiles {
		historyLen := len(profile.History)
		if historyLen <= historyLimiter.LimitLength {
			continue
		}
		profile.History = profile.History[(historyLen - historyLimiter.LimitLength):historyLen]
	}
}

type HistoryPrinter struct{}

func (historyPrinter *HistoryPrinter) Listen(profiles map[string]*k.Profile) {
	for _, profile := range profiles {
		historyLen := len(profile.History)
		fmt.Printf("[%s.History] len: %d; value: %v\n", profile.Name, historyLen, profile.History)
	}
}
