package system

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jgfranco17/aeternum/api/pkg/environment"
	"github.com/jgfranco17/aeternum/api/pkg/logging"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Aeternum API!",
	})
}

func ServiceInfoHandler(codebaseSpec *ProjectCodebase, startTime time.Time) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeSinceStart := time.Since(startTime)
		uptimeSeconds := fmt.Sprintf("%ds", int(timeSinceStart.Seconds()))
		c.JSON(http.StatusOK, ServiceInfo{
			Name:        "Aeternum API",
			Codebase:    *codebaseSpec,
			Environment: environment.GetApplicationEnv(),
			Uptime:      uptimeSeconds,
		})
	}
}

func HealthCheckHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthStatus{
			Timestamp: time.Now().Format(time.RFC822),
			Status:    "healthy",
		})
	}
}

func NotFoundHandler(c *gin.Context) {
	log := logging.FromContext(c)
	log.Errorf("Non-existent endpoint accessed: %s", c.Request.URL.Path)
	c.JSON(http.StatusNotFound, newMissingEndpoint(c.Request.URL.Path))
}

func newMissingEndpoint(endpoint string) BasicErrorInfo {
	return BasicErrorInfo{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("Endpoint '%s' does not exist", endpoint),
	}
}
