package tests

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"

	notion "github.com/kuekiko/NotionGO"
)

// 创建测试服务器
func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"object": "list", "results": []}`))
	}))
}

// TestClientPerformance 测试客户端性能
func TestClientPerformance(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := notion.NewClient("test-token")

	concurrency := 10
	requests := 1000

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requests/concurrency; j++ {
				_, err := client.Blocks.Get("test-block-id")
				if err != nil {
					t.Logf("请求错误: %v", err)
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("总请求数: %d", requests)
	t.Logf("并发数: %d", concurrency)
	t.Logf("总耗时: %v", duration)
	t.Logf("平均延迟: %v", duration/time.Duration(requests))
	t.Logf("QPS: %.2f", float64(requests)/duration.Seconds())
}

// TestObjectPoolPerformance 测试对象池性能
func TestObjectPoolPerformance(t *testing.T) {
	iterations := 10000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		block := &notion.Block{}
		block.Type = notion.TypeParagraph
		_ = block
	}

	duration := time.Since(start)
	t.Logf("对象池操作次数: %d", iterations)
	t.Logf("总耗时: %v", duration)
	t.Logf("平均每次操作耗时: %v", duration/time.Duration(iterations))
}

// TestMemoryUsage 测试内存使用
func TestMemoryUsage(t *testing.T) {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 执行一些操作
	iterations := 10000
	for i := 0; i < iterations; i++ {
		block := &notion.Block{}
		block.Type = notion.TypeParagraph
		_ = block
	}

	runtime.ReadMemStats(&m2)

	t.Logf("操作前堆内存: %d bytes", m1.HeapAlloc)
	t.Logf("操作后堆内存: %d bytes", m2.HeapAlloc)
	t.Logf("内存增长: %d bytes", m2.HeapAlloc-m1.HeapAlloc)
	t.Logf("GC 次数: %d", m2.NumGC-m1.NumGC)
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := notion.NewClient("test-token")

	concurrency := 100
	iterations := 100

	var wg sync.WaitGroup
	errors := make(chan error, concurrency*iterations)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, err := client.Blocks.Get("test-block-id")
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)
	errorCount := len(errors)

	t.Logf("总请求数: %d", concurrency*iterations)
	t.Logf("错误数: %d", errorCount)
	t.Logf("总耗时: %v", duration)
	t.Logf("QPS: %.2f", float64(concurrency*iterations)/duration.Seconds())
}
