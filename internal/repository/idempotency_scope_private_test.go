package repository_test

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
)

func TestIdempotencyKeyIsScopedByRoute(t *testing.T) {
	store, _ := testsupport.Open(t)
	user := testsupport.User(t, store, "operator@eco.local")
	r := repository.Idempotency{DB: store.DB}
	if err := r.Put(context.Background(), "same-key", user.ID, "/tasks", 201, "tasks-response"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := r.Get(context.Background(), "same-key", user.ID, "/projects"); err != sql.ErrNoRows {
		t.Fatalf("idempotency record leaked across routes: %v", err)
	}
}
