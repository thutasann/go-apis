package config

import (
	"encoding/json"
	"os"
	"runtime"
)

// Config is loaded once at startup and passed by pointer to all subsystems.
// Zero values in JSON fall back to defaults - users only configure what they care about.
type Config struct {
	Port            int    `json:"port"`
	PagesDir        string `json:"pagesDir"`
	PublicDir       string `json:"publicDir"`
	BuildDir        string `json:"buildDir"`
	WorkerPoolSize  int    `json:"workerPoolSize"`
	Dev             bool   `json:"dev"`
	CacheMaxEntries int    `json:"cacheMaxEntries"`
}

func DefaultConfig() *Config {
	return &Config{
		Port:            3000,
		PagesDir:        "pages",
		PublicDir:       "public",
		BuildDir:        ".build",
		WorkerPoolSize:  runtime.NumCPU(), // one V8 isolate per core - no oversubscription
		Dev:             false,
		CacheMaxEntries: 10000,
	}
}

// Load merges config file on top of defaults.
// Missing file is not an error - defaults are production-ready.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Zero means "auto" - fall back to CPU count
	if cfg.WorkerPoolSize == 0 {
		cfg.WorkerPoolSize = runtime.NumCPU()
	}

	return cfg, nil
}
