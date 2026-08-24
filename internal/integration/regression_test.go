package integration_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/filters"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/limits"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/policy"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/recovery"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
	"time"
)

func TestRepeatedDatabaseStartup(t *testing.T) {
	for i := 0; i < 5; i++ {
		store, a := testsupport.Open(t)
		if _, _, e := a.Login(context.Background(), "admin@eco.local", "admin123"); e != nil {
			t.Fatal(e)
		}
		if testsupport.Count(t, store, "schema_migrations") != 2 {
			t.Fatal("migration")
		}
		store.Close()
	}
}
func TestTaskFilterMatrix(t *testing.T) {
	for _, f := range []filters.TaskFilter{{}, {ProjectID: "p"}, {Status: "accepted"}, {Assignee: "u"}, {ProjectID: "p", Status: "submitted", Assignee: "u", Limit: 100, Offset: 9}} {
		q := filters.BuildTaskQuery(f)
		if q.SQL == "" || len(q.Args) < 2 {
			t.Fatal(q)
		}
	}
}
func TestPolicyMatrix(t *testing.T) {
	roles := []domain.Role{domain.RoleAdmin, domain.RoleInspector, domain.RoleOperator}
	actions := []policy.Action{policy.ActionCreateProject, policy.ActionCreateTask, policy.ActionInspectTask, policy.ActionReadAudit}
	for _, r := range roles {
		for _, a := range actions {
			_ = policy.Allowed(domain.User{Role: r}, a)
		}
	}
}
func TestRateWindowReset(t *testing.T) {
	c := limits.New(limits.Window{Max: 1, Duration: time.Second})
	now := time.Now()
	if !c.Allow(now) || c.Allow(now) {
		t.Fatal()
	}
	if !c.Allow(now.Add(time.Second)) {
		t.Fatal("window did not reset")
	}
}
func TestRetryContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	e := recovery.Run(ctx, recovery.Policy{Attempts: 5, BaseDelay: time.Millisecond}, func(context.Context) error { calls++; cancel(); return errors.New("transient") })
	if e == nil || calls > 2 {
		t.Fatal(e, calls)
	}
}
func TestSentinelsRemainDistinct(t *testing.T) {
	errs := []error{apperror.ErrInvalid, apperror.ErrConflict, apperror.ErrNotFound, apperror.ErrUnauthorized, apperror.ErrForbidden, apperror.ErrExpired, apperror.ErrCancelled}
	for i := range errs {
		for j := range errs {
			if i != j && errors.Is(errs[i], errs[j]) {
				t.Fatal("sentinel alias")
			}
		}
	}
}
func TestStatusTransitionsRejectCycles(t *testing.T) {
	for _, from := range []string{"planned", "assigned", "in_progress", "submitted", "accepted", "rejected"} {
		if domain.CanTransition(from, from) {
			t.Fatal("self transition", from)
		}
	}
}
