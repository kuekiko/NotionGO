package database

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuekiko/NotionGO/client"
)

// testClient 创建一个用于测试的客户端，返回客户端和测试服务器
func testClient(t *testing.T, fn http.HandlerFunc) (*client.Client, *httptest.Server) {
	server := httptest.NewServer(fn)
	c := client.NewClient("test-api-key")
	c.SetHTTPClient(&http.Client{
		Transport: &testTransport{
			server:        server,
			baseTransport: http.DefaultTransport,
		},
	})
	return c, server
}

// testTransport 是一个自定义的 http.RoundTripper，用于重写请求 URL
type testTransport struct {
	server        *httptest.Server
	baseTransport http.RoundTripper
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.server.Listener.Addr().String()
	return t.baseTransport.RoundTrip(req)
}

func TestDatabaseGet(t *testing.T) {
	expectedDB := &Database{
		ID:             "test-db-id",
		CreatedTime:    "2021-01-01T00:00:00.000Z",
		LastEditedTime: "2021-01-01T00:00:00.000Z",
		Title: []RichText{
			{
				Type:      "text",
				PlainText: "Test Database",
				Text: Text{
					Content: "Test Database",
				},
			},
		},
		Properties: map[string]Property{
			"Name": {
				ID:   "title",
				Type: "title",
				Name: "Name",
			},
		},
	}

	c, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/databases/test-db-id" {
			t.Errorf("Expected path /v1/databases/test-db-id, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(expectedDB)
	})
	defer server.Close()

	service := NewService(c)

	db, err := service.Get(context.Background(), "test-db-id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if db.ID != expectedDB.ID {
		t.Errorf("Expected database ID %s, got %s", expectedDB.ID, db.ID)
	}
}

func TestDatabaseQuery(t *testing.T) {
	expectedResp := &QueryResponse{
		Results: []Page{
			{
				ID:             "page-id",
				CreatedTime:    "2021-01-01T00:00:00.000Z",
				LastEditedTime: "2021-01-01T00:00:00.000Z",
				Properties: map[string]interface{}{
					"Name": map[string]interface{}{
						"title": []interface{}{
							map[string]interface{}{
								"text": map[string]interface{}{
									"content": "Test Page",
								},
							},
						},
					},
				},
			},
		},
		HasMore:    false,
		NextCursor: "",
	}

	c, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/databases/test-db-id/query" {
			t.Errorf("Expected path /v1/databases/test-db-id/query, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		var params QueryParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	service := NewService(c)

	params := &QueryParams{
		PageSize: 10,
		Filter: map[string]interface{}{
			"property": "Name",
			"text": map[string]interface{}{
				"contains": "Test",
			},
		},
	}

	resp, err := service.Query(context.Background(), "test-db-id", params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Results) != len(expectedResp.Results) {
		t.Errorf("Expected %d results, got %d", len(expectedResp.Results), len(resp.Results))
	}

	if resp.Results[0].ID != expectedResp.Results[0].ID {
		t.Errorf("Expected page ID %s, got %s", expectedResp.Results[0].ID, resp.Results[0].ID)
	}
}
