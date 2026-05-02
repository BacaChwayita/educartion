package config_test

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/config"
)

func TestGetConfig(t *testing.T) {
	t.Setenv("PORT", "8080")

	cfg := config.GetConfig()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port 8080, got %s", cfg.Port)
	}

	t.Setenv("PORT", "8081")

	cfg = config.GetConfig()

	if cfg.Port == "8081" {
		t.Error("Expected Port 8080, got 8081. Did not expect GetConfig to reload config (singleton instance)")
	} else if cfg.Port != "8080" {
		t.Errorf("Expected Port 8080, got %s", cfg.Port)
	}
}
