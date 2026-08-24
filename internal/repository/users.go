package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Users struct{ DB *sql.DB }

func (r Users) Create(ctx context.Context, u domain.User) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO users(id,email,name,role,password_hash,active,created_at) VALUES(?,?,?,?,?,?,?)`, u.ID, u.Email, u.Name, u.Role, u.PasswordHash, 1, u.CreatedAt.UTC().Format(time.RFC3339))
	return err
}
func (r Users) ByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	var active int
	var ts string
	err := r.DB.QueryRowContext(ctx, `SELECT id,email,name,role,password_hash,active,created_at FROM users WHERE email=?`, email).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.PasswordHash, &active, &ts)
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339, ts)
	return u, err
}
func (r Users) ByID(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	var active int
	var ts string
	err := r.DB.QueryRowContext(ctx, `SELECT id,email,name,role,password_hash,active,created_at FROM users WHERE id=?`, id).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.PasswordHash, &active, &ts)
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339, ts)
	return u, err
}
func (r Users) Count(ctx context.Context) (int, error) {
	var n int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&n)
	return n, err
}
