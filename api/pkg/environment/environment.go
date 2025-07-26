package environment

import (
	"os"
)

const (
	APPLICATION_ENV_LOCAL = "local"
	APPLICATION_ENV_DEV   = "dev"
	APPLICATION_ENV_STAGE = "stage"
	APPLICATION_ENV_PROD  = "prod"
)

const (
	ENV_KEY_ENVIRONMENT = "ENVIRONMENT"
	ENV_KEY_VERSION     = "APP_VERSION"
)

func IsLocalEnvironment() bool {
	return GetApplicationEnv() == APPLICATION_ENV_LOCAL
}

func GetEnvWithDefault(key string, defaultValue string) string {
	value, present := os.LookupEnv(key)
	if present {
		return value
	}
	return defaultValue
}

func GetApplicationEnv() string {
	return GetEnvWithDefault(ENV_KEY_ENVIRONMENT, APPLICATION_ENV_LOCAL)
}
