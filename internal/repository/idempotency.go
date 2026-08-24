package repository

import (
	"context"
	"database/sql"
	"time"
)

type Idempotency struct{ DB *sql.DB }

func (r Idempotency) Get(ctx context.Context, key, user, route string) (int, string, error) {
	var c int
	var b string
	e := r.DB.QueryRowContext(ctx, `SELECT response_code,response_body FROM idempotency_keys WHERE key=? AND user_id=? AND route=?`, key, user, route).Scan(&c, &b)
	return c, b, e
}
func (r Idempotency) Put(ctx context.Context, key, user, route string, code int, body string) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO idempotency_keys(key,user_id,route,response_code,response_body,created_at) VALUES(?,?,?,?,?,?)`, key, user, route, code, body, time.Now().UTC().Format(time.RFC3339))
	return e
}
