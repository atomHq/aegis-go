package aegis

import (
	"context"
	"net/http"
)

func (c *Client) ListAuditLogs(ctx context.Context, filter *AuditLogFilter) ([]AuditLog, error) {
	var logs []AuditLog
	err := c.do(ctx, http.MethodGet, "/api/v1/audit-logs", auditQuery(filter), nil, c.apiKeyAuth(), &logs)
	return logs, err
}
