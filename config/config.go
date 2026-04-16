package config

import (
	"os"
	"strings"
)

// Config holds all environment-driven configuration for the service.
// Add new fields here as the project grows — never read os.Getenv outside this package.
type Config struct {
	Port           string
	OriginsAllowed []string
	PostgresURL    string
}

// Load reads environment variables and returns a populated Config.
// It panics on missing required values so misconfiguration is caught at startup.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		panic("required environment variable PORT is not set")
	}

	return &Config{
		Port:           port,
		OriginsAllowed: splitCSV(os.Getenv("ORIGIN_ALLOWED")),
		PostgresURL:    os.Getenv("POSTGRES_URL"),
	}
}

func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
