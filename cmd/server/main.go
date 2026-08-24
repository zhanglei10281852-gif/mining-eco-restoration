package main

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/config"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/httpapi"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/migrations"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()
	store, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		logger.Error("database open failed", "error", err)
		return
	}
	defer store.Close()
	if err = migrations.Apply(context.Background(), store); err != nil {
		logger.Error("migration failed", "error", err)
		return
	}
	if err = httpapi.New(cfg, store, logger).Run(context.Background()); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
