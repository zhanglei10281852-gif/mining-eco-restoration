package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
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
	to := map[bool]string{true: "accepted", false: "rejected"}[result == "pass"]
	i := domain.Inspection{ID: ids.New("insp"), TaskID: taskID, InspectorID: actor.ID, Result: result, Notes: notes, Score: score, CreatedAt: now}
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return domain.Inspection{}, e
	}
	ok, e := s.Tasks.ReserveInspectionTx(ctx, tx, taskID, "submitted", to, t.Version, now)
	if e != nil {
		_ = tx.Rollback()
		return domain.Inspection{}, fmt.Errorf("reserve inspection: %w", e)
	}
	if !ok {
		_ = tx.Rollback()
		return domain.Inspection{}, apperror.ErrConflict
	}
	if e = s.Inspections.CreateTx(ctx, tx, i); e != nil {
		_ = tx.Rollback()
		return i, fmt.Errorf("inspection transaction: %w", e)
	}
	if e = s.Audit.WriteTx(ctx, tx, domain.AuditEntry{ID: ids.New("aud"), ActorID: actor.ID, Action: "task.inspected", EntityType: "task", EntityID: taskID, RequestID: requestID, Details: result, CreatedAt: now}); e != nil {
		_ = tx.Rollback()
		return i, fmt.Errorf("inspection transaction: %w", e)
	}
	if e = tx.Commit(); e != nil {
		return i, e
	}
	return i, nil
}
