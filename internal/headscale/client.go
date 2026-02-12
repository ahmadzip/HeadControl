package headscale

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"headcontrol/internal/model"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return data, resp.StatusCode, nil
}

func (c *Client) doGet(path string) ([]byte, error) {
	data, status, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, c.parseError(data, status)
	}
	return data, nil
}

func (c *Client) doPost(path string, body interface{}) ([]byte, error) {
	data, status, err := c.doRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, c.parseError(data, status)
	}
	return data, nil
}

func (c *Client) doDelete(path string) error {
	_, status, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("delete failed with status %d", status)
	}
	return nil
}

func (c *Client) parseError(data []byte, status int) error {
	var apiErr model.ErrorResponse
	if json.Unmarshal(data, &apiErr) == nil && apiErr.Message != "" {
		return fmt.Errorf("API error (%d): %s", apiErr.Code, apiErr.Message)
	}
	return fmt.Errorf("unexpected status %d: %s", status, string(data))
}

func (c *Client) TestConnection() error {
	data, status, err := c.doRequest("GET", "/api/v1/user", nil)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	body := strings.TrimSpace(string(data))
	log.Printf("[TestConnection] status=%d body=%q url=%s", status, body, c.BaseURL)

	if status == 401 || status == 403 {
		if body != "" && body != "Unauthorized" {
			return fmt.Errorf("authentication failed (HTTP %d): %s", status, body)
		}
		return fmt.Errorf("authentication failed (HTTP %d): the API key was rejected by the server", status)
	}
	if status != 200 {
		return c.parseError(data, status)
	}
	return nil
}

func (c *Client) ListUsers() ([]model.User, error) {
	data, err := c.doGet("/api/v1/user")
	if err != nil {
		return nil, err
	}
	var resp model.UsersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}
	return resp.Users, nil
}

func (c *Client) CreateUser(name, displayName, email, pictureURL string) (*model.User, error) {
	data, err := c.doPost("/api/v1/user", map[string]string{
		"name":        name,
		"displayName": displayName,
		"email":       email,
		"pictureUrl":  pictureURL,
	})
	if err != nil {
		return nil, err
	}
	var resp model.UserResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode user: %w", err)
	}
	return &resp.User, nil
}

func (c *Client) RenameUser(oldID, newName string) (*model.User, error) {
	data, err := c.doPost(fmt.Sprintf("/api/v1/user/%s/rename/%s", oldID, newName), nil)
	if err != nil {
		return nil, err
	}
	var resp model.UserResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode user: %w", err)
	}
	return &resp.User, nil
}

func (c *Client) DeleteUser(id string) error {
	return c.doDelete(fmt.Sprintf("/api/v1/user/%s", id))
}

func (c *Client) ListNodes() ([]model.Node, error) {
	data, err := c.doGet("/api/v1/node")
	if err != nil {
		return nil, err
	}
	var resp model.NodesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode nodes: %w", err)
	}
	return resp.Nodes, nil
}

func (c *Client) GetNode(nodeID string) (*model.Node, error) {
	data, err := c.doGet(fmt.Sprintf("/api/v1/node/%s", nodeID))
	if err != nil {
		return nil, err
	}
	var resp model.NodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode node: %w", err)
	}
	return &resp.Node, nil
}

func (c *Client) RenameNode(nodeID, newName string) (*model.Node, error) {
	data, err := c.doPost(fmt.Sprintf("/api/v1/node/%s/rename/%s", nodeID, newName), nil)
	if err != nil {
		return nil, err
	}
	var resp model.NodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode node: %w", err)
	}
	return &resp.Node, nil
}

func (c *Client) ExpireNode(nodeID string) (*model.Node, error) {
	data, err := c.doPost(fmt.Sprintf("/api/v1/node/%s/expire", nodeID), nil)
	if err != nil {
		return nil, err
	}
	var resp model.NodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode node: %w", err)
	}
	return &resp.Node, nil
}

func (c *Client) DeleteNode(nodeID string) error {
	return c.doDelete(fmt.Sprintf("/api/v1/node/%s", nodeID))
}

func (c *Client) SetNodeTags(nodeID string, tags []string) (*model.Node, error) {
	data, err := c.doPost(fmt.Sprintf("/api/v1/node/%s/tags", nodeID), map[string][]string{"tags": tags})
	if err != nil {
		return nil, err
	}
	var resp model.NodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode node: %w", err)
	}
	return &resp.Node, nil
}

func (c *Client) SetApprovedRoutes(nodeID string, routes []string) (*model.Node, error) {
	data, err := c.doPost(fmt.Sprintf("/api/v1/node/%s/approve_routes", nodeID), map[string][]string{"routes": routes})
	if err != nil {
		return nil, err
	}
	var resp model.NodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode node: %w", err)
	}
	return &resp.Node, nil
}
