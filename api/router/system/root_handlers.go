package system

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jgfranco17/aeternum/api/db"
	"github.com/jgfranco17/aeternum/api/environment"
	"github.com/jgfranco17/aeternum/api/logging"
	"github.com/supabase-community/gotrue-go/types"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterHandler handles user registration using Supabase Auth
func RegisterHandler() func(c *gin.Context) error {
	return func(c *gin.Context) error {
		var regReq RegisterRequest
		if err := c.ShouldBindJSON(&regReq); err != nil {
			return fmt.Errorf("Invalid request body: %w", err)
		}

		client := db.GetSupabaseClient()
		signupReq := types.SignupRequest{
			Email:    regReq.Email,
			Password: regReq.Password,
		}
		user, err := client.Auth.Signup(signupReq)
		if err != nil {
			return fmt.Errorf("Registration failed: %w", err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":    user.ID,
			"email": user.Email,
		})
		return nil
	}
}

// LoginHandler handles user authentication and returns JWT token
func LoginHandler() func(c *gin.Context) error {
	return func(c *gin.Context) error {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			return fmt.Errorf("Invalid request body: %w", err)
		}

		client := db.GetSupabaseClient()
		tokenReq := types.TokenRequest{
			Email:     loginReq.Email,
			Password:  loginReq.Password,
			GrantType: "password",
		}
		tokenResp, err := client.Auth.Token(tokenReq)
		if err != nil {
			return fmt.Errorf("Invalid credentials: %w", err)
		}

		token := tokenResp.AccessToken
		if token == "" {
			return fmt.Errorf("No access token returned from Supabase")
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
			Type:  "Bearer",
		})
		return nil
	}
}

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Aeternum API!",
	})
}

func ServiceInfoHandler(codebaseSpec *ProjectCodebase, startTime time.Time) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeSinceStart := time.Since(startTime)
		uptimeSeconds := fmt.Sprintf("%ds", int(timeSinceStart.Seconds()))
		c.JSON(http.StatusOK, ServiceInfo{
			Name:        "Aeternum API",
			Codebase:    *codebaseSpec,
			Environment: environment.GetApplicationEnv(),
			Uptime:      uptimeSeconds,
		})
	}
}

func HealthCheckHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthStatus{
			Timestamp: time.Now().Format(time.RFC822),
			Status:    "healthy",
		})
	}
}

func NotFoundHandler(c *gin.Context) {
	log := logging.FromContext(c)
	log.Errorf("Non-existent endpoint accessed: %s", c.Request.URL.Path)
	c.JSON(http.StatusNotFound, newMissingEndpoint(c.Request.URL.Path))
}

func newMissingEndpoint(endpoint string) BasicErrorInfo {
	return BasicErrorInfo{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("Endpoint '%s' does not exist", endpoint),
	}
}
