package config

import "os"

type Config struct{ Address, DatabaseURL, SessionTTL string }

func Load() Config {
	c := Config{Address: ":8080", DatabaseURL: "file:mining-eco.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", SessionTTL: "24h"}
	if v := os.Getenv("APP_ADDRESS"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		c.DatabaseURL = v
	}
	if v := os.Getenv("SESSION_TTL"); v != "" {
		c.SessionTTL = v
	}
	return c
}
