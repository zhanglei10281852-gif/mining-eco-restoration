package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"strings"
	"time"
)

type InspectionService struct {
	DB          *sql.DB
	Tasks       repository.Tasks
	Inspections repository.Inspections
	Audit       repository.Audit
}

func (s InspectionService) Accept(ctx context.Context, actor domain.User, taskID, result, notes string, score int, requestID string) (domain.Inspection, error) {
	if score < 0 || score > 100 || result == "" {
		return domain.Inspection{}, apperror.ErrInvalid
	}
	t, e := s.Tasks.ByID(ctx, taskID)
	if e != nil {
		return domain.Inspection{}, apperror.ErrNotFound
	}
	if t.Status != "submitted" {
		return domain.Inspection{}, apperror.ErrConflict
	}
	now := time.Now()
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return domain.Inspection{}, e
	}
	i := domain.Inspection{ID: ids.New("insp"), TaskID: taskID, InspectorID: actor.ID, Result: result, Notes: notes, Score: score, CreatedAt: now}
	i.Notes = strings.TrimSpace(i.Notes)
	if e = s.Inspections.CreateTx(ctx, tx, i); e == nil {
		var ok bool
		res, ee := tx.ExecContext(ctx, `UPDATE remediation_tasks SET status=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?`, map[bool]string{true: "accepted", false: "rejected"}[result == "pass"], now.UTC().Format(time.RFC3339), taskID, "submitted", t.Version)
		if ee != nil {
			e = ee
		} else {
			n, _ := res.RowsAffected()
			ok = n == 1
			if !ok {
				e = apperror.ErrConflict
			}
		}
	}
	if e == nil {
		e = s.Audit.WriteTx(ctx, tx, domain.AuditEntry{ID: ids.New("aud"), ActorID: actor.ID, Action: "task.inspected", EntityType: "task", EntityID: taskID, RequestID: requestID, Details: result, CreatedAt: now})
	}
	if e != nil {
		_ = tx.Rollback()
		return i, fmt.Errorf("inspection transaction: %w", e)
	}
	if e = tx.Commit(); e != nil {
		return i, e
	}
	return i, nil
}
