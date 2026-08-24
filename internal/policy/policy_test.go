package policy_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/policy"
	"testing"
)

func TestRoleActions(t *testing.T) {
	admin := domain.User{Role: domain.RoleAdmin}
	if !policy.Allowed(admin, policy.ActionReadAudit) {
		t.Fatal()
	}
	op := domain.User{Role: domain.RoleOperator}
	if policy.Allowed(op, policy.ActionReadAudit) || !policy.Allowed(op, policy.ActionCreateTask) {
		t.Fatal()
	}
	in := domain.User{Role: domain.RoleInspector}
	if !policy.Allowed(in, policy.ActionInspectTask) || policy.Allowed(in, policy.ActionCreateProject) {
		t.Fatal()
	}
}
func TestRequireAndActions(t *testing.T) {
	u := domain.User{Role: domain.RoleInspector}
	if e := policy.Require(context.Background(), u, policy.ActionCreateProject); e == nil {
		t.Fatal()
	}
	if len(policy.Actions(domain.RoleAdmin)) != 4 {
		t.Fatal()
	}
	if len(policy.Actions(domain.RoleOperator)) != 1 {
		t.Fatal()
	}
}
