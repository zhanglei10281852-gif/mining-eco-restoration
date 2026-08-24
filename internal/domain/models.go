package domain

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleInspector Role = "inspector"
	RoleOperator  Role = "operator"
)

type User struct {
	ID, Email, Name string
	Role            Role
	PasswordHash    string
	Active          bool
	CreatedAt       time.Time
}
type Session struct {
	ID, UserID, TokenHash string
	ExpiresAt             time.Time
	RevokedAt             *time.Time
	CreatedAt             time.Time
}
type Project struct {
	ID, Code, Name, Region, OwnerID string
	Status                          string
	CreatedAt, UpdatedAt            time.Time
	Version                         int64
}
type Plot struct {
	ID, ProjectID, Name, SoilType string
	AreaM2                        float64
	Status                        string
	CreatedAt                     time.Time
}
type Sample struct {
	ID, PlotID, CollectorID, Metric, Value, Unit string
	CollectedAt                                  time.Time
	CreatedAt                                    time.Time
}
type RemediationTask struct {
	ID, ProjectID, PlotID, AssigneeID, Title, Description, Status string
	DueAt                                                         *time.Time
	Version                                                       int64
	CreatedAt, UpdatedAt                                          time.Time
}
type Inspection struct {
	ID, TaskID, InspectorID, Result, Notes string
	Score                                  int
	CreatedAt                              time.Time
}
type AuditEntry struct {
	ID, ActorID, Action, EntityType, EntityID, RequestID, Details string
	CreatedAt                                                     time.Time
}

func ValidRole(r Role) bool { return r == RoleAdmin || r == RoleInspector || r == RoleOperator }
