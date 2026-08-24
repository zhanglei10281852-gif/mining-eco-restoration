package service_test

import (
	"context"
	"testing"
)

func TestInspectionPersistsNotes(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, err := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Wetland recovery", "Restore reed beds", "notes-private")
	if err != nil {
		t.Fatal(err)
	}
	task = progress(t, f, task)
	if _, err = f.inspections.Accept(context.Background(), f.inspector, task.ID, "pass", "reed beds stable", 91, "notes-private"); err != nil {
		t.Fatal(err)
	}
	items, err := f.inspections.Inspections.ByTask(context.Background(), task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Notes != "reed beds stable" {
		t.Fatalf("inspection notes lost: %#v", items)
	}
}
