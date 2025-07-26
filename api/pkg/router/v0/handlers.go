package v0

import (
	"context"
	"fmt"
	"net/http"

	exec "github.com/jgfranco17/aeternum/execution"

	"github.com/jgfranco17/aeternum/api/pkg/db"
	"github.com/jgfranco17/aeternum/api/pkg/httperror"
	"github.com/jgfranco17/aeternum/api/pkg/logging"

	"github.com/gin-gonic/gin"
)

type dbClient interface {
	Disconnect(ctx context.Context) error
	GetResult(ctx context.Context, id string) error
}

func runTests() func(c *gin.Context) error {
	return func(c *gin.Context) error {
		var req exec.TestExecutionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return fmt.Errorf("Invalid request body: %w", err)
		}

		response, err := exec.ExecuteTests(c, req)
		if err != nil {
			return fmt.Errorf("Failed to execute tests: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}
}

func getTestResultsById(username string, token string, uri string) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		log := logging.FromContext(c)
		resultId := c.Query("id")
		if resultId == "" {
			return httperror.NewInputError(c, "Empty ID parameter")
		}
		client, err := db.NewMongoClient(c, uri, username, token)
		defer client.Disconnect(c)
		if resultId == "" {
			return fmt.Errorf("Failed to create database client: %w", err)
		}

		result, err := client.GetResult(c, resultId)
		if err != nil {
			return fmt.Errorf("Failed to fetch data from collection: %w", err)
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
