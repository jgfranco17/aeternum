package execution

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTestsSuccess(t *testing.T) {
	mux := http.NewServeMux()
	for _, endpoint := range []string{"/home", "/healthz"} {
		mux.HandleFunc(endpoint, func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	}
	mockServer := httptest.NewServer(mux)
	request := TestExecutionRequest{
		BaseURL: mockServer.URL,
		Endpoints: []Endpoint{
			{
				Path:           "/home",
				ExpectedStatus: http.StatusOK,
			},
			{
				Path:           "/healthz",
				ExpectedStatus: http.StatusOK,
			},
		},
		MaxTimeoutSeconds: nil,
	}
	results, err := ExecuteTests(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, "PASS", results.Status)
}

func TestExecuteTestsStatusCodeDoNotMatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/index", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/will-fail", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})
	mockServer := httptest.NewServer(mux)
	request := TestExecutionRequest{
		BaseURL: mockServer.URL,
		Endpoints: []Endpoint{
			{
				Path:           "/index",
				ExpectedStatus: http.StatusOK,
			},
			{
				Path:           "/will-fail",
				ExpectedStatus: http.StatusOK,
			},
		},
		MaxTimeoutSeconds: nil,
	}
	results, err := ExecuteTests(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, "FAIL", results.Status)
}

func TestExecuteTestsFailedRequest(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewUnstartedServer(mux)
	request := TestExecutionRequest{
		BaseURL: mockServer.URL,
		Endpoints: []Endpoint{
			{
				Path:           "/home",
				ExpectedStatus: http.StatusOK,
			},
		},
		MaxTimeoutSeconds: nil,
	}
	_, err := ExecuteTests(context.Background(), request)
	assert.Errorf(t, err, "Failed to make 1 requests")
}
