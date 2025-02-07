package core

import (
	"os"
	"path"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func uniqueDir(t *testing.T) string {
	dir := t.TempDir()
	id, _ := uuid.NewV4()
	fullDir := path.Join(dir, id.String())
	err := os.Mkdir(fullDir, 0777)

	assert.NoError(t, err)
	return fullDir
}

func TestLoadConfigCompleteConfig(t *testing.T) {
	// Setup environment variables
	dir := uniqueDir(t)
	configFile := path.Join(dir, "config.yaml")
	t.Setenv("AETERNUM_GITHUB_TOKEN", "abcdefg4321")
	configFileContents := `
AETERNUM_GITHUB_URL: https://github.com
AETERNUM_LOG_LEVEL: WARN`
	err := os.WriteFile(configFile, []byte(configFileContents), 0666)
	assert.NoError(t, err)

	config, err := LoadConfig(dir)
	assert.NoError(t, err)
	assert.Equal(t, "abcdefg4321", config.GithubToken())
	assert.Equal(t, "https://github.com", config.GithubBaseUrl())
}

func TestLoadConfigFromFiles(t *testing.T) {
	dir := uniqueDir(t)
	t.Setenv("AETERNUM_GITHUB_TOKEN", "abcdefg4321")
	configFile := path.Join(dir, "config.yaml")
	configFileContents := `AETERNUM_GITHUB_URL: https://github.com
AETERNUM_LOG_LEVEL: WARN`
	err := os.WriteFile(configFile, []byte(configFileContents), 0666)
	assert.NoError(t, err)

	config, err := LoadConfig(dir)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com", config.GithubBaseUrl())
	assert.Equal(t, "WARN", config.LogLevel())
}

func TestLoadConfigFileNotExist(t *testing.T) {
	emptyDir := uniqueDir(t)
	_, err := LoadConfig(emptyDir)

	assert.ErrorContains(t, err, "Failed to read config file")
}

func TestLoadConfigSecretsNotSet(t *testing.T) {
	dir := uniqueDir(t)
	configFile := path.Join(dir, "config.yaml")
	configFileContents := `AETERNUM_GITHUB_URL: https://github.com
AETERNUM_LOG_LEVEL: WARN`
	err := os.WriteFile(configFile, []byte(configFileContents), 0666)
	assert.NoError(t, err)

	_, err = LoadConfig(dir)
	assert.ErrorContains(t, err, "Failed to load secrets from environment")
}
