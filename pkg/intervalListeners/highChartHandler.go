package intervalListeners

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type HighChartHandler struct {
	profiles map[string]*k.Profile
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
		fmt.Println("empty")
		return
	}

	data := make([][]string, len(handler.profiles)+1)
	sortedNames := make([]string, 0)
	for name, _ := range handler.profiles {
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
	for i := int64(0); i < seriesLength; i++ {
		seriesIdx := seriesLength - i
		timestampIdx := timestampLen - 1 - i*historyStep
		if timestampIdx < 0 {
			data[0][seriesIdx] = fmt.Sprintf("%d", timestampProfile.History[timestampLen-1].UnixNano/time.Millisecond.Nanoseconds()-i*historyStep*time.Second.Milliseconds())
		} else {
			data[0][seriesIdx] = fmt.Sprintf("%d", timestampProfile.History[timestampIdx].UnixNano/time.Millisecond.Nanoseconds())
		}
		for j, name := range sortedNames {
			profile := handler.profiles[name]
			profileHistoryIdx := int64(len(profile.History)) - 1 - i*historyStep
			if profileHistoryIdx < 0 {
				continue
			}
			data[j+1][seriesIdx] = fmt.Sprintf("%d", profile.History[profileHistoryIdx].RequestCounts)
		}
	}

	CsvHanlde(w, r, data)
}

func (handler *HighChartHandler) Listen(profiles map[string]*k.Profile, startTime time.Time) {
	handler.profiles = profiles
}
