package health_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/health"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
)

func TestReportReady(t *testing.T) {
	store, _ := testsupport.Open(t)
	r := (health.Checker{DB: store.DB}).Report(context.Background(), true)
	if r.Status != "ready" || !r.Database || !r.Worker {
		t.Fatal(r)
	}
}
