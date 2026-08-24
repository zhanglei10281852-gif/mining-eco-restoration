package middleware_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDGeneratedAndPreserved(t *testing.T) {
	h := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleware.GetRequestID(r.Context()) == "" {
			t.Fatal("missing context id")
		}
	}))
	for _, id := range []string{"", "fixed"} {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		if id != "" {
			r.Header.Set("X-Request-ID", id)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Header().Get("X-Request-ID") == "" {
			t.Fatal("missing header")
		}
		if id != "" && w.Header().Get("X-Request-ID") != id {
			t.Fatal("changed id")
		}
	}
}
