package notion

import (
	"encoding/json"
)

// Color 表示颜色
type Color string

const (
	ColorDefault Color = "default"
	ColorGray    Color = "gray"
	ColorBrown   Color = "brown"
	ColorOrange  Color = "orange"
	ColorYellow  Color = "yellow"
	ColorGreen   Color = "green"
	ColorBlue    Color = "blue"
	ColorPurple  Color = "purple"
	ColorPink    Color = "pink"
	ColorRed     Color = "red"

	ColorGrayBackground   Color = "gray_background"
	ColorBrownBackground  Color = "brown_background"
	ColorOrangeBackground Color = "orange_background"
	ColorYellowBackground Color = "yellow_background"
	ColorGreenBackground  Color = "green_background"
	ColorBlueBackground   Color = "blue_background"
	ColorPurpleBackground Color = "purple_background"
	ColorPinkBackground   Color = "pink_background"
	ColorRedBackground    Color = "red_background"
)

// BlockType 表示块类型
type BlockType string

const (
	// 基本块类型
	TypeParagraph        BlockType = "paragraph"
	TypeHeading1         BlockType = "heading_1"
	TypeHeading2         BlockType = "heading_2"
	TypeHeading3         BlockType = "heading_3"
	TypeBulletedListItem BlockType = "bulleted_list_item"
	TypeNumberedListItem BlockType = "numbered_list_item"
	TypeToDo             BlockType = "to_do"
	TypeToggle           BlockType = "toggle"
	TypeCode             BlockType = "code"
	TypeQuote            BlockType = "quote"
	TypeCallout          BlockType = "callout"
	TypeDivider          BlockType = "divider"

	// 媒体块类型
	TypeImage       BlockType = "image"
	TypeVideo       BlockType = "video"
	TypeFile        BlockType = "file"
	TypePDF         BlockType = "pdf"
	TypeBookmark    BlockType = "bookmark"
	TypeEmbed       BlockType = "embed"
	TypeLinkPreview BlockType = "link_preview"

	// 数据块类型
	TypeTable           BlockType = "table"
	TypeTableRow        BlockType = "table_row"
	TypeTableOfContents BlockType = "table_of_contents"
	TypeEquation        BlockType = "equation"

	// 布局块类型
	TypeColumnList BlockType = "column_list"
	TypeColumn     BlockType = "column"
	TypeBreadcrumb BlockType = "breadcrumb"

	// 特殊块类型
	TypeChildPage     BlockType = "child_page"
	TypeChildDatabase BlockType = "child_database"
	TypeSyncedBlock   BlockType = "synced_block"
	TypeTemplate      BlockType = "template"
)

// Text 表示文本内容
type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

// Link 表示链接
type Link struct {
	URL string `json:"url"`
}

// RichText 表示富文本内容
type RichText struct {
	Type        string      `json:"type"`
	Text        *Text       `json:"text,omitempty"`
	Annotations *Annotation `json:"annotations,omitempty"`
	PlainText   string      `json:"plain_text"`
	Href        string      `json:"href,omitempty"`
}

// Annotation 表示文本注释
type Annotation struct {
	Bold          bool  `json:"bold"`
	Italic        bool  `json:"italic"`
	Strikethrough bool  `json:"strikethrough"`
	Underline     bool  `json:"underline"`
	Code          bool  `json:"code"`
	Color         Color `json:"color"`
}

// Parent 表示父对象
type Parent struct {
	Type       string `json:"type"`                  // "page_id", "database_id", "block_id" 或 "workspace"
	PageID     string `json:"page_id,omitempty"`     // 当 Type 为 "page_id" 时
	DatabaseID string `json:"database_id,omitempty"` // 当 Type 为 "database_id" 时
	BlockID    string `json:"block_id,omitempty"`    // 当 Type 为 "block_id" 时
}

// User 表示用户
type User struct {
	Object    string  `json:"object"`               // 总是 "user"
	ID        string  `json:"id"`                   // 用户 ID
	Type      string  `json:"type"`                 // "person" 或 "bot"
	Name      string  `json:"name"`                 // 用户名称
	AvatarURL string  `json:"avatar_url,omitempty"` // 头像 URL
	Person    *Person `json:"person,omitempty"`     // 个人用户信息
	Bot       *Bot    `json:"bot,omitempty"`        // 机器人用户信息
}

// Person 表示个人用户信息
type Person struct {
	Email string `json:"email"`
}

// Bot 表示机器人用户信息
type Bot struct {
	Owner       *Owner `json:"owner"`
	WorkspaceID string `json:"workspace_id"`
}

// Owner 表示机器人所有者
type Owner struct {
	Type      string `json:"type"`
	Workspace bool   `json:"workspace"`
}

// File 表示文件
type File struct {
	Type     string     `json:"type"`               // "external" 或 "file"
	External *External  `json:"external,omitempty"` // 当 Type 为 "external" 时
	File     *FileInfo  `json:"file,omitempty"`     // 当 Type 为 "file" 时
	Caption  []RichText `json:"caption,omitempty"`
}

// External 表示外部文件
type External struct {
	URL string `json:"url"`
}

// FileInfo 表示文件信息
type FileInfo struct {
	URL        string `json:"url"`
	ExpiryTime string `json:"expiry_time"`
}

// Icon 表示图标
type Icon struct {
	Type     string    `json:"type"`               // "emoji" 或 "external" 或 "file"
	Emoji    string    `json:"emoji,omitempty"`    // 当 Type 为 "emoji" 时
	External *External `json:"external,omitempty"` // 当 Type 为 "external" 时
	File     *FileInfo `json:"file,omitempty"`     // 当 Type 为 "file" 时
}

