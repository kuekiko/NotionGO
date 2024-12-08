package tests

import (
	"context"
	"testing"

	"github.com/kuekiko/NotionGO/blocks"
	"github.com/kuekiko/NotionGO/client"
	"github.com/kuekiko/NotionGO/pages"
)

// FuzzCreatePage 模糊测试页面创建
func FuzzCreatePage(f *testing.F) {
	client := client.NewClient("test-token")
	ctx := context.Background()

	f.Add("Test Page", "Test Description")
	f.Fuzz(func(t *testing.T, title string, description string) {
		if len(title) > 2000 || len(description) > 2000 {
			t.Skip("输入太长")
		}

		params := &pages.CreatePageParams{
			Parent: pages.Parent{
				Type:   "page_id",
				PageID: "test-parent-id",
			},
			Properties: pages.Properties{
				"title": pages.TitleProperty{
					Title: []blocks.RichText{
						{
							Type: "text",
							Text: &blocks.Text{
								Content: title,
							},
						},
					},
				},
				"description": pages.RichTextProperty{
					RichText: []blocks.RichText{
						{
							Type: "text",
							Text: &blocks.Text{
								Content: description,
							},
						},
					},
				},
			},
		}

		page, err := client.Pages.Create(ctx, params)
		if err != nil {
			t.Logf("创建页面错误: %v", err)
		}
		if page != nil {
			t.Logf("创建页面成功: %s", page.ID)
		}
	})
}

// FuzzCreateDatabase 模糊测试数据库创建
func FuzzCreateDatabase(f *testing.F) {
	client := client.NewClient("test-token")
	ctx := context.Background()

	f.Add("Test Database")
	f.Fuzz(func(t *testing.T, title string) {
		if len(title) > 2000 {
			t.Skip("输入太长")
		}

		params := &pages.CreateDatabaseParams{
			Parent: pages.Parent{
				Type:   "page_id",
				PageID: "test-parent-id",
			},
			Title: []blocks.RichText{
				{
					Type: "text",
					Text: &blocks.Text{
						Content: title,
					},
				},
			},
			Properties: map[string]pages.PropertyConfig{
				"Name": {
					Title: &pages.EmptyObject{},
				},
			},
		}

		database, err := client.Pages.CreateDatabase(ctx, params)
		if err != nil {
			t.Logf("创建数据库错误: %v", err)
		}
		if database != nil {
			t.Logf("创建数据库成功: %s", database.ID)
		}
	})
}

// FuzzAppendBlocks 模糊测试添加块
func FuzzAppendBlocks(f *testing.F) {
	client := client.NewClient("test-token")
	ctx := context.Background()

	f.Add("Test Block")
	f.Fuzz(func(t *testing.T, content string) {
		if len(content) > 2000 {
			t.Skip("输入太长")
		}

		block := blocks.Block{
			Type: string(blocks.TypeParagraph),
			Paragraph: &blocks.ParagraphBlock{
				RichText: []blocks.RichText{
					{
						Type: "text",
						Text: &blocks.Text{
							Content: content,
						},
					},
				},
			},
		}

		blocks, err := client.Blocks.AppendChildren(ctx, "test-parent-id", []blocks.Block{block})
		if err != nil {
			t.Logf("添加块错误: %v", err)
		}
		if len(blocks) > 0 {
			t.Logf("添加块成功: %s", blocks[0].ID)
		}
	})
}

// FuzzRichText 模糊测试富文本
func FuzzRichText(f *testing.F) {
	f.Add("Test Text", true, true, true, true, "default")
	f.Fuzz(func(t *testing.T, content string, bold bool, italic bool, strikethrough bool, underline bool, color string) {
		if len(content) > 2000 {
			t.Skip("输入太长")
		}

		richText := blocks.RichText{
			Type: "text",
			Text: &blocks.Text{
				Content: content,
			},
			Annotations: &blocks.Annotation{
				Bold:          bold,
				Italic:        italic,
				Strikethrough: strikethrough,
				Underline:     underline,
				Color:         blocks.Color(color),
			},
		}

		if richText.Text.Content != content {
			t.Errorf("内容不匹配: 期望 %s, 实际 %s", content, richText.Text.Content)
		}
	})
}

// FuzzDatabaseQuery 模糊测试数据库查询
func FuzzDatabaseQuery(f *testing.F) {
	client := client.NewClient("test-token")
	ctx := context.Background()

	f.Add("Test Query", 10)
	f.Fuzz(func(t *testing.T, query string, pageSize int) {
		if pageSize < 1 || pageSize > 100 {
			t.Skip("页面大小无效")
		}

		params := &pages.QueryDatabaseParams{
			Filter: &pages.PropertyFilter{
				Property: "title",
				Text: &pages.TextCondition{
					Contains: query,
				},
			},
			PageSize: pageSize,
		}

		results, err := client.Pages.QueryDatabase(ctx, "test-database-id", params)
		if err != nil {
			t.Logf("查询错误: %v", err)
		}
		if results != nil {
			t.Logf("查询成功: 找到 %d 条结果", len(results.Results))
		}
	})
}

// FuzzBlockContent 模糊测试块内容
func FuzzBlockContent(f *testing.F) {
	client := client.NewClient("test-token")
	ctx := context.Background()

	f.Add("Test Content", "paragraph", "default")
	f.Fuzz(func(t *testing.T, content string, blockType string, color string) {
		if len(content) > 2000 {
			t.Skip("输入太长")
		}

		block := blocks.Block{
			Type: blockType,
		}

		switch blockType {
		case string(blocks.TypeParagraph):
			block.Paragraph = &blocks.ParagraphBlock{
				RichText: []blocks.RichText{
					{
						Type: "text",
						Text: &blocks.Text{
							Content: content,
						},
					},
				},
				Color: blocks.Color(color),
			}
		case string(blocks.TypeHeading1):
			block.Heading1 = &blocks.HeadingBlock{
				RichText: []blocks.RichText{
					{
						Type: "text",
						Text: &blocks.Text{
							Content: content,
						},
					},
				},
				Color: blocks.Color(color),
			}
		}

		blocks, err := client.Blocks.AppendChildren(ctx, "test-parent-id", []blocks.Block{block})
		if err != nil {
			t.Logf("添加块错误: %v", err)
		}
		if len(blocks) > 0 {
			t.Logf("添加块成功: %s", blocks[0].ID)
		}
	})
}
