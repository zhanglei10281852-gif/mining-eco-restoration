package clock_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/clock"
	"testing"
	"time"
)

func TestFixedClock(t *testing.T) {
	f := clock.At("2026-01-01T00:00:00Z")
	if f.Now().UTC().Year() != 2026 {
		t.Fatal(f.Now())
	}
	if f.Since(f.Now().Add(-time.Second)) != time.Second {
		t.Fatal()
	}
}
func TestRealClock(t *testing.T) {
	r := clock.Real{}
	before := time.Now()
	if r.Now().Before(before) {
		t.Fatal()
	}
	if r.Since(before) < 0 {
		t.Fatal()
	}
}
