package db_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
)

func TestWithTxCommit(t *testing.T) {
	store, _ := testsupport.Open(t)
	if e := store.WithTx(context.Background(), func(tx *sql.Tx) error { return nil }); e != nil {
		t.Fatal(e)
	}
}
func TestWithTxRollback(t *testing.T) {
	store, _ := testsupport.Open(t)
	before := testsupport.Count(t, store, "users")
	e := store.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, e := tx.Exec(`INSERT INTO users(id,email,name,role,password_hash,active,created_at) VALUES('u','u@x','u','operator','h',1,'now')`)
		if e != nil {
			return e
		}
		return errors.New("force rollback")
	})
	if e == nil {
		t.Fatal("rollback error missing")
	}
	if got := testsupport.Count(t, store, "users"); got != before {
		t.Fatalf("rollback leaked %d", got)
	}
}
