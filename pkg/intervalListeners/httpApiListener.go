package intervalListeners

import (
	"encoding/json"
	"fmt"
	"net/http"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

const ProcTimesApiPath = "/api/procTimes"
const SummaryApiPath = "/api/summary"

type HttpApi struct {
	CachedData
}

type CachedData struct {
	ProcTimes []byte
	Summary   []byte
}

func NewHttpApi() *HttpApi {
	return new(HttpApi)
}

func (httpApi *HttpApi) Serve(port int) *HttpApi {
	server := http.NewServeMux()
	httpApi.SetupHandler(server)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	}()
	return httpApi
}

func (httpApi *HttpApi) SetupHandler(server *http.ServeMux) {
	server.HandleFunc(ProcTimesApiPath, httpApi.PrintProcTimes)
	server.HandleFunc(SummaryApiPath, httpApi.PrintSummary)
}

func (httpApi *HttpApi) PrintProcTimes(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, httpApi.ProcTimes)
}

func (httpApi *HttpApi) PrintSummary(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, httpApi.Summary)
}

func JsonHandle(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (httpApi *HttpApi) Listen(profiles map[string]*k.Profile) {
	previousProcTimes := make(map[string][]int)
	summary := make(map[string]*ProfileSummary)
	for name, profile := range profiles {
		previousProcTimes[name] = profile.PreviousProcTimes

		currentIndex := len(profile.History) - 1
		profileSummary := &ProfileSummary{
			TotalRequestsCounts:   profile.TotalRequestCount,
			IntervalRequestCounts: profile.History[currentIndex].RequestCounts,
			IntervalProcTime:      profile.History[currentIndex].AvgProcTimes,
		}
		summary[name] = profileSummary
	}

	previousProcTimesJ, err := json.Marshal(previousProcTimes)
	if err == nil {
		httpApi.ProcTimes = previousProcTimesJ
	}

	summaryJ, err := json.Marshal(summary)
	if err == nil {
		httpApi.Summary = summaryJ
	}
}

type ProfileSummary struct {
	TotalRequestsCounts   int64 //`json:"totalRequestsCount"`
	IntervalRequestCounts int
	IntervalProcTime      int
}
