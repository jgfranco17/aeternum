package core

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type CliCommandFunction func() *cobra.Command

type CommandRunner func(cmd *cobra.Command, args []string)

type CliRunResult struct {
	ShellOutput string
	Error       error
}

func createTempDir(t *testing.T) string {
	t.Helper()
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Return the directory path for use in tests
	return tempDir
}

// Helper function to simulate CLI execution
func ExecuteTestCommand(cmdGetter CliCommandFunction, args ...string) CliRunResult {
	buf := new(bytes.Buffer)
	cmd := cmdGetter()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	_, err := cmd.ExecuteC()
	return CliRunResult{
		ShellOutput: buf.String(),
		Error:       err,
	}
}

func TestPingCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	output := ExecuteTestCommand(GetPingCommand, server.URL)
	assert.NoError(t, output.Error, "Unexpected error while executing ping command")
}

func TestPingCommandServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	output := ExecuteTestCommand(GetPingCommand, server.URL, "--timeout", "1")
	assert.NoError(t, output.Error, "Unexpected error while executing ping command")
}

func TestPingCommandUnreachable(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	output := ExecuteTestCommand(GetPingCommand, server.URL, "--timeout", "1")
	assert.ErrorContains(t, output.Error, "Failed to reach target")
}
