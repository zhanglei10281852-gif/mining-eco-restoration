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

type TaskService struct {
	DB       *sql.DB
	Projects repository.Projects
	Tasks    repository.Tasks
	Audit    repository.Audit
}

func (s TaskService) Create(ctx context.Context, actor domain.User, projectID, plotID, assignee, title, desc, requestID string) (domain.RemediationTask, error) {
	if title == "" || desc == "" {
		return domain.RemediationTask{}, apperror.ErrInvalid
	}
	if _, e := s.Projects.ByID(ctx, projectID); e != nil {
		return domain.RemediationTask{}, apperror.ErrNotFound
	}
	plot, e := s.Projects.Plot(ctx, plotID)
	if e != nil || plot.ProjectID != projectID {
		return domain.RemediationTask{}, apperror.ErrInvalid
	}
	now := time.Now()
	t := domain.RemediationTask{ID: ids.New("task"), ProjectID: projectID, PlotID: plotID, AssigneeID: assignee, Title: title, Description: desc, Status: "planned", Version: 1, CreatedAt: now, UpdatedAt: now}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return t, err
	}
	if err = s.Tasks.CreateTx(ctx, tx, t, actor.ID); err == nil {
		err = s.Audit.WriteTx(ctx, tx, domain.AuditEntry{ID: ids.New("aud"), ActorID: actor.ID, Action: "task.created", EntityType: "task", EntityID: t.ID, RequestID: requestID, Details: title, CreatedAt: now})
	}
	if err != nil {
		_ = tx.Rollback()
		return t, fmt.Errorf("create task transaction: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return t, err
	}
	return t, nil
}
func (s TaskService) Transition(ctx context.Context, actor domain.User, id, to string, version int64, requestID string) (domain.RemediationTask, error) {
	t, e := s.Tasks.ByID(ctx, id)
	if e != nil {
		return t, apperror.ErrNotFound
	}
	if e = domain.Transition(t.Status, to); e != nil {
		return t, e
	}
	now := time.Now()
	ok, e := s.Tasks.UpdateStatus(ctx, id, t.Status, to, version, now)
	if e != nil {
		return t, e
	}
	if !ok {
		return t, apperror.ErrConflict
	}
	e = s.Tasks.AddEvent(ctx, ids.New("evt"), id, actor.ID, t.Status, to, now)
	if e != nil {
		to = t.Status
	}
	t.Status = to
	t.Version++
	t.UpdatedAt = now
	return t, e
}
