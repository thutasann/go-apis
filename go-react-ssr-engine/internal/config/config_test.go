package config

import (
	"os"
	"runtime"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Port)
	}
	if cfg.WorkerPoolSize != runtime.NumCPU() {
		t.Errorf("expected %d workers, got %d", runtime.NumCPU(), cfg.WorkerPoolSize)
	}
	if cfg.Dev != false {
		t.Error("expected dev=false by default")
	}
}

func TestLoadMissing(t *testing.T) {
	// Missing config file should return defaults, not error
	cfg, err := Load("nonexistent.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected default port, got %d", cfg.Port)
	}
}

func TestLoadOverrides(t *testing.T) {
	// Write a temp config that overrides port only
	content := []byte(`{"port": 8080}`)
	tmpFile := "test_config.json"
	os.WriteFile(tmpFile, content, 0644)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Port overridden
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}

	// Other fields keep defaults
	if cfg.PagesDir != "pages" {
		t.Errorf("expected default pagesDir, got %s", cfg.PagesDir)
	}
}

func TestLoadZeroWorkers(t *testing.T) {
	// workerPoolSize=0 in JSON should fall back to NumCPU
	content := []byte(`{"workerPoolSize": 0}`)
	tmpFile := "test_zero_workers.json"
	os.WriteFile(tmpFile, content, 0644)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if cfg.WorkerPoolSize != runtime.NumCPU() {
		t.Errorf("expected %d workers for zero value, got %d", runtime.NumCPU(), cfg.WorkerPoolSize)
	}
}
