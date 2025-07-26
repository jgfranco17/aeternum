package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserClaimsKey = "user_claims"
)

// AuthMiddleware validates JWT tokens and adds user claims to context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set(UserClaimsKey, claims)
		c.Next()
	}
}

// GetUserClaims retrieves user claims from the gin context
func GetUserClaims(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get(UserClaimsKey)
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*Claims)
	return userClaims, ok
}
