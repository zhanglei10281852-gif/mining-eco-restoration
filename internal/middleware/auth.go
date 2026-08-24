package middleware

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/auth"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"net/http"
)

const userKey ctxKey = "user"
const sessionKey ctxKey = "session"

func Require(a auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, s, e := a.Authenticate(r.Context(), r.Header.Get("Authorization"))
		if e != nil {
			http.Error(w, apperror.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(context.WithValue(r.Context(), userKey, u), sessionKey, s)))
	})
}
func Role(roles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := r.Context().Value(userKey).(domain.User)
			if !ok {
				http.Error(w, apperror.ErrUnauthorized.Error(), 401)
				return
			}
			for _, role := range roles {
				if u.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, apperror.ErrForbidden.Error(), 403)
		})
	}
}
func CurrentUser(ctx context.Context) (domain.User, bool) {
	u, ok := ctx.Value(userKey).(domain.User)
	return u, ok
}
func CurrentSession(ctx context.Context) (domain.Session, bool) {
	s, ok := ctx.Value(sessionKey).(domain.Session)
	return s, ok
}
