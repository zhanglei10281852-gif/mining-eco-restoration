package service_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/pagination"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/service"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
	"time"
)

type fixture struct {
	store                      *db.Store
	admin, operator, inspector domain.User
	projects                   service.ProjectService
	tasks                      service.TaskService
	samples                    service.SampleService
	inspections                service.InspectionService
}

func newFixture(t *testing.T) fixture {
	store, _ := testsupport.Open(t)
	return fixture{store: store, admin: testsupport.User(t, store, "admin@eco.local"), operator: testsupport.User(t, store, "operator@eco.local"), inspector: testsupport.User(t, store, "inspector@eco.local"), projects: service.ProjectService{DB: store.DB, Projects: repository.Projects{DB: store.DB}, Audit: repository.Audit{DB: store.DB}}, tasks: service.TaskService{DB: store.DB, Projects: repository.Projects{DB: store.DB}, Tasks: repository.Tasks{DB: store.DB}, Audit: repository.Audit{DB: store.DB}}, samples: service.SampleService{Projects: repository.Projects{DB: store.DB}, Samples: repository.Samples{DB: store.DB}}, inspections: service.InspectionService{DB: store.DB, Tasks: repository.Tasks{DB: store.DB}, Inspections: repository.Inspections{DB: store.DB}, Audit: repository.Audit{DB: store.DB}}}
}
func (f fixture) project(t *testing.T) (domain.Project, domain.Plot) {
	p, e := f.projects.Create(context.Background(), f.admin, "P-1", "Demonstration Mine", "Ningxia")
	if e != nil {
		t.Fatal(e)
	}
	pl, e := f.projects.AddPlot(context.Background(), f.admin, p.ID, "North Plot", "loam", 500)
	if e != nil {
		t.Fatal(e)
	}
	return p, pl
}
func TestProjectValidationAndList(t *testing.T) {
	f := newFixture(t)
	if _, e := f.projects.Create(context.Background(), f.admin, "", "", " "); !errors.Is(e, apperror.ErrInvalid) {
		t.Fatalf("invalid accepted %v", e)
	}
	for i := 0; i < 3; i++ {
		_, _ = f.projects.Create(context.Background(), f.admin, ids.New("code"), "Project", "Region")
	}
	items, e := f.projects.List(context.Background(), pagination.Page{Limit: 2})
	if e != nil || len(items) != 2 {
		t.Fatal(len(items), e)
	}
}
func TestTaskLifecycleAndInvalidTransition(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, e := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Slope greening", "Seed native grass", "req-1")
	if e != nil {
		t.Fatal(e)
	}
	for _, to := range []string{"assigned", "in_progress", "submitted"} {
		task, e = f.tasks.Transition(context.Background(), f.operator, task.ID, to, task.Version, "req-2")
		if e != nil {
			t.Fatalf("%s %v", to, e)
		}
	}
	if _, e = f.tasks.Transition(context.Background(), f.operator, task.ID, "planned", task.Version, "req-3"); !errors.Is(e, apperror.ErrConflict) {
		t.Fatalf("illegal transition %v", e)
	}
}
func TestTaskOptimisticConflict(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, e := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Soil stabilization", "Cover soil", "r")
	if e != nil {
		t.Fatal(e)
	}
	first, e := f.tasks.Transition(context.Background(), f.operator, task.ID, "assigned", 1, "r")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.tasks.Transition(context.Background(), f.operator, task.ID, "in_progress", 1, "r"); !errors.Is(e, apperror.ErrConflict) {
		t.Fatalf("stale version accepted %v", e)
	}
	if first.Version != 2 {
		t.Fatal(first.Version)
	}
}
func progress(t *testing.T, f fixture, task domain.RemediationTask) domain.RemediationTask {
	var e error
	for _, to := range []string{"assigned", "in_progress", "submitted"} {
		task, e = f.tasks.Transition(context.Background(), f.operator, task.ID, to, task.Version, "r")
		if e != nil {
			t.Fatal(e)
		}
	}
	return task
}
func TestInspectionAccept(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	task, e := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Soil remediation", "Apply amendment", "r")
	if e != nil {
		t.Fatal(e)
	}
	task = progress(t, f, task)
	in, e := f.inspections.Accept(context.Background(), f.inspector, task.ID, "pass", "qualified", 92, "r")
	if e != nil || in.Score != 92 {
		t.Fatal(in, e)
	}
	got, e := f.tasks.Tasks.ByID(context.Background(), task.ID)
	if e != nil || got.Status != "accepted" {
		t.Fatal(got, e)
	}
}
func TestInspectionValidation(t *testing.T) {
	f := newFixture(t)
	if _, e := f.inspections.Accept(context.Background(), f.inspector, "missing", "pass", "", 101, "r"); !errors.Is(e, apperror.ErrInvalid) {
		t.Fatalf("bad score %v", e)
	}
}
func TestSamplesRecordList(t *testing.T) {
	f := newFixture(t)
	_, pl := f.project(t)
	for i := 0; i < 5; i++ {
		if _, e := f.samples.Record(context.Background(), f.inspector, pl.ID, "moisture", string(rune('1'+i)), "%", time.Now()); e != nil {
			t.Fatal(e)
		}
	}
	items, e := f.samples.List(context.Background(), pl.ID, 3)
	if e != nil || len(items) != 3 {
		t.Fatal(len(items), e)
	}
}
func TestTaskCreateRollbackOnAuditFailure(t *testing.T) {
	f := newFixture(t)
	p, pl := f.project(t)
	if _, e := f.store.Exec(`DROP TABLE audit_logs`); e != nil {
		t.Fatal(e)
	}
	if _, e := f.tasks.Create(context.Background(), f.operator, p.ID, pl.ID, f.operator.ID, "Task", "Description", "r"); e == nil {
		t.Fatal("audit failure ignored")
	}
	if n := testsupport.Count(t, f.store, "remediation_tasks"); n != 0 {
		t.Fatalf("task leaked %d", n)
	}
}
