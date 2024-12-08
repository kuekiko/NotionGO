package comments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuekiko/NotionGO/client"
)

// Service 处理评论相关的操作
type Service struct {
	client *client.Client
}

// NewService 创建一个新的评论服务
func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

// Comment 表示一条评论
type Comment struct {
	ID             string          `json:"id"`
	ParentID       string          `json:"parent_id"`
	ParentType     string          `json:"parent_type"`
	CreatedTime    string          `json:"created_time"`
	LastEditedTime string          `json:"last_edited_time"`
	CreatedBy      json.RawMessage `json:"created_by"`
	RichText       []RichText      `json:"rich_text"`
	Resolved       bool            `json:"resolved"`
	Discussion     *Discussion     `json:"discussion,omitempty"`
}

// Discussion 表示评论讨论
type Discussion struct {
	ID       string `json:"id"`
	ThreadID string `json:"thread_id"`
}

// RichText 表示富文本内容
type RichText struct {
	Type        string      `json:"type"`
	Text        *Text       `json:"text,omitempty"`
	Annotations *Annotation `json:"annotations,omitempty"`
	PlainText   string      `json:"plain_text"`
	Href        string      `json:"href,omitempty"`
}

// Text 表示文本内容
type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

// Link 表示链接
type Link struct {
	URL string `json:"url"`
}

// Annotation 表示文本注释
type Annotation struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

// ListResponse 表示评论列表响应
type ListResponse struct {
	Results    []Comment `json:"results"`
	HasMore    bool      `json:"has_more"`
	NextCursor string    `json:"next_cursor,omitempty"`
}

// ListParams 表示列出评论的参数
type ListParams struct {
	BlockID     string `json:"block_id,omitempty"`
	PageID      string `json:"page_id,omitempty"`
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

// CreateParams 表示创建评论的参数
type CreateParams struct {
	ParentID   string      `json:"parent_id"`
	ParentType string      `json:"parent_type"`
	RichText   []RichText  `json:"rich_text"`
	Discussion *Discussion `json:"discussion,omitempty"`
}

// List 列出评论
func (s *Service) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, "/comments", params)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	defer resp.Body.Close()

	var listResp ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode list response: %w", err)
	}

	return &listResp, nil
}

// Create 创建评论
func (s *Service) Create(ctx context.Context, params *CreateParams) (*Comment, error) {
	resp, err := s.client.Do(ctx, http.MethodPost, "/comments", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	defer resp.Body.Close()

	var comment Comment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("failed to decode comment response: %w", err)
	}

	return &comment, nil
}
