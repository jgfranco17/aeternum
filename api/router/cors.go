package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	API_URL_LOCAL = "http://localhost:3000"
	API_URL_DEV   = "https://dev.aeternum-ci.com"
	API_URL_STAGE = "https://stage.aeternum-ci.com"
	API_URL_BETA  = "https://beta.aeternum-ci.com"
	API_URL_PROD  = "https://www.aeternum-ci.com"
)

func GetCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
	})
}
