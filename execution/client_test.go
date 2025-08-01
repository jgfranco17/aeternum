package execution

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTests_Success(t *testing.T) {
	mux := http.NewServeMux()
	for _, endpoint := range []string{"/home", "/healthz"} {
		mux.HandleFunc(endpoint, func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	}
	mockServer := httptest.NewServer(mux)
	request := TargetDefinition{
		BaseURL: mockServer.URL,
		Endpoints: []Endpoint{
			{
				Path:           "/home",
				Method:         "GET",
				ExpectedStatus: http.StatusOK,
			},
			{
				Path:           "/healthz",
				Method:         "GET",
				ExpectedStatus: http.StatusOK,
			},
		},
		MaxTimeoutSeconds: nil,
	}
	results, err := Run(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, StatusPass, results.Status)
}

func TestExecuteTestsFail_StatusCodeDoNotMatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/index", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/will-fail", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})
	mockServer := httptest.NewServer(mux)
	request := TargetDefinition{
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
	results, err := Run(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, StatusFail, results.Status)
}

func TestExecuteTestsFail_UnreachedRequest(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewUnstartedServer(mux)
	request := TargetDefinition{
		BaseURL: mockServer.URL,
		Endpoints: []Endpoint{
			{
				Path:           "/home",
				ExpectedStatus: http.StatusOK,
			},
		},
		MaxTimeoutSeconds: nil,
	}
	_, err := Run(context.Background(), request)
	assert.Errorf(t, err, "Failed to make 1 requests")
}
