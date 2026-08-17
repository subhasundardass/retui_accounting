// config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration.
//
// Fields are tagged for both YAML and JSON so either config.yml or
// config.json can be dropped next to the binary and read without code
// changes.
type Config struct {
	AppName    string `json:"app_name"     yaml:"app_name"`
	Version    string `json:"version"      yaml:"version"`
	Debug      bool   `json:"debug"        yaml:"debug"`
	LogLevel   string `json:"log_level"    yaml:"log_level"`
	SQLitePath string `json:"sqlite_path"  yaml:"sqlite_path"`

	// App-level fields present in config.yml but previously dropped on
	// the floor because Load() never parsed YAML at all.
	DateFormat string `json:"date_format" yaml:"date_format"`
	Timezone   string `json:"timezone"    yaml:"timezone"`
}

var (
	config     *Config
	configPath string // absolute path actually loaded, "" if defaults only
	once       sync.Once
	mu         sync.RWMutex
)

// defaults returns a Config seeded with safe fallback values. It is the
// baseline that a discovered config file is merged on top of, and the
// final result if no config file is found at all.
func defaults() *Config {
	return &Config{
		AppName:    "Accountant",
		Version:    "1.0.0",
		Debug:      false,
		LogLevel:   "info",
		SQLitePath: "./data/retui.db",

		DateFormat: "DD/MM/YYYY",
		Timezone:   "UTC",
	}
}

// Load loads configuration once per process and caches the result.
// Safe for concurrent use. Call Reload (typically only in tests) to
// force re-reading from disk.
//
// Resolution order (first match wins):
//  1. $RETUI_CONFIG — explicit path to a .yml/.yaml/.json file
//  2. <executable dir>/config.yml (or .yaml)
//  3. <executable dir>/config.json
//  4. <working dir>/config.yml (or .yaml)   — convenient for `go run`
//  5. <working dir>/config.json
//  6. built-in defaults, no file required
//
// Anchoring the search on the executable's own directory (see
// ExecutableDir) is what makes `make build` → build/{binary,config.yml}
// work correctly no matter where the binary is invoked from.
//
// A handful of fields can be overridden via environment variables for
// deployment flexibility without editing the shipped config file:
//   - DEBUG              -> Debug
//   - RETUI_SQLITE_PATH  -> SQLitePath
func Load() *Config {
	mu.RLock()
	if config != nil {
		cfg := config
		mu.RUnlock()
		return cfg
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if config == nil {
		config, configPath = load()
	}

	return config
}

// Reload forces the configuration to be re-read from disk, bypassing the
// singleton cache. Intended for tests and hot-reload tooling — regular
// application code should use Load().
func Reload() *Config {
	mu.Lock()
	cfg, path := load()
	config = cfg
	configPath = path
	mu.Unlock()
	return cfg
}

// Path returns the absolute path of the config file that was actually
// loaded, or "" if no file was found and built-in defaults are in use.
func Path() string {
	mu.RLock()
	defer mu.RUnlock()
	return configPath
}

func load() (*Config, string) {
	cfg := defaults()

	path, format := resolveConfigPath()
	if path != "" {
		if err := loadInto(cfg, path, format); err != nil {
			// Config discovery found a file but couldn't parse it. Fall
			// back to defaults rather than crash the app on a typo'd
			// YAML file — but surface it loudly since the resulting
			// config being "wrong" is otherwise silent and confusing.
			fmt.Fprintf(os.Stderr, "config: failed to load %s: %v (falling back to defaults)\n", path, err)
			cfg = defaults()
			path = ""
		}
	}

	applyEnvOverrides(cfg)
	cfg.SQLitePath = resolveDataPath(cfg.SQLitePath, path)

	return cfg, path
}

// resolveConfigPath searches the well-known locations for a config file
// and returns its path plus a format hint ("yaml" or "json").
func resolveConfigPath() (path string, format string) {
	if override := os.Getenv("RETUI_CONFIG"); override != "" {
		return override, formatFromExt(override)
	}

	candidates := []struct {
		name   string
		format string
	}{
		{"config.yml", "yaml"},
		{"config.yaml", "yaml"},
		{"config.json", "json"},
	}

	for _, dir := range []string{ExecutableDir(), fallbackWd()} {
		for _, c := range candidates {
			p := filepath.Join(dir, c.name)
			if fileExists(p) {
				return p, c.format
			}
		}
	}

	return "", ""
}

func formatFromExt(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yml", ".yaml":
		return "yaml"
	case ".json":
		return "json"
	default:
		return "yaml"
	}
}

// loadInto reads path and unmarshals it onto cfg, so any field the file
// doesn't set keeps its current (default) value.
func loadInto(cfg *Config, path, format string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	switch format {
	case "json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("parse json: %w", err)
		}
	default:
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("parse yaml: %w", err)
		}
	}

	return nil
}

func applyEnvOverrides(cfg *Config) {
	if v, ok := os.LookupEnv("DEBUG"); ok {
		cfg.Debug = v == "true" || v == "1"
	}
	if v := os.Getenv("RETUI_SQLITE_PATH"); v != "" {
		cfg.SQLitePath = v
	}
}

// resolveDataPath anchors a relative SQLitePath to the directory the
// config file was loaded from (falling back to the executable's own
// directory when running purely on defaults). Absolute paths pass
// through unchanged.
//
// This is what keeps the sqlite file inside build/data/ when the app is
// launched as build/retui, instead of creating data/ wherever the
// caller's shell happened to have as its working directory.
func resolveDataPath(sqlitePath, configFilePath string) string {
	if filepath.IsAbs(sqlitePath) {
		return sqlitePath
	}

	anchor := ExecutableDir()
	if configFilePath != "" {
		anchor = filepath.Dir(configFilePath)
	}

	return filepath.Join(anchor, sqlitePath)
}
