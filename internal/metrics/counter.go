package metrics

import "sync/atomic"

type Counter struct {
	requests atomic.Uint64
	failures atomic.Uint64
}

func (c *Counter) Request()                   { c.requests.Add(1) }
func (c *Counter) Failure()                   { c.failures.Add(1) }
func (c *Counter) Snapshot() (uint64, uint64) { return c.requests.Load(), c.failures.Load() }
