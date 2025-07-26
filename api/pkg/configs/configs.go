package configs

import (
	"context"
	"fmt"
	"strings"

	"github.com/jgfranco17/aeternum/api/pkg/environment"
	"github.com/jgfranco17/aeternum/api/pkg/logging"
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
	missingEnvVars := []string{}
	username := environment.GetEnvWithDefault(EnvVarMongoUser, "")
	if username == "" {
		missingEnvVars = append(missingEnvVars, EnvVarMongoUser)
	}
	password := environment.GetEnvWithDefault(EnvVarMongoPassword, "")
	if password == "" {
		missingEnvVars = append(missingEnvVars, EnvVarMongoPassword)
	}
	uri := environment.GetEnvWithDefault(EnvVarMongoUri, "")
	if uri == "" {
		missingEnvVars = append(missingEnvVars, EnvVarMongoUri)
	}
	if len(missingEnvVars) > 0 {
		return "", "", "", fmt.Errorf("Missing %d environment variables: %s", len(missingEnvVars), strings.Join(missingEnvVars, ", "))
	}
	return username, password, uri, nil
}

func NewConfigFromSecrets() (*EnvironmentConfig, error) {
	log := logging.FromContext(context.Background())
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
	log.Infof("Loaded configs from environment")
	return &config, nil
}
