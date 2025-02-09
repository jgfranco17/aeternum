package configs

import (
	"context"
	"fmt"
	"os"
	"path"

	env "api/pkg/core/environment"
	logger "api/pkg/core/obs"

	"gopkg.in/yaml.v3"
)

const (
	EnvVarGithubToken string = "AETERNUM_GITHUB_TOKEN"
	ConfigFileName    string = "config.yaml"
)

type GithubConfig interface {
	GithubToken() string
	GithubBaseUrl() string
}

type EnvironmentConfig struct {
	EnvGithubBaseUrl string `yaml:"AETERNUM_GITHUB_URL"`
	EnvGithubToken   string `yaml:"AETERNUM_GITHUB_TOKEN"`
	EnvLogLevel      string `yaml:"AETERNUM_LOG_LEVEL"`
}

func (c *EnvironmentConfig) GithubBaseUrl() string {
	return c.EnvGithubBaseUrl
}

func (c *EnvironmentConfig) GithubToken() string {
	return c.EnvGithubToken
}

func (c *EnvironmentConfig) LogLevel() string {
	return c.EnvLogLevel
}

func loadFromFile(configPath string, config *EnvironmentConfig) error {
	log := logger.GetLoggerFromContext(context.Background())
	log.Infof("Loading configuration from %s", configPath)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file at %s: %w", configPath, err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config file: %w", err)
	}

	log.Info("Configuration was loaded successfully.")

	return nil
}

func loadSecrets(config *EnvironmentConfig) error {
	log := logger.GetLoggerFromContext(context.Background())
	log.Infof("Loading secrets from env")
	githubToken := env.GetEnvWithDefault(EnvVarGithubToken, "")
	if githubToken == "" {
		return fmt.Errorf("Github token was not set")
	}
	config.EnvGithubToken = githubToken
	log.Info("Configuration was loaded successfully.")
	return nil
}

// Load the application configuration
func LoadConfig(configDir string) (*EnvironmentConfig, error) {
	log := logger.GetLoggerFromContext(context.Background())
	config := EnvironmentConfig{}
	configFile := path.Join(configDir, ConfigFileName)
	log.Infof("Loading configuration from %s", configFile)
	err := loadFromFile(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to load configs from file: %w", err)
	}
	err = loadSecrets(&config)
	if err != nil {
		return nil, fmt.Errorf("Failed to load secrets from environment: %w", err)
	}
	return &config, err
}
