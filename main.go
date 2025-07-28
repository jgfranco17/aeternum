package main

import (
	"flag"

	"github.com/jgfranco17/aeternum/api/db"
	env "github.com/jgfranco17/aeternum/api/environment"
	"github.com/jgfranco17/aeternum/api/router"
	"github.com/jgfranco17/aeternum/api/router/system"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	port    = flag.Int("port", 8080, "Port to listen on")
	devMode = flag.Bool("dev", true, "Run server in debug mode")
)

func init() {
	logrus.SetReportCaller(true)

	if env.IsLocalEnvironment() {
		logrus.SetFormatter(&logrus.TextFormatter{})
		gin.SetMode(gin.DebugMode)
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		gin.SetMode(gin.ReleaseMode)
	}
	prometheus.Register(system.HttpLastRequestReceivedTime)
}

func main() {
	flag.Parse()
	if *devMode {
		logrus.Infof("Running API server on port %d in dev mode", *port)
	} else {
		logrus.Infof("Running API production server on port %d", *port)
		gin.SetMode(gin.ReleaseMode)
	}
	dbClient, err := db.NewClient()
	if err != nil {
		logrus.Fatalf("Error initializing database client: %v", err)
	}
	service, err := router.CreateNewService(*port, dbClient)
	if err != nil {
		logrus.Fatalf("Error creating the server: %v", err)
	}
	err = service.Run()
	if err != nil {
		logrus.Fatalf("Error starting the server: %v", err)
	}
}
