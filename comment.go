package notion

// CommentService 表示评论服务
type CommentService struct {
	client *Client
}

// NewCommentService 创建评论服务
func NewCommentService(client *Client) *CommentService {
	return &CommentService{client: client}
}

// Create 创建评论
func (s *CommentService) Create(params *CreateCommentParams) (*Comment, error) {
	comment := new(Comment)
	err := s.client.post("comments", params, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// List 列出评论
func (s *CommentService) List(blockID string, params *ListParams) (*ListResponse, error) {
	path := "comments"
	if blockID != "" {
		path += "?block_id=" + blockID
	}
	response := new(ListResponse)
	err := s.client.get(path, params, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
