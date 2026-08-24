package audittrail_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/audittrail"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"strings"
	"testing"
	"time"
)

func TestDescribeRedact(t *testing.T) {
	a := domain.AuditEntry{Action: "task.created", EntityType: "task", EntityID: "t1", ActorID: "u1", RequestID: "r1"}
	if !strings.Contains(audittrail.Describe(a), "task/t1") {
		t.Fatal()
	}
	if strings.Contains(audittrail.Redact("token secret password hunter"), "secret") {
		t.Fatal()
	}
}
func TestRecentWindow(t *testing.T) {
	now := time.Now()
	a := domain.AuditEntry{CreatedAt: now.Add(-time.Minute)}
	if !audittrail.IsRecent(a, now, 2*time.Minute) {
		t.Fatal()
	}
	if audittrail.IsRecent(a, now, 30*time.Second) {
		t.Fatal()
	}
	if audittrail.IsRecent(a, now.Add(-2*time.Minute), time.Minute) {
		t.Fatal()
	}
}
