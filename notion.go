package notion

import (
	"time"

	"github.com/kuekiko/NotionGO/blocks"
	"github.com/kuekiko/NotionGO/client"
	"github.com/kuekiko/NotionGO/comments"
	"github.com/kuekiko/NotionGO/database"
	"github.com/kuekiko/NotionGO/pages"
	"github.com/kuekiko/NotionGO/search"
	"github.com/kuekiko/NotionGO/users"
)

// Client 是 Notion API 的主要客户端
type Client struct {
	client    *client.Client
	Database  *database.Service
	Pages     *pages.Service
	Blocks    *blocks.Service
	Search    *search.Service
	Users     *users.Service
	Comments  *comments.Service
	RateLimit *RateLimiter
}

// RateLimiter 处理 API 速率限制
type RateLimiter struct {
	Remaining  int
	ResetAt    time.Time
	MaxRetries int
}

// ClientOption 定义客户端选项
type ClientOption func(*Client)

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(retries int) ClientOption {
	return func(c *Client) {
		c.RateLimit.MaxRetries = retries
	}
}

// WithRetryWaitTime 设置重试等待时间
func WithRetryWaitTime(min, max time.Duration) ClientOption {
	return func(c *Client) {
		c.client = client.NewClient(c.client.GetAPIKey(), client.WithRetryWaitTime(min, max))
	}
}

// NewClient 创建一个新的 Notion 客户端
func NewClient(apiKey string, opts ...ClientOption) *Client {
	baseClient := client.NewClient(apiKey)
	c := &Client{
		client: baseClient,
		RateLimit: &RateLimiter{
			MaxRetries: 3, // 默认重试次数
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(c)
	}

	// 初始化各个服务
	c.Database = database.NewService(baseClient)
	c.Pages = pages.NewService(baseClient)
	c.Blocks = blocks.NewService(baseClient)
	c.Search = search.NewService(baseClient)
	c.Users = users.NewService(baseClient)
	c.Comments = comments.NewService(baseClient)

	return c
}

// Version 返回 SDK 版本
func Version() string {
	return "1.0.0"
}

// BlockBuilder 创建一个新的块构建器
func BlockBuilder(blockType blocks.BlockType) *blocks.BlockBuilder {
	return blocks.NewBlockBuilder(blockType)
}

// SearchBuilder 创建一个新的搜索构建器
func SearchBuilder() *search.SearchBuilder {
	return search.NewSearchBuilder()
}

// GetRateLimit 获取当前速率限制状态
func (c *Client) GetRateLimit() *RateLimiter {
	return c.RateLimit
}

// WaitForRateLimit 等待直到速率限制重置
func (c *Client) WaitForRateLimit() {
	if c.RateLimit.Remaining <= 0 && c.RateLimit.ResetAt.After(time.Now()) {
		time.Sleep(time.Until(c.RateLimit.ResetAt))
	}
}
