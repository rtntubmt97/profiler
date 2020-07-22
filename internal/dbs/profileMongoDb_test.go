package dbs

import (
	"testing"

	"github.com/rtntubmt97/profiler/internal/defines"
)

func TestRetrieveProfile(t *testing.T) {
	id := int64(1)
	err, profile := MongoDbInstance().RetrieveProfile(id)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
	expectProfile := defines.Profile{Id: 1, Name: "Tu", Job: "SE"}
	if expectProfile != profile {
		t.Error("data is not match")
		t.Errorf("profile: %+v\n", profile)
		t.Errorf("expectProfile: %+v\n", expectProfile)
	}
}
