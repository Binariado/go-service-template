package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all environment-driven configuration for the service.
// Add new fields here as the project grows — never read os.Getenv outside this package.
type Config struct {
	Port                     string
	OriginsAllowed           []string
	PostgresURL              string
	GraphQLPlaygroundEnabled bool
}

// Load reads environment variables and returns a populated Config.
// It panics on missing or invalid required values so misconfiguration is caught at startup.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		panic("required environment variable PORT is not set")
	}

	postgresURL := strings.TrimSpace(os.Getenv("POSTGRES_URL"))
	if postgresURL == "" {
		panic("required environment variable POSTGRES_URL is not set")
	}

	return &Config{
		Port:                     port,
		OriginsAllowed:           splitCSV(os.Getenv("ORIGIN_ALLOWED")),
		PostgresURL:              postgresURL,
		GraphQLPlaygroundEnabled: parseBool("GRAPHQL_PLAYGROUND_ENABLED"),
	}
}

// splitCSV splits a comma-separated value, trimming spaces and dropping empty items.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			values = append(values, value)
		}
	}
	return values
}

// parseBool reads an optional boolean variable. Unset means false; an invalid
// value panics so typos never silently change behavior.
func parseBool(name string) bool {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return false
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		panic(name + " must be true or false")
	}
	return parsed
}
