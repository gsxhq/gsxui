// Package cli implements the gsxui command: init, add, list.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config is gsxui.json: where vendored Go packages, JS, and CSS live,
// relative to the module root. The module path itself is always read from
// go.mod, never stored.
type Config struct {
	UI  string `json:"ui"`
	JS  string `json:"js"`
	CSS string `json:"css"`
}

func DefaultConfig() Config {
	return Config{UI: "ui", JS: "web/gsxui", CSS: "web/gsxui.css"}
}

func LoadConfig(dir string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		return Config{}, fmt.Errorf("gsxui.json not found — run 'gsxui init' first")
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("parsing gsxui.json: %w", err)
	}
	return c, nil
}

func (c Config) Save(dir string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "gsxui.json"), append(data, '\n'), 0o644)
}
