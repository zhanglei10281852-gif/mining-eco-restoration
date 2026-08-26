package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"log/slog"
	"time"
)

type Reminder struct {
	Store    *db.Store
	Log      *slog.Logger
	Interval time.Duration
	MaxRetry int
}

func New(store *db.Store, log *slog.Logger) *Reminder {
	return &Reminder{Store: store, Log: log, Interval: 30 * time.Second, MaxRetry: 3}
}
func (w *Reminder) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	w.process(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}
func (w *Reminder) process(ctx context.Context) {
	tasks, e := repository.Tasks{DB: w.Store.DB}.Pending(ctx, 20)
	if e != nil {
		w.Log.Error("worker query failed", "error", e)
		return
	}
	for _, t := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
		}
		w.Log.Info("task reminder", "task_id", t.ID, "status", t.Status)
	}
}
