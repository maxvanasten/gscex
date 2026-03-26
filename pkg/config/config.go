package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type GameVersion string

const (
	GameT5 GameVersion = "t5"
	GameT6 GameVersion = "t6"
)

type GameConfig struct {
	ScriptsRepo   string `json:"scripts_repo"`
	ScriptsBranch string `json:"scripts_branch"`
}

type Config struct {
	Games        map[string]GameConfig `json:"games"`
	CacheDir     string                `json:"cache_dir"`
	AutoUpdate   bool                  `json:"auto_update"`
	MaxResults   int                   `json:"max_results"`
	ContextLines int                   `json:"context_lines"`
	DefaultGame  string                `json:"default_game"`
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Games: map[string]GameConfig{
			"t5": {
				ScriptsRepo:   "https://github.com/plutoniummod/t5-scripts",
				ScriptsBranch: "main",
			},
			"t6": {
				ScriptsRepo:   "https://github.com/plutoniummod/t6-scripts",
				ScriptsBranch: "main",
			},
		},
		CacheDir:     filepath.Join(home, ".gscex"),
		AutoUpdate:   false,
		MaxResults:   20,
		ContextLines: 3,
		DefaultGame:  "t6",
	}
}

func (c *Config) Path() string {
	return filepath.Join(c.CacheDir, "config.json")
}

func (c *Config) ScriptsPath(game string) string {
	return filepath.Join(c.CacheDir, "scripts-"+game)
}

func (c *Config) IndexPath(game string) string {
	return filepath.Join(c.CacheDir, "index-"+game+".json")
}

func (c *Config) GetGameRepo(game string) (string, string) {
	if gameCfg, ok := c.Games[game]; ok {
		return gameCfg.ScriptsRepo, gameCfg.ScriptsBranch
	}
	return "", ""
}

func (c *Config) Save() error {
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.Path(), data, 0644)
}

func Load() (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(cfg.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, cfg.Save()
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadFrom(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
