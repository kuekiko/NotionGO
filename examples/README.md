# Notion SDK 示例

本目录包含了一些使用 Notion SDK 的示例代码，展示了如何使用 SDK 实现各种常见场景。

## 配置说明

在运行示例之前，需要先配置必要的参数：

1. 首次运行时会自动在 `config` 目录下创建 `config.json` 文件
2. 编辑 `config.json` 文件，填写以下信息：
   ```json
   {
       "api_key": "your-api-key",      // 替换为你的 Notion API Key
       "page_id": "your-page-id",      // 替换为你的页面 ID
       "database_id": "your-database-id" // 替换为你的数据库 ID
   }
   ```

### 获取配置信息

1. API Key:
   - 访问 [Notion Developers](https://developers.notion.com/)
   - 创建一个新的集成
   - 复制生成的 API Key

2. Page ID:
   - 在 Notion 中打开目标页面
   - 从 URL 中复制页面 ID（形如：`https://notion.so/Your-Page-{page_id}`）

3. Database ID:
   - 在 Notion 中打开目标数据库
   - 从 URL 中复制数据库 ID（形如：`https://notion.so/Your-Database-{database_id}`）

## 日志说明

示例程序会自动在 `logs` 目录下生成日志文件：
- 文件名格式：`notion_YYYY-MM-DD.log`
- 包含三种级别的日志：
  - INFO: 普通信息
  - ERROR: 错误信息
  - DEBUG: 调试信息

## 示例列表

### 1. 集成示例 (integration)
展示了基本的 API 使用方法：
- 获取和查询数据库
- 创建和更新页面
- 设置页面属性

### 2. 项目管理 (project)
展示如何使用 Notion 创建和管理项目：
- 创建项目数据库
- 添加新项目
- 设置项目状态和优先级
- 查询进行中的项目

### 3. 知识库 (knowledge)
展示如何使用 Notion 构建知识库：
- 创建文档页面
- 添加富文本内容
- 使用不同类型的块
- 添加评论
- 搜索文档

### 4. 任务管理 (task)
展示如何使用 Notion 进行任务管理：
- 创建任务数据库
- 分配任务给用户
- 设置任务状态和截止日期
- 查询待处理的任务

## 使用方法

1. 克隆仓库
2. 首次运行会自动创建配置文件
3. 编辑 `config/config.json` 填写实际配置
4. 运行示例：

```bash
# 进入示例目录
cd examples/cmd

# 运行集成示例
go run main.go integration

# 运行项目管理示例
go run main.go project

# 运行知识库示例
go run main.go knowledge

# 运行任务管理示例
go run main.go task
```

## 错误处理

示例程序包含完整的错误处理：
1. 配置错误：
   - 配置文件不存在时会自动创建
   - 验证配置有效性
   - 详细的错误提示

2. API 错误：
   - 网络连接错误
   - 认证错误
   - 权限错误
   - 速率限制
   - 资源不存在

3. 数据处理错误：
   - 类型转换错误
   - 数据验证错误
   - 格式错误

所有错误都会：
- 在控制台显示
- 记录到日志文件
- 提供清晰的错误描述和处理建议

## 注意事项

1. 运行示例前请确保：
   - 获取了 Notion API Key
   - 将集成添加到目标页面
   - 有适当的权限

2. 示例中的数据库结构和属性都是可以根据需要修改的

3. 建议先在测试页面上运行示例，以免影响生产环境

4. 查看日志文件了解详细的执行过程和错误信息

## 更多信息

- [Notion API 文档](https://developers.notion.com/)
- [SDK 文档](../docs/API.md)
- [Notion 开发者社区](https://developers.notion.com/community) 