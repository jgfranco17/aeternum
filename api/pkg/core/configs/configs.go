package configs

import (
	"context"
	"fmt"

	"api/pkg/core/environment"
	logger "api/pkg/core/obs"
)

const (
	EnvVarLogLevel      string = "AETERNUM_LOG_LEVEL"
	EnvVarMongoUser     string = "AETERNUM_MONGO_USER"
	EnvVarMongoPassword string = "AETERNUM_MONGO_PASSWORD"
	EnvVarMongoUri      string = "AETERNUM_MONGO_URI"
)

type GithubConfig interface {
	GithubToken() string
	GithubBaseUrl() string
}

type EnvironmentConfig struct {
	EnvLogLevel      string
	EnvMongoUser     string
	EnvMongoPassword string
	EnvMongoUri      string
}

func (c *EnvironmentConfig) LogLevel() string {
	return c.EnvLogLevel
}

func (c *EnvironmentConfig) MongoUser() string {
	return c.EnvMongoUser
}

func (c *EnvironmentConfig) MongoPassword() string {
	return c.EnvMongoPassword
}

func (c *EnvironmentConfig) MongoUri() string {
	return c.EnvMongoUri
}

func getMongoConfigs() (string, string, string, error) {
	username := environment.GetEnvWithDefault(EnvVarMongoUser, "")
	if username == "" {
		return "", "", "", fmt.Errorf("Environment variable %s not set", EnvVarMongoUser)
	}
	password := environment.GetEnvWithDefault(EnvVarMongoPassword, "")
	if username == "" {
		return "", "", "", fmt.Errorf("Environment variable %s not set", EnvVarMongoPassword)
	}
	uri := environment.GetEnvWithDefault(EnvVarMongoUri, "")
	if username == "" {
		return "", "", "", fmt.Errorf("Environment variable %s not set", EnvVarMongoUri)
	}
	return username, password, uri, nil
}

func NewConfigFromSecrets() (*EnvironmentConfig, error) {
	log := logger.GetLoggerFromContext(context.Background())
	log.Infof("Loading secrets from env")
	mongoUser, mongoPassword, mongoUri, err := getMongoConfigs()
	if err != nil {
		return nil, fmt.Errorf("Failed to load configs: %w", err)
	}
	config := EnvironmentConfig{
		EnvLogLevel:      environment.GetEnvWithDefault(EnvVarLogLevel, "DEBUG"),
		EnvMongoUser:     mongoUser,
		EnvMongoPassword: mongoPassword,
		EnvMongoUri:      mongoUri,
	}
	return &config, nil
}
