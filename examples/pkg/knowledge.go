package pkg

import (
	notion "github.com/kuekiko/NotionGO"
	"github.com/kuekiko/NotionGO/examples/pkg/config"
	"github.com/kuekiko/NotionGO/examples/pkg/logger"
)

// RunKnowledgeBase 运行知识库示例
func RunKnowledgeBase(cfg *config.Config) {
	logger.Info("开始运行知识库示例...")

	// 创建客户端
	client := notion.NewClient(cfg.APIKey)
	logger.Debug("已创建 Notion 客户端")

	// 1. 创建文档页面
	logger.Info("正在创建文档页面...")
	createParams := &notion.PageCreateParams{
		Parent: notion.Parent{
			Type:   "page_id",
			PageID: cfg.PageID,
		},
		Properties: map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]interface{}{
						"content": "Go 编程指南",
					},
				},
			},
		},
		Children: []notion.Block{
			{
				Type: notion.TypeHeading1,
				Heading1: &notion.HeadingBlock{
					RichText: []notion.RichText{
						{
							Type: "text",
							Text: &notion.Text{
								Content: "Go 编程指南",
							},
						},
					},
				},
			},
			{
				Type: notion.TypeParagraph,
				Paragraph: &notion.ParagraphBlock{
					RichText: []notion.RichText{
						{
							Type: "text",
							Text: &notion.Text{
								Content: "本指南介绍 Go 语言的基础知识和最佳实践。",
							},
						},
					},
				},
			},
			{
				Type: notion.TypeHeading2,
				Heading2: &notion.HeadingBlock{
					RichText: []notion.RichText{
						{
							Type: "text",
							Text: &notion.Text{
								Content: "1. 安装",
							},
						},
					},
				},
			},
			{
				Type: notion.TypeParagraph,
				Paragraph: &notion.ParagraphBlock{
					RichText: []notion.RichText{
						{
							Type: "text",
							Text: &notion.Text{
								Content: "访问 ",
							},
						},
						{
							Type: "text",
							Text: &notion.Text{
								Content: "Go 官网",
								Link: &notion.Link{
									URL: "https://golang.org",
								},
							},
						},
						{
							Type: "text",
							Text: &notion.Text{
								Content: " 下载并安装 Go。",
							},
						},
					},
				},
			},
			{
				Type: notion.TypeCode,
				Code: &notion.CodeBlock{
					Language: "bash",
					RichText: []notion.RichText{
						{
							Type: "text",
							Text: &notion.Text{
								Content: "go version",
							},
						},
					},
				},
			},
		},
	}

	page, err := client.Pages.Create(createParams)
	if err != nil {
		logger.Error("创建页面失败: %v", err)
		return
	}

	logger.Info("创建文档成功: %s", page.ID)
	logger.Debug("页面属性: %+v", page.Properties)

	// 2. 添加评论
	logger.Info("正在添加评论...")
	commentParams := &notion.CreateCommentParams{
		ParentID:   page.ID,
		ParentType: "page_id",
		RichText: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{
					Content: "这是一个很好的入门指南！",
				},
			},
		},
	}

	comment, err := client.Comments.Create(commentParams)
	if err != nil {
		logger.Error("添加评论失败: %v", err)
		return
	}

	logger.Info("添加评论成功: %s", comment.ID)

	// 3. 搜索文档
	logger.Info("正在搜索文档...")
	searchParams := &notion.SearchParams{
		Query: "Go 编程",
		Filter: &notion.SearchFilter{
			Property: "object",
			Value:    "page",
		},
		Sort: &notion.SearchSort{
			Direction: "descending",
			Timestamp: "last_edited_time",
		},
		PageSize: 10,
	}

	results, err := client.Search.Search(searchParams)
	if err != nil {
		logger.Error("搜索失败: %v", err)
		return
	}

	logger.Info("找到 %d 个相关文档", len(results.Results))
	for _, result := range results.Results {
		if page, ok := result.(*notion.Page); ok {
			if title, ok := page.Properties["title"].(map[string]interface{}); ok {
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

	logger.Info("知识库示例运行完成")
}
