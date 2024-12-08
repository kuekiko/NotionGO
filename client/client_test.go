package client

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

const (
	defaultRetryCount   = 3
	defaultRetryWaitMin = 1
	defaultRetryWaitMax = 30
)

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

func setupTestServer(t *testing.T, handler fasthttp.RequestHandler) *fasthttputil.InmemoryListener {
	ln := fasthttputil.NewInmemoryListener()
	go func() {
		if err := fasthttp.Serve(ln, handler); err != nil {
			t.Errorf("Error in test server: %v", err)
		}
	}()
	return ln
}

func TestClientRetry(t *testing.T) {
	attempts := 0
	ln := setupTestServer(t, func(ctx *fasthttp.RequestCtx) {
		attempts++
		if attempts <= 2 {
			ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	defer ln.Close()

	client := NewClient("test-token")
	client.httpClient = &fasthttp.HostClient{
		Addr: "localhost",
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	client.retryWaitMin = 1
	client.retryWaitMax = 5

	resp, err := client.Do(context.Background(), "GET", "test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected status code %d, got %d", fasthttp.StatusOK, resp.StatusCode())
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestClientError(t *testing.T) {
	ln := setupTestServer(t, func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "Invalid request"}`))
	})
	defer ln.Close()

	client := NewClient("test-token")
	client.httpClient = &fasthttp.HostClient{
		Addr: "localhost",
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	_, err := client.Do(context.Background(), "GET", "test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "API错误 400: {\"message\": \"Invalid request\"}"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, err.Error())
	}
}

func TestClientContext(t *testing.T) {
	ln := setupTestServer(t, func(ctx *fasthttp.RequestCtx) {
		time.Sleep(100 * time.Millisecond)
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write([]byte(`{"result": "ok"}`))
	})
	defer ln.Close()

	client := NewClient("test-token")
	client.httpClient = &fasthttp.HostClient{
		Addr: "localhost",
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Do(ctx, "GET", "test", nil)
	if err == nil {
		t.Fatal("Expected context deadline exceeded error, got nil")
	}
}
