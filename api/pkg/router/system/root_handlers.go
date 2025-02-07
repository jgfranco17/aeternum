package system

import (
	"fmt"
	"net/http"
	"time"

	"api/pkg/core"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Aeternum API page!",
	})
}

func ServiceInfoHandler(c *gin.Context) {
	// Send the parsed JSON core as a response
	c.JSON(http.StatusOK, core.AboutInfo{
		Name:        "Aeternum API",
		Author:      "Joaquin Franco",
		Repository:  "https://github.com/jgfranco17/aeternum-api",
		Environment: core.GetApplicationEnv(),
		License:     "MIT",
		Languages:   []string{"Go"},
	})
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, core.HealthStatus{
		Timestamp: time.Now().Format(time.RFC822),
		Status:    "healthy",
	})
}

func NotFoundHandler(c *gin.Context) {
	log := core.FromContext(c)
	log.Errorf("Non-existent endpoint accessed: %s", c.Request.URL.Path)
	c.JSON(http.StatusNotFound, newMissingEndpoint(c.Request.URL.Path))
}

func newMissingEndpoint(endpoint string) core.BasicErrorInfo {
	return core.BasicErrorInfo{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("Endpoint '%s' does not exist", endpoint),
	}
}
