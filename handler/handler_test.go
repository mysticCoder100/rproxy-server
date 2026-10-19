package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mysticCoder100/rproxy-server/balancer"
	"github.com/mysticCoder100/rproxy-server/handler"
)

func TestHandleRequest_Success(t *testing.T) {
	// 1. Spin up a local mock backend server
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello from mock backend"))
	}))
	defer mockBackend.Close()

	backendURL, err := url.Parse(mockBackend.URL)
	if err != nil {
		t.Fatalf("failed to parse mock backend URL: %v", err)
	}

	// 2. Initialize the server pool with our mock backend
	// (Assumes Server has an exported Url field to match target.Url)
	pool := balancer.ServerPool{
		&balancer.Server{
			Url: backendURL,
		},
	}

	// 3. Instantiate the proxy handler closure
	proxyHandler := handler.HandleRequest(pool)

	// 4. Create a test request and a response recorder
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rec := httptest.NewRecorder()

	// 5. Execute the handler
	proxyHandler(rec, req)

	// 6. Assert the response
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expectedBody := "hello from mock backend"
	if string(bodyBytes) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, string(bodyBytes))
	}
}

func TestHandleRequest_NoServersAvailable(t *testing.T) {
	// Test edge case: Empty server pool should return 502 Bad Gateway
	emptyPool := balancer.ServerPool{}

	proxyHandler := handler.HandleRequest(emptyPool)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	proxyHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway for empty pool, got %d", resp.StatusCode)
	}
}