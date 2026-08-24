package health

import (
	"context"
	"time"
)

type Report struct {
	Status    string
	CheckedAt time.Time
	Database  bool
	Worker    bool
}

func (c Checker) Report(ctx context.Context, worker bool) Report {
	ok := c.Ready(ctx)
	status := "ready"
	if !ok {
		status = "not_ready"
	}
	return Report{Status: status, CheckedAt: time.Now(), Database: ok, Worker: worker}
}
