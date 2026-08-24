package httpapi

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/pagination"
	"net/http"
)

type projectInput struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Region string `json:"region"`
}
type plotInput struct {
	ProjectID string  `json:"project_id"`
	Name      string  `json:"name"`
	SoilType  string  `json:"soil_type"`
	AreaM2    float64 `json:"area_m2"`
}

func (s *Server) projects(w http.ResponseWriter, r *http.Request) {
	u, _ := middleware.CurrentUser(r.Context())
	if r.Method == "GET" {
		ps, e := s.Projects.List(r.Context(), pagination.Parse(r.URL.Query().Get("limit"), r.URL.Query().Get("offset")))
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 200, map[string]any{"items": ps})
		return
	}
	if r.Method == "POST" {
		var in projectInput
		if e := decode(r, &in); e != nil {
			writeErr(w, r, e)
			return
		}
		p, e := s.Projects.Create(r.Context(), u, in.Code, in.Name, in.Region)
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 201, p)
		return
	}
	if r.Method == "PUT" && r.URL.Query().Get("plot") != "" {
		var in plotInput
		if e := decode(r, &in); e != nil {
			writeErr(w, r, e)
			return
		}
		p, e := s.Projects.AddPlot(r.Context(), u, in.ProjectID, in.Name, in.SoilType, in.AreaM2)
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 201, p)
		return
	}
	writeErr(w, r, apperror.ErrInvalid)
}
