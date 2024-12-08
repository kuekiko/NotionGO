package notion

// SearchService 表示搜索服务
type SearchService struct {
	client *Client
}

// NewSearchService 创建搜索服务
func NewSearchService(client *Client) *SearchService {
	return &SearchService{client: client}
}

// SearchParams 表示搜索参数
type SearchParams struct {
	Query       string        `json:"query"`                  // 搜索关键词
	Filter      *SearchFilter `json:"filter,omitempty"`       // 过滤条件
	Sort        *SearchSort   `json:"sort,omitempty"`         // 排序条件
	StartCursor string        `json:"start_cursor,omitempty"` // 起始游标
	PageSize    int           `json:"page_size,omitempty"`    // 页面大小
}

// SearchFilter 表示搜索过滤条件
type SearchFilter struct {
	Value      string `json:"value"`                 // 过滤值
	Property   string `json:"property"`              // 属性
	Timestamps string `json:"timestamps,omitempty"`  // 时间戳
	IsDatabase bool   `json:"is_database,omitempty"` // 是否为数据库
	IsPage     bool   `json:"is_page,omitempty"`     // 是否为页面
}

// SearchSort 表示搜索排序条件
type SearchSort struct {
	Direction string `json:"direction"`           // 排序方向："ascending" 或 "descending"
	Timestamp string `json:"timestamp,omitempty"` // 时间戳字段："created_time" 或 "last_edited_time"
}

// Search 搜索
func (s *SearchService) Search(params *SearchParams) (*ListResponse, error) {
	response := new(ListResponse)
	err := s.client.post("search", params, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SearchBuilder 表示搜索构建器
type SearchBuilder struct {
	params SearchParams
}

// NewSearchBuilder 创建搜索构建器
func NewSearchBuilder() *SearchBuilder {
	return &SearchBuilder{}
}

// Query 设置搜索关键词
func (b *SearchBuilder) Query(query string) *SearchBuilder {
	b.params.Query = query
	return b
}

// Filter 设置过滤条件
func (b *SearchBuilder) Filter(filter *SearchFilter) *SearchBuilder {
	b.params.Filter = filter
	return b
}

// Sort 设置排序条件
func (b *SearchBuilder) Sort(sort *SearchSort) *SearchBuilder {
	b.params.Sort = sort
	return b
}

// StartCursor 设置起始游标
func (b *SearchBuilder) StartCursor(cursor string) *SearchBuilder {
	b.params.StartCursor = cursor
	return b
}

// PageSize 设置页面大小
func (b *SearchBuilder) PageSize(size int) *SearchBuilder {
	b.params.PageSize = size
	return b
}

// Build 构建搜索参数
func (b *SearchBuilder) Build() *SearchParams {
	return &b.params
}
