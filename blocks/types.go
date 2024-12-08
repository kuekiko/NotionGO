package blocks

// BlockType 表示块类型
type BlockType string

const (
	TypeParagraph       BlockType = "paragraph"
	TypeHeading1        BlockType = "heading_1"
	TypeHeading2        BlockType = "heading_2"
	TypeHeading3        BlockType = "heading_3"
	TypeBulletedList    BlockType = "bulleted_list_item"
	TypeNumberedList    BlockType = "numbered_list_item"
	TypeToDo            BlockType = "to_do"
	TypeToggle          BlockType = "toggle"
	TypeChildPage       BlockType = "child_page"
	TypeChildDatabase   BlockType = "child_database"
	TypeEmbed           BlockType = "embed"
	TypeImage           BlockType = "image"
	TypeVideo           BlockType = "video"
	TypeFile            BlockType = "file"
	TypePDF             BlockType = "pdf"
	TypeBookmark        BlockType = "bookmark"
	TypeCallout         BlockType = "callout"
	TypeQuote           BlockType = "quote"
	TypeEquation        BlockType = "equation"
	TypeDivider         BlockType = "divider"
	TypeTableOfContents BlockType = "table_of_contents"
	TypeBreadcrumb      BlockType = "breadcrumb"
	TypeColumnList      BlockType = "column_list"
	TypeColumn          BlockType = "column"
	TypeLinkPreview     BlockType = "link_preview"
	TypeTemplate        BlockType = "template"
	TypeSyncedBlock     BlockType = "synced_block"
	TypeTable           BlockType = "table"
	TypeTableRow        BlockType = "table_row"
	TypeCode            BlockType = "code"
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
)

// Language 表示代码块的语言
type Language string

const (
	LangPlainText    Language = "plain text"
	LangABAP         Language = "abap"
	LangArduino      Language = "arduino"
	LangBash         Language = "bash"
	LangBasic        Language = "basic"
	LangC            Language = "c"
	LangClojure      Language = "clojure"
	LangCoffeescript Language = "coffeescript"
	LangCPP          Language = "c++"
	LangCSharp       Language = "c#"
	LangCSS          Language = "css"
	LangDart         Language = "dart"
	LangDiff         Language = "diff"
	LangDocker       Language = "docker"
	LangElixir       Language = "elixir"
	LangElm          Language = "elm"
	LangErlang       Language = "erlang"
	LangFlow         Language = "flow"
	LangFortran      Language = "fortran"
	LangFSharp       Language = "f#"
	LangGherkin      Language = "gherkin"
	LangGLSL         Language = "glsl"
	LangGo           Language = "go"
	LangGraphQL      Language = "graphql"
	LangGroovy       Language = "groovy"
	LangHaskell      Language = "haskell"
	LangHTML         Language = "html"
	LangJava         Language = "java"
	LangJavaScript   Language = "javascript"
	LangJSON         Language = "json"
	LangJulia        Language = "julia"
	LangKotlin       Language = "kotlin"
	LangLatex        Language = "latex"
	LangLess         Language = "less"
	LangLisp         Language = "lisp"
	LangLivescript   Language = "livescript"
	LangLua          Language = "lua"
	LangMakefile     Language = "makefile"
	LangMarkdown     Language = "markdown"
	LangMatlab       Language = "matlab"
	LangMermaid      Language = "mermaid"
	LangNix          Language = "nix"
	LangObjectiveC   Language = "objective-c"
	LangOCaml        Language = "ocaml"
	LangPascal       Language = "pascal"
	LangPerl         Language = "perl"
	LangPHP          Language = "php"
	LangPLSQL        Language = "plsql"
	LangPowershell   Language = "powershell"
	LangPython       Language = "python"
	LangR            Language = "r"
	LangRacket       Language = "racket"
	LangReason       Language = "reason"
	LangRuby         Language = "ruby"
	LangRust         Language = "rust"
	LangSASS         Language = "sass"
	LangScala        Language = "scala"
	LangScheme       Language = "scheme"
	LangSCSS         Language = "scss"
	LangShell        Language = "shell"
	LangSQL          Language = "sql"
	LangSwift        Language = "swift"
	LangTypeScript   Language = "typescript"
	LangVBNet        Language = "vb.net"
	LangVerilog      Language = "verilog"
	LangVHDL         Language = "vhdl"
	LangVisualBasic  Language = "visual basic"
	LangWebAssembly  Language = "webassembly"
	LangXML          Language = "xml"
	LangYAML         Language = "yaml"
)

