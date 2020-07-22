package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Temp struct {
	Foo int    `json:"foo"`
	Bar string `json:"bar"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	data := Temp{Foo: 2, Bar: "bar"}
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, data)
	})

	r.Run(":9080")
}
