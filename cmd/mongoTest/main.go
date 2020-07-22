package main

import (
	"fmt"

	"github.com/rtntubmt97/profiler/internal/dbs"
)

func main() {
	err, profile := dbs.MongoDbInstance().RetrieveProfile(1)

	fmt.Println(err)
	fmt.Println(profile)
}
