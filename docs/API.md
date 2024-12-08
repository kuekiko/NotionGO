# Notion SDK for Go API 文档

## 安装

```bash
go get github.com/kuekiko/NotionGO
```

## 快速开始

```go
package main

import (
    "fmt"
    notion "github.com/kuekiko/NotionGO"
)

func main() {
    // 创建客户端
    client := notion.NewClient("your-api-key")

    // 获取数据库
    db, err := client.Database.Get("database-id")
    if err != nil {
        panic(err)
    }

    fmt.Printf("数据库标题: %s\n", db.Title[0].PlainText)
}
```

## API 参考

### 客户端

```go
// 创建客户端
client := notion.NewClient("your-api-key")
```

### 数据库操作

```go
// 获取数据库
db, err := client.Database.Get("database-id")

// 创建数据库
createParams := &notion.DatabaseCreateParams{
    Parent: notion.Parent{
        Type:   "page_id",
        PageID: "page-id",
    },
    Title: []notion.RichText{
        {
            Type: "text",
            Text: &notion.Text{
                Content: "数据库标题",
            },
        },
    },
    Properties: map[string]notion.Property{
        "Name": {
            Type:  "title",
            Title: &notion.EmptyObject{},
        },
        "Status": {
            Type: "select",
            Select: &notion.SelectConfig{
                Options: []notion.Option{
                    {Name: "未开始", Color: "gray"},
                    {Name: "进行中", Color: "blue"},
                    {Name: "已完成", Color: "green"},
                },
            },
        },
    },
}
db, err := client.Database.Create(createParams)

// 查询数据库
queryParams := &notion.DatabaseQueryParams{
    Filter: map[string]interface{}{
        "property": "Status",
        "select": map[string]interface{}{
            "equals": "进行中",
        },
    },
    PageSize: 10,
}
results, err := client.Database.Query("database-id", queryParams)
```

### 页面操作

```go
// 获取页面
page, err := client.Pages.Get("page-id")

// 创建页面
createParams := &notion.PageCreateParams{
    Parent: notion.Parent{
        Type:       "database_id",
        DatabaseID: "database-id",
    },
    Properties: map[string]interface{}{
        "Name": map[string]interface{}{
            "title": []map[string]interface{}{
                {
                    "text": map[string]interface{}{
                        "content": "页面标题",
                    },
                },
            },
        },
    },
}
page, err := client.Pages.Create(createParams)

// 更新页面
updatePage := &notion.Page{
    Properties: map[string]interface{}{
        "Name": map[string]interface{}{
            "title": []map[string]interface{}{
                {
                    "text": map[string]interface{}{
                        "content": "更新后的标题",
                    },
                },
            },
        },
    },
}
page, err := client.Pages.Update("page-id", updatePage)
```

### 块操作

```go
// 获取块
block, err := client.Blocks.Get("block-id")

// 更新块
updateBlock := &notion.Block{
    Type: notion.TypeParagraph,
    Paragraph: &notion.ParagraphBlock{
        RichText: []notion.RichText{
            {
                Type: "text",
                Text: &notion.Text{
                    Content: "更新后的内容",
                },
            },
        },
    },
}
block, err := client.Blocks.Update("block-id", updateBlock)

// 获取子块
children, err := client.Blocks.ListChildren("block-id", &notion.ListParams{
    PageSize: 10,
})

// 追加子块
children := []notion.Block{
    {
        Type: notion.TypeParagraph,
        Paragraph: &notion.ParagraphBlock{
            RichText: []notion.RichText{
                {
                    Type: "text",
                    Text: &notion.Text{
                        Content: "新段落",
                    },
                },
            },
        },
    },
}
result, err := client.Blocks.AppendChildren("block-id", children)
```

### 搜索操作

```go
// 搜索
params := &notion.SearchParams{
    Query: "搜索关键词",
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
results, err := client.Search.Search(params)
```

### 用户操作

```go
// 获取当前用户
me, err := client.Users.Me()

// 获取用户
user, err := client.Users.Get("user-id")

// 列出用户
users, err := client.Users.List(&notion.ListParams{
    PageSize: 10,
})
```

### 评论操作

```go
// 创建评论
createParams := &notion.CreateCommentParams{
    ParentID:   "block-id",
    ParentType: "block_id",
    RichText: []notion.RichText{
        {
            Type: "text",
            Text: &notion.Text{
                Content: "评论内容",
            },
        },
    },
}
comment, err := client.Comments.Create(createParams)

// 列出评论
comments, err := client.Comments.List("block-id", &notion.ListParams{
    PageSize: 10,
})
```

## 错误处理

SDK 使用自定义的错误类型 `notion.Error`，可以通过以下方式处理错误：

```go
if err != nil {
    if notionErr, ok := err.(*notion.Error); ok {
        fmt.Printf("错误代码: %s\n", notionErr.Code)
        fmt.Printf("错误消息: %s\n", notionErr.Message)
        fmt.Printf("HTTP 状态码: %d\n", notionErr.Status)
    }
}
```

## 类型定义

### 基本类型

