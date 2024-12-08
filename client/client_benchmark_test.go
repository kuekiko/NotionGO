package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/klauspost/compress/gzip"
)

func BenchmarkClientDo(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "test-id", "type": "page"}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.httpClient = &http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}

func BenchmarkClientDoWithBody(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "test-id", "type": "page"}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.httpClient = &http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	body := map[string]interface{}{
		"title": "Test Page",
		"properties": map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]interface{}{
						"content": "Test Page",
					},
				},
			},
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(ctx, http.MethodPost, "/test", body)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}

func BenchmarkClientDoWithRetry(b *testing.B) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts%3 == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "test-id", "type": "page"}`))
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Header().Set("Retry-After", "1")
		}
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.httpClient = &http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	}
	client.retryWaitMin = 1
	client.retryWaitMax = 2

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}

func BenchmarkClientDoWithLargeBody(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "test-id", "type": "page"}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.httpClient = &http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	body := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		body[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(ctx, http.MethodPost, "/test", body)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}

func BenchmarkClientDoWithCompression(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		w.WriteHeader(http.StatusOK)
		gz.Write([]byte(`{"id": "test-id", "type": "page"}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.httpClient = &http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}
