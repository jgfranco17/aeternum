package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigCompleteConfig(t *testing.T) {
	// Setup environment variables
	t.Setenv("AETERNUM_LOG_LEVEL", "WARN")
	t.Setenv("AETERNUM_MONGO_USER", "my-user")
	t.Setenv("AETERNUM_MONGO_PASSWORD", "mongo-token")
	t.Setenv("AETERNUM_MONGO_URI", "some-uri")

	config, err := NewConfigFromSecrets()
	assert.NoError(t, err)
	assert.Equal(t, "WARN", config.LogLevel())
	assert.Equal(t, "my-user", config.MongoUser())
	assert.Equal(t, "mongo-token", config.MongoPassword())
	assert.Equal(t, "some-uri", config.MongoUri())
}
