package domain_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"testing"
)

func TestTransitionMatrix(t *testing.T) {
	cases := []struct {
		from, to string
		ok       bool
	}{{"planned", "assigned", true}, {"assigned", "in_progress", true}, {"in_progress", "submitted", true}, {"submitted", "accepted", true}, {"submitted", "rejected", true}, {"rejected", "in_progress", true}, {"accepted", "planned", false}, {"planned", "accepted", false}, {"", "planned", false}}
	for _, tc := range cases {
		if got := domain.CanTransition(tc.from, tc.to); got != tc.ok {
			t.Fatalf("%s->%s got %v", tc.from, tc.to, got)
		}
		if err := domain.Transition(tc.from, tc.to); (err == nil) != tc.ok {
			t.Fatalf("transition error mismatch %s->%s", tc.from, tc.to)
		}
	}
}
func TestRoles(t *testing.T) {
	for _, r := range []domain.Role{domain.RoleAdmin, domain.RoleInspector, domain.RoleOperator} {
		if !domain.ValidRole(r) {
			t.Fatalf("role %q invalid", r)
		}
	}
	if domain.ValidRole("guest") {
		t.Fatal("guest accepted")
	}
}
func TestStatusesStable(t *testing.T) {
	s := domain.TaskStatuses()
	if len(s) != 6 {
		t.Fatalf("statuses=%d", len(s))
	}
	seen := map[string]bool{}
	for _, v := range s {
		if seen[v] {
			t.Fatalf("duplicate %s", v)
		}
		seen[v] = true
	}
}
