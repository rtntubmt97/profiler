package intervalListeners

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type DataTableHandler struct {
	cachedData []byte
}

type DataTableResponse struct {
	ITotalRecords        int        `json:"iTotalRecords"`
	ITotalDisplayRecords int        `json:"iTotalDisplayRecords"`
	SEcho                string     `json:"sEcho"`
	SColumns             string     `json:"sColumns"`
	AaData               [][]string `json:"aaData"`
}

var defaultDTRsp []byte
var columnNames []string = []string{
	"apiName",
	"totalRequestsCounts",
	"intervalRequestCounts",
	"intervalProcTime",
	"intervalProcRate"}

func init() {
	dataRsp := DataTableResponse{}
	dataRsp.AaData = make([][]string, 0)
	data, err := json.Marshal(dataRsp)
	if err == nil {
		defaultDTRsp = data
	} else {
		defaultDTRsp = nil
	}

}

var count int

// parameters document: http://legacy.datatables.net/usage/server-side
func (handler *DataTableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	JsonHandle(w, r, handler.cachedData)
}

func (handler *DataTableHandler) Listen(profiles map[string]*k.Profile, startTime time.Time) {
	dataRsp := DataTableResponse{}
	dataRsp.AaData = make([][]string, 0)

	for name, profile := range profiles {
		currentIndex := len(profile.History) - 1
		avgProcTimes := profile.History[currentIndex].AvgProcTimes
		procRate := 0
		if avgProcTimes != 0 {
			procRate = 1000000 / avgProcTimes
		}

		record := make([]string, 5)
		record[0] = name
		record[1] = fmt.Sprintf("%d", profile.TotalRequestCount)
		record[2] = fmt.Sprintf("%d", profile.History[currentIndex].RequestCounts)
		record[3] = fmt.Sprintf("%d", avgProcTimes)
		record[4] = fmt.Sprintf("%d", procRate)

		dataRsp.AaData = append(dataRsp.AaData, record)
	}
	dataRsp.ITotalRecords = len(dataRsp.AaData)
	dataRsp.ITotalDisplayRecords = len(dataRsp.AaData)

	data, err := json.Marshal(dataRsp)
	if err != nil {
		data = defaultDTRsp
	}

	handler.cachedData = data
}
