package notion

// Database 表示数据库对象
type Database struct {
	Object         string              `json:"object"`           // 总是 "database"
	ID             string              `json:"id"`               // 数据库 ID
	CreatedTime    string              `json:"created_time"`     // 创建时间
	LastEditedTime string              `json:"last_edited_time"` // 最后编辑时间
	CreatedBy      User                `json:"created_by"`       // 创建者
	LastEditedBy   User                `json:"last_edited_by"`   // 最后编辑者
	Title          []RichText          `json:"title"`            // 标题
	Description    []RichText          `json:"description"`      // 描述
	Icon           *Icon               `json:"icon,omitempty"`   // 图标
	Cover          *File               `json:"cover,omitempty"`  // 封面
	Properties     map[string]Property `json:"properties"`       // 属性
	Parent         Parent              `json:"parent"`           // 父对象
	URL            string              `json:"url"`              // URL
	Archived       bool                `json:"archived"`         // 是否已归档
	IsInline       bool                `json:"is_inline"`        // 是否内联
}

// DatabaseService 表示数据库服务
type DatabaseService struct {
	client *Client
}

// NewDatabaseService 创建数据库服务
func NewDatabaseService(client *Client) *DatabaseService {
	return &DatabaseService{client: client}
}

// DatabaseCreateParams 表示创建数据库的参数
type DatabaseCreateParams struct {
	Parent     Parent              `json:"parent"`              // 父对象
	Title      []RichText          `json:"title"`               // 标题
	Properties map[string]Property `json:"properties"`          // 属性
	Icon       *Icon               `json:"icon,omitempty"`      // 图标
	Cover      *File               `json:"cover,omitempty"`     // 封面
	IsInline   bool                `json:"is_inline,omitempty"` // 是否内联
}

// Create 创建数据库
func (s *DatabaseService) Create(params *DatabaseCreateParams) (*Database, error) {
	database := new(Database)
	err := s.client.post("databases", params, database)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// Get 获取数据库
func (s *DatabaseService) Get(databaseID string) (*Database, error) {
	path := "databases/" + databaseID
	database := new(Database)
	err := s.client.get(path, nil, database)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// Update 更新数据库
func (s *DatabaseService) Update(databaseID string, params *Database) (*Database, error) {
	path := "databases/" + databaseID
	database := new(Database)
	err := s.client.patch(path, params, database)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// Delete 删除数据库
func (s *DatabaseService) Delete(databaseID string) error {
	path := "databases/" + databaseID
	return s.client.delete(path)
}

// List 列出数据库
func (s *DatabaseService) List(params *ListParams) (*ListResponse, error) {
	response := new(ListResponse)
	err := s.client.get("databases", params, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// DatabaseQueryParams 表示查询数据库的参数
type DatabaseQueryParams struct {
	Filter      interface{} `json:"filter,omitempty"`       // 过滤条件
	Sorts       []Sort      `json:"sorts,omitempty"`        // 排序条件
	StartCursor string      `json:"start_cursor,omitempty"` // 起始游标
	PageSize    int         `json:"page_size,omitempty"`    // 页面大小
}

// Sort 表示排序条件
type Sort struct {
	Property  string `json:"property"`            // 属性名称
	Direction string `json:"direction"`           // 排序方向："ascending" 或 "descending"
	Timestamp string `json:"timestamp,omitempty"` // 时间戳字段："created_time" 或 "last_edited_time"
}

// Query 查询数据库
func (s *DatabaseService) Query(databaseID string, params *DatabaseQueryParams) (*ListResponse, error) {
	path := "databases/" + databaseID + "/query"
	response := new(ListResponse)
	err := s.client.post(path, params, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
