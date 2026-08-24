package recovery

import (
	"context"
	"time"
)

type Policy struct {
	Attempts  int
	BaseDelay time.Duration
}

func Run(ctx context.Context, p Policy, fn func(context.Context) error) error {
	if p.Attempts < 1 {
		p.Attempts = 1
	}
	if p.BaseDelay <= 0 {
		p.BaseDelay = 10 * time.Millisecond
	}
	var err error
	for i := 0; i < p.Attempts; i++ {
		if e := ctx.Err(); e != nil {
			return e
		}
		if err = fn(ctx); err == nil {
			return nil
		}
		if i+1 < p.Attempts {
			d := p.BaseDelay * time.Duration(1<<i)
			timer := time.NewTimer(d)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return err
}
func Backoff(base time.Duration, attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 10 {
		attempt = 10
	}
	return base * time.Duration(1<<attempt)
}
