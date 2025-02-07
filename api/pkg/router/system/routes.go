package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func SetSystemRoutes(route *gin.Engine) {
	route.GET("/", HomeHandler)
	route.GET("/service-info", ServiceInfoHandler)
	route.GET("/healthz", HealthCheckHandler)
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	route.NoRoute(NotFoundHandler)
}
