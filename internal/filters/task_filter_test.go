package filters_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/filters"
	"strings"
	"testing"
)

func TestBuildQueryClauses(t *testing.T) {
	q := filters.BuildTaskQuery(filters.TaskFilter{ProjectID: "p", Status: "submitted", Assignee: "u", Limit: 10, Offset: 3})
	if !strings.Contains(q.SQL, "project_id=?") || !strings.Contains(q.SQL, "status=?") || !strings.Contains(q.SQL, "assignee_id=?") {
		t.Fatal(q.SQL)
	}
	if len(q.Args) != 5 || q.Args[3] != 10 || q.Args[4] != 3 {
		t.Fatal(q.Args)
	}
}
func TestBuildQueryDefaults(t *testing.T) {
	q := filters.BuildTaskQuery(filters.TaskFilter{Limit: 200, Offset: -1})
	if !strings.HasSuffix(q.SQL, "LIMIT ? OFFSET ?") {
		t.Fatal(q.SQL)
	}
	if q.Args[len(q.Args)-2] != 20 || q.Args[len(q.Args)-1] != 0 {
		t.Fatal(q.Args)
	}
}
func TestNormalizeStatus(t *testing.T) {
	for _, v := range []string{"planned", "ASSIGNED", "in_progress", "submitted", "accepted", "rejected"} {
		if filters.NormalizeStatus(v) == "" {
			t.Fatal(v)
		}
	}
	if filters.NormalizeStatus("unknown") != "" {
		t.Fatal()
	}
}
