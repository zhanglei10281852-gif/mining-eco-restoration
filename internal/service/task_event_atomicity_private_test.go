package service_test

import (
	"context"
	"testing"
)

func TestTaskTransitionFailureKeepsState(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, err := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Slope drainage", "Clear blocked channels", "event-private")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.store.Exec("DROP TABLE task_events"); err != nil {
		t.Fatal(err)
	}
	if _, err = f.tasks.Transition(context.Background(), f.operator, task.ID, "assigned", task.Version, "event-private"); err == nil {
		t.Fatal("transition unexpectedly succeeded")
	}
	current, err := f.tasks.Tasks.ByID(context.Background(), task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.Status != "planned" || current.Version != task.Version {
		t.Fatalf("failed transition changed task: status=%s version=%d", current.Status, current.Version)
	}
}
