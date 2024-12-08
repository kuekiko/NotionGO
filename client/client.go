package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// BaseURL 是 Notion API 的基础 URL
	BaseURL = "https://api.notion.com/v1/"
	// APIVersion 是 Notion API 的版本
	APIVersion = "2022-06-28"
)

// Client 表示 Notion API 客户端
type Client struct {
	apiKey       string
	httpClient   *fasthttp.Client
	retryCount   int
	retryWaitMin time.Duration
	retryWaitMax time.Duration
}

// NewClient 创建一个新的客户端
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &fasthttp.Client{
			Name:                          "NotionGO",
			MaxConnsPerHost:               10000,
			MaxIdleConnDuration:           10 * time.Second,
			ReadTimeout:                   10 * time.Second,
			WriteTimeout:                  10 * time.Second,
			MaxResponseBodySize:           50 * 1024 * 1024, // 50MB
			DisableHeaderNamesNormalizing: true,
			MaxConnWaitTimeout:            30 * time.Second,
		},
		retryCount:   3,
		retryWaitMin: 1 * time.Second,
		retryWaitMax: 30 * time.Second,
	}
}

// Do 执行 HTTP 请求
func (c *Client) Do(ctx context.Context, method, path string, body interface{}) (*fasthttp.Response, error) {
	// 创建请求和响应对象
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	// 设置 URL
	req.SetRequestURI(BaseURL + path)
	req.Header.SetMethod(method)

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Notion-Version", APIVersion)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 如果有请求体，编码为 JSON
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("编码请求体失败: %v", err)
		}
		req.SetBody(jsonBody)
		fmt.Printf("请求体: %s\n", string(jsonBody))
	}

	fmt.Printf("请求 URL: %s\n", req.URI().String())
	fmt.Printf("请求头: %+v\n", req.Header.String())

	// 发送请求
	if err := c.httpClient.Do(req, resp); err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	return resp, nil
}

// Get 发送 GET 请求
func (c *Client) Get(path string, params interface{}, v interface{}) error {
	ctx := context.Background()
	resp, err := c.Do(ctx, "GET", path, params)
	if err != nil {
		return err
	}

	// 检查状态码
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("API错误 %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// 解码响应
	if err := json.Unmarshal(resp.Body(), v); err != nil {
		return fmt.Errorf("解码响应失败: %v", err)
	}

	return nil
}

// Post 发送 POST 请求
func (c *Client) Post(path string, body interface{}, v interface{}) error {
	ctx := context.Background()
	resp, err := c.Do(ctx, "POST", path, body)
	if err != nil {
		return err
	}

	// 检查状态码
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("API错误 %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// 解码响应
	if err := json.Unmarshal(resp.Body(), v); err != nil {
		return fmt.Errorf("解码响应失败: %v", err)
	}

	return nil
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(path string, body interface{}, v interface{}) error {
	ctx := context.Background()
	resp, err := c.Do(ctx, "PATCH", path, body)
	if err != nil {
		return err
	}

	// 检查状态码
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("API错误 %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// 解码响应
	if err := json.Unmarshal(resp.Body(), v); err != nil {
		return fmt.Errorf("解码响应失败: %v", err)
	}

	return nil
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(path string) error {
	ctx := context.Background()
	resp, err := c.Do(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	// 检查状态码
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("API错误 %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
