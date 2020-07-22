package handlers

import (
	"github.com/rtntubmt97/profiler/internal/dbs"
	"github.com/rtntubmt97/profiler/internal/defines"
)

var db defines.ProfileDb

func init() {
	db = dbs.MongoDbInstance()
}
