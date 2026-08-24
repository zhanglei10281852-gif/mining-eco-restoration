package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/config"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/httpapi"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/migrations"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/testsupport"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func server(t *testing.T) *httpapi.Server {
	store, _ := testsupport.Open(t)
	_ = migrations.Apply(context.Background(), store)
	return httpapi.New(config.Config{Address: ":0", DatabaseURL: ""}, store, slog.Default())
}
func request(t *testing.T, h http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	var b bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&b).Encode(body)
	}
	r := httptest.NewRequest(method, path, &b)
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func TestHealthReady(t *testing.T) {
	s := server(t)
	for _, path := range []string{"/healthz", "/readyz"} {
		w := request(t, s.Mux, http.MethodGet, path, nil, "")
		if w.Code != 200 {
			t.Fatalf("%s=%d", path, w.Code)
		}
	}
}
func login(t *testing.T, s *httpapi.Server, email, pwd string) string {
	w := request(t, s.Mux, http.MethodPost, "/api/auth/login", map[string]string{"email": email, "password": pwd}, "")
	if w.Code != 200 {
		t.Fatalf("login=%d %s", w.Code, w.Body.String())
	}
	var out struct {
		Token string `json:"token"`
	}
	if e := json.Unmarshal(w.Body.Bytes(), &out); e != nil {
		t.Fatal(e)
	}
	return out.Token
}
func TestLoginLogoutFlow(t *testing.T) {
	s := server(t)
	tok := login(t, s, "admin@eco.local", "admin123")
	w := request(t, s.Mux, http.MethodGet, "/api/projects", nil, tok)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	w = request(t, s.Mux, http.MethodPost, "/api/auth/logout", nil, tok)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	w = request(t, s.Mux, http.MethodGet, "/api/projects", nil, tok)
	if w.Code != 401 {
		t.Fatalf("revoked=%d", w.Code)
	}
}
func TestProjectAndPlotContract(t *testing.T) {
	s := server(t)
	tok := login(t, s, "admin@eco.local", "admin123")
	w := request(t, s.Mux, http.MethodPost, "/api/projects", map[string]any{"code": "P-1", "name": "西沟", "region": "陕西"}, tok)
	if w.Code != 201 {
		t.Fatal(w.Code, w.Body.String())
	}
	var p struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &p)
	w = request(t, s.Mux, http.MethodPut, "/api/projects?plot=1", map[string]any{"project_id": p.ID, "name": "一号地块", "soil_type": "砂", "area_m2": 10}, tok)
	if w.Code != 201 {
		t.Fatal(w.Code, w.Body.String())
	}
}
func TestRoleForbidden(t *testing.T) {
	s := server(t)
	tok := login(t, s, "operator@eco.local", "operate123")
	w := request(t, s.Mux, http.MethodPost, "/api/inspections", map[string]any{"task_id": "x", "result": "pass", "score": 90}, tok)
	if w.Code != 403 {
		t.Fatalf("operator inspection=%d", w.Code)
	}
}
func TestInvalidContentType(t *testing.T) {
	s := server(t)
	r := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"admin@eco.local"}`))
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
}
func TestRequestIDPropagates(t *testing.T) {
	s := server(t)
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.Header.Set("X-Request-ID", "req-fixed")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, r)
	if w.Header().Get("X-Request-ID") != "req-fixed" {
		t.Fatal(w.Header())
	}
}
func TestMalformedJSON(t *testing.T) {
	s := server(t)
	tok := login(t, s, "admin@eco.local", "admin123")
	r := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewBufferString("{"))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
	_, _ = io.Copy(io.Discard, w.Result().Body)
}
