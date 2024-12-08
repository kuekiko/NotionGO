package main

import (
	"fmt"
	"os"

	"github.com/kuekiko/NotionGO/examples/pkg"
	"github.com/kuekiko/NotionGO/examples/pkg/config"
	"github.com/kuekiko/NotionGO/examples/pkg/logger"
)

func main() {
	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 加载配置
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Error("加载配置失败: %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		logger.Info("请指定要运行的示例：")
		logger.Info("  integration  - 运行集成示例")
		logger.Info("  project      - 运行项目管理示例")
		logger.Info("  knowledge    - 运行知识库示例")
		logger.Info("  task         - 运行任务管理示例")
		logger.Info("  database     - 读取数据库示例")
		logger.Info("  page         - 读取页面示例")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "integration":
		pkg.RunIntegrationExample(cfg)
	case "project":
		pkg.RunProjectManagement(cfg)
	case "knowledge":
		pkg.RunKnowledgeBase(cfg)
	case "task":
		pkg.RunTaskManagement(cfg)
	case "database":
		pkg.RunDatabaseReader(cfg)
	case "page":
		pkg.RunPageReader(cfg)
	default:
		logger.Error("未知的示例: %s", os.Args[1])
		os.Exit(1)
	}
}
