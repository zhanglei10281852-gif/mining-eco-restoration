package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"time"
)

type Projects struct{ DB *sql.DB }

func (r Projects) Create(ctx context.Context, p domain.Project) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO projects(id,code,name,region,owner_id,status,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, p.ID, p.Code, p.Name, p.Region, p.OwnerID, p.Status, p.Version, p.CreatedAt.UTC().Format(time.RFC3339), p.UpdatedAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("persist project: %w", err)
	}
	return err
}
func (r Projects) ByID(ctx context.Context, id string) (domain.Project, error) {
	var p domain.Project
	var c, u string
	err := r.DB.QueryRowContext(ctx, `SELECT id,code,name,region,owner_id,status,version,created_at,updated_at FROM projects WHERE id=?`, id).Scan(&p.ID, &p.Code, &p.Name, &p.Region, &p.OwnerID, &p.Status, &p.Version, &c, &u)
	p.CreatedAt, _ = time.Parse(time.RFC3339, c)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, u)
	return p, err
}
func (r Projects) List(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id,code,name,region,owner_id,status,version,created_at,updated_at FROM projects ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Project{}
	for rows.Next() {
		var p domain.Project
		var c, u string
		if err = rows.Scan(&p.ID, &p.Code, &p.Name, &p.Region, &p.OwnerID, &p.Status, &p.Version, &c, &u); err != nil {
			return nil, err
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, c)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, u)
		out = append(out, p)
	}
	return out, rows.Err()
}
func (r Projects) AddPlot(ctx context.Context, p domain.Plot) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO plots(id,project_id,name,soil_type,area_m2,status,created_at) VALUES(?,?,?,?,?,?,?)`, p.ID, p.ProjectID, p.Name, p.SoilType, p.AreaM2, p.Status, p.CreatedAt.UTC().Format(time.RFC3339))
	return err
}
func (r Projects) Plot(ctx context.Context, id string) (domain.Plot, error) {
	var p domain.Plot
	var ts string
	err := r.DB.QueryRowContext(ctx, `SELECT id,project_id,name,soil_type,area_m2,status,created_at FROM plots WHERE id=?`, id).Scan(&p.ID, &p.ProjectID, &p.Name, &p.SoilType, &p.AreaM2, &p.Status, &ts)
	p.CreatedAt, _ = time.Parse(time.RFC3339, ts)
	return p, err
}
