package testsupport

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/auth"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/migrations"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"testing"
	"time"
)

func Open(t *testing.T) (*db.Store, auth.Service) {
	t.Helper()
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared&_pragma=foreign_keys(1)"
	s, e := db.Open(dsn)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.Close() })
	if e = migrations.Apply(context.Background(), s); e != nil {
		t.Fatal(e)
	}
	a := auth.Service{Users: repository.Users{DB: s.DB}, Sessions: repository.Sessions{DB: s.DB}, TTL: time.Hour}
	if e = a.EnsureDefaults(context.Background()); e != nil {
		t.Fatal(e)
	}
	return s, a
}
func User(t *testing.T, s *db.Store, email string) domain.User {
	t.Helper()
	u, e := repository.Users{DB: s.DB}.ByEmail(context.Background(), email)
	if e != nil {
		t.Fatal(e)
	}
	return u
}
func MustExec(t *testing.T, s *db.Store, q string, args ...any) {
	t.Helper()
	if _, e := s.Exec(q, args...); e != nil {
		t.Fatal(e)
	}
}
func Count(t *testing.T, s *db.Store, table string) int {
	t.Helper()
	var n int
	if e := s.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n); e != nil {
		t.Fatal(e)
	}
	return n
}
