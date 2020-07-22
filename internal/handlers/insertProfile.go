package handlers

import (
	"github.com/rtntubmt97/profiler/internal/defines"
	"github.com/rtntubmt97/profiler/internal/utils"
)

func InsertProfile(profile defines.Profile) error {
	err := db.CreateProfile(profile)
	return utils.WrapError(pkgName, "InsertProfile failed", err)
}
