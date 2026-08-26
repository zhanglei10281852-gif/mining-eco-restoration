package httpapi

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var in loginRequest
	if e := decode(r, &in); e != nil {
		writeErr(w, r, e)
		return
	}
	u, sess, e := s.Auth.Login(r.Context(), in.Email, in.Password)
	if e != nil {
		writeErr(w, r, e)
		return
	}
	writeJSON(w, 200, map[string]any{"token": sess.TokenHash, "expires_at": sess.ExpiresAt, "user": map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role}})
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := middleware.CurrentSession(r.Context())
	if e := s.Auth.Logout(r.Context(), sess); e != nil {
		writeErr(w, r, e)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "revoked"})
}
