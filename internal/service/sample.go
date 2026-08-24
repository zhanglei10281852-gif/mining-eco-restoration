package service

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"time"
)

type SampleService struct {
	Projects repository.Projects
	Samples  repository.Samples
}

func (s SampleService) Record(ctx context.Context, actor domain.User, plot, metric, value, unit string, at time.Time) (domain.Sample, error) {
	if metric == "" || value == "" || unit == "" {
		return domain.Sample{}, apperror.ErrInvalid
	}
	if _, e := s.Projects.Plot(ctx, plot); e != nil {
		return domain.Sample{}, apperror.ErrNotFound
	}
	x := domain.Sample{ID: ids.New("smp"), PlotID: plot, CollectorID: actor.ID, Metric: metric, Value: value, Unit: unit, CollectedAt: at, CreatedAt: time.Now()}
	if e := s.Samples.Create(ctx, x); e != nil {
		return x, e
	}
	return x, nil
}
func (s SampleService) List(ctx context.Context, plot string, limit int) ([]domain.Sample, error) {
	return s.Samples.List(ctx, plot, limit)
}
