package system

import (
	"context"
	"time"

	"github.com/jgfranco17/aeternum/api/logging"
	v0 "github.com/jgfranco17/aeternum/api/router/v0"

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
	route.POST("/register", v0.WithErrorHandling(RegisterHandler()))
	route.POST("/login", v0.WithErrorHandling(LoginHandler()))
	route.GET("/healthz", HealthCheckHandler())
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	route.NoRoute(NotFoundHandler)
}
