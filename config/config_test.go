package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("POSTGRES_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("ORIGIN_ALLOWED", " http://a.com , ,http://b.com ")
	t.Setenv("GRAPHQL_PLAYGROUND_ENABLED", "true")

	cfg := Load()

	if cfg.Port != "3000" {
		t.Errorf("Port = %q", cfg.Port)
	}
	if want := []string{"http://a.com", "http://b.com"}; !reflect.DeepEqual(cfg.OriginsAllowed, want) {
		t.Errorf("OriginsAllowed = %v, want %v", cfg.OriginsAllowed, want)
	}
	if !cfg.GraphQLPlaygroundEnabled {
		t.Error("GraphQLPlaygroundEnabled = false, want true")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("POSTGRES_URL", "postgres://localhost/db")
	t.Setenv("ORIGIN_ALLOWED", "")
	t.Setenv("GRAPHQL_PLAYGROUND_ENABLED", "")

	cfg := Load()

	if len(cfg.OriginsAllowed) != 0 {
		t.Errorf("OriginsAllowed = %v, want empty", cfg.OriginsAllowed)
	}
	if cfg.GraphQLPlaygroundEnabled {
		t.Error("GraphQLPlaygroundEnabled = true, want false")
	}
}

func TestLoadPanics(t *testing.T) {
	tests := map[string]map[string]string{
		"missing port":     {"PORT": "", "POSTGRES_URL": "postgres://localhost/db"},
		"missing postgres": {"PORT": "3000", "POSTGRES_URL": ""},
		"invalid bool":     {"PORT": "3000", "POSTGRES_URL": "postgres://localhost/db", "GRAPHQL_PLAYGROUND_ENABLED": "maybe"},
	}
	for name, env := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv("GRAPHQL_PLAYGROUND_ENABLED", "")
			for k, v := range env {
				t.Setenv(k, v)
			}
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			Load()
		})
	}
}
