package v0

import (
	"github.com/jgfranco17/aeternum/api/auth"

	"github.com/gin-gonic/gin"
)

// Adds v0 routes to the router.
func SetRoutes(route *gin.Engine) error {
	v0 := route.Group("/v0")
	// Apply authentication middleware to all v0 routes
	v0.Use(auth.AuthMiddleware())
	{
		testExecutionRoutes := v0.Group("/tests")
		{
			testExecutionRoutes.POST("/run", WithErrorHandling(runTests()))
			testExecutionRoutes.GET("/results", WithErrorHandling(getTestResultsById("", "", "")))
			testExecutionRoutes.GET("/history", WithErrorHandling(getUserTestResults()))
		}
	}
	return nil
}
