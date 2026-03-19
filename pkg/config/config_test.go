package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.ScriptsRepo != "https://github.com/plutoniummod/t6-scripts" {
		t.Errorf("wrong scripts repo: %s", cfg.ScriptsRepo)
	}
	if cfg.ScriptsBranch != "main" {
		t.Errorf("wrong branch: %s", cfg.ScriptsBranch)
	}
	if cfg.MaxResults != 20 {
		t.Errorf("wrong max results: %d", cfg.MaxResults)
	}
}

func TestConfigPaths(t *testing.T) {
	home, _ := os.UserHomeDir()
	cfg := Default()

	wantCache := filepath.Join(home, ".gscex")
	if cfg.CacheDir != wantCache {
		t.Errorf("CacheDir = %s, want %s", cfg.CacheDir, wantCache)
	}

	wantConfig := filepath.Join(wantCache, "config.json")
	if cfg.Path() != wantConfig {
		t.Errorf("Path() = %s, want %s", cfg.Path(), wantConfig)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{
		CacheDir:    tmpDir,
		MaxResults:  50,
		ScriptsRepo: "https://example.com/test",
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Test LoadFrom with the same path
	cfg2, err := LoadFrom(cfg.Path())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg2.MaxResults != 50 {
		t.Errorf("MaxResults = %d, want 50", cfg2.MaxResults)
	}

	if cfg2.ScriptsRepo != "https://example.com/test" {
		t.Errorf("ScriptsRepo = %s, want https://example.com/test", cfg2.ScriptsRepo)
	}
}
