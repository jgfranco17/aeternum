package system

import (
	"context"
	"time"

	"api/pkg/core/obs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func SetSystemRoutes(route *gin.Engine) {
	log := obs.GetLoggerFromContext(context.Background())
	startTime = time.Now()
	specs, err := GetCodebaseSpecFromFile("specs.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, homeRoute := range []string{"", "/home"} {
		route.GET(homeRoute, HomeHandler)
	}
	route.GET("/service-info", ServiceInfoHandler(specs, startTime))
	route.GET("/healthz", HealthCheckHandler())
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	route.NoRoute(NotFoundHandler)
}
