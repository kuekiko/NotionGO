package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuekiko/NotionGO/client"
)

// Page 表示一个 Notion 页面
type Page struct {
	ID             string                 `json:"id"`
	CreatedTime    string                 `json:"created_time"`
	LastEditedTime string                 `json:"last_edited_time"`
	Parent         Parent                 `json:"parent"`
	Properties     map[string]interface{} `json:"properties"`
}

// Parent 表示页面的父级
type Parent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id,omitempty"`
	PageID     string `json:"page_id,omitempty"`
}

// Service 处理页面相关的操作
type Service struct {
	client *client.Client
}

// NewService 创建一个新的页面服务
func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

// Get 获取页面信息
func (s *Service) Get(ctx context.Context, pageID string) (*Page, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/pages/%s", pageID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	defer resp.Body.Close()

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to decode page response: %w", err)
	}

	return &page, nil
}

// CreateParams 表示创建页面的参数
type CreateParams struct {
	Parent     Parent                 `json:"parent"`
	Properties map[string]interface{} `json:"properties"`
	Children   []Block                `json:"children,omitempty"`
}

// Create 创建新页面
func (s *Service) Create(ctx context.Context, params *CreateParams) (*Page, error) {
	resp, err := s.client.Do(ctx, http.MethodPost, "/pages", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer resp.Body.Close()

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to decode create page response: %w", err)
	}

	return &page, nil
}

// UpdateParams 表示更新页面的参数
type UpdateParams struct {
	Properties map[string]interface{} `json:"properties"`
	Archived   bool                   `json:"archived,omitempty"`
}

// Update 更新页面
func (s *Service) Update(ctx context.Context, pageID string, params *UpdateParams) (*Page, error) {
	resp, err := s.client.Do(ctx, http.MethodPatch, fmt.Sprintf("/pages/%s", pageID), params)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}
	defer resp.Body.Close()

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to decode update page response: %w", err)
	}

	return &page, nil
}

// Block 表示页面中的块
type Block struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// PropertyItem 表示页面属性项
type PropertyItem struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	NextURL       string          `json:"next_url,omitempty"`
	HasMore       bool            `json:"has_more"`
	PropertyValue json.RawMessage `json:"property_value"`
}

// PropertyResponse 表示属性响应
type PropertyResponse struct {
	Object         string          `json:"object"`
	Results        []PropertyItem  `json:"results"`
	NextURL        string          `json:"next_url,omitempty"`
	HasMore        bool            `json:"has_more"`
	PropertyItem   *PropertyItem   `json:"property_item,omitempty"`
	PropertyValues json.RawMessage `json:"property_values,omitempty"`
}

// GetProperty 获取页面属性
func (s *Service) GetProperty(ctx context.Context, pageID, propertyID string) (*PropertyResponse, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/pages/%s/properties/%s", pageID, propertyID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get page property: %w", err)
	}
	defer resp.Body.Close()

	var property PropertyResponse
	if err := json.NewDecoder(resp.Body).Decode(&property); err != nil {
		return nil, fmt.Errorf("failed to decode property response: %w", err)
	}

	return &property, nil
}

// GetAllProperties 获取页面所有属性
func (s *Service) GetAllProperties(ctx context.Context, pageID string) (map[string]*PropertyResponse, error) {
	page, err := s.Get(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	properties := make(map[string]*PropertyResponse)
	for propertyID := range page.Properties {
		property, err := s.GetProperty(ctx, pageID, propertyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get property %s: %w", propertyID, err)
		}
		properties[propertyID] = property
	}

	return properties, nil
}
