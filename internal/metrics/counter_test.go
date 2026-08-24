package metrics_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/metrics"
	"sync"
	"testing"
)

func TestCounterConcurrent(t *testing.T) {
	var c metrics.Counter
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Request()
				if j%3 == 0 {
					c.Failure()
				}
			}
		}()
	}
	wg.Wait()
	r, f := c.Snapshot()
	if r != 2000 || f != 680 {
		t.Fatalf("snapshot %d %d", r, f)
	}
}
func TestCounterZero(t *testing.T) {
	var c metrics.Counter
	r, f := c.Snapshot()
	if r != 0 || f != 0 {
		t.Fatal(r, f)
	}
}
