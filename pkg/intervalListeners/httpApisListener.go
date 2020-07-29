package intervalListeners

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

const ProfilerInfoApiPath = "/api/profiler-info"
const ProcTimesApiPath = "/api/proc-times"

const RqRateChartApiPath = "/api/highchart/request-rate"
const PrRateChartApiPath = "/api/highchart/process-rate"
const DataTableApiPath = "/api/data-table"

type HttpApis struct {
	port int
	name string
	CachedData

	dataTableHandler     *DataTableHandler
	rqCountChartHandler  *HighChartHandler
	procRateChartHandler *HighChartHandler
}

type CachedData struct {
	ProfilerInfo []byte
	ProcTimes    []byte
	Summary      []byte

	Histories [][]string
}

func NewHttpApi(port int, name string) *HttpApis {
	return &HttpApis{dataTableHandler: new(DataTableHandler),
		rqCountChartHandler:  NewRqRateChartHandler(),
		procRateChartHandler: NewProcRateChartHandler(),
		port:                 port,
		name:                 name,
	}
}

func (httpApi *HttpApis) Serve(port int, name string) *HttpApis {
	server := http.NewServeMux()
	httpApi.SetupHandler(server)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	}()
	return httpApi
}

func (httpApi *HttpApis) SetupHandler(server *http.ServeMux) {
	server.HandleFunc(ProcTimesApiPath, httpApi.PrintProcTimes)
	server.HandleFunc(ProfilerInfoApiPath, httpApi.PrintProfilerInfo)

	server.Handle(RqRateChartApiPath, httpApi.rqCountChartHandler)
	server.Handle(PrRateChartApiPath, httpApi.procRateChartHandler)
	server.Handle(DataTableApiPath, httpApi.dataTableHandler)
}

func (httpApi *HttpApis) PrintProfilerInfo(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, httpApi.CachedData.ProfilerInfo)
}

func (httpApi *HttpApis) PrintProcTimes(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, httpApi.ProcTimes)
}

func JsonHandle(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CsvHanlde(w http.ResponseWriter, r *http.Request, data [][]string) {
	// fmt.Printf("duration: %s\n", r.URL.Query().Get("duration"))
	// fmt.Printf("interval: %s\n", r.URL.Query().Get("interval"))

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=Histories_%d.csv", time.Now().UnixNano()))
	w.WriteHeader(http.StatusOK)

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.WriteAll(data)
	w.Write(b.Bytes())
}

func (httpApi *HttpApis) Listen(profiles map[string]*k.Profile, startTime time.Time) {
	httpApi.dataTableHandler.Listen(profiles, startTime)
	httpApi.rqCountChartHandler.Listen(profiles, startTime)
	httpApi.procRateChartHandler.Listen(profiles, startTime)

	previousProcTimes := make(map[string][]int)
	for name, profile := range profiles {
		previousProcTimes[name] = profile.PreviousProcTimes
	}

	previousProcTimesJ, err := json.Marshal(previousProcTimes)
	if err == nil {
		httpApi.ProcTimes = previousProcTimesJ
	}

	if httpApi.CachedData.ProfilerInfo == nil {
		profilerInfo := ProfilerInfo{
			Name:      httpApi.name,
			Port:      httpApi.port,
			StartTime: startTime.Format("2006-01-02 15:04:05"),
		}
		profilerInfoJ, err := json.Marshal(profilerInfo)
		if err == nil {
			httpApi.CachedData.ProfilerInfo = profilerInfoJ
		}
	}
}

type ProfilerInfo struct {
	Name      string `json:"name"`
	Port      int    `json:"port"`
	StartTime string `json:"startTime"`
}

type ProfileSummary struct {
	TotalRequestsCount   int64 `json:"totalRequestsCount"`
	IntervalRequestCount int   `json:"intervalRequestCount"`
	IntervalProcTime     int   `json:"intervalProcTime"`
}
