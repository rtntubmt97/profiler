package kernel

import (
	"math/rand"
	"sync"
	"testing"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func TestRecord(t *testing.T) {
	profile := Profile{}

	n := 10
	sample := randomIntSlice(n)

	ret := make([]int, n)
	for _, record := range sample {
		ret = append(ret, profile.Record(int64(record)))
	}

	errorIndex := findDifIndex(ret, profile.CurrentProcTimes)
	if errorIndex != MaxInt {
		t.Errorf("index %d got different value %d (sample) vs %d (record)", errorIndex, sample[errorIndex], profile.CurrentProcTimes[errorIndex])
	}
}

func TestConcurrentRecord(t *testing.T) {
	profile := Profile{}

	var wg sync.WaitGroup
	sample1 := randomIntSlice(100)
	ret1 := make([]int, 100)
	sample2 := randomIntSlice(200)
	ret2 := make([]int, 200)
	loopRecord := func(records, ret []int) {
		defer wg.Done()
		for i, record := range records {
			ret[i] = profile.Record(int64(record))
		}
	}
	wg.Add(2)
	go loopRecord(sample1, ret1)
	go loopRecord(sample2, ret2)
	wg.Wait()

	// sample1[0] = 0
	// sample2[0] = 0
	for _, v := range ret1 {
		if !intSliceContains(profile.CurrentProcTimes, v) {
			t.Errorf("records dont contain %d while ret1 does\n", v)
		}
	}
	for _, v := range ret2 {
		if !intSliceContains(profile.CurrentProcTimes, v) {
			t.Errorf("records dont contain %d while ret2 does\n", v)
		}
	}

	// sample1 = append(sample1, 0)
	if len(sample1)+len(sample2) != len(profile.CurrentProcTimes) {
		t.Errorf("len(sample1) + len(sample2 ) != len(profile.currentProcTimes)\n")
		t.Errorf("%d + %d != %d\n", len(sample1), len(sample2), len(profile.CurrentProcTimes))
	}
}

func TestIntervalUpdate(t *testing.T) {
	profile := Profile{}

	n1 := 100
	sample1 := randomIntSlice(n1)
	ret1 := make([]int, n1)
	for _, record := range sample1 {
		ret1 = append(ret1, profile.Record(int64(record)))
	}
	profile.intervalUpdate()

	errorIndex := findDifIndex(ret1, profile.PreviousProcTimes)
	if errorIndex != MaxInt {
		t.Errorf("index %d got different value %d (ret) vs %d (record)", errorIndex, sample1[errorIndex], profile.PreviousProcTimes[errorIndex])
	}

	n2 := 200
	sample2 := randomIntSlice(n2)
	ret2 := make([]int, n2)
	for _, record := range sample2 {
		ret2 = append(ret2, profile.Record(int64(record)))
	}
	profile.intervalUpdate()
	profile.intervalUpdate()

	errorIndex = findDifIndex(ret2, profile.PreviousProcTimes)
	if errorIndex != MaxInt {
		t.Errorf("index %d got different value %d (ret) vs %d (record)", errorIndex, sample2[errorIndex], profile.PreviousProcTimes[errorIndex])
	}

	if profile.History[0].RequestCounts != n1 ||
		profile.History[1].RequestCounts != n2 ||
		profile.History[2].RequestCounts != 0 {
		t.Errorf("requestCounts mismatches\n")
	}

	if sumIntSlice(ret1)/n1 != profile.History[0].AvgProcTimes ||
		sumIntSlice(ret2)/n2 != profile.History[1].AvgProcTimes ||
		0 != profile.History[2].AvgProcTimes {
		t.Errorf("avgProcTimes mismatches\n")
	}
}

func intSliceContains(slice []int, num int) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

func sumIntSlice(slice []int) int {
	ret := 0
	for _, v := range slice {
		ret += v
	}
	return ret
}

func randomIntSlice(n int) []int {
	ret := make([]int, n)
	for i := range ret {
		ret[i] = rand.Int()
	}
	return ret
}

func findDifIndex(sample1, sample2 []int) int {
	for i := range sample2 {
		if sample2[i] != sample2[i] {
			return i
		}
	}
	return MaxInt
}
