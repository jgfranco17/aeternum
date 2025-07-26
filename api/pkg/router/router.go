package router

import (
	"fmt"
	"os"

	env "github.com/jgfranco17/aeternum/api/pkg/environment"
	"github.com/jgfranco17/aeternum/api/pkg/logging"
	"github.com/jgfranco17/aeternum/api/pkg/router/headers"
	system "github.com/jgfranco17/aeternum/api/pkg/router/system"
	v0 "github.com/jgfranco17/aeternum/api/pkg/router/v0"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	Router *gin.Engine
	Port   int
}

func (s *Service) Run() error {
	err := s.Router.Run(fmt.Sprintf(":%v", s.Port))
	if err != nil {
		return fmt.Errorf("Failed to start service on port %v: %w", s.Port, err)
	}
	return nil
}

// Add the fields we want to expose in the logger to the request context
func addLoggerFields() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !env.IsLocalEnvironment() {
			requestID := uuid.NewString()
			environment := os.Getenv(env.ENV_KEY_ENVIRONMENT)
			version := os.Getenv(env.ENV_KEY_VERSION)

			// Golang recommends contexts to use custom types instead
			// of strings, but gin defines key as a string.
			c.Set(string(logging.RequestId), requestID)
			c.Set(string(logging.Environment), environment)
			c.Set(string(logging.Version), version)

			originInfo, err := headers.CreateOriginInfoHeader(c)

			if err == nil && originInfo.Origin != "" {
				c.Set(string(logging.Origin), fmt.Sprintf("%s@%s", originInfo.Origin, originInfo.Version))
			}
		}
		c.Next()
	}
}

// Log the start and completion of a request
func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logging.FromContext(c)

		origin := c.Request.Header.Get("Origin")
		log.Infof("Request Started: [%s] %s from %s", c.Request.Method, c.Request.URL, origin)
		c.Next()
		log.Infof("Request Completed: [%s] %s", c.Request.Method, c.Request.URL)
	}
}

// Configure the router adding routes and middlewares
func getRouter(withSystemInfo bool) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(addLoggerFields())
	router.Use(logRequest())
	router.Use(GetCors())
	router.Use(system.PrometheusMiddleware())
	system.SetSystemRoutes(router, withSystemInfo)
	err := v0.SetRoutes(router)
	if err != nil {
		return nil, fmt.Errorf("Failed to set v0 routes: %w", err)
	}

	return router, nil
}

/*
Create a backend service instance.

[IN] port: server port to listen on

[OUT] *Service: new backend service instance
*/
func CreateNewService(port int) (*Service, error) {
	router, err := getRouter(true)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new service instance: %w", err)
	}
	return &Service{
		Router: router,
		Port:   port,
	}, nil
}
