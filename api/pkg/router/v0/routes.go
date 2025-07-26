package v0

import (
	"github.com/jgfranco17/aeternum/api/pkg/configs"

	"github.com/gin-gonic/gin"
)

// Adds v0 routes to the router.
func SetRoutes(route *gin.Engine) error {
	appConfig, err := configs.NewConfigFromSecrets()
	if err != nil {
		return err
	}
	v0 := route.Group("/v0")
	{
		testExecutionRoutes := v0.Group("/tests")
		{
			testExecutionRoutes.POST("/run", WithErrorHandling(runTests()))
			testExecutionRoutes.GET("/results", WithErrorHandling(getTestResultsById(appConfig.MongoUser(), appConfig.MongoPassword(), appConfig.MongoUri())))
		}
	}
	return nil
}
