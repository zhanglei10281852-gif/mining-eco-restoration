package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Audit struct{ DB *sql.DB }

func (r Audit) WriteTx(ctx context.Context, tx *sql.Tx, a domain.AuditEntry) error {
	_, e := tx.ExecContext(ctx, `INSERT INTO audit_logs(id,actor_id,action,entity_type,entity_id,request_id,details,created_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.ActorID, a.Action, a.EntityType, a.EntityID, a.RequestID, a.Details, a.CreatedAt.UTC().Format(time.RFC3339))
	return e
}
func (r Audit) Recent(ctx context.Context, limit int) ([]domain.AuditEntry, error) {
	rows, e := r.DB.QueryContext(ctx, `SELECT id,actor_id,action,entity_type,entity_id,request_id,details,created_at FROM audit_logs ORDER BY created_at DESC LIMIT ?`, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.AuditEntry{}
	for rows.Next() {
		var a domain.AuditEntry
		var actor, ts sql.NullString
		if e = rows.Scan(&a.ID, &actor, &a.Action, &a.EntityType, &a.EntityID, &a.RequestID, &a.Details, &ts); e != nil {
			return nil, e
		}
		a.ActorID = actor.String
		a.CreatedAt, _ = time.Parse(time.RFC3339, ts.String)
		out = append(out, a)
	}
	return out, rows.Err()
}
