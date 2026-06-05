package config_test

import (
	"os"
	"testing"

	"github.com/franciscomxs/simulator/internal/infrastructure/config"
)

func TestLoad_DefaultPort(t *testing.T) {
	os.Unsetenv("PORT")
	cfg := config.Load()
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
}

func TestLoad_ValidPort(t *testing.T) {
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("PORT")
	cfg := config.Load()
	if cfg.Port != "3000" {
		t.Errorf("expected port 3000, got %s", cfg.Port)
	}
}

func TestLoad_InvalidPortFallsBack(t *testing.T) {
	os.Setenv("PORT", "99999")
	defer os.Unsetenv("PORT")
	cfg := config.Load()
	if cfg.Port != "8080" {
		t.Errorf("expected fallback to 8080 for out-of-range port, got %s", cfg.Port)
	}
}

func TestLoad_NonNumericPortFallsBack(t *testing.T) {
	os.Setenv("PORT", "abc")
	defer os.Unsetenv("PORT")
	cfg := config.Load()
	if cfg.Port != "8080" {
		t.Errorf("expected fallback to 8080 for non-numeric port, got %s", cfg.Port)
	}
}
