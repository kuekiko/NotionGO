# Notion SDK for Go API 文档

## 目录

- [概述](#概述)
- [安装](#安装)
- [快速开始](#快速开始)
- [客户端配置](#客户端配置)
- [错误处理](#错误处理)
- [速率限制](#速率限制)
- [API 参考](#api-参考)
- [最佳实践](#最佳实践)
- [性能优化](#性能优化)

## 概述

Notion SDK for Go 是一个用于访问 Notion API v1 的 Go 语言客户端库。它提供了以下功能：

- 完整的 Notion API v1 支持
- 类型安全的 API 调用
- 自动重试和错误处理
- 速率限制处理
- 并发安全
- 性能优化

## 安装

```bash
go get github.com/kuekiko/NotionGO
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    notion "github.com/kuekiko/NotionGO"
)

func main() {
    // 创建客户端
    client := notion.NewClient("your-api-key")

    // 获取数据库
    db, err := client.Database.Get(context.Background(), "database-id")
    if err != nil {
        panic(err)
    }

    fmt.Printf("数据库标题: %s\n", db.Title[0].PlainText)
}
```

## 客户端配置

客户端支持多种配置选项：

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithMaxRetries(5),                    // 设置最大重试次数
    notion.WithRetryWaitTime(1*time.Second, 30*time.Second), // 设置重试等待时间
    notion.WithTimeout(30*time.Second),          // 设置超时时间
    notion.WithHTTPClient(customHTTPClient),     // 设置自定义 HTTP 客户端
)
```

## 错误处理

SDK 定义了多种错误类型：

```go
if err != nil {
    switch {
    case errors.IsNotFound(err):
        // 处理 404 错误
    case errors.IsRateLimited(err):
        // 处理速率限制错误
    case errors.IsSizeLimitExceeded(err):
        // 处理大小限制错误
    case errors.IsValidationError(err):
        // 处理验证错误
    default:
        // 处理其他错误
    }
}
```

## 速率限制

SDK 自动处理速率限制：

```go
// 获取当前速率限制状态
rateLimit := client.GetRateLimit()
fmt.Printf("剩余请求数: %d\n", rateLimit.Remaining)
fmt.Printf("重置时间: %s\n", rateLimit.ResetAt)

// 等待直到速率限制重置
client.WaitForRateLimit()
```

## API 参考

### 数据库操作

```go
// 获取数据库
db, err := client.Database.Get(ctx, "database-id")

// 查询数据库
results, err := client.Database.Query(ctx, "database-id", &database.QueryParams{
    Filter: map[string]interface{}{
        "property": "Status",
        "select": map[string]interface{}{
            "equals": "Done",
        },
    },
    Sorts: []database.Sort{
        {
            Property:  "LastEdited",
            Direction: "descending",
        },
    },
    PageSize: 100,
})

// 创建数据库
newDB, err := client.Database.Create(ctx, &database.CreateParams{
    Parent: database.Parent{
        Type:   "page_id",
        PageID: "parent-page-id",
    },
    Title: []database.RichText{
        {
            Type: "text",
            Text: database.Text{
                Content: "项目跟踪",
            },
        },
    },
    Properties: map[string]database.Property{
        "Name": {
            Type: "title",
            Name: "Name",
        },
        "Status": {
            Type: "select",
            Select: &database.SelectOptions{
                Options: []database.Option{
                    {Name: "Not Started", Color: "gray"},
                    {Name: "In Progress", Color: "blue"},
                    {Name: "Done", Color: "green"},
                },
            },
        },
    },
})
```

### 页面操作

```go
// 获取页面
page, err := client.Pages.Get(ctx, "page-id")

// 创建页面
newPage, err := client.Pages.Create(ctx, &pages.CreateParams{
    Parent: pages.Parent{
        Type:       "database_id",
        DatabaseID: "database-id",
    },
    Properties: map[string]interface{}{
        "Name": map[string]interface{}{
            "title": []map[string]interface{}{
                {
                    "text": map[string]interface{}{
                        "content": "新任务",
                    },
                },
            },
        },
        "Status": map[string]interface{}{
            "select": map[string]interface{}{
                "name": "Not Started",
            },
        },
    },
})

// 更新页面
updatedPage, err := client.Pages.Update(ctx, "page-id", &pages.UpdateParams{
    Properties: map[string]interface{}{
        "Status": map[string]interface{}{
            "select": map[string]interface{}{
                "name": "In Progress",
            },
        },
    },
})
```

### 块操作

```go
// 获取块
block, err := client.Blocks.Get(ctx, "block-id")

// 获取子块
children, err := client.Blocks.GetChildren(ctx, "block-id")

// 添加块
newBlocks, err := client.Blocks.AppendChildren(ctx, "block-id", []blocks.Block{
    notion.BlockBuilder(blocks.TypeParagraph).
        WithParagraph("这是一个段落").
        Build(),
    notion.BlockBuilder(blocks.TypeHeading1).
        WithHeading("这是一个标题").
        Build(),
})

// 删除块
err := client.Blocks.Delete(ctx, "block-id")
```

### 搜索操作

```go
// 搜索
results, err := client.Search.Search(ctx, notion.SearchBuilder().
    Query("项目").
    FilterByType("page").
    SortBy("descending", "last_edited_time").
    PageSize(10).
    Build())
```

### 用户操作

```go
// 获取当前用户
me, err := client.Users.Me(ctx)

// 获取用户
user, err := client.Users.Get(ctx, "user-id")

// 列出用户
users, err := client.Users.List(ctx, &users.ListParams{
    PageSize: 100,
})
```

### 评论操作

```go
// 创建评论
comment, err := client.Comments.Create(ctx, &comments.CreateParams{
    ParentID:   "page-id",
    ParentType: "page_id",
    RichText: []comments.RichText{
        {
            Type: "text",
            Text: &comments.Text{
                Content: "这是一条评论",
            },
        },
    },
})

// 获取评论
comments, err := client.Comments.List(ctx, &comments.ListParams{
    BlockID: "block-id",
})
```

## 最佳实践

1. 使用上下文控制请求超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

db, err := client.Database.Get(ctx, "database-id")
```

2. 批量操作时使用分页：

```go
var allResults []pages.Page
startCursor := ""

for {
    results, err := client.Database.Query(ctx, "database-id", &database.QueryParams{
        StartCursor: startCursor,
        PageSize:    100,
    })
    if err != nil {
        return err
    }

    allResults = append(allResults, results.Results...)

    if !results.HasMore {
        break
    }
    startCursor = results.NextCursor
}
```

3. 处理大型响应时使用流式处理：

```go
err := client.Database.QueryStream(ctx, "database-id", &database.QueryParams{}, func(page *pages.Page) error {
    // 处理每个页面
    return nil
})
```

4. 使用重试机制处理临时错误：

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithMaxRetries(5),
    notion.WithRetryWaitTime(1*time.Second, 30*time.Second),
)
```

5. 使用验证器检查输入：

```go
params := &database.CreateParams{
    // ...
}
if err := params.Validate(); err != nil {
    return err
}
```

## 性能优化

1. 使用对象池减少内存分配：

```go
// 获取对象
richText := pool.Get[blocks.RichText](&pool.RichTextPool)
defer pool.Put(&pool.RichTextPool, richText)
```

2. 使用压缩传输：

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithCompression(true),
)
```

3. 使用连接池：

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithMaxIdleConns(100),
    notion.WithMaxIdleConnsPerHost(10),
)
```

4. 使用缓存：

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithCache(cache.NewMemoryCache()),
)
```

5. 使用并发控制：

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithMaxConcurrentRequests(100),
)