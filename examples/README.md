# Notion SDK 示例

本目录包含了一些使用 Notion SDK 的示例代码，展示了如何使用 SDK 实现各种常见场景。

## 示例列表

### 1. 项目管理 (project_management.go)

展示如何使用 Notion 创建和管理项目，包括：
- 创建项目数据库
- 添加新项目
- 设置项目状态和优先级
- 查询进行中的项目

### 2. 知识库 (knowledge_base.go)

展示如何使用 Notion 构建知识库，包括：
- 创建文档页面
- 添加富文本内容
- 使用不同类型的块
- 添加评论
- 搜索文档

### 3. 任务管理 (task_management.go)

展示如何使用 Notion 进行任务管理，包括：
- 创建任务数据库
- 分配任务给用户
- 设置任务状态和截止日期
- 查询待处理的任务

## 使用方法

1. 克隆仓库
2. 在每个示例文件中替换 `your-api-key` 为你的 Notion API Key
3. 在每个示例文件中替换 `your-page-id` 为你的目标页面 ID
4. 运行示例：

```bash
# 运行项目管理示例
go run cmd/main.go project

# 运行知识库示例
go run cmd/main.go knowledge

# 运行任务管理示例
go run cmd/main.go task
```

## 注意事项

1. 运行示例前请确保已经：
   - 获取了 Notion API Key
   - 将集成添加到目标页面
   - 有适当的权限

2. 示例中的数据库结构和属性都是可以根据需要修改的

3. 建议先在测试页面上运行示例，以免影响生产环境

## 更多信息

- [Notion API 文档](https://developers.notion.com/)
- [SDK 文档](../README.md)
- [Notion 开发者社区](https://developers.notion.com/community) 