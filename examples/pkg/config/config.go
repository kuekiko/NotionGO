package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 表示配置信息
type Config struct {
	APIKey     string `json:"api_key"`     // Notion API Key
	PageID     string `json:"page_id"`     // 页面 ID
	DatabaseID string `json:"database_id"` // 数据库 ID
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	// 获取配置文件的绝对路径
	configPath := filepath.Join("config", filename)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果配置文件不存在，创建示例配置
		config := &Config{
			APIKey:     "your-api-key",
			PageID:     "your-page-id",
			DatabaseID: "your-database-id",
		}
		// 创建配置目录
		if err := os.MkdirAll("config", 0755); err != nil {
			return nil, fmt.Errorf("创建配置目录失败: %v", err)
		}
		// 写入示例配置
		file, err := os.Create(configPath)
		if err != nil {
			return nil, fmt.Errorf("创建配置文件失败: %v", err)
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(config); err != nil {
			return nil, fmt.Errorf("写入配置文件失败: %v", err)
		}
		return nil, fmt.Errorf("已创建示例配置文件 %s，请填写实际的配置信息", configPath)
	}

	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	config := new(Config)
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if config.APIKey == "your-api-key" {
		return nil, fmt.Errorf("请在配置文件中设置实际的 API Key")
	}

	return config, nil
}
