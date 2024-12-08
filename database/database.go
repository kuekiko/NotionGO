package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuekiko/NotionGO/client"
)

// Database 表示一个 Notion 数据库
type Database struct {
	ID             string              `json:"id"`
	CreatedTime    string              `json:"created_time"`
	LastEditedTime string              `json:"last_edited_time"`
	Title          []RichText          `json:"title"`
	Properties     map[string]Property `json:"properties"`
}

// RichText 表示富文本内容
type RichText struct {
	Type      string `json:"type"`
	PlainText string `json:"plain_text"`
	Text      Text   `json:"text,omitempty"`
}

// Text 表示文本内容
type Text struct {
	Content string `json:"content"`
}

// Property 表示数据库的属性
type Property struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Name   string         `json:"name"`
	Select *SelectOptions `json:"select,omitempty"`
	Date   *DateOptions   `json:"date,omitempty"`
	People *PeopleOptions `json:"people,omitempty"`
	Status *StatusOptions `json:"status,omitempty"`
}

// SelectOptions 表示选择选项
type SelectOptions struct {
	Options []Option `json:"options"`
}

// Option 表示选项
type Option struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// DateOptions 表示日期选项
type DateOptions struct {
	Format string `json:"format,omitempty"`
}

// PeopleOptions 表示人员选项
type PeopleOptions struct {
}

// StatusOptions 表示状态选项
type StatusOptions struct {
	Options []Option `json:"options"`
}

// Service 处理数据库相关的操作
type Service struct {
	client *client.Client
}

// NewService 创建一个新的数据库服务
func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

// Get 获取数据库信息
func (s *Service) Get(ctx context.Context, databaseID string) (*Database, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/databases/%s", databaseID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}
	defer resp.Body.Close()

	var database Database
	if err := json.NewDecoder(resp.Body).Decode(&database); err != nil {
		return nil, fmt.Errorf("failed to decode database response: %w", err)
	}

	return &database, nil
}

// CreateParams 表示创建数据库的参数
type CreateParams struct {
	Parent     Parent              `json:"parent"`
	Title      []RichText          `json:"title"`
	Properties map[string]Property `json:"properties"`
	Icon       *Icon               `json:"icon,omitempty"`
	Cover      *Cover              `json:"cover,omitempty"`
}

// Parent 表示父级
type Parent struct {
	Type   string `json:"type"`
	PageID string `json:"page_id,omitempty"`
}

// Icon 表示图标
type Icon struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji,omitempty"`
	File  *File  `json:"file,omitempty"`
}

// Cover 表示封面
type Cover struct {
	Type string `json:"type"`
	File *File  `json:"file,omitempty"`
}

// File 表示文件
type File struct {
	URL      string `json:"url"`
	ExpireAt string `json:"expire_at,omitempty"`
}

// Create 创建数据库
func (s *Service) Create(ctx context.Context, params *CreateParams) (*Database, error) {
	resp, err := s.client.Do(ctx, http.MethodPost, "/databases", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	defer resp.Body.Close()

	var database Database
	if err := json.NewDecoder(resp.Body).Decode(&database); err != nil {
		return nil, fmt.Errorf("failed to decode database response: %w", err)
	}

	return &database, nil
}

// QueryParams 表示数据库查询参数
type QueryParams struct {
	Filter      map[string]interface{} `json:"filter,omitempty"`
	Sorts       []Sort                 `json:"sorts,omitempty"`
	StartCursor string                 `json:"start_cursor,omitempty"`
	PageSize    int                    `json:"page_size,omitempty"`
}

// Sort 表示排序参数
type Sort struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}

// Query 查询数据库
func (s *Service) Query(ctx context.Context, databaseID string, params *QueryParams) (*QueryResponse, error) {
	resp, err := s.client.Do(ctx, http.MethodPost, fmt.Sprintf("/databases/%s/query", databaseID), params)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer resp.Body.Close()

	var queryResp QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode query response: %w", err)
	}

	return &queryResp, nil
}

// QueryResponse 表示数据库查询响应
type QueryResponse struct {
	Results    []Page `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

// Page 表示数据库中的页面
type Page struct {
	ID             string                 `json:"id"`
	CreatedTime    string                 `json:"created_time"`
	LastEditedTime string                 `json:"last_edited_time"`
	Properties     map[string]interface{} `json:"properties"`
}
