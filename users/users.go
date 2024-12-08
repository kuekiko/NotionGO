package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuekiko/NotionGO/client"
)

// Service 处理用户相关的操作
type Service struct {
	client *client.Client
}

// NewService 创建一个新的用户服务
func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

// User 表示 Notion 用户
type User struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Person    *Person `json:"person,omitempty"`
	Bot       *Bot    `json:"bot,omitempty"`
}

// Person 表示个人用户信息
type Person struct {
	Email string `json:"email"`
}

// Bot 表示机器人用户信息
type Bot struct {
	Owner       *Owner `json:"owner"`
	WorkspaceID string `json:"workspace_id"`
}

// Owner 表示机器人所有者
type Owner struct {
	Type      string `json:"type"`
	Workspace bool   `json:"workspace"`
}

// ListResponse 表示用户列表响应
type ListResponse struct {
	Results    []User `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// ListParams 表示列出用户的参数
type ListParams struct {
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

// Get 获取用户信息
func (s *Service) Get(ctx context.Context, userID string) (*User, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/users/%s", userID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil
}

// List 列出所有用户
func (s *Service) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, "/users", params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer resp.Body.Close()

	var listResp ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode list response: %w", err)
	}

	return &listResp, nil
}

// Me 获取当前用户信息
func (s *Service) Me(ctx context.Context) (*User, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, "/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode current user response: %w", err)
	}

	return &user, nil
}
