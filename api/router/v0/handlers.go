package v0

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/jgfranco17/aeternum/api/auth"
	"github.com/jgfranco17/aeternum/api/db"
	"github.com/jgfranco17/aeternum/api/httperror"
	"github.com/jgfranco17/aeternum/api/logging"
	"github.com/jgfranco17/aeternum/execution"

	"github.com/gin-gonic/gin"
)

func runTests(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		log := logging.FromContext(c)
		userClaims, exists := auth.GetUserClaims(c)
		if !exists {
			return httperror.New(c, http.StatusUnauthorized, "user claims not found in context")
		}

		var req execution.TargetDefinition
		if err := c.ShouldBindJSON(&req); err != nil {
			return httperror.New(c, http.StatusBadRequest, "invalid request body: %w", err)
		}

		log.Debugf("Running %d tests for %s", len(req.Endpoints), req.BaseURL)
		response, err := execution.Run(c, req)
		if err != nil {
			return fmt.Errorf("failed to execute tests: %w", err)
		}
		err = dbClient.StoreTestResult(c, userClaims.UserID, response)
		if err != nil {
			// Log the error but don't fail the request
			log.Errorf("Failed to store test result: %v", err)
		}

		c.JSON(http.StatusOK, response)

		log.WithFields(logrus.Fields{
			"id":        userClaims.UserID,
			"target":    req.BaseURL,
			"endpoints": len(req.Endpoints),
			"timeout":   req.MaxTimeout,
		}).Info("Test execution completed")
		return nil
	}
}

func getTestResultsById(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		// Get user claims from context
		userClaims, exists := auth.GetUserClaims(c)
		if !exists {
			return httperror.New(c, http.StatusUnauthorized, "user claims not found in context")
		}

		log := logging.FromContext(c)
		resultID := c.Query("id")
		if resultID == "" {
			return httperror.New(c, http.StatusBadRequest, "Empty ID parameter")
		}

		result, err := dbClient.GetTestResult(c, userClaims.UserID, resultID)
		if err != nil {
			return fmt.Errorf("Failed to fetch test result: %w", err)
		}
		if result == nil {
			return httperror.New(c, http.StatusNotFound, "No result found for ID %s", resultID)
		}
		log.WithFields(logrus.Fields{
			"id": resultID,
		}).Infof("Found requested results for ID")
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
			return httperror.New(c, http.StatusUnauthorized, "user claims not found in request context")
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
