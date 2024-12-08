package notion

import (
	"github.com/kuekiko/NotionGO/client"
)

// Client 表示 Notion API 客户端
type Client struct {
	client *client.Client

	// 服务
	Blocks   *BlockService
	Pages    *PageService
	Database *DatabaseService
	Users    *UserService
	Search   *SearchService
	Comments *CommentService
}

// NewClient 创建 Notion API 客户端
func NewClient(token string) *Client {
	c := &Client{
		client: client.NewClient(token),
	}

	c.Blocks = NewBlockService(c)
	c.Pages = NewPageService(c)
	c.Database = NewDatabaseService(c)
	c.Users = NewUserService(c)
	c.Search = NewSearchService(c)
	c.Comments = NewCommentService(c)

	return c
}

// get 发送 GET 请求
func (c *Client) get(path string, params interface{}, v interface{}) error {
	return c.client.Get(path, params, v)
}

// post 发送 POST 请求
func (c *Client) post(path string, body interface{}, v interface{}) error {
	return c.client.Post(path, body, v)
}

// patch 发送 PATCH 请求
func (c *Client) patch(path string, body interface{}, v interface{}) error {
	return c.client.Patch(path, body, v)
}

// delete 发送 DELETE 请求
func (c *Client) delete(path string) error {
	return c.client.Delete(path)
}
