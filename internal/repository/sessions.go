package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Sessions struct{ DB *sql.DB }

func (r Sessions) Create(ctx context.Context, s domain.Session) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO sessions(id,user_id,token_hash,expires_at,created_at) VALUES(?,?,?,?,?)`, s.ID, s.UserID, s.TokenHash, s.ExpiresAt.UTC().Format(time.RFC3339), s.CreatedAt.UTC().Format(time.RFC3339))
	return err
}
func (r Sessions) ActiveByHash(ctx context.Context, hash string, now time.Time) (domain.Session, error) {
	var s domain.Session
	var exp, created string
	var revoked sql.NullString
	err := r.DB.QueryRowContext(ctx, `SELECT id,user_id,token_hash,expires_at,revoked_at,created_at FROM sessions WHERE token_hash=?`, hash).Scan(&s.ID, &s.UserID, &s.TokenHash, &exp, &revoked, &created)
	s.ExpiresAt, _ = time.Parse(time.RFC3339, exp)
	s.CreatedAt, _ = time.Parse(time.RFC3339, created)
	if revoked.Valid {
		t, _ := time.Parse(time.RFC3339, revoked.String)
		s.RevokedAt = &t
	}
	if err == nil && (s.RevokedAt != nil || !s.ExpiresAt.After(now)) {
		return domain.Session{}, sql.ErrNoRows
	}
	return s, err
}
func (r Sessions) Revoke(ctx context.Context, id string, at time.Time) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE id=? AND revoked_at IS NULL`, at.UTC().Format(time.RFC3339), id)
	return err
}
func (r Sessions) RevokeUser(ctx context.Context, userID string, at time.Time) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE user_id=? AND revoked_at IS NULL`, at.UTC().Format(time.RFC3339), userID)
	return err
}
