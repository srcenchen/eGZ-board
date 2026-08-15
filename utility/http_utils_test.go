package utility

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write([]byte("weather"))
	}))
	defer server.Close()

	body, err := HttpGet(server.URL)
	if err != nil {
		t.Fatalf("HttpGet: %v", err)
	}
	if body != "weather" {
		t.Fatalf("body = %q, want weather", body)
	}
}

func TestHttpGetRejectsErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Error(response, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := HttpGet(server.URL)
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatalf("error = %v, want status 503", err)
	}
}

func TestHttpGetTimesOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		time.Sleep(100 * time.Millisecond)
		response.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	originalClient := HTTPClient
	HTTPClient = &http.Client{Timeout: 10 * time.Millisecond}
	t.Cleanup(func() { HTTPClient = originalClient })
	if _, err := HttpGet(server.URL); err == nil {
		t.Fatal("HttpGet returned nil error after timeout")
	}
}
