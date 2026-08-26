package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/auth"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/config"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/db"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/health"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/middleware"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/repository"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	Cfg         config.Config
	DB          *db.Store
	Log         *slog.Logger
	Auth        auth.Service
	Projects    service.ProjectService
	Tasks       service.TaskService
	Samples     service.SampleService
	Inspections service.InspectionService
	Mux         *http.ServeMux
	HTTP        *http.Server
}

func New(cfg config.Config, store *db.Store, log *slog.Logger) *Server {
	a := auth.Service{Users: repository.Users{DB: store.DB}, Sessions: repository.Sessions{DB: store.DB}, TTL: 24 * time.Hour}
	_ = a.EnsureDefaults(context.Background())
	s := &Server{Cfg: cfg, DB: store, Log: log, Auth: a, Projects: service.ProjectService{DB: store.DB, Projects: repository.Projects{DB: store.DB}, Audit: repository.Audit{DB: store.DB}}, Tasks: service.TaskService{DB: store.DB, Projects: repository.Projects{DB: store.DB}, Tasks: repository.Tasks{DB: store.DB}, Audit: repository.Audit{DB: store.DB}}, Samples: service.SampleService{Projects: repository.Projects{DB: store.DB}, Samples: repository.Samples{DB: store.DB}}, Inspections: service.InspectionService{DB: store.DB, Tasks: repository.Tasks{DB: store.DB}, Inspections: repository.Inspections{DB: store.DB}, Audit: repository.Audit{DB: store.DB}}, Mux: http.NewServeMux()}
	s.routes()
	return s
}
func (s *Server) routes() {
	s.Mux.Handle("/healthz", middleware.RequestID(http.HandlerFunc(s.health)))
	s.Mux.Handle("/readyz", middleware.RequestID(http.HandlerFunc(s.ready)))
	s.Mux.HandleFunc("/api/auth/login", s.login)
	s.Mux.Handle("/api/auth/logout", middleware.RequestID(middleware.Require(s.Auth, http.HandlerFunc(s.logout))))
	protected := middleware.RequestID(middleware.Logging(s.Log, middleware.Require(s.Auth, s.Mux)))
	_ = protected
	s.Mux.Handle("/api/projects", middleware.RequestID(middleware.Logging(s.Log, middleware.Require(s.Auth, http.HandlerFunc(s.projects)))))
	s.Mux.Handle("/api/tasks", middleware.RequestID(middleware.Logging(s.Log, middleware.Require(s.Auth, http.HandlerFunc(s.tasks)))))
	s.Mux.Handle("/api/samples", middleware.RequestID(middleware.Logging(s.Log, middleware.Require(s.Auth, http.HandlerFunc(s.samples)))))
	s.Mux.Handle("/api/inspections", middleware.RequestID(middleware.Logging(s.Log, middleware.Require(s.Auth, middleware.Role(domain.RoleInspector, domain.RoleAdmin)(http.HandlerFunc(s.inspections))))))
}
func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{Addr: s.Cfg.Address, Handler: s.Mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	s.HTTP = srv
	go func() { <-ctx.Done(); _ = srv.Shutdown(context.Background()) }()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-ch:
			s.Log.Info("shutdown signal", "signal", sig)
			_ = srv.Shutdown(context.Background())
		case <-ctx.Done():
		}
	}()
	s.Log.Info("server listening", "address", s.Cfg.Address)
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, r *http.Request, e error) {
	status := 500
	code := "internal_error"
	switch {
	case errors.Is(e, apperror.ErrInvalid):
		status = 400
		code = "invalid_request"
	case errors.Is(e, apperror.ErrUnauthorized):
		status = 401
		code = "unauthorized"
	case errors.Is(e, apperror.ErrForbidden):
		status = 403
		code = "forbidden"
	case errors.Is(e, apperror.ErrNotFound), errors.Is(e, sql.ErrNoRows):
		status = 404
		code = "not_found"
	case errors.Is(e, apperror.ErrConflict):
		status = 409
		code = "conflict"
	}
	writeJSON(w, status, map[string]any{"code": code, "message": e.Error(), "request_id": middleware.GetRequestID(r.Context())})
}
func decode(r *http.Request, v any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return apperror.ErrInvalid
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("%w: %v", apperror.ErrInvalid, err)
	}
	return nil
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok"})
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if (health.Checker{DB: s.DB.DB}).Ready(r.Context()) {
		writeJSON(w, 200, map[string]any{"status": "ready"})
		return
	}
	writeJSON(w, 503, map[string]any{"status": "not_ready"})
}
