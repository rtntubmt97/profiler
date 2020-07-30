package intervalListeners

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type ExtractChartData func(profile *k.Profile, startIdx, step int) interface{}

type HighChartHandler struct {
	profiles        map[string]*k.Profile
	displayProfiles map[string]bool
	hideProfiles    map[string]bool

	extractChartData ExtractChartData
}

const maxIntervalTick int = 30

// parameters document: http://legacy.HighCharts.net/usage/server-side
func (handler *HighChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	duration, err1 := strconv.Atoi(r.URL.Query().Get("duration"))
	if err1 != nil {
		return
	}

	durationMillis := int64(duration) * time.Minute.Milliseconds()
	intervalMillis := durationMillis / int64(maxIntervalTick)

	if handler.profiles == nil || handler.profiles["LastIntervalUpdate"] == nil {
		return
	}

	data := make([][]string, len(handler.profiles)+1)
	sortedNames := make([]string, 0)
	for name, _ := range handler.profiles {
		if handler.hideProfiles != nil && handler.hideProfiles[name] {
			continue
		}
		if handler.displayProfiles != nil && !handler.displayProfiles[name] {
			continue
		}
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	historyStep := intervalMillis / time.Second.Milliseconds()
	if historyStep == 0 {
		historyStep = 1
	}

	seriesLength := int64(durationMillis / intervalMillis)
	for i := 0; i < len(handler.profiles)+1; i++ {
		data[i] = make([]string, seriesLength+1)
		if i == 0 {
			data[i][0] = "null"
		} else {
			data[i][0] = sortedNames[i-1]
		}
	}

	timestampProfile := handler.profiles["LastIntervalUpdate"]
	timestampLen := int64(len(timestampProfile.History))
	padding := timestampLen % historyStep // ignore some lasted history for the consistency time series
	if historyStep < padding {
		return
	}
	for i := int64(0); i < seriesLength; i++ {
		seriesIdx := seriesLength - i
		timestampIdx := timestampLen - padding - 1 - i*historyStep
		timestampValue := timestampProfile.History[timestampLen-1].UnixNano/time.Millisecond.Nanoseconds() - i*historyStep*time.Second.Milliseconds()
		if timestampIdx >= 0 {
			timestampValue = timestampProfile.History[timestampIdx].UnixNano / time.Millisecond.Nanoseconds()
		}
		data[0][seriesIdx] = fmt.Sprintf("%d", timestampValue)
		for j, name := range sortedNames {
			profile := handler.profiles[name]
			historyEndIdx := int64(len(profile.History)) - padding - i*historyStep
			if historyEndIdx-historyStep < 0 {
				continue
			}
			// handler.extractChartData = sumRqCounts
			avgCount := handler.extractChartData(profile, int(historyEndIdx-historyStep), int(historyStep))
			data[j+1][seriesIdx] = fmt.Sprintf("%d", avgCount)
		}
	}

	CsvHanlde(w, r, data)
}

func (handler *HighChartHandler) Listen(profiles map[string]*k.Profile, startTime time.Time) {
	handler.profiles = profiles
}

func NewRqRateChartHandler() *HighChartHandler {
	return &HighChartHandler{extractChartData: intervalRqRate}
}

func intervalRqRate(profile *k.Profile, start, step int) interface{} {
	ret := 0
	for i := start; i < start+step; i++ {
		ret += profile.History[i].RequestCounts
	}
	return ret
}

func NewProcRateChartHandler() *HighChartHandler {
	return &HighChartHandler{extractChartData: intervalProcRate}
}

func intervalProcRate(profile *k.Profile, start, step int) interface{} {
	intervalProcRate := int64(0)
	for i := start; i < start+step; i++ {
		if profile.History[i].AvgProcTimes == 0 {
			continue
		}
		intervalProcRate += time.Second.Microseconds() / int64(profile.History[i].AvgProcTimes)
	}
	return intervalProcRate
}
