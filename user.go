package notion

// UserService 表示用户服务
type UserService struct {
	client *Client
}

// NewUserService 创建用户服务
func NewUserService(client *Client) *UserService {
	return &UserService{client: client}
}

// Get 获取用户
func (s *UserService) Get(userID string) (*User, error) {
	path := "users/" + userID
	user := new(User)
	err := s.client.get(path, nil, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List 列出用户
func (s *UserService) List(params *ListParams) (*ListResponse, error) {
	response := new(ListResponse)
	err := s.client.get("users", params, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Me 获取当前用户
func (s *UserService) Me() (*User, error) {
	user := new(User)
	err := s.client.get("users/me", nil, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
