package clock_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/clock"
	"testing"
	"time"
)

func TestTickerStop(t *testing.T) { tk := clock.NewTicker(time.Millisecond); tk.Stop() }
func TestEveryStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	n := 0
	done := make(chan struct{})
	go func() {
		clock.Every(ctx, time.Millisecond, func(time.Time) {
			n++
			if n == 2 {
				cancel()
			}
		})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("ticker did not stop")
	}
}
