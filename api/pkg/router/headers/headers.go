package headers

import (
	core_errors "api/pkg/core/errors"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type OriginInfo struct {
	Origin  string `json:"origin"`
	Version string `json:"version"`
}

// Gets the origin info based on the received headers
func CreateOriginInfoHeader(c *gin.Context) (OriginInfo, error) {
	header := c.Request.Header["X-Origin-Info"]

	json_header := OriginInfo{}

	if len(header) == 0 {
		return json_header, core_errors.NewInputError(c, "X-Origin-Info header not found.")
	}

	err := json.Unmarshal([]byte(header[0]), &json_header)
	if err != nil {
		return json_header, core_errors.NewInputError(c, "Header schema validation: %s", err.Error())
	}
	return json_header, nil
}
