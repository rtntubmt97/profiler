package intervalListeners

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/rtntubmt97/profiler/pkg/kernel"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

func genProfiles(n int) map[string]*k.Profile {
	ret := make(map[string]*k.Profile)

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("api %d", i)
		profile := new(k.Profile)
		historySnap := k.HistorySnap{
			UnixNano:      time.Now().UnixNano(),
			RequestCounts: rand.Int() % 1000,
			AvgProcTimes:  rand.Int() % 1000,
		}
		profile.History = append(profile.History, historySnap)

		ret[name] = profile
	}

	return ret
}

func TestChainLog(t *testing.T) {
	consoleLog := NewConsolLog()
	fileLog, file := NewFileLog("/home/tumd/golang-repositories/profiler/test/out.txt")
	chainLog := make(IntervalListeners, 2)
	chainLog[0] = consoleLog
	chainLog[1] = fileLog

	chainLog.Listen(genProfiles(0))
	time.Sleep(time.Second * 2)
	chainLog.Listen(genProfiles(10))
	time.Sleep(time.Second * 2)

	file.Sync()
}

func TestHistoryLimiter(t *testing.T) {
	historyLimiter := &HistoryLimiter{LimitLength: 5}
	listeners := IntervalListeners{historyLimiter, new(HistoryPrinter)}
	kernel.NewProfiler(1000, listeners)
	time.Sleep(time.Second * 11)
}