// Block 表示一个块
type Block struct {
	Object         string `json:"object"`
	ID             string `json:"id"`
	Type           string `json:"type"`
	CreatedTime    string `json:"created_time"`
	LastEditedTime string `json:"last_edited_time"`
	HasChildren    bool   `json:"has_children"`
	Archived       bool   `json:"archived"`

	// 块类型特定的内容
	Paragraph       *ParagraphBlock       `json:"paragraph,omitempty"`
	Heading1        *HeadingBlock         `json:"heading_1,omitempty"`
	Heading2        *HeadingBlock         `json:"heading_2,omitempty"`
	Heading3        *HeadingBlock         `json:"heading_3,omitempty"`
	BulletedList    *ListBlock            `json:"bulleted_list_item,omitempty"`
	NumberedList    *ListBlock            `json:"numbered_list_item,omitempty"`
	ToDo            *ToDoBlock            `json:"to_do,omitempty"`
	Toggle          *ToggleBlock          `json:"toggle,omitempty"`
	ChildPage       *ChildPageBlock       `json:"child_page,omitempty"`
	ChildDatabase   *ChildDatabaseBlock   `json:"child_database,omitempty"`
	Embed           *EmbedBlock           `json:"embed,omitempty"`
	Image           *FileBlock            `json:"image,omitempty"`
	Video           *FileBlock            `json:"video,omitempty"`
	File            *FileBlock            `json:"file,omitempty"`
	PDF             *FileBlock            `json:"pdf,omitempty"`
	Bookmark        *BookmarkBlock        `json:"bookmark,omitempty"`
	Callout         *CalloutBlock         `json:"callout,omitempty"`
	Quote           *QuoteBlock           `json:"quote,omitempty"`
	Equation        *EquationBlock        `json:"equation,omitempty"`
	Divider         *DividerBlock         `json:"divider,omitempty"`
	TableOfContents *TableOfContentsBlock `json:"table_of_contents,omitempty"`
	Breadcrumb      *BreadcrumbBlock      `json:"breadcrumb,omitempty"`
	ColumnList      *ColumnListBlock      `json:"column_list,omitempty"`
	Column          *ColumnBlock          `json:"column,omitempty"`
	LinkPreview     *LinkPreviewBlock     `json:"link_preview,omitempty"`
	Template        *TemplateBlock        `json:"template,omitempty"`
	SyncedBlock     *SyncedBlockBlock     `json:"synced_block,omitempty"`
	Table           *TableBlock           `json:"table,omitempty"`
	TableRow        *TableRowBlock        `json:"table_row,omitempty"`
	Code            *CodeBlock            `json:"code,omitempty"`
}

