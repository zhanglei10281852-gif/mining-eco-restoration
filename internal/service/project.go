package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/pagination"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"time"
)

type ProjectService struct {
	DB       *sql.DB
	Projects repository.Projects
	Audit    repository.Audit
}

func (s ProjectService) Create(ctx context.Context, actor domain.User, code, name, region string) (domain.Project, error) {
	if code == "" || name == "" || region == "" {
		return domain.Project{}, apperror.ErrInvalid
	}
	now := time.Now()
	p := domain.Project{ID: ids.New("prj"), Code: code, Name: name, Region: region, OwnerID: actor.ID, Status: "active", Version: 1, CreatedAt: now, UpdatedAt: now}
	e := s.Projects.Create(ctx, p)
	if e != nil {
		return p, fmt.Errorf("create project: %w", e)
	}
	return p, nil
}
func (s ProjectService) List(ctx context.Context, page pagination.Page) ([]domain.Project, error) {
	return s.Projects.List(ctx, page.Limit, page.Offset)
}
func (s ProjectService) AddPlot(ctx context.Context, actor domain.User, projectID, name, soil string, area float64) (domain.Plot, error) {
	if area <= 0 {
		return domain.Plot{}, apperror.ErrInvalid
	}
	if _, e := s.Projects.ByID(ctx, projectID); e != nil {
		return domain.Plot{}, apperror.ErrNotFound
	}
	p := domain.Plot{ID: ids.New("plot"), ProjectID: projectID, Name: name, SoilType: soil, AreaM2: area, Status: "open", CreatedAt: time.Now()}
	if e := s.Projects.AddPlot(ctx, p); e != nil {
		return p, fmt.Errorf("add plot: %w", e)
	}
	return p, nil
}
