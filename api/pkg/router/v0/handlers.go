package v0

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "api/pkg/core/errors"
	exec "api/pkg/core/execution"

	"github.com/gin-gonic/gin"
)

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

func getTestResultsById() func(c *gin.Context) error {
	return func(c *gin.Context) error {
		value := c.Query("id")
		number, err := strconv.Atoi(value)
		if err != nil {
			return core_errors.NewInputError(c, "Failed to parse ID: %w", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Test results for ID %v", number),
			"results": []string{"result1", "result2"},
		})
		return nil
	}
}
