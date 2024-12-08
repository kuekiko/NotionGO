package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testClient 创建一个用于测试的客户端，返回客户端和测试服务器
func testClient(t *testing.T, fn http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(fn)
	client := NewClient("test-api-key")
	// 替换客户端中的 httpClient，使其请求指向测试服务器
	client.httpClient = &http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	}
	return client, server
}

// testTransport 是一个自定义的 http.RoundTripper，用于重写请求 URL
type testTransport struct {
	server        *httptest.Server
	baseTransport http.RoundTripper
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 将请求重定向到测试服务器
	req.URL.Scheme = "http"
	req.URL.Host = t.server.Listener.Addr().String()
	// 移除 /v1 前缀，如果存在的话
	if len(req.URL.Path) > 3 && req.URL.Path[:3] == "/v1" {
		req.URL.Path = req.URL.Path[3:]
	}
	return t.baseTransport.RoundTrip(req)
}

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.apiKey)
	}

	if client.retryCount != defaultRetryCount {
		t.Errorf("Expected retry count %d, got %d", defaultRetryCount, client.retryCount)
	}
}

func TestClientRetry(t *testing.T) {
	attempts := 0
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	client.retryWaitMin = 1 * time.Millisecond
	client.retryWaitMax = 5 * time.Millisecond

	resp, err := client.Do(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestClientRateLimit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "5")
		w.Header().Set("X-RateLimit-Reset", "1609459200") // 2021-01-01 00:00:00 UTC
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	resp, err := client.Do(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	rateLimit := client.parseRateLimit(resp)
	if rateLimit.Remaining != 5 {
		t.Errorf("Expected remaining rate limit 5, got %d", rateLimit.Remaining)
	}
}

func TestClientError(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request"}`))
	})
	defer server.Close()

	_, err := client.Do(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "HTTP 400: Invalid request"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, err.Error())
	}
}

func TestClientContext(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Do(ctx, http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("Expected context deadline exceeded error, got nil")
	}
}
