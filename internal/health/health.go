package health

import (
	"context"
	"database/sql"
)

type Checker struct{ DB *sql.DB }

func (c Checker) Live() bool                     { return true }
func (c Checker) Ready(ctx context.Context) bool { return c.DB.PingContext(ctx) == nil }