// RichText 表示富文本内容
type RichText struct {
	Type        string      `json:"type"`
	Text        *Text       `json:"text,omitempty"`
	Mention     *Mention    `json:"mention,omitempty"`
	Equation    *Equation   `json:"equation,omitempty"`
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
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Mention 表示提及
type Mention struct {
	Type     string           `json:"type"`
	User     *UserMention     `json:"user,omitempty"`
	Page     *PageMention     `json:"page,omitempty"`
	Database *DBMention       `json:"database,omitempty"`
	Date     *DateMention     `json:"date,omitempty"`
	Template *TemplateMention `json:"template_mention,omitempty"`
}

// UserMention 表示用户提及
type UserMention struct {
	ID string `json:"id"`
}

// PageMention 表示页面提及
type PageMention struct {
	ID string `json:"id"`
}

// DBMention 表示数据库提及
type DBMention struct {
	ID string `json:"id"`
}

// DateMention 表示日期提及
type DateMention struct {
	Start string `json:"start"`
	End   string `json:"end,omitempty"`
}

// TemplateMention 表示模板提及
type TemplateMention struct {
	Type       string `json:"type"`
	TemplateID string `json:"template_mention_date,omitempty"`
}

// Equation 表示公式
type Equation struct {
	Expression string `json:"expression"`
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

// FileBlock 表示文件块
type FileBlock struct {
	Type     string `json:"type"`
	URL      string `json:"url,omitempty"`
	ExpireAt string `json:"expire_time,omitempty"`
}

// Icon 表示图标
type Icon struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji,omitempty"`
	File  *File  `json:"file,omitempty"`
}

// File 表示文件
type File struct {
	Type     string `json:"type"`
	URL      string `json:"url,omitempty"`
	ExpireAt string `json:"expire_time,omitempty"`
}

// Parent 表示父对象
type Parent struct {
	Type       string `json:"type"`
	PageID     string `json:"page_id,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
	BlockID    string `json:"block_id,omitempty"`
	Workspace  bool   `json:"workspace,omitempty"`
}

// ParagraphBlock 表示段落块
type ParagraphBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

// HeadingBlock 表示标题块
type HeadingBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

// ListBlock 表示列表块
type ListBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

// ToDoBlock 表示待办事项块
type ToDoBlock struct {
	RichText []RichText `json:"rich_text"`
	Checked  bool       `json:"checked"`
	Color    Color      `json:"color"`
}

// ToggleBlock 表示折叠块
type ToggleBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

// ChildPageBlock 表示子页面块
type ChildPageBlock struct {
	Title string `json:"title"`
}

// ChildDatabaseBlock 表示子数据库块
type ChildDatabaseBlock struct {
	Title string `json:"title"`
}

// EmbedBlock 表示嵌入块
type EmbedBlock struct {
	URL string `json:"url"`
}

// BookmarkBlock 表示书签块
type BookmarkBlock struct {
	URL     string     `json:"url"`
	Caption []RichText `json:"caption"`
}

// CalloutBlock 表示标注块
type CalloutBlock struct {
	RichText []RichText `json:"rich_text"`
	Icon     *Icon      `json:"icon"`
	Color    Color      `json:"color"`
}

// QuoteBlock 表示引用块
type QuoteBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

// EquationBlock 表示公式块
type EquationBlock struct {
	Expression string `json:"expression"`
}

// DividerBlock 表示分割线块
type DividerBlock struct{}

// TableOfContentsBlock 表示目录块
type TableOfContentsBlock struct {
	Color Color `json:"color"`
}

// BreadcrumbBlock 表示面包屑块
type BreadcrumbBlock struct{}

// ColumnListBlock 表示列表块
type ColumnListBlock struct{}

// ColumnBlock 表示列块
type ColumnBlock struct{}

// LinkPreviewBlock 表示链接预览块
type LinkPreviewBlock struct {
	URL string `json:"url"`
}

// TemplateBlock 表示模板块
type TemplateBlock struct {
	RichText []RichText `json:"rich_text"`
}

// SyncedBlockBlock 表示同步块
type SyncedBlockBlock struct {
	SyncedFrom *SyncedFrom `json:"synced_from"`
}

// SyncedFrom 表示同步来源
type SyncedFrom struct {
	BlockID string `json:"block_id"`
}

// TableBlock 表示表格块
type TableBlock struct {
	TableWidth      int  `json:"table_width"`
	HasColumnHeader bool `json:"has_column_header"`
	HasRowHeader    bool `json:"has_row_header"`
}

// TableRowBlock 表示表格行块
type TableRowBlock struct {
	Cells [][]RichText `json:"cells"`
}

// CodeBlock 表示代码块
type CodeBlock struct {
	RichText []RichText `json:"rich_text"`
	Caption  []RichText `json:"caption"`
	Language Language   `json:"language"`
}

// BlockBuilder 用于构建块内容
type BlockBuilder struct {
	block Block
}

// NewBlockBuilder 创建新的块构建器
func NewBlockBuilder(blockType BlockType) *BlockBuilder {
	return &BlockBuilder{
		block: Block{
			Type: string(blockType),
		},
	}
}

// WithParagraph 添加段落内容
func (b *BlockBuilder) WithParagraph(text string) *BlockBuilder {
	b.block.Paragraph = &ParagraphBlock{
		RichText: []RichText{
			{
				Type: "text",
				Text: &Text{
					Content: text,
				},
			},
		},
	}
	return b
}

// WithHeading 添加标题内容
func (b *BlockBuilder) WithHeading(text string) *BlockBuilder {
	switch b.block.Type {
	case string(TypeHeading1):
		b.block.Heading1 = &HeadingBlock{
			RichText: []RichText{
				{
					Type: "text",
					Text: &Text{
						Content: text,
					},
				},
			},
		}
	case string(TypeHeading2):
		b.block.Heading2 = &HeadingBlock{
			RichText: []RichText{
				{
					Type: "text",
					Text: &Text{
						Content: text,
					},
				},
			},
		}
	case string(TypeHeading3):
		b.block.Heading3 = &HeadingBlock{
			RichText: []RichText{
				{
					Type: "text",
					Text: &Text{
						Content: text,
					},
				},
			},
		}
	}
	return b
}

// WithCode 添加代码块内容
func (b *BlockBuilder) WithCode(code string, language Language) *BlockBuilder {
	b.block.Code = &CodeBlock{
		RichText: []RichText{
			{
				Type: "text",
				Text: &Text{
					Content: code,
				},
			},
		},
		Language: language,
	}
	return b
}

// WithCallout 添加标注块内容
func (b *BlockBuilder) WithCallout(text string, icon *Icon, color Color) *BlockBuilder {
	b.block.Callout = &CalloutBlock{
		RichText: []RichText{
			{
				Type: "text",
				Text: &Text{
					Content: text,
				},
			},
		},
		Icon:  icon,
		Color: color,
	}
	return b
}

// Build 构建块
func (b *BlockBuilder) Build() Block {
	return b.block
}
