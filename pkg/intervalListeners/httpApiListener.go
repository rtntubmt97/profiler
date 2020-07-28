package intervalListeners

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
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
	highChartHandler *HighChartHandler
}

type CachedData struct {
	ProcTimes []byte
	Summary   []byte

	Histories [][]string
}

func NewHttpApi() *HttpApi {
	return &HttpApi{dataTableHandler: new(DataTableHandler), highChartHandler: new(HighChartHandler)}
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

	server.Handle(HistoriesApiPath, httpApi.highChartHandler)
	server.Handle(DataTableApiPath, httpApi.dataTableHandler)
}

func (httpApi *HttpApi) PrintProcTimes(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, httpApi.ProcTimes)
}

func (httpApi *HttpApi) PrintHistories(w http.ResponseWriter, r *http.Request) {
	// histories := [][]string{{"api1", "api2"}, {"1", "1,1"}, {"2", "5"}}
	CsvHanlde(w, r, seriesByRowHistories)
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

func (httpApi *HttpApi) Listen(profiles map[string]*k.Profile, startTime time.Time) {
	httpApi.dataTableHandler.Listen(profiles, startTime)
	httpApi.highChartHandler.Listen(profiles, startTime)
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

var seriesByColumnHistories [][]string
var seriesByRowHistories [][]string

func init() {
	seriesByColumnHistories = make([][]string, 100)

	seriesByColumnHistories[0] = []string{"unixSec", "api1", "api2"}
	currentSec := time.Now().Unix()
	for i := 1; i < 100; i++ {
		row := make([]string, 3)
		row[0] = fmt.Sprintf("%d", (currentSec+int64(i))*time.Second.Milliseconds())
		row[1] = fmt.Sprintf("%d", rand.Int()%100)
		row[2] = fmt.Sprintf("%d", rand.Int()%100)
		seriesByColumnHistories[i] = row
	}

	seriesByRowHistories = make([][]string, 3)
	seriesByRowHistories[0] = make([]string, 100)
	seriesByRowHistories[1] = make([]string, 100)
	seriesByRowHistories[2] = make([]string, 100)
	seriesByRowHistories[0][0] = "null"
	seriesByRowHistories[1][0] = "api1"
	seriesByRowHistories[2][0] = "api2"
	for i := 1; i < 100; i++ {
		seriesByRowHistories[0][i] = fmt.Sprintf("%d", (currentSec+int64(i))*time.Second.Milliseconds())
		seriesByRowHistories[1][i] = fmt.Sprintf("%d", rand.Int()%100)
		seriesByRowHistories[2][i] = fmt.Sprintf("%d", rand.Int()%100)
	}
}
