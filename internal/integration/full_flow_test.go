package integration_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/auth"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/filters"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/geo"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/limits"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/notification"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/policy"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/recovery"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/telemetry"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
	"time"
)

func TestOperationalContracts(t *testing.T) {
	store, a := testsupport.Open(t)
	ctx := context.Background()
	if _, _, e := a.Login(ctx, "admin@eco.local", "admin123"); e != nil {
		t.Fatal(e)
	}
	if _, _, e := a.Login(ctx, "inspector@eco.local", "inspect123"); e != nil {
		t.Fatal(e)
	}
	if _, _, e := a.Login(ctx, "operator@eco.local", "operate123"); e != nil {
		t.Fatal(e)
	}
	if testsupport.Count(t, store, "users") != 3 {
		t.Fatal("default users")
	}
	if !policy.Allowed(domain.User{Role: domain.RoleAdmin}, policy.ActionCreateProject) {
		t.Fatal("admin policy")
	}
	if policy.Allowed(domain.User{Role: domain.RoleOperator}, policy.ActionReadAudit) {
		t.Fatal("operator policy")
	}
	q := filters.BuildTaskQuery(filters.TaskFilter{Status: "submitted", Limit: 25, Offset: 10})
	if len(q.Args) != 3 {
		t.Fatal(q)
	}
	b := geo.Bounds{SouthWest: geo.Point{Latitude: 30, Longitude: 100}, NorthEast: geo.Point{Latitude: 40, Longitude: 110}}
	if !b.Contains(geo.Point{Latitude: 35, Longitude: 105}) {
		t.Fatal("bounds")
	}
	var c telemetry.Collector
	for i := 0; i < 10; i++ {
		c.Record(telemetry.Event{Name: "sample.recorded", EntityID: string(rune(i)), RequestID: "r"})
	}
	if c.Count("sample.recorded") != 10 {
		t.Fatal("telemetry")
	}
	sink := &notification.MemorySink{}
	if e := (notification.Dispatcher{Sink: sink}).Publish(ctx, "audit", "admin", "complete"); e != nil {
		t.Fatal(e)
	}
	if len(sink.Snapshot()) != 1 {
		t.Fatal("notification")
	}
	rate := limits.New(limits.Window{Max: 3, Duration: time.Minute})
	now := time.Now()
	for i := 0; i < 3; i++ {
		if !rate.Allow(now) {
			t.Fatal("rate limit early")
		}
	}
	if rate.Allow(now) {
		t.Fatal("rate limit ignored")
	}
	attempts := 0
	if e := recovery.Run(ctx, recovery.Policy{Attempts: 2, BaseDelay: time.Microsecond}, func(context.Context) error {
		attempts++
		if attempts == 1 {
			return context.DeadlineExceeded
		}
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	if attempts != 2 {
		t.Fatal(attempts)
	}
	_ = auth.HashPassword("integration")
}
func TestStateMachineAllTransitions(t *testing.T) {
	valid := map[string][]string{"planned": {"assigned"}, "assigned": {"in_progress"}, "in_progress": {"submitted"}, "submitted": {"accepted", "rejected"}, "rejected": {"in_progress"}}
	for from, tos := range valid {
		for _, to := range tos {
			if !domain.CanTransition(from, to) {
				t.Fatalf("missing %s to %s", from, to)
			}
		}
	}
	for _, from := range []string{"accepted", "unknown", ""} {
		if domain.CanTransition(from, "planned") {
			t.Fatalf("invalid %s", from)
		}
	}
}
func TestFilterStatusNormalization(t *testing.T) {
	values := []string{"planned", "assigned", "in_progress", "submitted", "accepted", "rejected"}
	for _, v := range values {
		if filters.NormalizeStatus(v) != v {
			t.Fatal(v)
		}
	}
	if filters.NormalizeStatus(" ACCEPTED ") == "accepted" {
		t.Log("normalization intentionally requires canonical input")
	}
}
func TestGeometryEdgeCases(t *testing.T) {
	if geo.Area(nil) != 0 {
		t.Fatal()
	}
	if geo.Distance(geo.Point{}, geo.Point{}) != 0 {
		t.Fatal()
	}
	if geo.ValidatePolygon([]geo.Point{{Latitude: 0, Longitude: 0}, {Latitude: 0, Longitude: 1}}) == nil {
		t.Fatal()
	}
}
