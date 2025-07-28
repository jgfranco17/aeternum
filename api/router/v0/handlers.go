package v0

import (
	"fmt"
	"net/http"

	exec "github.com/jgfranco17/aeternum/execution"

	"github.com/jgfranco17/aeternum/api/auth"
	"github.com/jgfranco17/aeternum/api/db"
	"github.com/jgfranco17/aeternum/api/httperror"
	"github.com/jgfranco17/aeternum/api/logging"

	"github.com/gin-gonic/gin"
)

func runTests(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		// Get user claims from context
		userClaims, exists := auth.GetUserClaims(c)
		if !exists {
			return fmt.Errorf("user claims not found in context")
		}

		var req exec.TestExecutionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return fmt.Errorf("Invalid request body: %w", err)
		}

		response, err := exec.ExecuteTests(c, req)
		if err != nil {
			return fmt.Errorf("Failed to execute tests: %w", err)
		}

		err = dbClient.StoreTestResult(c, userClaims.UserID, response)
		if err != nil {
			// Log the error but don't fail the request
			log := logging.FromContext(c)
			log.Errorf("Failed to store test result: %v", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}
}

func getTestResultsById(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		// Get user claims from context
		userClaims, exists := auth.GetUserClaims(c)
		if !exists {
			return fmt.Errorf("user claims not found in context")
		}

		log := logging.FromContext(c)
		resultId := c.Query("id")
		if resultId == "" {
			return httperror.New(c, http.StatusBadRequest, "Empty ID parameter")
		}

		result, err := dbClient.GetTestResult(c, userClaims.UserID, resultId)
		if err != nil {
			return fmt.Errorf("Failed to fetch test result: %w", err)
		}
		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No result found for ID %s", resultId),
			})
			return fmt.Errorf("No result found for ID %s", resultId)
		}
		log.Infof("Found results for ID %s", resultId)
		c.JSON(http.StatusOK, result)
		return nil
	}
}

// New handler to get all test results for a user
func getUserTestResults(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		// Get user claims from context
		userClaims, exists := auth.GetUserClaims(c)
		if !exists {
			return httperror.New(c, http.StatusBadRequest, "user claims not found in request context")
		}

		// Get limit from query parameter, default to 10
		limit := 10
		if limitStr := c.Query("limit"); limitStr != "" {
			if parsed, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil || parsed != 1 {
				return httperror.New(c, http.StatusBadRequest, "Invalid limit parameter")
			}
		}

		results, err := dbClient.GetUserTestResults(c, userClaims.UserID, limit)
		if err != nil {
			return fmt.Errorf("Failed to fetch user test results: %w", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"results": results,
			"count":   len(results),
		})
		return nil
	}
}
