package notion

import (
	"testing"
)

const (
	testAPIKey     = ""
	testDatabaseID = ""
)

func TestIntegrationDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := NewClient(testAPIKey)

	// 获取数据库
	db, err := client.Database.Get(testDatabaseID)
	if err != nil {
		t.Fatalf("获取数据库失败: %v", err)
	}

	t.Logf("数据库 ID: %s", db.ID)
	if len(db.Title) > 0 {
		t.Logf("数据库标题: %s", db.Title[0].PlainText)
	}

	// 查询数据库
	queryParams := &DatabaseQueryParams{
		PageSize: 10,
	}

	results, err := client.Database.Query(testDatabaseID, queryParams)
	if err != nil {
		t.Fatalf("查询数据库失败: %v", err)
	}

	t.Logf("查询到 %d 条记录", len(results.Results))
}

func TestIntegrationPage(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := NewClient(testAPIKey)

	// 创建页面
	createParams := &PageCreateParams{
		Parent: Parent{
			Type:       "database_id",
			DatabaseID: testDatabaseID,
		},
		Properties: map[string]interface{}{
			"Name": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": "Test Page",
						},
					},
				},
			},
		},
	}

	page, err := client.Pages.Create(createParams)
	if err != nil {
		t.Fatalf("创建页面失败: %v", err)
	}

	t.Logf("创建的页面 ID: %s", page.ID)

	// 获取页面
	retrievedPage, err := client.Pages.Get(page.ID)
	if err != nil {
		t.Fatalf("获取页面失败: %v", err)
	}

	// 安全地获取标题
	if title, ok := retrievedPage.Properties["Name"].(map[string]interface{}); ok {
		if titleArr, ok := title["title"].([]interface{}); ok && len(titleArr) > 0 {
			if titleObj, ok := titleArr[0].(map[string]interface{}); ok {
				if textObj, ok := titleObj["text"].(map[string]interface{}); ok {
					if content, ok := textObj["content"].(string); ok {
						t.Logf("页面标题: %s", content)
					}
				}
			}
		}
	}

	// 更新页面
	updatePage := &Page{
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
		t.Fatalf("更新页面失败: %v", err)
	}

	t.Logf("页面已更新: %s", updatedPage.ID)
}
