package auth_test

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/auth"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"strings"
	"testing"
	"time"
)

func TestDefaultUsersAndLogin(t *testing.T) {
	store, a := testsupport.Open(t)
	u, s, e := a.Login(context.Background(), "admin@eco.local", "admin123")
	if e != nil {
		t.Fatal(e)
	}
	if u.Role != "admin" {
		t.Fatal(u.Role)
	}
	if s.TokenHash == "" || !strings.HasPrefix(s.TokenHash, "tok_") {
		t.Fatal("token missing")
	}
	got, _, e := a.Authenticate(context.Background(), "Bearer "+s.TokenHash)
	if e != nil || got.ID != u.ID {
		t.Fatalf("auth %v %v", got, e)
	}
	if testsupport.Count(t, store, "sessions") != 1 {
		t.Fatal("session not persisted")
	}
}
func TestWrongPasswordAndMalformedBearer(t *testing.T) {
	_, a := testsupport.Open(t)
	if _, _, e := a.Login(context.Background(), "admin@eco.local", "bad"); e != apperror.ErrUnauthorized {
		t.Fatalf("wrong password %v", e)
	}
	for _, v := range []string{"", "Token abc", "Bearer", "Bearer a b"} {
		if _, _, e := a.Authenticate(context.Background(), v); e != apperror.ErrUnauthorized {
			t.Fatalf("%q accepted: %v", v, e)
		}
	}
}
func TestLogoutRevokes(t *testing.T) {
	_, a := testsupport.Open(t)
	_, s, e := a.Login(context.Background(), "operator@eco.local", "operate123")
	if e != nil {
		t.Fatal(e)
	}
	if e = a.Logout(context.Background(), s); e != nil {
		t.Fatal(e)
	}
	if _, _, e = a.Authenticate(context.Background(), "Bearer "+s.TokenHash); e != apperror.ErrUnauthorized {
		t.Fatalf("revoked token accepted: %v", e)
	}
}
func TestHashDeterminism(t *testing.T) {
	if auth.HashPassword("x") != auth.HashPassword("x") {
		t.Fatal("hash changed")
	}
	if auth.HashPassword("x") == auth.HashPassword("y") {
		t.Fatal("hash collision")
	}
}
func TestEnsureDefaultsIdempotent(t *testing.T) {
	store, a := testsupport.Open(t)
	if e := a.EnsureDefaults(context.Background()); e != nil {
		t.Fatal(e)
	}
	if n := testsupport.Count(t, store, "users"); n != 3 {
		t.Fatalf("users=%d", n)
	}
}
func TestSessionExpiryWithShortTTL(t *testing.T) {
	store, a := testsupport.Open(t)
	a.TTL = time.Millisecond
	_, s, e := a.Login(context.Background(), "inspector@eco.local", "inspect123")
	if e != nil {
		t.Fatal(e)
	}
	time.Sleep(5 * time.Millisecond)
	if _, _, e = a.Authenticate(context.Background(), "Bearer "+s.TokenHash); e != apperror.ErrUnauthorized {
		t.Fatalf("expired token accepted: %v", e)
	}
	_ = store
}
