package config_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/config"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("PORT", "8080")

	cfg := config.LoadConfig()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port 8080, got %s", cfg.Port)
	}
}
