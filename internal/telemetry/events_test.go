package telemetry_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/telemetry"
	"testing"
)

func TestCollectorRecordRecentCount(t *testing.T) {
	var c telemetry.Collector
	for i := 0; i < 5; i++ {
		c.Record(telemetry.Event{Name: "task", EntityID: string(rune('a' + i))})
	}
	if c.Count("task") != 5 || c.Count("other") != 0 {
		t.Fatal()
	}
	if len(c.Recent(2)) != 2 || len(c.Recent(99)) != 5 {
		t.Fatal()
	}
	if _, e := c.Recent(1)[0].JSON(); e != nil {
		t.Fatal(e)
	}
}
func TestNilAttributesInitialized(t *testing.T) {
	var c telemetry.Collector
	c.Record(telemetry.Event{Name: "x"})
	if c.Recent(1)[0].Attributes == nil {
		t.Fatal()
	}
}
