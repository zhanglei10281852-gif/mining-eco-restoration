package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Samples struct{ DB *sql.DB }

func (r Samples) Create(ctx context.Context, s domain.Sample) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO monitoring_samples(id,plot_id,collector_id,metric,value,unit,collected_at,created_at) VALUES(?,?,?,?,?,?,?,?)`, s.ID, s.PlotID, s.CollectorID, s.Metric, s.Value, s.Unit, s.CollectedAt.UTC().Format(time.RFC3339), s.CreatedAt.UTC().Format(time.RFC3339))
	return e
}
func (r Samples) List(ctx context.Context, plot string, limit int) ([]domain.Sample, error) {
	rows, e := r.DB.QueryContext(ctx, `SELECT id,plot_id,collector_id,metric,value,unit,collected_at,created_at FROM monitoring_samples WHERE plot_id=? ORDER BY collected_at DESC LIMIT ?`, plot, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Sample{}
	for rows.Next() {
		var s domain.Sample
		var c, cr string
		if e = rows.Scan(&s.ID, &s.PlotID, &s.CollectorID, &s.Metric, &s.Value, &s.Unit, &c, &cr); e != nil {
			return nil, e
		}
		s.CollectedAt, _ = time.Parse(time.RFC3339, c)
		s.CreatedAt, _ = time.Parse(time.RFC3339, cr)
		out = append(out, s)
	}
	return out, rows.Err()
}
