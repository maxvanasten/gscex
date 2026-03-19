package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ScriptsRepo   string `json:"scripts_repo"`
	ScriptsBranch string `json:"scripts_branch"`
	CacheDir      string `json:"cache_dir"`
	AutoUpdate    bool   `json:"auto_update"`
	MaxResults    int    `json:"max_results"`
	ContextLines  int    `json:"context_lines"`
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		ScriptsRepo:   "https://github.com/plutoniummod/t6-scripts",
		ScriptsBranch: "main",
		CacheDir:      filepath.Join(home, ".gscex"),
		AutoUpdate:    false,
		MaxResults:    20,
		ContextLines:  3,
	}
}

func (c *Config) Path() string {
	return filepath.Join(c.CacheDir, "config.json")
}

func (c *Config) ScriptsPath() string {
	return filepath.Join(c.CacheDir, "scripts")
}

func (c *Config) IndexPath() string {
	return filepath.Join(c.CacheDir, "index.json")
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
