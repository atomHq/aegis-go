package aegis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "http://localhost:8080"

type Client struct {
	baseURL    string
	apiKey     string
	authToken  string
	httpClient *http.Client
	userAgent  string
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func WithAuthToken(token string) Option {
	return func(c *Client) {
		c.authToken = token
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		if userAgent != "" {
			c.userAgent = userAgent
		}
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "aegis-go/0.1.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type apiSuccess[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
	Meta   Meta   `json:"meta"`
}

type apiErrorEnvelope struct {
	Status string   `json:"status"`
	Error  APIError `json:"error"`
	Meta   Meta     `json:"meta"`
}

type Meta struct {
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

type APIError struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Status    int         `json:"-"`
	RequestID string      `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("aegis: %s (%s)", e.Message, e.Code)
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, authToken string, out any) error {
	var requestBody io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return err
		}
		requestBody = buf
	}

	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var envelope apiErrorEnvelope
		if err := json.Unmarshal(raw, &envelope); err == nil && envelope.Error.Code != "" {
			envelope.Error.Status = res.StatusCode
			envelope.Error.RequestID = envelope.Meta.RequestID
			return &envelope.Error
		}
		return &APIError{Code: "UNEXPECTED_RESPONSE", Message: strings.TrimSpace(string(raw)), Status: res.StatusCode}
	}

	if out == nil || len(raw) == 0 {
		return nil
	}

	var envelope apiSuccess[json.RawMessage]
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return err
	}
	if envelope.Status != "success" {
		return &APIError{Code: "UNEXPECTED_RESPONSE", Message: "unexpected response status", Status: res.StatusCode, RequestID: envelope.Meta.RequestID}
	}
	return json.Unmarshal(envelope.Data, out)
}

func (c *Client) apiKeyAuth() string {
	return c.apiKey
}

func (c *Client) jwtAuth() (string, error) {
	if c.authToken == "" {
		return "", fmt.Errorf("aegis: JWT auth token required for API key management endpoints")
	}
	return c.authToken, nil
}

func listQuery(opts *ListOptions) url.Values {
	q := url.Values{}
	if opts == nil {
		return q
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Cursor != nil {
		q.Set("cursor", opts.Cursor.UTC().Format(time.RFC3339))
	}
	return q
}

func auditQuery(filter *AuditLogFilter) url.Values {
	q := url.Values{}
	if filter == nil {
		return q
	}
	if filter.Action != "" {
		q.Set("action", filter.Action)
	}
	if filter.ResourceType != "" {
		q.Set("resource_type", filter.ResourceType)
	}
	if filter.StartTime != nil {
		q.Set("start_time", filter.StartTime.UTC().Format(time.RFC3339))
	}
	if filter.EndTime != nil {
		q.Set("end_time", filter.EndTime.UTC().Format(time.RFC3339))
	}
	if filter.Limit > 0 {
		q.Set("limit", strconv.Itoa(filter.Limit))
	}
	if filter.Cursor != nil {
		q.Set("cursor", filter.Cursor.UTC().Format(time.RFC3339))
	}
	return q
}

func projectPath(projectID string, suffix string) string {
	return "/api/v1/projects/" + url.PathEscape(projectID) + suffix
}

func secretPath(projectID, key, suffix string) string {
	return projectPath(projectID, "/secrets/"+url.PathEscape(key)+suffix)
}
