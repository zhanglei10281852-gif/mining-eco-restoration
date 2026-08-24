package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
	"time"
)

func TestSampleRecordHonorsCanceledContext(t *testing.T) {
	f := newFixture(t)
	_, pl := f.project(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	x := domain.Sample{ID: ids.New("smp"), PlotID: pl.ID, CollectorID: f.operator.ID, Metric: "water_ph", Value: "7.2", Unit: "pH", CollectedAt: time.Now(), CreatedAt: time.Now()}
	if err := f.samples.Samples.Create(ctx, x); err == nil {
		t.Fatal("sample record unexpectedly ignored canceled context")
	}
	if count := testsupport.Count(t, f.store, "monitoring_samples"); count != 0 {
		t.Fatalf("canceled request persisted %d samples", count)
	}
}
