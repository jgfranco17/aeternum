package system

import (
	"context"
	"time"

	"github.com/jgfranco17/aeternum/api/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func SetSystemRoutes(route *gin.Engine, includeSystemInfo bool) {
	log := logging.FromContext(context.Background())
	startTime = time.Now()
	if includeSystemInfo {
		specs, err := GetCodebaseSpecFromFile("specs.json")
		if err != nil {
			log.Fatal(err)
		}
		route.GET("/service-info", ServiceInfoHandler(specs, startTime))
	}
	for _, homeRoute := range []string{"", "/home"} {
		route.GET(homeRoute, HomeHandler)
	}
	route.GET("/healthz", HealthCheckHandler())
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	route.NoRoute(NotFoundHandler)
}
