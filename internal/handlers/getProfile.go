package handlers

import (
	"time"

	"github.com/rtntubmt97/profiler/internal/utils"
	app "github.com/rtntubmt97/profiler/pkg/applications"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

const pkgName = "handlers"

// var profiler k.Profiler = app.NewFileProfiler(1000, "/home/tumd/golang-repositories/profiler/test/out.txt")
var profiler k.Profiler = app.HttpPageProfiler()

func GetProfile(id int64) (error, interface{}) {
	startTime := time.Now().UnixNano()
	defer profiler.Record("GetProfile", startTime)
	err, profile := db.RetrieveProfile(id)
	return utils.WrapError(pkgName, "GetProfile failed", err), profile
}
