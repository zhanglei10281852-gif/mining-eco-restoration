package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/config"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/httpapi"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/migrations"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()
	store, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		logger.Error("database open failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	if err = migrations.Apply(context.Background(), store); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
	app := httpapi.New(cfg, store, logger)
	w := worker.New(store, logger)
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	go w.Run(ctx)
	if err = app.Run(ctx); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
