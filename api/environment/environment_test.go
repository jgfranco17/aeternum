package environment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentRetrieval(t *testing.T) {
	environments := []string{
		"local",
		"dev",
		"stage",
		"prod",
		"ci",
	}
	for _, sampleEnv := range environments {
		t.Run(fmt.Sprintf("Check environment [%s]", sampleEnv), func(t *testing.T) {
			t.Setenv(ENV_KEY_ENVIRONMENT, sampleEnv)
			valueFromEnv := GetEnvWithDefault(ENV_KEY_ENVIRONMENT, "")
			assert.Equal(t, sampleEnv, valueFromEnv)
		})
	}
}

func TestCheckLocalEnvironment(t *testing.T) {
	t.Setenv(ENV_KEY_ENVIRONMENT, APPLICATION_ENV_LOCAL)
	assert.True(t, IsLocalEnvironment())
}
