package limits

import (
	"errors"
	"sync"
	"time"
)

type Window struct {
	Max      int
	Duration time.Duration
}
type Counter struct {
	mu      sync.Mutex
	started time.Time
	count   int
	window  Window
}

func New(w Window) *Counter {
	if w.Max < 1 {
		w.Max = 1
	}
	if w.Duration <= 0 {
		w.Duration = time.Minute
	}
	return &Counter{window: w}
}
func (c *Counter) Allow(now time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.started.IsZero() || now.Sub(c.started) >= c.window.Duration {
		c.started = now
		c.count = 0
	}
	if c.count >= c.window.Max {
		return false
	}
	c.count++
	return true
}
func (c *Counter) Remaining(now time.Time) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.started.IsZero() || now.Sub(c.started) >= c.window.Duration {
		return c.window.Max
	}
	if c.count >= c.window.Max {
		return 0
	}
	return c.window.Max - c.count
}
func (c *Counter) Validate() error {
	if c.window.Max <= 0 {
		return errors.New("max must be positive")
	}
	return nil
}
