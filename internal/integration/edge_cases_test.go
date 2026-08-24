package integration_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/clock"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/codec"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/geo"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/notification"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/telemetry"
	"testing"
	"time"
)

func TestCodecAndClockTogether(t *testing.T) {
	value := map[string]int{"projects": 4, "tasks": 10}
	encoded, e := codec.Encode(value)
	if e != nil {
		t.Fatal(e)
	}
	var decoded map[string]int
	if e = codec.Decode(encoded, &decoded); e != nil || decoded["tasks"] != 10 {
		t.Fatal(decoded, e)
	}
	fixed := clock.Fixed{Value: time.Unix(10, 0)}
	if fixed.Since(time.Unix(9, 0)) != time.Second {
		t.Fatal()
	}
}
func TestNotificationCancellation(t *testing.T) {
	sink := &notification.MemorySink{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := (notification.Dispatcher{Sink: sink}).Publish(ctx, "x", "y", "z"); e == nil {
		t.Fatal()
	}
	if len(sink.Snapshot()) != 0 {
		t.Fatal()
	}
}
func TestTelemetryJSON(t *testing.T) {
	event := telemetry.Event{Name: "audit", RequestID: "r", EntityID: "e", Attributes: map[string]string{"source": "api"}}
	raw, e := event.JSON()
	if e != nil || len(raw) == 0 {
		t.Fatal(e)
	}
	var c telemetry.Collector
	c.Record(event)
	if c.Count("audit") != 1 {
		t.Fatal()
	}
}
func TestGeoInvalidCoordinates(t *testing.T) {
	for _, p := range []geo.Point{{Latitude: -91}, {Longitude: 181}, {Latitude: 100, Longitude: -181}} {
		if p.Valid() {
			t.Fatal(p)
		}
	}
}
