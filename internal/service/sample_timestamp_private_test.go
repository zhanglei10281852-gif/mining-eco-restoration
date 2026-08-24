package service_test

import (
	"context"
	"testing"
	"time"
)

func TestSampleRecordPreservesCollectedAt(t *testing.T) {
	f := newFixture(t)
	_, pl := f.project(t)
	wanted := time.Date(2024, 6, 1, 8, 30, 0, 0, time.UTC)
	sample, err := f.samples.Record(context.Background(), f.inspector, pl.ID, "moisture", "42", "%", wanted)
	if err != nil {
		t.Fatal(err)
	}
	if !sample.CollectedAt.Equal(wanted) {
		t.Fatalf("collected timestamp changed: got %s want %s", sample.CollectedAt, wanted)
	}
	items, err := f.samples.List(context.Background(), pl.ID, 10)
	if err != nil || len(items) != 1 || !items[0].CollectedAt.Equal(wanted) {
		t.Fatalf("stored timestamp mismatch: %#v %v", items, err)
	}
}
