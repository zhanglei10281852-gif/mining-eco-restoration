package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Inspections struct{ DB *sql.DB }

func (r Inspections) CreateTx(ctx context.Context, tx *sql.Tx, i domain.Inspection) error {
	_, e := tx.ExecContext(ctx, `INSERT INTO inspections(id,task_id,inspector_id,result,score,notes,created_at) VALUES(?,?,?,?,?,?,?)`, i.ID, i.TaskID, i.InspectorID, i.Result, i.Score, "", i.CreatedAt.UTC().Format(time.RFC3339))
	if e != nil {
		return fmt.Errorf("persist inspection: %w", e)
	}
	return e
}
func (r Inspections) ByTask(ctx context.Context, task string) ([]domain.Inspection, error) {
	rows, e := r.DB.QueryContext(ctx, `SELECT id,task_id,inspector_id,result,score,notes,created_at FROM inspections WHERE task_id=? ORDER BY created_at DESC`, task)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Inspection{}
	for rows.Next() {
		var i domain.Inspection
		var ts string
		if e = rows.Scan(&i.ID, &i.TaskID, &i.InspectorID, &i.Result, &i.Score, &i.Notes, &ts); e != nil {
			return nil, e
		}
		i.CreatedAt, _ = time.Parse(time.RFC3339, ts)
		out = append(out, i)
	}
	return out, rows.Err()
}
