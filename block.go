package notion

// Block 表示块对象
type Block struct {
	Object         string    `json:"object"`           // 总是 "block"
	ID             string    `json:"id"`               // 块 ID
	Parent         Parent    `json:"parent"`           // 父对象
	Type           BlockType `json:"type"`             // 块类型
	CreatedTime    string    `json:"created_time"`     // 创建时间
	LastEditedTime string    `json:"last_edited_time"` // 最后编辑时间
	CreatedBy      User      `json:"created_by"`       // 创建者
	LastEditedBy   User      `json:"last_edited_by"`   // 最后编辑者
	HasChildren    bool      `json:"has_children"`     // 是否有子块
	Archived       bool      `json:"archived"`         // 是否已归档

	// 不同类型的块具有不同的属性
	Paragraph        *ParagraphBlock     `json:"paragraph,omitempty"`
	Heading1         *HeadingBlock       `json:"heading_1,omitempty"`
	Heading2         *HeadingBlock       `json:"heading_2,omitempty"`
	Heading3         *HeadingBlock       `json:"heading_3,omitempty"`
	BulletedListItem *ListItemBlock      `json:"bulleted_list_item,omitempty"`
	NumberedListItem *ListItemBlock      `json:"numbered_list_item,omitempty"`
	ToDo             *ToDoBlock          `json:"to_do,omitempty"`
	Toggle           *ToggleBlock        `json:"toggle,omitempty"`
	ChildPage        *ChildPageBlock     `json:"child_page,omitempty"`
	ChildDatabase    *ChildDatabaseBlock `json:"child_database,omitempty"`
	Embed            *EmbedBlock         `json:"embed,omitempty"`
	Image            *File               `json:"image,omitempty"`
	Video            *File               `json:"video,omitempty"`
	File             *File               `json:"file,omitempty"`
	PDF              *File               `json:"pdf,omitempty"`
	Bookmark         *BookmarkBlock      `json:"bookmark,omitempty"`
	Callout          *CalloutBlock       `json:"callout,omitempty"`
	Quote            *QuoteBlock         `json:"quote,omitempty"`
	Equation         *EquationBlock      `json:"equation,omitempty"`
	Divider          *EmptyObject        `json:"divider,omitempty"`
	TableOfContents  *EmptyObject        `json:"table_of_contents,omitempty"`
	Breadcrumb       *EmptyObject        `json:"breadcrumb,omitempty"`
	ColumnList       *ColumnListBlock    `json:"column_list,omitempty"`
	Column           *ColumnBlock        `json:"column,omitempty"`
	LinkPreview      *LinkPreviewBlock   `json:"link_preview,omitempty"`
	Template         *TemplateBlock      `json:"template,omitempty"`
	SyncedBlock      *SyncedBlock        `json:"synced_block,omitempty"`
	Table            *TableBlock         `json:"table,omitempty"`
	TableRow         *TableRowBlock      `json:"table_row,omitempty"`
	Code             *CodeBlock          `json:"code,omitempty"`
}

// ParagraphBlock 表示段落块
type ParagraphBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children,omitempty"`
}

// HeadingBlock 表示标题块
type HeadingBlock struct {
	RichText     []RichText `json:"rich_text"`
	Color        Color      `json:"color"`
	IsToggleable bool       `json:"is_toggleable"`
}

// ListItemBlock 表示列表项块
type ListItemBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children,omitempty"`
}

// ToDoBlock 表示待办事项块
type ToDoBlock struct {
	RichText []RichText `json:"rich_text"`
	Checked  bool       `json:"checked"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children,omitempty"`
}

// ToggleBlock 表示折叠块
type ToggleBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children,omitempty"`
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
	Children []Block    `json:"children,omitempty"`
}

// QuoteBlock 表示引用块
type QuoteBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children,omitempty"`
}

// EquationBlock 表示公式块
type EquationBlock struct {
	Expression string `json:"expression"`
}

// ColumnListBlock 表示列表块
type ColumnListBlock struct {
	Children []Block `json:"children"`
}

// ColumnBlock 表示列块
type ColumnBlock struct {
	Children []Block `json:"children"`
}

// LinkPreviewBlock 表示链接预览块
type LinkPreviewBlock struct {
	URL string `json:"url"`
}

// TemplateBlock 表示模板块
type TemplateBlock struct {
	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children"`
}

// SyncedBlock 表示同步块
type SyncedBlock struct {
	SyncedFrom *SyncedFrom `json:"synced_from"`
	Children   []Block     `json:"children,omitempty"`
}

// SyncedFrom 表示同步来源
type SyncedFrom struct {
	BlockID string `json:"block_id"`
}

// TableBlock 表示表格块
type TableBlock struct {
	TableWidth      int     `json:"table_width"`
	HasColumnHeader bool    `json:"has_column_header"`
	HasRowHeader    bool    `json:"has_row_header"`
	Children        []Block `json:"children"`
}

// TableRowBlock 表示表格行块
type TableRowBlock struct {
	Cells [][]RichText `json:"cells"`
}

// CodeBlock 表示代码块
type CodeBlock struct {
	RichText []RichText `json:"rich_text"`
	Caption  []RichText `json:"caption"`
	Language string     `json:"language"`
}

// BlockService 表示块服务
type BlockService struct {
	client *Client
}

// NewBlockService 创建块服务
func NewBlockService(client *Client) *BlockService {
	return &BlockService{client: client}
}

// Get 获取块
func (s *BlockService) Get(blockID string) (*Block, error) {
	path := "blocks/" + blockID
	block := new(Block)
	err := s.client.get(path, nil, block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// Update 更新块
func (s *BlockService) Update(blockID string, block *Block) (*Block, error) {
	path := "blocks/" + blockID
	response := new(Block)
	err := s.client.patch(path, block, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Delete 删除块
func (s *BlockService) Delete(blockID string) error {
	path := "blocks/" + blockID
	return s.client.delete(path)
}

// ListChildren 列出子块
func (s *BlockService) ListChildren(blockID string, params *ListParams) (*ListResponse, error) {
	path := "blocks/" + blockID + "/children"
	response := new(ListResponse)
	err := s.client.get(path, nil, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AppendChildren 追加子块
func (s *BlockService) AppendChildren(blockID string, children []Block) (*ListResponse, error) {
	path := "blocks/" + blockID + "/children"
	response := new(ListResponse)
	err := s.client.patch(path, map[string]interface{}{
		"children": children,
	}, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
