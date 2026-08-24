package service_test

import (
	"context"
	"testing"
)

func TestProjectCreationStartsActive(t *testing.T) {
	f := newFixture(t)
	p, err := f.projects.Create(context.Background(), f.admin, "P-ACTIVE", "Green Corridor", "Inner Mongolia")
	if err != nil {
		t.Fatal(err)
	}
	if p.Status != "active" {
		t.Fatalf("project created in %s state", p.Status)
	}
	stored, err := f.projects.Projects.ByID(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != "active" {
		t.Fatalf("persisted project state is %s", stored.Status)
	}
}
