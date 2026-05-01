package aegis

import (
	"context"
	"net/http"
)

func (c *Client) CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error) {
	var project Project
	err := c.do(ctx, http.MethodPost, "/api/v1/projects", nil, input, c.apiKeyAuth(), &project)
	return &project, err
}

func (c *Client) ListProjects(ctx context.Context, opts *ListOptions) ([]Project, error) {
	var projects []Project
	err := c.do(ctx, http.MethodGet, "/api/v1/projects", listQuery(opts), nil, c.apiKeyAuth(), &projects)
	return projects, err
}

func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	var project Project
	err := c.do(ctx, http.MethodGet, projectPath(projectID, ""), nil, nil, c.apiKeyAuth(), &project)
	return &project, err
}

func (c *Client) UpdateProject(ctx context.Context, projectID string, input UpdateProjectInput) (*Project, error) {
	var project Project
	err := c.do(ctx, http.MethodPatch, projectPath(projectID, ""), nil, input, c.apiKeyAuth(), &project)
	return &project, err
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	return c.do(ctx, http.MethodDelete, projectPath(projectID, ""), nil, nil, c.apiKeyAuth(), nil)
}