// SelectOptions 表示选择选项配置
type SelectOptions struct {
	Options []Option `json:"options"`
}

// Option 表示选择选项
type Option struct {
	Name  string `json:"name"`
	Color Color  `json:"color,omitempty"`
}

// EmptyObject 表示空对象
type EmptyObject struct{}

// NumberConfig 表示数字属性配置
type NumberConfig struct {
	Format string `json:"format"` // 数字格式
}

// SelectConfig 表示选择属性配置
type SelectConfig struct {
	Options []Option `json:"options"` // 选项列表
}

// FormulaConfig 表示公式属性配置
type FormulaConfig struct {
	Expression string `json:"expression"` // 公式表达式
}

// RelationConfig 表示关联属性配置
type RelationConfig struct {
	DatabaseID         string `json:"database_id"`          // 关联的数据库 ID
	SyncedPropertyID   string `json:"synced_property_id"`   // 同步的属性 ID
	SyncedPropertyName string `json:"synced_property_name"` // 同步的属性名称
}

// RollupConfig 表示汇总属性配置
type RollupConfig struct {
	RelationPropertyName string `json:"relation_property_name"` // 关联属性名称
	RelationPropertyID   string `json:"relation_property_id"`   // 关联属性 ID
	RollupPropertyName   string `json:"rollup_property_name"`   // 汇总属性名称
	RollupPropertyID     string `json:"rollup_property_id"`     // 汇总属性 ID
	Function             string `json:"function"`               // 汇总函数
}

// StatusConfig 表示状态属性配置
type StatusConfig struct {
	Options []StatusOption `json:"options"` // 状态选项列表
	Groups  []StatusGroup  `json:"groups"`  // 状态组列表
}

// StatusOption 表示状态选项
type StatusOption struct {
	ID    string `json:"id"`    // 选项 ID
	Name  string `json:"name"`  // 选项名称
	Color Color  `json:"color"` // 选项颜色
}

// StatusGroup 表示状态组
type StatusGroup struct {
	ID      string   `json:"id"`      // 组 ID
	Name    string   `json:"name"`    // 组名称
	Color   string   `json:"color"`   // 组颜色
	Options []string `json:"options"` // 组中的选项 ID 列表
}

// Property 表示数据库属性
type Property struct {
	ID             string          `json:"id"`                         // 属性 ID
	Type           string          `json:"type"`                       // 属性类型
	Name           string          `json:"name"`                       // 属性名称
	Title          *EmptyObject    `json:"title,omitempty"`            // 标题属性
	RichText       *EmptyObject    `json:"rich_text,omitempty"`        // 富文本属性
	Number         *NumberConfig   `json:"number,omitempty"`           // 数字属性
	Select         *SelectConfig   `json:"select,omitempty"`           // 选择属性
	MultiSelect    *SelectConfig   `json:"multi_select,omitempty"`     // 多选属性
	Date           *EmptyObject    `json:"date,omitempty"`             // 日期属性
	People         *EmptyObject    `json:"people,omitempty"`           // 人员属性
	Files          *EmptyObject    `json:"files,omitempty"`            // 文件属性
	Checkbox       *EmptyObject    `json:"checkbox,omitempty"`         // 复选框属性
	URL            *EmptyObject    `json:"url,omitempty"`              // URL 属性
	Email          *EmptyObject    `json:"email,omitempty"`            // 邮箱属性
	PhoneNumber    *EmptyObject    `json:"phone_number,omitempty"`     // 电话属性
	Formula        *FormulaConfig  `json:"formula,omitempty"`          // 公式属性
	Relation       *RelationConfig `json:"relation,omitempty"`         // 关联属性
	Rollup         *RollupConfig   `json:"rollup,omitempty"`           // 汇总属性
	CreatedTime    *EmptyObject    `json:"created_time,omitempty"`     // 创建时间属性
	CreatedBy      *EmptyObject    `json:"created_by,omitempty"`       // 创建者属性
	LastEditedTime *EmptyObject    `json:"last_edited_time,omitempty"` // 最后编辑时间属性
	LastEditedBy   *EmptyObject    `json:"last_edited_by,omitempty"`   // 最后编辑者属性
	Status         *StatusConfig   `json:"status,omitempty"`           // 状态属性
	Unique         *EmptyObject    `json:"unique,omitempty"`           // 唯一属性
}

// ListParams 表示列出资源的参数
type ListParams struct {
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

// ListResponse 表示列出资源的响应
type ListResponse struct {
	Results    []interface{} `json:"results"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// Comment 表示评论
type Comment struct {
	ID             string          `json:"id"`
	ParentID       string          `json:"parent_id"`
	ParentType     string          `json:"parent_type"`
	CreatedTime    string          `json:"created_time"`
	LastEditedTime string          `json:"last_edited_time"`
	CreatedBy      json.RawMessage `json:"created_by"`
	RichText       []RichText      `json:"rich_text"`
	Resolved       bool            `json:"resolved"`
	Discussion     *Discussion     `json:"discussion,omitempty"`
}

// Discussion 表示评论讨论
type Discussion struct {
	ID       string `json:"id"`
	ThreadID string `json:"thread_id"`
}

// CreateCommentParams 表示创建评论的参数
type CreateCommentParams struct {
	ParentID   string      `json:"parent_id"`
	ParentType string      `json:"parent_type"`
	RichText   []RichText  `json:"rich_text"`
	Discussion *Discussion `json:"discussion,omitempty"`
}
