package repository_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"testing"
)

func TestPendingHonorsCanceledContext(t *testing.T) {
	store, _ := testsupport.Open(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (repository.Tasks{DB: store.DB}).Pending(ctx, 20); err == nil {
		t.Fatal("pending query ignored canceled context")
	}
}
