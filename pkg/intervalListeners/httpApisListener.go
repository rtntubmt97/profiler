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

const ProcTimesApiPath = "/api/proc-times"

const RqRateChartApiPath = "/api/highchart/request-rate"
const PrRateChartApiPath = "/api/highchart/process-rate"
const DataTableApiPath = "/api/data-table"

type HttpApis struct {
	CachedData

	dataTableHandler    *DataTableHandler
	rqCountChartHandler *HighChartHandler
	prCountChartHandler *HighChartHandler
}

type CachedData struct {
	ProcTimes []byte
	Summary   []byte

	Histories [][]string
}

func NewHttpApi() *HttpApis {
	return &HttpApis{dataTableHandler: new(DataTableHandler),
		rqCountChartHandler: NewRqRateChartHandler(),
		prCountChartHandler: NewProcRateChartHandler(),
	}
}

func (httpApi *HttpApis) Serve(port int) *HttpApis {
	server := http.NewServeMux()
	httpApi.SetupHandler(server)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	}()
	return httpApi
}

func (httpApi *HttpApis) SetupHandler(server *http.ServeMux) {
	server.HandleFunc(ProcTimesApiPath, httpApi.PrintProcTimes)

	server.Handle(RqRateChartApiPath, httpApi.rqCountChartHandler)
	server.Handle(PrRateChartApiPath, httpApi.prCountChartHandler)
	server.Handle(DataTableApiPath, httpApi.dataTableHandler)
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
	httpApi.prCountChartHandler.Listen(profiles, startTime)
	previousProcTimes := make(map[string][]int)
	for name, profile := range profiles {
		previousProcTimes[name] = profile.PreviousProcTimes
	}

	previousProcTimesJ, err := json.Marshal(previousProcTimes)
	if err == nil {
		httpApi.ProcTimes = previousProcTimesJ
	}
}

type ProfileSummary struct {
	TotalRequestsCount   int64 `json:"totalRequestsCount"`
	IntervalRequestCount int   `json:"intervalRequestCount"`
	IntervalProcTime     int   `json:"intervalProcTime"`
}
