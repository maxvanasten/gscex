package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	// Check T6 game config
	t6Game, ok := cfg.Games["t6"]
	if !ok {
		t.Error("t6 game not found in config")
	}
	if t6Game.ScriptsRepo != "https://github.com/plutoniummod/t6-scripts" {
		t.Errorf("wrong t6 scripts repo: %s", t6Game.ScriptsRepo)
	}
	if t6Game.ScriptsBranch != "main" {
		t.Errorf("wrong t6 branch: %s", t6Game.ScriptsBranch)
	}

	// Check T5 game config
	t5Game, ok := cfg.Games["t5"]
	if !ok {
		t.Error("t5 game not found in config")
	}
	if t5Game.ScriptsRepo != "https://github.com/plutoniummod/t5-scripts" {
		t.Errorf("wrong t5 scripts repo: %s", t5Game.ScriptsRepo)
	}
	if t5Game.ScriptsBranch != "main" {
		t.Errorf("wrong t5 branch: %s", t5Game.ScriptsBranch)
	}

	if cfg.MaxResults != 10000 {
		t.Errorf("wrong max results: %d", cfg.MaxResults)
	}

	if cfg.DefaultGame != "t6" {
		t.Errorf("wrong default game: %s", cfg.DefaultGame)
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

	// Test game-specific paths
	wantScriptsT6 := filepath.Join(wantCache, "scripts-t6")
	if cfg.ScriptsPath("t6") != wantScriptsT6 {
		t.Errorf("ScriptsPath(t6) = %s, want %s", cfg.ScriptsPath("t6"), wantScriptsT6)
	}

	wantIndexT6 := filepath.Join(wantCache, "index-t6.json")
	if cfg.IndexPath("t6") != wantIndexT6 {
		t.Errorf("IndexPath(t6) = %s, want %s", cfg.IndexPath("t6"), wantIndexT6)
	}
}

func TestGetGameRepo(t *testing.T) {
	cfg := Default()

	repo, branch := cfg.GetGameRepo("t6")
	if repo != "https://github.com/plutoniummod/t6-scripts" {
		t.Errorf("GetGameRepo(t6) repo = %s, want t6-scripts", repo)
	}
	if branch != "main" {
		t.Errorf("GetGameRepo(t6) branch = %s, want main", branch)
	}

	repo, branch = cfg.GetGameRepo("t5")
	if repo != "https://github.com/plutoniummod/t5-scripts" {
		t.Errorf("GetGameRepo(t5) repo = %s, want t5-scripts", repo)
	}
	if branch != "main" {
		t.Errorf("GetGameRepo(t5) branch = %s, want main", branch)
	}

	repo, branch = cfg.GetGameRepo("unknown")
	if repo != "" || branch != "" {
		t.Error("GetGameRepo(unknown) should return empty strings")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{
		CacheDir:   tmpDir,
		MaxResults: 50,
		Games: map[string]GameConfig{
			"t6": {
				ScriptsRepo:   "https://example.com/t6-test",
				ScriptsBranch: "test-branch",
			},
		},
		DefaultGame: "t6",
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

	t6Game, ok := cfg2.Games["t6"]
	if !ok {
		t.Fatal("t6 game not found in loaded config")
	}

	if t6Game.ScriptsRepo != "https://example.com/t6-test" {
		t.Errorf("ScriptsRepo = %s, want https://example.com/t6-test", t6Game.ScriptsRepo)
	}

	if t6Game.ScriptsBranch != "test-branch" {
		t.Errorf("ScriptsBranch = %s, want test-branch", t6Game.ScriptsBranch)
	}

	if cfg2.DefaultGame != "t6" {
		t.Errorf("DefaultGame = %s, want t6", cfg2.DefaultGame)
	}
}
