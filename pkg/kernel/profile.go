package kernel

import (
	"sync"
	"time"
)

type Profile struct {
	Name              string
	TotalRequestCount int64
	CurrentProcTimes  []int
	PreviousProcTimes []int
	History
	procTimesLock sync.Mutex
}

type HistorySnap struct {
	UnixNano      int64
	RequestCounts int
	AvgProcTimes  int //microsec
}

type History []HistorySnap

type Mark int64

func CreateMark() Mark {
	return Mark(time.Now().UnixNano())
}

func (pf *Profile) Record(mark Mark) int {
	startNanos := int64(mark)
	pf.procTimesLock.Lock()
	procTimeMicros := time.Duration(time.Now().UnixNano()-startNanos) / (time.Microsecond)
	pf.CurrentProcTimes = append(pf.CurrentProcTimes, int(procTimeMicros))
	pf.procTimesLock.Unlock()
	return int(procTimeMicros)
}

func (pf *Profile) intervalUpdate() {
	pf.procTimesLock.Lock()
	pf.PreviousProcTimes = pf.CurrentProcTimes
	pf.CurrentProcTimes = make([]int, 0, len(pf.PreviousProcTimes))
	pf.procTimesLock.Unlock()

	previousRequestCounts := len(pf.PreviousProcTimes)
	totalProcTime := 0
	for _, v := range pf.PreviousProcTimes {
		totalProcTime += v
	}

	avgProcTime := 0
	if previousRequestCounts != 0 {
		avgProcTime = totalProcTime / previousRequestCounts
	}
	historySlice := HistorySnap{
		UnixNano:      time.Now().UnixNano(),
		RequestCounts: previousRequestCounts,
		AvgProcTimes:  avgProcTime}

	pf.History = append(pf.History, historySlice)
	pf.TotalRequestCount += int64(previousRequestCounts)
}
