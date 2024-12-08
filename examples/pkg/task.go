package pkg

import (
	"time"

	notion "github.com/kuekiko/NotionGO"
	"github.com/kuekiko/NotionGO/examples/pkg/config"
	"github.com/kuekiko/NotionGO/examples/pkg/logger"
)

// RunTaskManagement 运行任务管理示例
func RunTaskManagement(cfg *config.Config) {
	logger.Info("开始运行任务管理示例...")

	// 创建客户端
	client := notion.NewClient(cfg.APIKey)
	logger.Debug("已创建 Notion 客户端")

	// 1. 创建任务数据库
	logger.Info("正在创建任务数据库...")
	createParams := &notion.DatabaseCreateParams{
		Parent: notion.Parent{
			Type:   "page_id",
			PageID: cfg.PageID,
		},
		Title: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{
					Content: "任务管理",
				},
			},
		},
		Properties: map[string]notion.Property{
			"任务名称": {
				Type:  "title",
				Title: &notion.EmptyObject{},
			},
			"状态": {
				Type: "select",
				Select: &notion.SelectConfig{
					Options: []notion.Option{
						{Name: "待处理", Color: "gray"},
						{Name: "进行中", Color: "blue"},
						{Name: "已完成", Color: "green"},
						{Name: "已取消", Color: "red"},
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
			"完成度": {
				Type:   "number",
				Number: &notion.NumberConfig{Format: "percent"},
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

	// 2. 添加新任务
	logger.Info("正在添加新任务...")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	createPageParams := &notion.PageCreateParams{
		Parent: notion.Parent{
			Type:       "database_id",
			DatabaseID: db.ID,
		},
		Properties: map[string]interface{}{
			"任务名称": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": "完成 SDK 文档",
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
			"截止日期": map[string]interface{}{
				"date": map[string]interface{}{
					"start": tomorrow,
				},
			},
			"完成度": map[string]interface{}{
				"number": 30,
			},
		},
	}

	page, err := client.Pages.Create(createPageParams)
	if err != nil {
		logger.Error("创建任务失败: %v", err)
		return
	}

	logger.Info("创建任务成功: %s", page.ID)
	logger.Debug("任务属性: %+v", page.Properties)

	// 3. 查询待处理的任务
	logger.Info("正在查询待处理的任务...")
	queryParams := &notion.DatabaseQueryParams{
		Filter: map[string]interface{}{
			"and": []map[string]interface{}{
				{
					"property": "状态",
					"select": map[string]interface{}{
						"equals": "待处理",
					},
				},
				{
					"property": "截止日期",
					"date": map[string]interface{}{
						"on_or_before": tomorrow,
					},
				},
			},
		},
		Sorts: []notion.Sort{
			{
				Property:  "优先级",
				Direction: "descending",
			},
			{
				Property:  "截止日期",
				Direction: "ascending",
			},
		},
		PageSize: 10,
	}

	results, err := client.Database.Query(db.ID, queryParams)
	if err != nil {
		logger.Error("查询任务失败: %v", err)
		return
	}

	logger.Info("找到 %d 个待处理的任务", len(results.Results))
	for _, result := range results.Results {
		if page, ok := result.(*notion.Page); ok {
			if title, ok := page.Properties["任务名称"].(map[string]interface{}); ok {
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

	logger.Info("任务管理示例运行完成")
}
