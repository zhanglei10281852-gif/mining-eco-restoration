package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"sort"
	"strings"
)

//go:embed sql/*.sql
var files embed.FS

func Apply(ctx context.Context, store *db.Store) error {
	if _, err := store.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		return err
	}
	entries, err := files.ReadDir("sql")
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		var n int
		_, _ = fmt.Sscanf(name, "%d", &n)
		var exists int
		if err = store.QueryRowContext(ctx, "SELECT COUNT(1) FROM schema_migrations WHERE version=?", name).Scan(&exists); err != nil {
			return err
		}
		if exists > 0 {
			continue
		}
		data, e := files.ReadFile("sql/" + name)
		if e != nil {
			return e
		}
		if err = store.WithTx(ctx, func(tx *sql.Tx) error {
			for _, stmt := range strings.Split(string(data), ";") {
				if strings.TrimSpace(stmt) != "" {
					if _, e := tx.ExecContext(ctx, stmt); e != nil {
						return e
					}
				}
			}
			_, e := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, applied_at) VALUES(?, datetime('now'))", name)
			return e
		}); err != nil {
			return err
		}
	}
	return nil
}
