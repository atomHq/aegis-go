package aegis

import (
	"context"
	"net/http"
	"strconv"
)

func (c *Client) PutSecret(ctx context.Context, projectID string, input PutSecretInput) (*SecretWriteResult, error) {
	var result SecretWriteResult
	err := c.do(ctx, http.MethodPut, projectPath(projectID, "/secrets"), nil, input, c.apiKeyAuth(), &result)
	return &result, err
}

func (c *Client) BulkPutSecrets(ctx context.Context, projectID string, secrets []PutSecretInput) ([]SecretWriteResult, error) {
	var result []SecretWriteResult
	err := c.do(ctx, http.MethodPut, projectPath(projectID, "/secrets/bulk"), nil, map[string]any{"secrets": secrets}, c.apiKeyAuth(), &result)
	return result, err
}

func (c *Client) GetSecret(ctx context.Context, projectID, key string) (*Secret, error) {
	var secret Secret
	err := c.do(ctx, http.MethodGet, secretPath(projectID, key, ""), nil, nil, c.apiKeyAuth(), &secret)
	return &secret, err
}

func (c *Client) BulkGetSecrets(ctx context.Context, projectID string, keys []string) (map[string]string, error) {
	var secrets map[string]string
	err := c.do(ctx, http.MethodPost, projectPath(projectID, "/secrets/bulk"), nil, map[string]any{"keys": keys}, c.apiKeyAuth(), &secrets)
	return secrets, err
}

func (c *Client) ListSecretKeys(ctx context.Context, projectID string, opts *ListOptions) ([]SecretKey, error) {
	var keys []SecretKey
	err := c.do(ctx, http.MethodGet, projectPath(projectID, "/secrets"), listQuery(opts), nil, c.apiKeyAuth(), &keys)
	return keys, err
}

func (c *Client) DeleteSecret(ctx context.Context, projectID, key string) error {
	return c.do(ctx, http.MethodDelete, secretPath(projectID, key, ""), nil, nil, c.apiKeyAuth(), nil)
}

func (c *Client) ListSecretVersions(ctx context.Context, projectID, key string) ([]SecretVersion, error) {
	var versions []SecretVersion
	err := c.do(ctx, http.MethodGet, secretPath(projectID, key, "/versions"), nil, nil, c.apiKeyAuth(), &versions)
	return versions, err
}

func (c *Client) GetSecretVersion(ctx context.Context, projectID, key string, version int) (*Secret, error) {
	var secret Secret
	err := c.do(ctx, http.MethodGet, secretPath(projectID, key, "/versions/"+strconv.Itoa(version)), nil, nil, c.apiKeyAuth(), &secret)
	return &secret, err
}
