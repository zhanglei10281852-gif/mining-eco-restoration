package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
)

func TestInspectionFailurePreservesSubmittedTask(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, err := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Terrace repair", "Reinforce the drainage terrace", "inspection-private")
	if err != nil {
		t.Fatal(err)
	}
	task = progress(t, f, task)
	if _, err = f.store.Exec("DROP TABLE audit_logs"); err != nil {
		t.Fatal(err)
	}
	if _, err = f.inspections.Accept(context.Background(), f.inspector, task.ID, "pass", "audit unavailable", 88, "inspection-private"); err == nil {
		t.Fatal("inspection unexpectedly succeeded")
	}
	current, err := f.tasks.Tasks.ByID(context.Background(), task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.Status != "submitted" || current.Version != task.Version {
		t.Fatalf("failed inspection changed task state: status=%s version=%d", current.Status, current.Version)
	}
	if count := testsupport.Count(t, f.store, "inspections"); count != 0 {
		t.Fatalf("failed inspection persisted %d records", count)
	}
}
