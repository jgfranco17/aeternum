package system

import (
	"fmt"
	"net/http"
	"time"

	"api/pkg/core/environment"
	"api/pkg/core/obs"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Aeternum API!",
	})
}

func ServiceInfoHandler(startTime time.Time) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, ServiceInfo{
			Name:        "Aeternum API",
			Author:      "Joaquin Gabriel Franco",
			Repository:  "https://github.com/jgfranco17/aeternum-api",
			Environment: environment.GetApplicationEnv(),
			Uptime:      time.Since(startTime),
			License:     "MIT",
			Languages:   []string{"Go"},
		})
	}
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthStatus{
		Timestamp: time.Now().Format(time.RFC822),
		Status:    "healthy",
	})
}

func NotFoundHandler(c *gin.Context) {
	log := obs.GetLoggerFromContext(c)
	log.Errorf("Non-existent endpoint accessed: %s", c.Request.URL.Path)
	c.JSON(http.StatusNotFound, newMissingEndpoint(c.Request.URL.Path))
}

func newMissingEndpoint(endpoint string) BasicErrorInfo {
	return BasicErrorInfo{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("Endpoint '%s' does not exist", endpoint),
	}
}
