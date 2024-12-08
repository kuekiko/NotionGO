package client

import (
	"context"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func BenchmarkClient(b *testing.B) {
	// 创建测试服务器
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.ListenAndServe(ln.Addr().String(), func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("application/json")
			ctx.Write([]byte(`{
				"object": "list",
				"results": [
					{
						"object": "block",
						"id": "test-id",
						"type": "paragraph",
						"paragraph": {
							"rich_text": [
								{
									"type": "text",
									"text": {
										"content": "Hello, World!",
										"link": null
									}
								}
							]
						}
					}
				]
			}`))
		})
		if err != nil {
			b.Fatal(err)
		}
	}()

	// 创建客户端
	client := NewClient("test-token")

	// 运行基准测试
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(context.Background(), "GET", "test", nil)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode() != fasthttp.StatusOK {
				b.Fatalf("unexpected status code: %d", resp.StatusCode())
			}
		}
	})
}

func BenchmarkClientWithLargePayload(b *testing.B) {
	// 创建大型响应数据
	largeResponse := make([]byte, 1024*1024) // 1MB
	for i := range largeResponse {
		largeResponse[i] = 'a'
	}

	// 创建测试服务器
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.ListenAndServe(ln.Addr().String(), func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("application/json")
			ctx.Write(largeResponse)
		})
		if err != nil {
			b.Fatal(err)
		}
	}()

	// 创建客户端
	client := NewClient("test-token")

	// 运行基准测试
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(context.Background(), "GET", "test", nil)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode() != fasthttp.StatusOK {
				b.Fatalf("unexpected status code: %d", resp.StatusCode())
			}
		}
	})
}

func BenchmarkClientWithConcurrency(b *testing.B) {
	// 创建测试服务器
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.ListenAndServe(ln.Addr().String(), func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("application/json")
			ctx.Write([]byte(`{"result": "ok"}`))
		})
		if err != nil {
			b.Fatal(err)
		}
	}()

	// 创建客户端
	client := NewClient("test-token")

	// 运行基准测试
	b.ResetTimer()
	b.SetParallelism(100) // 设置并发数
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(context.Background(), "GET", "test", nil)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode() != fasthttp.StatusOK {
				b.Fatalf("unexpected status code: %d", resp.StatusCode())
			}
		}
	})
}
