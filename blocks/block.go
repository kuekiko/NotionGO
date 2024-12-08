package blocks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuekiko/NotionGO/client"
)

// Service 处理块相关的操作
type Service struct {
	client *client.Client
}

// NewService 创建一个新的块服务
func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

// Get 获取块信息
func (s *Service) Get(ctx context.Context, blockID string) (*Block, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/blocks/%s", blockID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}
	defer resp.Body.Close()

	var block Block
	if err := json.NewDecoder(resp.Body).Decode(&block); err != nil {
		return nil, fmt.Errorf("failed to decode block response: %w", err)
	}

	return &block, nil
}

// GetChildren 获取块的子块
func (s *Service) GetChildren(ctx context.Context, blockID string) ([]Block, error) {
	resp, err := s.client.Do(ctx, http.MethodGet, fmt.Sprintf("/blocks/%s/children", blockID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get block children: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Results []Block `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode block children response: %w", err)
	}

	return response.Results, nil
}

// AppendChildren 向块添加子块
func (s *Service) AppendChildren(ctx context.Context, blockID string, children []Block) ([]Block, error) {
	params := struct {
		Children []Block `json:"children"`
	}{
		Children: children,
	}

	resp, err := s.client.Do(ctx, http.MethodPatch, fmt.Sprintf("/blocks/%s/children", blockID), params)
	if err != nil {
		return nil, fmt.Errorf("failed to append block children: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Results []Block `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode append children response: %w", err)
	}

	return response.Results, nil
}

// Delete 删除块
func (s *Service) Delete(ctx context.Context, blockID string) error {
	resp, err := s.client.Do(ctx, http.MethodDelete, fmt.Sprintf("/blocks/%s", blockID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete block: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// Update 更新块
func (s *Service) Update(ctx context.Context, blockID string, content interface{}) (*Block, error) {
	params := map[string]interface{}{
		"content": content,
	}

	resp, err := s.client.Do(ctx, http.MethodPatch, fmt.Sprintf("/blocks/%s", blockID), params)
	if err != nil {
		return nil, fmt.Errorf("failed to update block: %w", err)
	}
	defer resp.Body.Close()

	var block Block
	if err := json.NewDecoder(resp.Body).Decode(&block); err != nil {
		return nil, fmt.Errorf("failed to decode update block response: %w", err)
	}

	return &block, nil
}
