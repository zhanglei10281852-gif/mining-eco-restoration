package migrations_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/migrations"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
)

func TestMigrationVersionsAndTables(t *testing.T) {
	store, _ := testsupport.Open(t)
	var n int
	if e := store.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 2 {
		t.Fatalf("versions=%d", n)
	}
	for _, table := range []string{"users", "sessions", "projects", "plots", "monitoring_samples", "remediation_tasks", "task_events", "inspections", "audit_logs", "idempotency_keys"} {
		var exists int
		if e := store.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&exists); e != nil || exists != 1 {
			t.Fatalf("table %s missing %v", table, e)
		}
	}
}
func TestMigrationIdempotent(t *testing.T) {
	store, _ := testsupport.Open(t)
	before := testsupport.Count(t, store, "schema_migrations")
	if e := migrations.Apply(context.Background(), store); e != nil {
		t.Fatal(e)
	}
	after := testsupport.Count(t, store, "schema_migrations")
	if before != after {
		t.Fatalf("migration count changed %d %d", before, after)
	}
}
func TestForeignKeyEnforced(t *testing.T) {
	store, _ := testsupport.Open(t)
	if _, e := store.Exec(`INSERT INTO plots(id,project_id,name,soil_type,area_m2,status,created_at) VALUES('p','missing','x','s',1,'open','now')`); e == nil {
		t.Fatal("foreign key disabled")
	}
}
func TestUniqueConstraints(t *testing.T) {
	store, _ := testsupport.Open(t)
	u := testsupport.User(t, store, "admin@eco.local")
	args := []any{"p1", "C1", "项目", "区域", u.ID, "active", 1, "now", "now"}
	if _, e := store.Exec(`INSERT INTO projects(id,code,name,region,owner_id,status,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, args...); e != nil {
		t.Fatal(e)
	}
	if _, e := store.Exec(`INSERT INTO projects(id,code,name,region,owner_id,status,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, "p2", "C1", "项目2", "区域", u.ID, "active", 1, "now", "now"); e == nil {
		t.Fatal("duplicate code accepted")
	}
}
