package examples_test

import (
	"context"
	"testing"

	"github.com/kuekiko/NotionGO/client"
	"github.com/kuekiko/NotionGO/database"
	"github.com/kuekiko/NotionGO/pages"
	"github.com/kuekiko/NotionGO/search"
	"github.com/kuekiko/NotionGO/users"
)

func TestProjectManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	ctx := context.Background()
	c := client.NewClient("your-api-key")
	svc := database.NewService(c)

	// 创建项目数据库
	createParams := &database.CreateParams{
		Parent: database.Parent{
			Type:   "page_id",
			PageID: "your-page-id",
		},
		Title: []database.RichText{
			{
				Type: "text",
				Text: database.Text{
					Content: "项目管理",
				},
			},
		},
		Properties: map[string]database.Property{
			"项目名称": {
				Type: "title",
				Name: "项目名称",
			},
			"状态": {
				Type: "select",
				Select: &database.SelectOptions{
					Options: []database.Option{
						{Name: "未开始", Color: "gray"},
						{Name: "进行中", Color: "blue"},
						{Name: "已完成", Color: "green"},
					},
				},
			},
		},
	}

	db, err := svc.Create(ctx, createParams)
	if err != nil {
		t.Fatalf("创建数据库失败: %v", err)
	}

	// 查询数据库
	queryParams := &database.QueryParams{
		Filter: map[string]interface{}{
			"property": "状态",
			"select": map[string]interface{}{
				"equals": "进行中",
			},
		},
	}

	results, err := svc.Query(ctx, db.ID, queryParams)
	if err != nil {
		t.Fatalf("查询数据库失败: %v", err)
	}

	t.Logf("查询到 %d 个进行中的项目", len(results.Results))
}

func TestPageOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	ctx := context.Background()
	c := client.NewClient("your-api-key")
	svc := pages.NewService(c)

	// 创建页面
	createParams := &pages.CreateParams{
		Parent: pages.Parent{
			Type:   "page_id",
			PageID: "your-page-id",
		},
		Properties: map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]interface{}{
						"content": "测试页面",
					},
				},
			},
		},
	}

	page, err := svc.Create(ctx, createParams)
	if err != nil {
		t.Fatalf("创建页面失败: %v", err)
	}

	// 更新页面
	updateParams := &pages.UpdateParams{
		Properties: map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]interface{}{
						"content": "更新后的测试页面",
					},
				},
			},
		},
	}

	updatedPage, err := svc.Update(ctx, page.ID, updateParams)
	if err != nil {
		t.Fatalf("更新页面失败: %v", err)
	}

	// 获取页面
	retrievedPage, err := svc.Get(ctx, updatedPage.ID)
	if err != nil {
		t.Fatalf("获取页面失败: %v", err)
	}

	title := retrievedPage.Properties["title"].(map[string]interface{})
	titleText := title["title"].([]interface{})[0].(map[string]interface{})
	content := titleText["text"].(map[string]interface{})["content"].(string)
	t.Logf("页面标题: %s", content)
}

func TestSearchOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	ctx := context.Background()
	c := client.NewClient("your-api-key")
	searchSvc := search.NewService(c)

	// 搜索页面
	params := &search.SearchParams{
		Query: "测试页面",
		Filter: map[string]interface{}{
			"property": "object",
			"value":    "page",
		},
		Sort: &search.SearchSort{
			Direction: "descending",
			Timestamp: "last_edited_time",
		},
		PageSize: 10,
	}

	results, err := searchSvc.Search(ctx, params)
	if err != nil {
		t.Fatalf("搜索失败: %v", err)
	}

	t.Logf("找到 %d 个匹配的页面", len(results.Results))
}

func TestUserOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	ctx := context.Background()
	c := client.NewClient("your-api-key")
	userSvc := users.NewService(c)

	// 获取当前用户
	me, err := userSvc.Me(ctx)
	if err != nil {
		t.Fatalf("获取当前用户失败: %v", err)
	}

	t.Logf("当前用户: %s (%s)", me.Name, me.ID)

	// 列出所有用户
	users, err := userSvc.List(ctx, nil)
	if err != nil {
		t.Fatalf("列出用户失败: %v", err)
	}

	t.Logf("共有 %d 个用户", len(users.Results))
}
