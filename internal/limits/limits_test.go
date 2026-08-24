package limits_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/limits"
	"testing"
	"time"
)

func TestWindowAllowsAndBlocks(t *testing.T) {
	c := limits.New(limits.Window{Max: 2, Duration: time.Minute})
	now := time.Unix(100, 0)
	if !c.Allow(now) || !c.Allow(now) || c.Allow(now) {
		t.Fatal()
	}
	if c.Remaining(now) != 0 {
		t.Fatal()
	}
	if !c.Allow(now.Add(time.Minute)) {
		t.Fatal()
	}
	if e := c.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestWindowDefaultDuration(t *testing.T) {
	c := limits.New(limits.Window{Max: 0})
	if c.Remaining(time.Now()) != 1 {
		t.Fatal()
	}
}
