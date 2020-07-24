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

const HistoriesApiPath = "/api/histories"
const DataTableApiPath = "/api/data-table"

type HttpApi struct {
	CachedData

	dataTableHandler *DataTableHandler
}

type CachedData struct {
	ProcTimes []byte
	Summary   []byte

	Histories [][]string
}

func NewHttpApi() *HttpApi {
	return &HttpApi{dataTableHandler: new(DataTableHandler)}
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
	server.HandleFunc(HistoriesApiPath, httpApi.PrintHistories)

	server.Handle(DataTableApiPath, httpApi.dataTableHandler)
}

func (httpApi *HttpApi) PrintProcTimes(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, httpApi.ProcTimes)
}

func (httpApi *HttpApi) PrintHistories(w http.ResponseWriter, r *http.Request) {
	histories := [][]string{{"api1", "api2"}, {"1", "1,1"}, {"2", "5"}}
	CsvHanlde(w, r, histories)
}

func JsonHandle(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CsvHanlde(w http.ResponseWriter, r *http.Request, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=Histories_%d.csv", time.Now().UnixNano()))
	w.WriteHeader(http.StatusOK)

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.WriteAll(data)
	w.Write(b.Bytes())
}

func (httpApi *HttpApi) Listen(profiles map[string]*k.Profile, startTime time.Time, intervalTimeMillis int) {
	httpApi.dataTableHandler.Listen(profiles, startTime, intervalTimeMillis)
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