```go
// RichText 表示富文本内容
type RichText struct {
    Type        string      `json:"type"`
    Text        *Text       `json:"text,omitempty"`
    Annotations *Annotation `json:"annotations,omitempty"`
    PlainText   string      `json:"plain_text"`
    Href        string      `json:"href,omitempty"`
}

// Text 表示文本内容
type Text struct {
    Content string `json:"content"`
    Link    *Link  `json:"link,omitempty"`
}

// Link 表示链接
type Link struct {
    URL string `json:"url"`
}

// Annotation 表示文本注释
type Annotation struct {
    Bold          bool   `json:"bold"`
    Italic        bool   `json:"italic"`
    Strikethrough bool   `json:"strikethrough"`
    Underline     bool   `json:"underline"`
    Code          bool   `json:"code"`
    Color         Color  `json:"color"`
}
```

### 块类型

```go
// Block 表示块对象
type Block struct {
    Object         string          `json:"object"`
    ID            string          `json:"id"`
    Parent        Parent          `json:"parent"`
    Type          BlockType       `json:"type"`
    CreatedTime   string          `json:"created_time"`
    LastEditedTime string         `json:"last_edited_time"`
    CreatedBy     User           `json:"created_by"`
    LastEditedBy  User           `json:"last_edited_by"`
    HasChildren   bool           `json:"has_children"`
    Archived      bool           `json:"archived"`

    // 不同类型的块具有不同的属性
    Paragraph       *ParagraphBlock       `json:"paragraph,omitempty"`
    Heading1        *HeadingBlock         `json:"heading_1,omitempty"`
    Heading2        *HeadingBlock         `json:"heading_2,omitempty"`
    Heading3        *HeadingBlock         `json:"heading_3,omitempty"`
    BulletedListItem *ListItemBlock        `json:"bulleted_list_item,omitempty"`
    NumberedListItem *ListItemBlock        `json:"numbered_list_item,omitempty"`
    ToDo            *ToDoBlock            `json:"to_do,omitempty"`
    Toggle          *ToggleBlock          `json:"toggle,omitempty"`
    ChildPage       *ChildPageBlock       `json:"child_page,omitempty"`
    ChildDatabase   *ChildDatabaseBlock   `json:"child_database,omitempty"`
    Embed           *EmbedBlock           `json:"embed,omitempty"`
    Image           *File                 `json:"image,omitempty"`
    Video           *File                 `json:"video,omitempty"`
    File            *File                 `json:"file,omitempty"`
    PDF             *File                 `json:"pdf,omitempty"`
    Bookmark        *BookmarkBlock        `json:"bookmark,omitempty"`
    Callout         *CalloutBlock         `json:"callout,omitempty"`
    Quote           *QuoteBlock           `json:"quote,omitempty"`
    Equation        *EquationBlock        `json:"equation,omitempty"`
    Divider         *EmptyObject          `json:"divider,omitempty"`
    TableOfContents *EmptyObject          `json:"table_of_contents,omitempty"`
    Breadcrumb      *EmptyObject          `json:"breadcrumb,omitempty"`
    ColumnList      *ColumnListBlock      `json:"column_list,omitempty"`
    Column          *ColumnBlock          `json:"column,omitempty"`
    LinkPreview     *LinkPreviewBlock     `json:"link_preview,omitempty"`
    Template        *TemplateBlock        `json:"template,omitempty"`
    SyncedBlock     *SyncedBlock          `json:"synced_block,omitempty"`
    Table           *TableBlock           `json:"table,omitempty"`
    TableRow        *TableRowBlock        `json:"table_row,omitempty"`
    Code            *CodeBlock            `json:"code,omitempty"`
}
```

## 最佳实践

1. 使用 Builder 模式构建复杂的搜索参数：

```go
params := notion.NewSearchBuilder().
    Query("搜索关键词").
    Filter(&notion.SearchFilter{
        Property: "object",
        Value:    "page",
    }).
    Sort(&notion.SearchSort{
        Direction: "descending",
        Timestamp: "last_edited_time",
    }).
    PageSize(10).
    Build()

results, err := client.Search.Search(params)
```

2. 使用类型断言安全地处理属性：

```go
if title, ok := page.Properties["Name"].(map[string]interface{}); ok {
    if titleArr, ok := title["title"].([]interface{}); ok && len(titleArr) > 0 {
        if titleObj, ok := titleArr[0].(map[string]interface{}); ok {
            if textObj, ok := titleObj["text"].(map[string]interface{}); ok {
                if content, ok := textObj["content"].(string); ok {
                    fmt.Printf("页面标题: %s\n", content)
                }
            }
        }
    }
}
```

3. 使用常量定义块类型和颜色：

```go
block := &notion.Block{
    Type: notion.TypeParagraph,
    Paragraph: &notion.ParagraphBlock{
        RichText: []notion.RichText{
            {
                Type: "text",
                Text: &notion.Text{
                    Content: "段落内容",
                },
                Annotations: &notion.Annotation{
                    Color: notion.ColorBlue,
                },
            },
        },
    },
}