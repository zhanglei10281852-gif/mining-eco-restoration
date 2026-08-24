package filters

import (
	"fmt"
	"strings"
)

type TaskFilter struct {
	ProjectID, Status, Assignee string
	Limit, Offset               int
}
type Query struct {
	SQL  string
	Args []any
}

func BuildTaskQuery(f TaskFilter) Query {
	clauses := []string{"1=1"}
	args := []any{}
	if f.ProjectID != "" {
		clauses = append(clauses, "project_id=?")
		args = append(args, f.ProjectID)
	}
	if f.Status != "" {
		clauses = append(clauses, "status=?")
		args = append(args, f.Status)
	}
	if f.Assignee != "" {
		clauses = append(clauses, "assignee_id=?")
		args = append(args, f.Assignee)
	}
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	return Query{SQL: fmt.Sprintf("SELECT id,project_id,plot_id,assignee_id,title,description,status,due_at,version,created_at,updated_at FROM remediation_tasks WHERE %s ORDER BY updated_at DESC LIMIT ? OFFSET ?", strings.Join(clauses, " AND ")), Args: append(args, limit, offset)}
}
func NormalizeStatus(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "planned", "assigned", "in_progress", "submitted", "accepted", "rejected":
		return strings.ToLower(strings.TrimSpace(s))
	default:
		return ""
	}
}
