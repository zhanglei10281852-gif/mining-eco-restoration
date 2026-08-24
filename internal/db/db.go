package db

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
)

type Store struct{ *sql.DB }

func Open(dsn string) (*Store, error) {
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(8)
	d.SetMaxIdleConns(4)
	if err = d.Ping(); err != nil {
		d.Close()
		return nil, err
	}
	return &Store{d}, nil
}
func (s *Store) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
