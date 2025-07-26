package headers

import (
	"encoding/json"

	"github.com/jgfranco17/aeternum/api/pkg/httperror"

	"github.com/gin-gonic/gin"
)

type OriginInfo struct {
	Origin  string `json:"origin"`
	Version string `json:"version"`
}

// Gets the origin info based on the received headers
func CreateOriginInfoHeader(c *gin.Context) (OriginInfo, error) {
	header := c.Request.Header["X-Origin-Info"]

	jsonHeader := OriginInfo{}

	if len(header) == 0 {
		return jsonHeader, httperror.NewInputError(c, "X-Origin-Info header not found.")
	}

	err := json.Unmarshal([]byte(header[0]), &jsonHeader)
	if err != nil {
		return jsonHeader, httperror.NewInputError(c, "Header schema validation: %s", err.Error())
	}
	return jsonHeader, nil
}
