package telemetry

import (
	"encoding/json"
	"sync"
	"time"
)

type Event struct {
	Name, RequestID, EntityID string
	Attributes                map[string]string
	At                        time.Time
}
type Collector struct {
	mu     sync.RWMutex
	events []Event
}

func (c *Collector) Record(e Event) {
	if e.At.IsZero() {
		e.At = time.Now()
	}
	if e.Attributes == nil {
		e.Attributes = map[string]string{}
	}
	c.mu.Lock()
	c.events = append(c.events, e)
	c.mu.Unlock()
}
func (c *Collector) Recent(n int) []Event {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if n <= 0 || n > len(c.events) {
		n = len(c.events)
	}
	out := make([]Event, n)
	copy(out, c.events[len(c.events)-n:])
	return out
}
func (e Event) JSON() ([]byte, error) { return json.Marshal(e) }
func (c *Collector) Count(name string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	n := 0
	for _, e := range c.events {
		if e.Name == name {
			n++
		}
	}
	return n
}
