# Notion SDK for Go

这是一个用 Golang 编写的非官方 Notion API SDK，提供了对 Notion API v1 的完整支持。

## 功能特性

- 完整支持 Notion API v1
- 类型安全的 API 调用
- 自动重试和错误处理
- 速率限制处理
- 并发安全
- 性能优化
  - 使用 fasthttp 替代标准库
  - 使用对象池减少内存分配
  - 使用压缩传输
  - 使用连接池
- 完整的测试覆盖
  - 单元测试
  - 集成测试
  - 性能测试
  - 模糊测试

## 安装

```bash
go get github.com/kuekiko/NotionGO
```

## 快速开始

1. 首先在 [Notion Developers](https://developers.notion.com/) 获取你的 API Key

2. 基本使用示例：

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
    
    // 处理结果
    fmt.Printf("数据库标题: %s\n", db.Title[0].PlainText)
}
```

## 高级配置

```go
// 创建带配置的客户端
client := notion.NewClient(
    "your-api-key",
    notion.WithMaxRetries(5),                    // 设置最大重试次数
    notion.WithRetryWaitTime(1*time.Second, 30*time.Second), // 设置重试等待时间
    notion.WithTimeout(30*time.Second),          // 设置超时时间
    notion.WithCompression(true),                // 启用压缩传输
    notion.WithMaxIdleConns(100),                // 设置最大空闲连接数
    notion.WithMaxIdleConnsPerHost(10),          // 设置每个主机的最大空闲连接数
    notion.WithCache(cache.NewMemoryCache()),    // 启用缓存
    notion.WithMaxConcurrentRequests(100),       // 设置最大并发请求数
)
```

## 错误处理

```go
db, err := client.Database.Get(ctx, "database-id")
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

## 支持的 API

### 数据库操作

- 查询数据库
- 创建数据库
- 更新数据库
- 获取数据库列表

### 页面操作

- 创建页面
- 更新页面
- 获取页面
- 获取页面属性

### 块操作

- 获取块
- 更新块
- 删除块
- 获取子块
- 追加子块

### 用户操作

- 获取用户
- 列出用户
- 获取机器人用户

### 搜索操作

- 搜索页面
- 搜索数据库
- 搜索块

### 评论操作

- 创建评论
- 获取评论列表

## 性能优化

### 使用对象池

```go
// 从对象池获取富文本对象
richText := pool.Get[common.RichText](&pool.RichTextPool)
defer pool.Put(&pool.RichTextPool, richText)

// 设置富文本内容
richText.Type = "text"
richText.Text = &common.Text{
	Content: "Hello World",
}
```

### 使用压缩传输

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithCompression(true),
)
```

### 使用连接池

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithMaxIdleConns(100),
    notion.WithMaxIdleConnsPerHost(10),
)
```

### 使用缓存

```go
client := notion.NewClient(
    "your-api-key",
    notion.WithCache(cache.NewMemoryCache()),
)
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

4. 使用验证器检查输入：

```go
params := &database.CreateParams{
    // ...
}
if err := params.Validate(); err != nil {
    return err
}
```

## 贡献

欢迎提交 Pull Request 和 Issue！

1. Fork 项目
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 测试

```bash
# 运行所有测试
go test ./...

# 运行基准测试
go test -bench=. ./...

# 运行模糊测试
go test -fuzz=. ./...

# 运行性能测试
go test -run=^$ -bench=. -benchmem ./...
```

## 许可证

MIT License

## 相关链接

- [Notion API 文档](https://developers.notion.com/)
- [SDK 文档](./docs/API.md)
- [示例代码](./examples)
- [性能测试报告](./docs/PERFORMANCE.md)
- [更新日志](./CHANGELOG.md)
