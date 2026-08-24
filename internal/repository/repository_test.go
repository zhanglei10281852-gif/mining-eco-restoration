package repository_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
	"time"
)

func TestProjectsAndPlots(t *testing.T) {
	store, _ := testsupport.Open(t)
	u := testsupport.User(t, store, "admin@eco.local")
	r := repository.Projects{DB: store.DB}
	now := time.Now()
	p := domain.Project{ID: ids.New("prj"), Code: "K-01", Name: "北山修复", Region: "山西", OwnerID: u.ID, Status: "active", Version: 1, CreatedAt: now, UpdatedAt: now}
	if e := r.Create(context.Background(), p); e != nil {
		t.Fatal(e)
	}
	got, e := r.ByID(context.Background(), p.ID)
	if e != nil || got.Code != p.Code {
		t.Fatal(got, e)
	}
	plot := domain.Plot{ID: ids.New("plot"), ProjectID: p.ID, Name: "一号地块", SoilType: "黄土", AreaM2: 1200, Status: "open", CreatedAt: now}
	if e = r.AddPlot(context.Background(), plot); e != nil {
		t.Fatal(e)
	}
	gotPlot, e := r.Plot(context.Background(), plot.ID)
	if e != nil || gotPlot.ProjectID != p.ID {
		t.Fatal(gotPlot, e)
	}
	list, e := r.List(context.Background(), 10, 0)
	if e != nil || len(list) != 1 {
		t.Fatalf("list %d %v", len(list), e)
	}
}
func TestTaskVersionAndEvents(t *testing.T) {
	store, _ := testsupport.Open(t)
	u := testsupport.User(t, store, "operator@eco.local")
	r := repository.Tasks{DB: store.DB}
	p := repository.Projects{DB: store.DB}
	now := time.Now()
	pr := domain.Project{ID: ids.New("prj"), Code: "K-02", Name: "南坡", Region: "内蒙", OwnerID: u.ID, Status: "active", Version: 1, CreatedAt: now, UpdatedAt: now}
	if e := p.Create(context.Background(), pr); e != nil {
		t.Fatal(e)
	}
	pl := domain.Plot{ID: ids.New("plot"), ProjectID: pr.ID, Name: "地块", SoilType: "砂", AreaM2: 2, Status: "open", CreatedAt: now}
	if e := p.AddPlot(context.Background(), pl); e != nil {
		t.Fatal(e)
	}
	tx, e := store.BeginTx(context.Background(), nil)
	if e != nil {
		t.Fatal(e)
	}
	task := domain.RemediationTask{ID: ids.New("task"), ProjectID: pr.ID, PlotID: pl.ID, AssigneeID: u.ID, Title: "覆土", Description: "完成覆土", Status: "planned", Version: 1, CreatedAt: now, UpdatedAt: now}
	if e = r.CreateTx(context.Background(), tx, task, u.ID); e != nil {
		t.Fatal(e)
	}
	if e = tx.Commit(); e != nil {
		t.Fatal(e)
	}
	ok, e := r.UpdateStatus(context.Background(), task.ID, "planned", "assigned", 1, now)
	if e != nil || !ok {
		t.Fatalf("update %v %v", ok, e)
	}
	if e = r.AddEvent(context.Background(), ids.New("evt"), task.ID, u.ID, "planned", "assigned", now); e != nil {
		t.Fatal(e)
	}
	got, e := r.ByID(context.Background(), task.ID)
	if e != nil || got.Version != 2 || got.Status != "assigned" {
		t.Fatal(got, e)
	}
}
func TestSamplesAndInspections(t *testing.T) {
	store, _ := testsupport.Open(t)
	u := testsupport.User(t, store, "inspector@eco.local")
	samples := repository.Samples{DB: store.DB}
	now := time.Now()
	if e := samples.Create(context.Background(), domain.Sample{ID: ids.New("smp"), PlotID: "plot-x", CollectorID: u.ID, Metric: "ph", Value: "7", Unit: "pH", CollectedAt: now, CreatedAt: now}); e == nil {
		t.Fatal("missing plot accepted")
	}
}
func TestIdempotencyRoundTrip(t *testing.T) {
	store, _ := testsupport.Open(t)
	u := testsupport.User(t, store, "operator@eco.local")
	r := repository.Idempotency{DB: store.DB}
	if e := r.Put(context.Background(), "k", u.ID, "/api/tasks", 201, "{}"); e != nil {
		t.Fatal(e)
	}
	code, body, e := r.Get(context.Background(), "k", u.ID, "/api/tasks")
	if e != nil || code != 201 || body != "{}" {
		t.Fatal(code, body, e)
	}
}
func TestAuditRecent(t *testing.T) {
	store, _ := testsupport.Open(t)
	u := testsupport.User(t, store, "admin@eco.local")
	r := repository.Audit{DB: store.DB}
	tx, _ := store.Begin()
	e := r.WriteTx(context.Background(), tx, domain.AuditEntry{ID: ids.New("aud"), ActorID: u.ID, Action: "test", EntityType: "project", EntityID: "p", RequestID: "r", Details: "d", CreatedAt: time.Now()})
	if e != nil {
		t.Fatal(e)
	}
	if e = tx.Commit(); e != nil {
		t.Fatal(e)
	}
	items, e := r.Recent(context.Background(), 5)
	if e != nil || len(items) != 1 {
		t.Fatal(items, e)
	}
}
