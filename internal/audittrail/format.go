package audittrail

import (
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"strings"
	"time"
)

func Describe(a domain.AuditEntry) string {
	return fmt.Sprintf("%s %s/%s by %s request=%s", a.Action, a.EntityType, a.EntityID, a.ActorID, a.RequestID)
}
func Redact(v string) string {
	parts := strings.Fields(v)
	redactNext := false
	for i, p := range parts {
		lower := strings.ToLower(p)
		if redactNext {
			parts[i] = "[redacted]"
			redactNext = false
			continue
		}
		if strings.Contains(lower, "token") || strings.Contains(lower, "password") {
			parts[i] = "[redacted]"
			redactNext = true
		}
	}
	return strings.Join(parts, " ")
}
func IsRecent(a domain.AuditEntry, now time.Time, window time.Duration) bool {
	return !a.CreatedAt.Before(now.Add(-window)) && !a.CreatedAt.After(now)
}
