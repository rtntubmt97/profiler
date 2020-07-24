package intervalListeners

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type Log struct {
	profiles map[string]*k.Profile
	writer   io.Writer
}

func (log *Log) Listen(profiles map[string]*k.Profile, startTime time.Time, intervalTimeMillis int) {
	// currentTime := time.Now().Format("2006-01-02 15:04:05.000000")
	// fmt.Fprintf(log.writer, "------%s------\n", currentTime)

	sampleProfile, ok := profiles["LastIntervalUpdate"]
	if ok {
		logTime := sampleProfile.History[len(sampleProfile.History)-1].UnixNano
		fmt.Fprintf(log.writer, "---logTime: %s---\n", time.Unix(0, logTime).Format("2006-01-02 15:04:05.000000"))
	} else {
		currentTime := time.Now().UnixNano()
		fmt.Fprintf(log.writer, "---printTime: %s---\n", time.Unix(0, currentTime).Format("2006-01-02 15:04:05.000000"))
	}

	names := make([]string, 0, len(profiles))
	for k := range profiles {
		// if k == "LastIntervalUpdate" {
		// 	continue
		// }
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		profile := profiles[name]
		lastIndex := len(profile.History) - 1
		fmt.Fprintf(log.writer, "[%s] ReqCount: %d, AvgProcTime: %d\n",
			name, profile.History[lastIndex].RequestCounts, profile.History[lastIndex].AvgProcTimes)
	}
}

func NewConsolLog() *Log {
	ret := new(Log)
	ret.writer = os.Stdout
	return ret
}

func NewFileLog(filePath string) (*Log, *os.File) {
	ret := new(Log)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	ret.writer = file

	return ret, file
}
