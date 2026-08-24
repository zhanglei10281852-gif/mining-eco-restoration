package service_test

import (
	"context"
	"testing"
)

func TestTaskTransitionRecordsActualActor(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, err := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.inspector.ID, "Revegetation", "Plant native grass", "actor-private")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.tasks.Transition(context.Background(), f.operator, task.ID, "assigned", task.Version, "actor-private"); err != nil {
		t.Fatal(err)
	}
	var actor string
	if err = f.store.QueryRow("SELECT actor_id FROM task_events WHERE task_id=? AND to_status='assigned'", task.ID).Scan(&actor); err != nil {
		t.Fatal(err)
	}
	if actor != f.operator.ID {
		t.Fatalf("transition event attributed to %s instead of %s", actor, f.operator.ID)
	}
}
