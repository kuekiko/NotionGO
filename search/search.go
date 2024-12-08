package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuekiko/NotionGO/client"
)

// Service 处理搜索相关的操作
type Service struct {
	client *client.Client
}

// NewService 创建一个新的搜索服务
func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

// SearchParams 定义搜索参数
type SearchParams struct {
	Query       string                 `json:"query,omitempty"`
	Filter      map[string]interface{} `json:"filter,omitempty"`
	Sort        *SearchSort            `json:"sort,omitempty"`
	StartCursor string                 `json:"start_cursor,omitempty"`
	PageSize    int                    `json:"page_size,omitempty"`
}

// SearchSort 定义搜索排序
type SearchSort struct {
	Direction string `json:"direction"`
	Timestamp string `json:"timestamp"`
}

// SearchResponse 定义搜索响应
type SearchResponse struct {
	Results    []interface{} `json:"results"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// Search 执行搜索
func (s *Service) Search(ctx context.Context, params *SearchParams) (*SearchResponse, error) {
	resp, err := s.client.Do(ctx, http.MethodPost, "/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &searchResp, nil
}

// SearchBuilder 用于构建搜索请求
type SearchBuilder struct {
	params SearchParams
}

// NewSearchBuilder 创建新的搜索构建器
func NewSearchBuilder() *SearchBuilder {
	return &SearchBuilder{}
}

// Query 设置搜索查询
func (b *SearchBuilder) Query(query string) *SearchBuilder {
	b.params.Query = query
	return b
}

// FilterByType 按类型过滤
func (b *SearchBuilder) FilterByType(objectType string) *SearchBuilder {
	if b.params.Filter == nil {
		b.params.Filter = make(map[string]interface{})
	}
	b.params.Filter["object"] = objectType
	return b
}

// SortBy 设置排序
func (b *SearchBuilder) SortBy(direction, timestamp string) *SearchBuilder {
	b.params.Sort = &SearchSort{
		Direction: direction,
		Timestamp: timestamp,
	}
	return b
}

// PageSize 设置页面大小
func (b *SearchBuilder) PageSize(size int) *SearchBuilder {
	b.params.PageSize = size
	return b
}

// StartCursor 设置起始游标
func (b *SearchBuilder) StartCursor(cursor string) *SearchBuilder {
	b.params.StartCursor = cursor
	return b
}

// Build 构建搜索参数
func (b *SearchBuilder) Build() *SearchParams {
	return &b.params
}
