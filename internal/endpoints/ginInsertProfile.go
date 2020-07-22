package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	def "github.com/rtntubmt97/profiler/internal/defines"
	"github.com/rtntubmt97/profiler/internal/handlers"
	"github.com/rtntubmt97/profiler/internal/utils"
)

const pkgName = "endpoints"

func ginInsertProfile(c *gin.Context) {
	profile, err := extractProfileData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, def.ResponseData{Description: def.InvalidParam, Data: nil})
		err = utils.WrapError(pkgName, "ginInsertProfile failed", err)
		utils.PrintError(err)
		return
	}

	err = handlers.InsertProfile(profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.ResponseData{Description: def.HandleFailed, Data: nil})
		err = utils.WrapError(pkgName, "ginInsertProfile failed", err)
		utils.PrintError(err)
		return
	}

	fmt.Println(profile)
	c.JSON(http.StatusOK, def.ResponseData{Description: def.Success, Data: nil})
}

func extractProfileData(c *gin.Context) (def.Profile, error) {
	var profile def.Profile

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return profile, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}
