package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/menloresearch/cli/internal/config"
)

var ErrNoAPIKey = errors.New("API key not set. Run 'menlo config apikey' to set it")

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		if config.IsNotExist(err) {
			cfg = config.DefaultConfig()
		} else {
			return nil, err
		}
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    cfg.PlatformURL,
		apiKey:     cfg.APIKey,
	}, nil
}

// Battery represents battery status
type Battery struct {
	Level  int  `json:"level"`
	Charging bool `json:"charging"`
}

// RobotStatus represents the robot's current status
type RobotStatus struct {
	RobotID   string  `json:"robot_id"`
	Battery   Battery `json:"battery"`
	Timestamp int64   `json:"timestamp"`
}

// LocalTimestamp returns the timestamp converted to local time
func (r *RobotStatus) LocalTimestamp() time.Time {
	return time.Unix(r.Timestamp, 0).Local()
}

// RobotResponse represents a robot in the API response
type RobotResponse struct {
	ID     string       `json:"id"`
	Model  string       `json:"model"`
	Type   string       `json:"type"`
	Name   string       `json:"name"`
	Status *RobotStatus `json:"status,omitempty"`
}

// ListResponse is a generic list response
type ListResponse[T any] struct {
	Code   string          `json:"code"`
	Result json.RawMessage `json:"result"`
	NextID *string         `json:"next_id"`
	Total  int64           `json:"total"`
}

func (r *ListResponse[T]) ParseResult(target *[]T) error {
	if len(r.Result) == 0 {
		*target = nil
		return nil
	}
	return json.Unmarshal(r.Result, target)
}

// GeneralResponse is a generic single item response
type GeneralResponse[T any] struct {
	Code   string          `json:"code"`
	Result json.RawMessage `json:"result"`
}

func (r *GeneralResponse[T]) ParseResult(target *T) error {
	if len(r.Result) == 0 {
		return nil
	}
	return json.Unmarshal(r.Result, target)
}

// doRequest makes an authenticated HTTP request
func (c *Client) doRequest(method, path string, queryParams map[string]string) (*http.Response, error) {
	return c.doRequestBody(method, path, queryParams, nil)
}

// doRequestBody makes an authenticated HTTP request with a body
func (c *Client) doRequestBody(method, path string, queryParams map[string]string, body []byte) (*http.Response, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("%w. Get your API key from: https://platform.menlo.ai/account/api-keys", ErrNoAPIKey)
	}

	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	reqURL.Path = path

	q := reqURL.Query()
	for k, v := range queryParams {
		if v != "" {
			q.Set(k, v)
		}
	}
	if len(q) > 0 {
		reqURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, reqURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if !isOK(resp.StatusCode) {
		return resp, getErrorMessage(resp.StatusCode)
	}

	return resp, nil
}

func isOK(code int) bool {
	return code >= 200 && code < 300
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

func getErrorMessage(statusCode int) error {
	switch statusCode {
	case 401:
		return fmt.Errorf("invalid API key, visit https://platform.menlo.ai/account/api-keys to generate a new one")
	case 403:
		return fmt.Errorf("access forbidden")
	case 404:
		return fmt.Errorf("robot not found or you don't have access to it")
	case 500:
		return fmt.Errorf("internal server error")
	default:
		return fmt.Errorf("request failed with status: %d", statusCode)
	}
}

// ListRobots fetches a list of robots
func (c *Client) ListRobots(limit int, afterPublicID string) (*ListResponse[RobotResponse], error) {
	resp, err := c.doRequest("GET", "v1/robots", map[string]string{
		"limit":           fmt.Sprintf("%d", limit),
		"after_public_id": afterPublicID,
	})
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	var result ListResponse[RobotResponse]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ListResponse[RobotResponse]) Robots() ([]RobotResponse, error) {
	var robots []RobotResponse
	err := r.ParseResult(&robots)
	return robots, err
}

// GetRobot fetches a robot by ID
func (c *Client) GetRobot(robotID string) (*RobotResponse, error) {
	resp, err := c.doRequest("GET", "v1/robots/"+robotID, nil)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	var result GeneralResponse[RobotResponse]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var robot RobotResponse
	if err := result.ParseResult(&robot); err != nil {
		return nil, err
	}

	return &robot, nil
}

// SessionResponse represents the WebRTC session info
type SessionResponse struct {
	SFUEndpoint string `json:"sfu_endpoint"`
	WebRTCToken string `json:"webrtc_token"`
}

// ValidSemanticCommands are the supported semantic commands
var ValidSemanticCommands = []string{
	"forward",
	"backward",
	"left",
	"right",
	"turn-left",
	"turn-right",
}

// SendSemanticCommand sends a semantic command to a robot
func (c *Client) SendSemanticCommand(robotID, command string) error {
	body, err := json.Marshal(map[string]string{"command": command})
	if err != nil {
		return err
	}

	resp, err := c.doRequestBody("POST", "v1/robots/"+robotID+"/semantic-command", nil, body)
	if err != nil {
		return err
	}
	defer closeBody(resp)

	return nil
}

// GetSnapshot downloads the latest snapshot image for a robot
func (c *Client) GetSnapshot(robotID string) (string, error) {
	resp, err := c.doRequest("GET", "v1/robots/"+robotID+"/snapshot", nil)
	if err != nil {
		return "", err
	}
	defer closeBody(resp)

	// Get config dir for snapshot storage
	configDir, err := config.ConfigDir()
	if err != nil {
		return "", err
	}
	snapshotDir := filepath.Join(configDir, "snapshot", robotID)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return "", err
	}

	// Save the image
	imagePath := filepath.Join(snapshotDir, "latest.jpeg")
	outFile, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = outFile.Close() }()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

// CreateSession creates a new session for a robot and returns WebRTC credentials
func (c *Client) CreateSession(robotID string) (*SessionResponse, error) {
	resp, err := c.doRequest("POST", "v1/robots/"+robotID+"/session", nil)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	var result GeneralResponse[SessionResponse]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var session SessionResponse
	if err := result.ParseResult(&session); err != nil {
		return nil, err
	}

	return &session, nil
}
