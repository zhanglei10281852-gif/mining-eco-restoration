package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"strings"
	"time"
)

type Service struct {
	Users    repository.Users
	Sessions repository.Sessions
	TTL      time.Duration
}

func HashPassword(p string) string {
	h := sha256.Sum256([]byte("mining-eco:" + p))
	return hex.EncodeToString(h[:])
}
func tokenHash(t string) string { h := sha256.Sum256([]byte(t)); return hex.EncodeToString(h[:]) }
func (s Service) EnsureDefaults(ctx context.Context) error {
	n, e := s.Users.Count(ctx)
	if e != nil || n > 0 {
		return e
	}
	now := time.Now()
	for _, u := range []domain.User{{ID: ids.New("usr"), Email: "admin@eco.local", Name: "平台管理员", Role: domain.RoleAdmin, PasswordHash: HashPassword("admin123"), CreatedAt: now}, {ID: ids.New("usr"), Email: "inspector@eco.local", Name: "生态验收员", Role: domain.RoleInspector, PasswordHash: HashPassword("inspect123"), CreatedAt: now}, {ID: ids.New("usr"), Email: "operator@eco.local", Name: "治理执行员", Role: domain.RoleOperator, PasswordHash: HashPassword("operate123"), CreatedAt: now}} {
		if e = s.Users.Create(ctx, u); e != nil {
			return e
		}
	}
	return nil
}
func (s Service) Login(ctx context.Context, email, password string) (domain.User, domain.Session, error) {
	u, e := s.Users.ByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if e != nil {
		return u, domain.Session{}, apperror.ErrUnauthorized
	}
	if !u.Active || u.PasswordHash != HashPassword(password) {
		return u, domain.Session{}, apperror.ErrUnauthorized
	}
	now := time.Now()
	tok := ids.New("tok")
	sess := domain.Session{ID: ids.New("ses"), UserID: u.ID, TokenHash: tokenHash(tok), ExpiresAt: now.Add(s.TTL), CreatedAt: now}
	if e = s.Sessions.Create(ctx, sess); e != nil {
		return u, domain.Session{}, e
	}
	sess.TokenHash = tok
	return u, sess, nil
}
func (s Service) Authenticate(ctx context.Context, bearer string) (domain.User, domain.Session, error) {
	parts := strings.Fields(bearer)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return domain.User{}, domain.Session{}, apperror.ErrUnauthorized
	}
	sess, e := s.Sessions.ActiveByHash(ctx, tokenHash(parts[1]), time.Now())
	if e != nil {
		return domain.User{}, domain.Session{}, apperror.ErrUnauthorized
	}
	u, e := s.Users.ByID(ctx, sess.UserID)
	if e != nil {
		return domain.User{}, domain.Session{}, e
	}
	return u, sess, nil
}
func (s Service) Logout(ctx context.Context, sess domain.Session) error {
	if sess.ID == "" {
		return errors.New("empty session")
	}
	return s.Sessions.Revoke(ctx, sess.ID, time.Now())
}
