package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Tasks struct{ DB *sql.DB }

func scanTask(row interface{ Scan(...any) error }) (domain.RemediationTask, error) {
	var t domain.RemediationTask
	var due, created, updated sql.NullString
	err := row.Scan(&t.ID, &t.ProjectID, &t.PlotID, &t.AssigneeID, &t.Title, &t.Description, &t.Status, &due, &t.Version, &created, &updated)
	if due.Valid {
		x, _ := time.Parse(time.RFC3339, due.String)
		t.DueAt = &x
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, created.String)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updated.String)
	return t, err
}
func (r Tasks) CreateTx(ctx context.Context, tx *sql.Tx, t domain.RemediationTask, actor string) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO remediation_tasks(id,project_id,plot_id,assignee_id,title,description,status,due_at,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, t.ID, t.ProjectID, t.PlotID, t.AssigneeID, t.Title, t.Description, t.Status, nil, t.Version, t.CreatedAt.UTC().Format(time.RFC3339), t.UpdatedAt.UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO task_events(id,task_id,actor_id,from_status,to_status,created_at) VALUES(?,?,?,?,?,?)`, t.ID+"_event", t.ID, actor, "", t.Status, t.CreatedAt.UTC().Format(time.RFC3339))
	return err
}
func (r Tasks) ByID(ctx context.Context, id string) (domain.RemediationTask, error) {
	return scanTask(r.DB.QueryRowContext(ctx, `SELECT id,project_id,plot_id,assignee_id,title,description,status,due_at,version,created_at,updated_at FROM remediation_tasks WHERE id=?`, id))
}
func (r Tasks) UpdateStatus(ctx context.Context, id, from, to string, version int64, at time.Time) (bool, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE remediation_tasks SET status=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?`, to, at.UTC().Format(time.RFC3339), id, from, version)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}
func (r Tasks) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id, from, to string, version int64, at time.Time) (bool, error) {
	res, err := tx.ExecContext(ctx, `UPDATE remediation_tasks SET status=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?`, to, at.UTC().Format(time.RFC3339), id, from, version)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}
func (r Tasks) AddEvent(ctx context.Context, eid, tid, actor, from, to string, at time.Time) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO task_events(id,task_id,actor_id,from_status,to_status,created_at) VALUES(?,?,?,?,?,?)`, eid, tid, actor, from, to, at.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("record task transition event: %w", err)
	}
	return err
}
func (r Tasks) AddEventTx(ctx context.Context, tx *sql.Tx, eid, tid, actor, from, to string, at time.Time) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO task_events(id,task_id,actor_id,from_status,to_status,created_at) VALUES(?,?,?,?,?,?)`, eid, tid, actor, from, to, at.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("record task transition event: %w", err)
	}
	return err
}
func (r Tasks) Pending(ctx context.Context, limit int) ([]domain.RemediationTask, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id,project_id,plot_id,assignee_id,title,description,status,due_at,version,created_at,updated_at FROM remediation_tasks WHERE status IN ('submitted','rejected') ORDER BY updated_at LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.RemediationTask{}
	for rows.Next() {
		t, e := scanTask(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
