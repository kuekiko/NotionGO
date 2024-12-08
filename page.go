package notion

// Page 表示页面对象
type Page struct {
	Object         string                 `json:"object"`           // 总是 "page"
	ID             string                 `json:"id"`               // 页面 ID
	CreatedTime    string                 `json:"created_time"`     // 创建时间
	LastEditedTime string                 `json:"last_edited_time"` // 最后编辑时间
	CreatedBy      User                   `json:"created_by"`       // 创建者
	LastEditedBy   User                   `json:"last_edited_by"`   // 最后编辑者
	Parent         Parent                 `json:"parent"`           // 父对象
	Archived       bool                   `json:"archived"`         // 是否已归档
	Properties     map[string]interface{} `json:"properties"`       // 属性
	URL            string                 `json:"url"`              // URL
	Icon           *Icon                  `json:"icon,omitempty"`   // 图标
	Cover          *File                  `json:"cover,omitempty"`  // 封面
}

// PageService 表示页面服务
type PageService struct {
	client *Client
}

// NewPageService 创建页面服务
func NewPageService(client *Client) *PageService {
	return &PageService{client: client}
}

// PageCreateParams 表示创建页面的参数
type PageCreateParams struct {
	Parent     Parent                 `json:"parent"`             // 父对象
	Properties map[string]interface{} `json:"properties"`         // 属性
	Children   []Block                `json:"children,omitempty"` // 子块
	Icon       *Icon                  `json:"icon,omitempty"`     // 图标
	Cover      *File                  `json:"cover,omitempty"`    // 封面
}

// Create 创建页面
func (s *PageService) Create(params *PageCreateParams) (*Page, error) {
	page := new(Page)
	err := s.client.post("pages", params, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// Get 获取页面
func (s *PageService) Get(pageID string) (*Page, error) {
	path := "pages/" + pageID
	page := new(Page)
	err := s.client.get(path, nil, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// Update 更新页面
func (s *PageService) Update(pageID string, params *Page) (*Page, error) {
	path := "pages/" + pageID
	page := new(Page)
	err := s.client.patch(path, params, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// Delete 删除页面
func (s *PageService) Delete(pageID string) error {
	path := "pages/" + pageID
	return s.client.delete(path)
}

// PropertyItem 表示属性项
type PropertyItem struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	NextURL  string      `json:"next_url,omitempty"`
	Property interface{} `json:"property"`
}

// GetProperty 获取页面属性
func (s *PageService) GetProperty(pageID, propertyID string) (*PropertyItem, error) {
	path := "pages/" + pageID + "/properties/" + propertyID
	property := new(PropertyItem)
	err := s.client.get(path, nil, property)
	if err != nil {
		return nil, err
	}
	return property, nil
}

// GetPropertyList 获取页面属性列表
func (s *PageService) GetPropertyList(pageID string, params *ListParams) (*ListResponse, error) {
	path := "pages/" + pageID + "/properties"
	response := new(ListResponse)
	err := s.client.get(path, params, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
