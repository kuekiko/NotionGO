package pkg

import (
	"encoding/json"

	notion "github.com/kuekiko/NotionGO"
	"github.com/kuekiko/NotionGO/examples/pkg/config"
	"github.com/kuekiko/NotionGO/examples/pkg/logger"
)

// RunDatabaseReader 运行数据库读取示例
func RunDatabaseReader(cfg *config.Config) {
	logger.Info("开始运行数据库读取示例...")

	// 验证配置
	if cfg.DatabaseID == "" {
		logger.Error("数据库 ID 不能为空")
		return
	}
	if len(cfg.DatabaseID) != 32 {
		logger.Error("数据库 ID 格式不正确，应为32位字符")
		return
	}

	// 创建客户端
	client := notion.NewClient(cfg.APIKey)
	logger.Debug("已创建 Notion 客户端")

	// 1. 获取数据库信息
	logger.Info("正在获取数据库信息...")
	logger.Debug("数据库 ID: %s", cfg.DatabaseID)
	db, err := client.Database.Get(cfg.DatabaseID)
	if err != nil {
		logger.Error("获取数据库失败: %v", err)
		return
	}

	// 打印数据库基本信息
	logger.Info("数据库基本信息:")
	logger.Info("- ID: %s", db.ID)
	logger.Info("- 创建时间: %s", db.CreatedTime)
	logger.Info("- 最后编辑时间: %s", db.LastEditedTime)
	logger.Info("- 创建者: %s (%s)", db.CreatedBy.Name, db.CreatedBy.ID)
	logger.Info("- URL: %s", db.URL)

	// 打印数据库标题
	if len(db.Title) > 0 {
		var title string
		for _, t := range db.Title {
			title += t.PlainText
		}
		logger.Info("- 标题: %s", title)
	}

	// 打印数据库属性
	logger.Info("\n数据库属性:")
	for name, prop := range db.Properties {
		propJSON, _ := json.MarshalIndent(prop, "  ", "  ")
		logger.Info("- %s: %s", name, string(propJSON))
	}

	// 2. 查询数据库内容
	logger.Info("\n正在查询数据库内容...")
	queryParams := &notion.DatabaseQueryParams{
		PageSize: 100, // 获取最多100条记录
		Filter:   nil, // 不使用过滤器
		Sorts:    nil, // 不使用排序
	}

	results, err := client.Database.Query(cfg.DatabaseID, queryParams)
	if err != nil {
		logger.Error("查询数据库失败: %v", err)
		return
	}

	// 打印所有记录
	logger.Info("\n数据库记录 (共 %d 条):", len(results.Results))
	for i, result := range results.Results {
		if page, ok := result.(*notion.Page); ok {
			logger.Info("\n记录 #%d:", i+1)
			logger.Info("- ID: %s", page.ID)
			logger.Info("- 创建时间: %s", page.CreatedTime)
			logger.Info("- 最后编辑时间: %s", page.LastEditedTime)
			logger.Info("- URL: %s", page.URL)

			// 打印所有属性
			logger.Info("- 属性:")
			for propName, propValue := range page.Properties {
				propJSON, _ := json.MarshalIndent(propValue, "    ", "  ")
				logger.Info("  %s: %s", propName, string(propJSON))
			}
		}
	}

	logger.Info("\n数据库读取示例运行完成")
}

