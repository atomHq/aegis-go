package aegis

import "time"

const (
	ScopeSecretsRead    = "secrets:read"
	ScopeSecretsWrite   = "secrets:write"
	ScopeSecretsAdmin   = "secrets:admin"
	ScopeProjectsManage = "projects:manage"
	ScopeAPIKeysManage  = "api_keys:manage"
	ScopeAuditRead      = "audit:read"
)

type Project struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Environment string    `json:"environment"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectInput struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	Environment string `json:"environment"`
}

type UpdateProjectInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type Secret struct {
	ID        string                 `json:"id"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	Version   int                    `json:"version"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
	CreatedBy string                 `json:"created_by,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type SecretWriteResult struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type SecretKey struct {
	Key           string    `json:"key"`
	LatestVersion int       `json:"latest_version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SecretVersion struct {
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	IsActive  bool      `json:"is_active"`
	CreatedBy string    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PutSecretInput struct {
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
}

type APIKey struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Scopes     []string   `json:"scopes"`
	ProjectIDs []string   `json:"project_ids,omitempty"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type CreateAPIKeyInput struct {
	Name       string     `json:"name"`
	Scopes     []string   `json:"scopes"`
	ProjectIDs []string   `json:"project_ids,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type CreateAPIKeyResult struct {
	Key          APIKey `json:"key"`
	PlaintextKey string `json:"plaintext_key"`
	Warning      string `json:"warning"`
}

type AuditLog struct {
	ID           string                 `json:"id"`
	TenantID     string                 `json:"tenant_id"`
	Actor        string                 `json:"actor"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

type ListOptions struct {
	Limit  int
	Cursor *time.Time
}

type AuditLogFilter struct {
	Action       string
	ResourceType string
	StartTime    *time.Time
	EndTime      *time.Time
	Limit        int
	Cursor       *time.Time
}
