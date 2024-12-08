package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/kuekiko/NotionGO/errors"
)

const (
	BaseURL    = "https://api.notion.com/v1"
	APIVersion = "2022-06-28"

	defaultRetryCount   = 3
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultTimeout      = 30 * time.Second
	defaultIdleTimeout  = 90 * time.Second
	defaultMaxIdleConns = 100
)

// Client 是 Notion API 的主要客户端
type Client struct {
	apiKey        string
	httpClient    *http.Client
	retryCount    int
	retryWaitMin  time.Duration
	retryWaitMax  time.Duration
	validateInput bool
}

// RateLimit 表示 API 速率限制信息
type RateLimit struct {
	Remaining int
	ResetAt   time.Time
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

// WithHTTPClient 设置自定义的 HTTP 客户端
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithInputValidation 设置是否验证输入
func WithInputValidation(validate bool) ClientOption {
	return func(c *Client) {
		c.validateInput = validate
	}
}

// NewClient 创建一个新的 Notion 客户端
func NewClient(apiKey string, opts ...ClientOption) *Client {
	// 创建默认的传输层
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

	// 创建默认的 HTTP 客户端
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   defaultTimeout,
	}

	c := &Client{
		apiKey:        apiKey,
		httpClient:    httpClient,
		retryCount:    defaultRetryCount,
		retryWaitMin:  defaultRetryWaitMin,
		retryWaitMax:  defaultRetryWaitMax,
		validateInput: true,
	}

	// 应用选项
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// validateRequestSize 验证请求大小
func (c *Client) validateRequestSize(body []byte) error {
	if len(body) > int(errors.SizeLimits.MaxPayloadSize) {
		return errors.NewError(errors.ErrSizeLimitExceeded,
			fmt.Sprintf("request size %d bytes exceeds maximum of %d bytes",
				len(body), errors.SizeLimits.MaxPayloadSize),
			http.StatusBadRequest)
	}
	return nil
}

// Do 执行 HTTP 请求
func (c *Client) Do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", BaseURL, path)

	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, errors.NewError(errors.ErrInvalidJSON,
				fmt.Sprintf("marshal request body: %v", err),
				http.StatusBadRequest)
		}

		// 验证请求大小
		if c.validateInput {
			if err := c.validateRequestSize(bodyBytes); err != nil {
				return nil, err
			}
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.NewError(errors.ErrInvalidRequest,
			fmt.Sprintf("create request: %v", err),
			http.StatusBadRequest)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Notion-Version", APIVersion)
	req.Header.Set("Content-Type", "application/json")
	if bodyBytes != nil {
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	}

	// 执行请求（带重试）
	var resp *http.Response
	var rateLimit *RateLimit

	for retries := 0; ; retries++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if retries >= c.retryCount {
				return nil, errors.NewError(errors.ErrRequestTimeout,
					fmt.Sprintf("execute request (after %d retries): %v", retries, err),
					http.StatusGatewayTimeout)
			}
			continue
		}

		// 解析速率限制信息
		rateLimit = c.parseRateLimit(resp)

		// 处理响应
		if resp.StatusCode >= 400 {
			if retries >= c.retryCount {
				defer resp.Body.Close()
				var apiErr struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
					return nil, errors.NewError(errors.ErrInvalidJSON,
						fmt.Sprintf("failed to decode error response: %v", err),
						resp.StatusCode)
				}
				return nil, errors.NewError(errors.ErrorCode(apiErr.Code),
					apiErr.Message,
					resp.StatusCode)
			}

			// 如果是速率限制错误，等待适当的时间
			if resp.StatusCode == http.StatusTooManyRequests {
				waitTime := c.getRetryWaitTime(retries, rateLimit)
				select {
				case <-ctx.Done():
					return nil, errors.NewError(errors.ErrContextCanceled,
						"context canceled while waiting for rate limit",
						http.StatusRequestTimeout)
				case <-time.After(waitTime):
					continue
				}
			}

			// 如果是服务器错误，重试
			if resp.StatusCode >= 500 {
				continue
			}
		}

		break
	}

	return resp, nil
}

// parseRateLimit 解析响应头中的速率限制信息
func (c *Client) parseRateLimit(resp *http.Response) *RateLimit {
	remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	resetAt, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)

	return &RateLimit{
		Remaining: remaining,
		ResetAt:   time.Unix(resetAt, 0),
	}
}

// getRetryWaitTime 计算重试等待时间
func (c *Client) getRetryWaitTime(retries int, rateLimit *RateLimit) time.Duration {
	// 如果有速率限制信息，使用重置时间
	if rateLimit != nil && rateLimit.ResetAt.After(time.Now()) {
		return time.Until(rateLimit.ResetAt)
	}

	// 否则使用指数退避
	wait := c.retryWaitMin * time.Duration(1<<uint(retries))
	if wait > c.retryWaitMax {
		wait = c.retryWaitMax
	}
	return wait
}

// GetAPIKey 返回 API Key
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// SetHTTPClient 设置客户端的 HTTP 客户端（主要用于测试）
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}
