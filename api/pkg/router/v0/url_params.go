package v0

import (
	"context"
	"encoding/base64"

	"github.com/jgfranco17/aeternum/api/pkg/httperror"

	"github.com/gin-gonic/gin"
)

// Decode a url param formatted in base64.
func decodeBase64Param(ctx context.Context, base64Param string) (string, error) {
	if base64Param == "" {
		return "", httperror.NewInputError(ctx, "the base 64 param cannot be empty")
	}
	value, err := base64.RawURLEncoding.DecodeString(base64Param)
	if err != nil {
		value, err = base64.StdEncoding.DecodeString(base64Param)
		if err != nil {
			return "", httperror.NewInputError(ctx, "Cannot decode the param with base64. Original value: %s. Error: %w", base64Param, err)
		}

	}
	return string(value), nil
}

func getBase64Param(c *gin.Context, key string) (string, error) {
	b64Uri := c.Param(key)
	return decodeBase64Param(c, b64Uri)
}
