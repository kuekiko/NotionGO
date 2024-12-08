package pkg

import (
	notion "github.com/kuekiko/NotionGO"
	"github.com/kuekiko/NotionGO/examples/pkg/config"
	"github.com/kuekiko/NotionGO/examples/pkg/logger"
)

// RunIntegrationExample 运行集成示例
func RunIntegrationExample(cfg *config.Config) {
	logger.Info("开始运行集成示例...")

	// 创建客户端
	client := notion.NewClient(cfg.APIKey)
	logger.Debug("已创建 Notion 客户端")

	// 1. 获取数据库
	logger.Info("正在获取数据库...")
	db, err := client.Database.Get(cfg.DatabaseID)
	if err != nil {
		logger.Error("获取数据库失败: %v", err)
		return
	}

	logger.Info("数据库 ID: %s", db.ID)
	if len(db.Title) > 0 {
		logger.Info("数据库标题: %s", db.Title[0].PlainText)
	}

	// 2. 查询数据库
	logger.Info("正在查询数据库...")
	queryParams := &notion.DatabaseQueryParams{
		PageSize: 10,
	}

	results, err := client.Database.Query(cfg.DatabaseID, queryParams)
	if err != nil {
		logger.Error("查询数据库失败: %v", err)
		return
	}

	logger.Info("查询到 %d 条记录", len(results.Results))

	// 3. 创建页面
	logger.Info("正在创建页面...")
	createParams := &notion.PageCreateParams{
		Parent: notion.Parent{
			Type:       "database_id",
			DatabaseID: cfg.DatabaseID,
		},
		Properties: map[string]interface{}{
			"Name": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": "测试页面",
						},
					},
				},
			},
			"Status": map[string]interface{}{
				"select": map[string]interface{}{
					"name": "进行中",
				},
			},
		},
	}

	page, err := client.Pages.Create(createParams)
	if err != nil {
		logger.Error("创建页面失败: %v", err)
		return
	}

	logger.Info("创建的页面 ID: %s", page.ID)
	logger.Debug("页面属性: %+v", page.Properties)

	// 4. 更新页面
	logger.Info("正在更新页面...")
	updatePage := &notion.Page{
		Properties: map[string]interface{}{
			"Name": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": "更新后的测试页面",
						},
					},
				},
			},
			"Status": map[string]interface{}{
				"select": map[string]interface{}{
					"name": "已完成",
				},
			},
		},
	}

	updatedPage, err := client.Pages.Update(page.ID, updatePage)
	if err != nil {
		logger.Error("更新页面失败: %v", err)
		return
	}

	logger.Info("页面已更新: %s", updatedPage.ID)
	logger.Info("集成示例运行完成")
}
