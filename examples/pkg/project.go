package pkg

import (
	notion "github.com/kuekiko/NotionGO"
	"github.com/kuekiko/NotionGO/examples/pkg/config"
	"github.com/kuekiko/NotionGO/examples/pkg/logger"
)

// RunProjectManagement 运行项目管理示例
func RunProjectManagement(cfg *config.Config) {
	logger.Info("开始运行项目管理示例...")

	// 创建客户端
	client := notion.NewClient(cfg.APIKey)
	logger.Debug("已创建 Notion 客户端")

	// 1. 创建项目数据库
	logger.Info("正在创建项目数据库...")
	createParams := &notion.DatabaseCreateParams{
		Parent: notion.Parent{
			Type:   "page_id",
			PageID: cfg.PageID,
		},
		Title: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{
					Content: "项目管理",
				},
			},
		},
		Properties: map[string]notion.Property{
			"名称": {
				Type:  "title",
				Title: &notion.EmptyObject{},
			},
			"状态": {
				Type: "select",
				Select: &notion.SelectConfig{
					Options: []notion.Option{
						{Name: "未开始", Color: "gray"},
						{Name: "进行中", Color: "blue"},
						{Name: "已完成", Color: "green"},
					},
				},
			},
			"优先级": {
				Type: "select",
				Select: &notion.SelectConfig{
					Options: []notion.Option{
						{Name: "低", Color: "gray"},
						{Name: "中", Color: "yellow"},
						{Name: "高", Color: "red"},
					},
				},
			},
			"负责人": {
				Type:   "people",
				People: &notion.EmptyObject{},
			},
			"截止日期": {
				Type: "date",
				Date: &notion.EmptyObject{},
			},
		},
	}

	db, err := client.Database.Create(createParams)
	if err != nil {
		logger.Error("创建数据库失败: %v", err)
		return
	}

	logger.Info("创建数据库成功: %s", db.ID)
	logger.Debug("数据库属性: %+v", db.Properties)

	// 2. 添加新项目
	logger.Info("正在添加新项目...")
	createPageParams := &notion.PageCreateParams{
		Parent: notion.Parent{
			Type:       "database_id",
			DatabaseID: db.ID,
		},
		Properties: map[string]interface{}{
			"名称": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": "示例项目",
						},
					},
				},
			},
			"状态": map[string]interface{}{
				"select": map[string]interface{}{
					"name": "进行中",
				},
			},
			"优先级": map[string]interface{}{
				"select": map[string]interface{}{
					"name": "高",
				},
			},
		},
	}

	page, err := client.Pages.Create(createPageParams)
	if err != nil {
		logger.Error("创建页面失败: %v", err)
		return
	}

	logger.Info("创建项目成功: %s", page.ID)
	logger.Debug("项目属性: %+v", page.Properties)

	// 3. 查询进行中的项目
	logger.Info("正在查询进行中的项目...")
	queryParams := &notion.DatabaseQueryParams{
		Filter: map[string]interface{}{
			"property": "状态",
			"select": map[string]interface{}{
				"equals": "进行中",
			},
		},
		Sorts: []notion.Sort{
			{
				Property:  "优先级",
				Direction: "descending",
			},
		},
		PageSize: 10,
	}

	results, err := client.Database.Query(db.ID, queryParams)
	if err != nil {
		logger.Error("查询数据库失败: %v", err)
		return
	}

	logger.Info("找到 %d 个进行中的项目", len(results.Results))
	for _, result := range results.Results {
		if page, ok := result.(*notion.Page); ok {
			if title, ok := page.Properties["名称"].(map[string]interface{}); ok {
				if titleArr, ok := title["title"].([]interface{}); ok && len(titleArr) > 0 {
					if titleObj, ok := titleArr[0].(map[string]interface{}); ok {
						if textObj, ok := titleObj["text"].(map[string]interface{}); ok {
							if content, ok := textObj["content"].(string); ok {
								logger.Info("- %s", content)
							}
						}
					}
				}
			}
		}
	}

	logger.Info("项目管理示例运行完成")
}
