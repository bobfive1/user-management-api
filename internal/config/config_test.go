package config

import (
	"testing"
	"time"
)

func TestLoadReadsYAMLIntoStruct(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.App.Name != "github.com/bobfive1/user-management-api" {
		t.Fatalf("app name = %q", cfg.App.Name)
	}
	if cfg.ServerAPI.Port != "0.0.0.0:8080" {
		t.Fatalf("server port = %q", cfg.ServerAPI.Port)
	}
	if cfg.ServerAPI.ReadTimeout != 20*time.Second {
		t.Fatalf("read timeout = %s", cfg.ServerAPI.ReadTimeout)
	}
	if cfg.Logging.Level != "info" {
		t.Fatalf("logging level = %q", cfg.Logging.Level)
	}
	if cfg.PostgresConfig.Port != 5432 {
		t.Fatalf("postgres port = %d", cfg.PostgresConfig.Port)
	}

}
