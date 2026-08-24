package middleware

import (
	"context"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"net/http"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = ids.New("req")
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIDKey, id)))
	})
}
func GetRequestID(ctx context.Context) string { v, _ := ctx.Value(RequestIDKey).(string); return v }
