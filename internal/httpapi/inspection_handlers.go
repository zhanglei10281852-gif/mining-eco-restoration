package httpapi

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"net/http"
)

type inspectionInput struct {
	TaskID string `json:"task_id"`
	Result string `json:"result"`
	Score  int    `json:"score"`
	Notes  string `json:"notes"`
}

func (s *Server) inspections(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeErr(w, r, apperror.ErrInvalid)
		return
	}
	u, _ := middleware.CurrentUser(r.Context())
	var in inspectionInput
	if e := decode(r, &in); e != nil {
		writeErr(w, r, e)
		return
	}
	x, e := s.Inspections.Accept(r.Context(), u, in.TaskID, in.Result, in.Notes, in.Score, middleware.GetRequestID(r.Context()))
	if e != nil {
		writeErr(w, r, e)
		return
	}
	writeJSON(w, 201, x)
}
