package config_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/config"
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("APP_ADDRESS")
	os.Unsetenv("DATABASE_URL")
	c := config.Load()
	if c.Address != ":8080" {
		t.Fatal(c.Address)
	}
	if c.DatabaseURL == "" || c.SessionTTL == "" {
		t.Fatal(c)
	}
}
func TestLoadEnvironment(t *testing.T) {
	t.Setenv("APP_ADDRESS", ":9999")
	t.Setenv("DATABASE_URL", "file:test.db")
	t.Setenv("SESSION_TTL", "5m")
	c := config.Load()
	if c.Address != ":9999" || c.DatabaseURL != "file:test.db" || c.SessionTTL != "5m" {
		t.Fatal(c)
	}
}
