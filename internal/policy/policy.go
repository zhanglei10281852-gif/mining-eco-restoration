package policy

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
)

type Action string

const (
	ActionCreateProject Action = "project.create"
	ActionCreateTask    Action = "task.create"
	ActionInspectTask   Action = "task.inspect"
	ActionReadAudit     Action = "audit.read"
)

func Allowed(u domain.User, a Action) bool {
	switch u.Role {
	case domain.RoleAdmin:
		return true
	case domain.RoleInspector:
		return a == ActionInspectTask || a == ActionCreateTask
	case domain.RoleOperator:
		return a == ActionCreateTask
	}
	return false
}
func Require(ctx context.Context, u domain.User, a Action) error {
	if !Allowed(u, a) {
		return context.Canceled
	}
	return nil
}
func Actions(r domain.Role) []Action {
	u := domain.User{Role: r}
	out := []Action{}
	for _, a := range []Action{ActionCreateProject, ActionCreateTask, ActionInspectTask, ActionReadAudit} {
		if Allowed(u, a) {
			out = append(out, a)
		}
	}
	return out
}
