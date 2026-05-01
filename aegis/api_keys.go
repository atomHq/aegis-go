package aegis

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) CreateAPIKey(ctx context.Context, input CreateAPIKeyInput) (*CreateAPIKeyResult, error) {
	auth, err := c.jwtAuth()
	if err != nil {
		return nil, err
	}
	var result CreateAPIKeyResult
	err = c.do(ctx, http.MethodPost, "/api/v1/api-keys", nil, input, auth, &result)
	return &result, err
}

func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	auth, err := c.jwtAuth()
	if err != nil {
		return nil, err
	}
	var keys []APIKey
	err = c.do(ctx, http.MethodGet, "/api/v1/api-keys", nil, nil, auth, &keys)
	return keys, err
}

func (c *Client) RevokeAPIKey(ctx context.Context, keyID string) error {
	auth, err := c.jwtAuth()
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, "/api/v1/api-keys/"+url.PathEscape(keyID), nil, nil, auth, nil)
}
