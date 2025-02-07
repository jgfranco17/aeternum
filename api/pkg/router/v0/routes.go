package v0

import (
	"github.com/gin-gonic/gin"
)

// Adds v0 routes to the router.
func SetRoutes(route *gin.Engine) {
	v0 := route.Group("/v0")
	{
		testExecutionRoutes := v0.Group("/tests")
		{
			testExecutionRoutes.POST("/run", WithErrorHandling(runTests()))
			testExecutionRoutes.GET("/results", WithErrorHandling(getTestResultsById()))
		}
	}
}
