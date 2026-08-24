package httpapi

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"net/http"
	"time"
)

type sampleInput struct {
	PlotID      string `json:"plot_id"`
	Metric      string `json:"metric"`
	Value       string `json:"value"`
	Unit        string `json:"unit"`
	CollectedAt string `json:"collected_at"`
}

func (s *Server) samples(w http.ResponseWriter, r *http.Request) {
	u, _ := middleware.CurrentUser(r.Context())
	plot := r.URL.Query().Get("plot_id")
	if r.Method == "GET" {
		items, e := s.Samples.List(r.Context(), plot, 100)
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 200, map[string]any{"items": items})
		return
	}
	if r.Method == "POST" {
		var in sampleInput
		if e := decode(r, &in); e != nil {
			writeErr(w, r, e)
			return
		}
		at := time.Now()
		if in.CollectedAt != "" {
			at, _ = time.Parse(time.RFC3339, in.CollectedAt)
		}
		x, e := s.Samples.Record(r.Context(), u, in.PlotID, in.Metric, in.Value, in.Unit, at)
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 201, x)
		return
	}
	writeErr(w, r, apperror.ErrInvalid)
}
