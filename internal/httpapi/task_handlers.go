package httpapi

import (
	"encoding/json"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"net/http"
	"strconv"
)

type taskInput struct {
	ProjectID   string `json:"project_id"`
	PlotID      string `json:"plot_id"`
	AssigneeID  string `json:"assignee_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Version     int64  `json:"version"`
}

func (s *Server) tasks(w http.ResponseWriter, r *http.Request) {
	u, _ := middleware.CurrentUser(r.Context())
	if r.Method == "POST" {
		var in taskInput
		if e := decode(r, &in); e != nil {
			writeErr(w, r, e)
			return
		}
		key := r.Header.Get("Idempotency-Key")
		if key != "" {
			if code, body, e := (repository.Idempotency{DB: s.DB.DB}).Get(r.Context(), key, u.ID, r.URL.Path); e == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
				_, _ = w.Write([]byte(body))
				return
			}
		}
		t, e := s.Tasks.Create(r.Context(), u, in.ProjectID, in.PlotID, in.AssigneeID, in.Title, in.Description, middleware.GetRequestID(r.Context()))
		if e != nil {
			writeErr(w, r, e)
			return
		}
		body, _ := json.Marshal(t)
		if key != "" {
			_ = (repository.Idempotency{DB: s.DB.DB}).Put(r.Context(), key, u.ID, r.URL.Path, 201, string(body))
		}
		writeJSON(w, 201, t)
		return
	}
	if r.Method == "PATCH" {
		id := r.URL.Query().Get("id")
		v, _ := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
		var in taskInput
		if e := decode(r, &in); e != nil {
			writeErr(w, r, e)
			return
		}
		t, e := s.Tasks.Transition(r.Context(), u, id, in.Status, v, middleware.GetRequestID(r.Context()))
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 200, t)
		return
	}
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		t, e := s.Tasks.Tasks.ByID(r.Context(), id)
		if e != nil {
			writeErr(w, r, e)
			return
		}
		writeJSON(w, 200, t)
		return
	}
	writeErr(w, r, apperror.ErrInvalid)
}
