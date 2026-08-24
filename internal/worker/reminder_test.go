package worker_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/worker"
	"log/slog"
	"testing"
	"time"
)

func TestWorkerStopsOnContext(t *testing.T) {
	store, _ := testsupport.Open(t)
	w := worker.New(store, slog.Default())
	w.Interval = time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { w.Run(ctx); close(done) }()
	time.Sleep(3 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("worker did not stop")
	}
}