// RunPageReader 运行页面读取示例
func RunPageReader(cfg *config.Config) {
	logger.Info("开始运行页面读取示例...")

	// 验证配置
	if cfg.PageID == "" {
		logger.Error("页面 ID 不能为空")
		return
	}

	// 创建客户端
	client := notion.NewClient(cfg.APIKey)
	logger.Debug("已创建 Notion 客户端")

	// 1. 获取页面信息
	logger.Info("正在获取页面信息...")
	page, err := client.Pages.Get(cfg.PageID)
	if err != nil {
		logger.Error("获取页面失败: %v", err)
		return
	}

	// 打印页面基本信息
	logger.Info("\n页面基本信息:")
	logger.Info("- ID: %s", page.ID)
	logger.Info("- 创建时间: %s", page.CreatedTime)
	logger.Info("- 最后编辑时间: %s", page.LastEditedTime)
	logger.Info("- 创建者: %s (%s)", page.CreatedBy.Name, page.CreatedBy.ID)
	logger.Info("- URL: %s", page.URL)

	// 打印页面属性
	logger.Info("\n页面属性:")
	for name, prop := range page.Properties {
		propJSON, _ := json.MarshalIndent(prop, "  ", "  ")
		logger.Info("- %s: %s", name, string(propJSON))
	}

	// 2. 获取页面内容（块）
	logger.Info("\n正在获取页面内容...")
	blocks, err := client.Blocks.ListChildren(page.ID, nil)
	if err != nil {
		logger.Error("获取页面内容失败: %v", err)
		return
	}

	// 打印所有块
	logger.Info("\n页面内容 (共 %d 个块):", len(blocks.Results))
	for i, result := range blocks.Results {
		block, ok := result.(*notion.Block)
		if !ok {
			logger.Error("块类型转换失败")
			continue
		}

		logger.Info("\n块 #%d:", i+1)
		logger.Info("- 类型: %s", block.Type)

		// 根据块类型打印内容
		switch block.Type {
		case notion.TypeParagraph:
			if block.Paragraph != nil && len(block.Paragraph.RichText) > 0 {
				var text string
				for _, rt := range block.Paragraph.RichText {
					text += rt.PlainText
				}
				logger.Info("- 内容: %s", text)
			}
		case notion.TypeHeading1:
			if block.Heading1 != nil && len(block.Heading1.RichText) > 0 {
				var text string
				for _, rt := range block.Heading1.RichText {
					text += rt.PlainText
				}
				logger.Info("- 内容: %s", text)
			}
		case notion.TypeHeading2:
			if block.Heading2 != nil && len(block.Heading2.RichText) > 0 {
				var text string
				for _, rt := range block.Heading2.RichText {
					text += rt.PlainText
				}
				logger.Info("- 内容: %s", text)
			}
		case notion.TypeHeading3:
			if block.Heading3 != nil && len(block.Heading3.RichText) > 0 {
				var text string
				for _, rt := range block.Heading3.RichText {
					text += rt.PlainText
				}
				logger.Info("- 内容: %s", text)
			}
		case notion.TypeBulletedListItem:
			if block.BulletedListItem != nil && len(block.BulletedListItem.RichText) > 0 {
				var text string
				for _, rt := range block.BulletedListItem.RichText {
					text += rt.PlainText
				}
				logger.Info("- 内容: • %s", text)
			}
		case notion.TypeNumberedListItem:
			if block.NumberedListItem != nil && len(block.NumberedListItem.RichText) > 0 {
				var text string
				for _, rt := range block.NumberedListItem.RichText {
					text += rt.PlainText
				}
				logger.Info("- 内容: %d. %s", i+1, text)
			}
		case notion.TypeToDo:
			if block.ToDo != nil && len(block.ToDo.RichText) > 0 {
				var text string
				for _, rt := range block.ToDo.RichText {
					text += rt.PlainText
				}
				checked := "[ ]"
				if block.ToDo.Checked {
					checked = "[x]"
				}
				logger.Info("- 内容: %s %s", checked, text)
			}
		case notion.TypeCode:
			if block.Code != nil && len(block.Code.RichText) > 0 {
				var text string
				for _, rt := range block.Code.RichText {
					text += rt.PlainText
				}
				logger.Info("- 语言: %s", block.Code.Language)
				logger.Info("- 代码:\n%s", text)
			}
		case notion.TypeImage:
			if block.Image != nil {
				if block.Image.File != nil {
					logger.Info("- 图片 URL: %s", block.Image.File.URL)
				} else if block.Image.External != nil {
					logger.Info("- 图片 URL: %s", block.Image.External.URL)
				}
			}
		default:
			blockJSON, _ := json.MarshalIndent(block, "  ", "  ")
			logger.Info("- 原始内容: %s", string(blockJSON))
		}
	}

	logger.Info("\n页面读取示例运行完成")
}

// printRichText 打印富文本内容
func printRichText(blockType string, richText []notion.RichText) {
	var text string
	for _, rt := range richText {
		text += rt.PlainText
	}
	logger.Info("- %s: %s", blockType, text)

	// 打印详细的富文本信息
	for i, rt := range richText {
		logger.Debug("  富文本 #%d:", i+1)
		logger.Debug("  - 类型: %s", rt.Type)
		logger.Debug("  - 文本: %s", rt.PlainText)
		if rt.Annotations != nil {
			logger.Debug("  - 样式: 粗体=%v, 斜体=%v, 删除线=%v, 下划线=%v, 代码=%v, 颜色=%s",
				rt.Annotations.Bold,
				rt.Annotations.Italic,
				rt.Annotations.Strikethrough,
				rt.Annotations.Underline,
				rt.Annotations.Code,
				rt.Annotations.Color)
		}
		if rt.Text != nil && rt.Text.Link != nil {
			logger.Debug("  - 链接: %s", rt.Text.Link.URL)
		}
	}
}
