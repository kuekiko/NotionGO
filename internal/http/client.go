package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/kuekiko/NotionGO/errors"
	"github.com/kuekiko/NotionGO/internal/pool"
	"github.com/valyala/fasthttp"
)

const (
	defaultTimeout      = 30 * time.Second
	defaultIdleTimeout  = 90 * time.Second
	defaultMaxIdleConns = 100
)

// Client 是优化的 HTTP 客户端
type Client struct {
	client        *fasthttp.Client
	retryCount    int
	retryWaitMin  time.Duration
	retryWaitMax  time.Duration
	validateInput bool
}

// ClientOption 定义客户端选项
type ClientOption func(*Client)

// WithRetryCount 设置重试次数
func WithRetryCount(count int) ClientOption {
	return func(c *Client) {
		c.retryCount = count
	}
}

// WithRetryWaitTime 设置重试等待时间
func WithRetryWaitTime(min, max time.Duration) ClientOption {
	return func(c *Client) {
		c.retryWaitMin = min
		c.retryWaitMax = max
	}
}

// WithInputValidation 设置是否验证输入
func WithInputValidation(validate bool) ClientOption {
	return func(c *Client) {
		c.validateInput = validate
	}
}

// NewClient 创建一个新的 HTTP 客户端
func NewClient(opts ...ClientOption) *Client {
	// 创建 fasthttp 客户端
	client := &fasthttp.Client{
		Name: "NotionSDK",
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
		MaxIdleConnDuration: defaultIdleTimeout,
		MaxConnsPerHost:     defaultMaxIdleConns,
		ReadTimeout:         defaultTimeout,
		WriteTimeout:        defaultTimeout,
		MaxResponseBodySize: 10 * 1024 * 1024, // 10MB
	}

	c := &Client{
		client:        client,
		retryCount:    3,
		retryWaitMin:  1 * time.Second,
		retryWaitMax:  30 * time.Second,
		validateInput: true,
	}

	// 应用选项
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Do 执行 HTTP 请求
func (c *Client) Do(ctx context.Context, method, url string, body interface{}, headers map[string]string) ([]byte, error) {
	// 获取请求和响应对象
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 设置方法和 URL
	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	// 设置请求体
	if body != nil {
		// 使用对象池获取缓冲区
		buf := pool.GetBytes()
		defer pool.PutBytes(buf)

		// 编码 JSON
		if err := json.NewEncoder(bytes.NewBuffer(*buf)).Encode(body); err != nil {
			return nil, errors.NewError(errors.ErrInvalidJSON,
				fmt.Sprintf("marshal request body: %v", err),
				400)
		}

		// 验证请求大小
		if c.validateInput {
			if len(*buf) > int(errors.SizeLimits.MaxPayloadSize) {
				return nil, errors.NewError(errors.ErrSizeLimitExceeded,
					fmt.Sprintf("request size %d bytes exceeds maximum of %d bytes",
						len(*buf), errors.SizeLimits.MaxPayloadSize),
					400)
			}
		}

		req.SetBody(*buf)
	}

	// 执行请求（带重试）
	var err error
	for retries := 0; ; retries++ {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return nil, errors.NewError(errors.ErrContextCanceled,
				"context canceled",
				499)
		default:
		}

		err = c.client.Do(req, resp)
		if err != nil {
			if retries >= c.retryCount {
				return nil, errors.NewError(errors.ErrRequestTimeout,
					fmt.Sprintf("execute request (after %d retries): %v", retries, err),
					504)
			}
			time.Sleep(c.getRetryWaitTime(retries))
			continue
		}

		// 处理响应
		statusCode := resp.StatusCode()
		if statusCode >= 400 {
			if retries >= c.retryCount {
				var apiErr struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(resp.Body(), &apiErr); err != nil {
					return nil, errors.NewError(errors.ErrInvalidJSON,
						fmt.Sprintf("failed to decode error response: %v", err),
						statusCode)
				}
				return nil, errors.NewError(errors.ErrorCode(apiErr.Code),
					apiErr.Message,
					statusCode)
			}

			// 如果是速率限制错误，等待适当的时间
			if statusCode == 429 {
				retryAfter := resp.Header.Peek("Retry-After")
				if len(retryAfter) > 0 {
					seconds, _ := time.ParseDuration(string(retryAfter) + "s")
					time.Sleep(seconds)
				} else {
					time.Sleep(c.getRetryWaitTime(retries))
				}
				continue
			}

			// 如果是服务器错误，重试
			if statusCode >= 500 {
				time.Sleep(c.getRetryWaitTime(retries))
				continue
			}
		}

		// 处理压缩响应
		var body []byte
		if bytes.Equal(resp.Header.Peek("Content-Encoding"), []byte("gzip")) {
			reader, err := gzip.NewReader(bytes.NewReader(resp.Body()))
			if err != nil {
				return nil, errors.NewError(errors.ErrInvalidJSON,
					fmt.Sprintf("failed to create gzip reader: %v", err),
					statusCode)
			}
			defer reader.Close()

			body, err = io.ReadAll(reader)
			if err != nil {
				return nil, errors.NewError(errors.ErrInvalidJSON,
					fmt.Sprintf("failed to read gzipped response: %v", err),
					statusCode)
			}
		} else {
			body = resp.Body()
		}

		return body, nil
	}
}

// getRetryWaitTime 计算重试等待时间
func (c *Client) getRetryWaitTime(retries int) time.Duration {
	wait := c.retryWaitMin * time.Duration(1<<uint(retries))
	if wait > c.retryWaitMax {
		wait = c.retryWaitMax
	}
	return wait
}

// Close 关闭客户端
func (c *Client) Close() error {
	c.client.CloseIdleConnections()
	return nil
}
