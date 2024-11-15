package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Api interface {
	HelloWorldHandler(c *gin.Context)
}

func HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}
