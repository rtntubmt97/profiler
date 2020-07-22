package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	def "github.com/rtntubmt97/profiler/internal/defines"
	"github.com/rtntubmt97/profiler/internal/handlers"
	// "github.com/rtntubmt97/profiler/internal/utils"
)

func ginGetProfile(c *gin.Context) {
	idStr := c.Request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, def.ResponseData{Description: def.InvalidParam, Data: nil})
		// utils.PrintError(err)
		return
	}

	err, profile := handlers.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusOK, def.ResponseData{Description: def.HandleFailed, Data: nil})
		// utils.PrintError(err)
		return
	}

	c.JSON(http.StatusOK, def.ResponseData{Description: def.Success, Data: profile})
}
