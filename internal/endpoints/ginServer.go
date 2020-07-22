package endpoints

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ginServer *gin.Engine

func init() {
	// ginServer = gin.Default()
	ginServer = gin.New()
	ginServer.GET("/getProfile", ginGetProfile)
	ginServer.POST("/insertProfile", ginInsertProfile)
	ginServer.POST("/printData", printDataHandler)
}

func StartGinServer() {
	// ginServer.GET("/getProfile", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "hi")
	// })
	ginServer.Run(":9080")
}

func printDataHandler(c *gin.Context) {
	fmt.Println(c.Request.Body)
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(c.Request.FormValue("f"))
	fmt.Println(body)
	fmt.Println(string(body))
	fmt.Println(c.Request.Header.Get("Host"))
	fmt.Println(c.Request.Header.Get("User-Agent"))
	c.String(http.StatusOK, "done")
}
