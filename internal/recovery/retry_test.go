package recovery_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/recovery"
	"testing"
	"time"
)

func TestRunRetriesUntilSuccess(t *testing.T) {
	n := 0
	e := recovery.Run(context.Background(), recovery.Policy{Attempts: 4, BaseDelay: time.Microsecond}, func(context.Context) error {
		n++
		if n < 3 {
			return errors.New("retry")
		}
		return nil
	})
	if e != nil || n != 3 {
		t.Fatal(e, n)
	}
}
func TestRunReturnsLastError(t *testing.T) {
	n := 0
	want := errors.New("fail")
	e := recovery.Run(context.Background(), recovery.Policy{Attempts: 3, BaseDelay: time.Microsecond}, func(context.Context) error { n++; return want })
	if !errors.Is(e, want) || n != 3 {
		t.Fatal(e, n)
	}
}
func TestRunCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := recovery.Run(ctx, recovery.Policy{Attempts: 3}, func(context.Context) error { return nil }); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestBackoffCap(t *testing.T) {
	if recovery.Backoff(time.Second, 0) != time.Second {
		t.Fatal()
	}
	if recovery.Backoff(time.Second, 2) != 4*time.Second {
		t.Fatal()
	}
	if recovery.Backoff(time.Second, 20) != 1024*time.Second {
		t.Fatal()
	}
}
