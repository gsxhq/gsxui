// Package cli implements the gsxui command: init, add, list.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

// errConfigNotFound sentinels a true absence of gsxui.json, distinguishing
// it from other read failures (permissions) or a parse error. Callers that
// need to tell "not initialized yet" apart from "broken" use errors.Is
// against this.
var errConfigNotFound = fmt.Errorf("gsxui.json not found — run 'gsxui init' first")

// Config is gsxui.json: where vendored Go packages, JS, and CSS live,
// relative to the module root. The module path itself is always read from
// go.mod, never stored.
type Config struct {
	UI      string            `json:"ui"`
	JS      string            `json:"js"`
	CSS     string            `json:"css"`
	Managed map[string]string `json:"managed,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		UI:      "ui",
		JS:      "web/gsxui",
		CSS:     "web/gsxui/index.css",
		Managed: make(map[string]string),
	}
}

func LoadConfig(dir string) (Config, error) {
	configPath, err := artifactPath(dir, "gsxui.json")
	if err != nil {
		return Config{}, fmt.Errorf("reading gsxui.json: %w", err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, errConfigNotFound
		}
		return Config{}, fmt.Errorf("reading gsxui.json: %w", err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("parsing gsxui.json: %w", err)
	}
	if c.Managed == nil {
		c.Managed = make(map[string]string)
	}
	if err := validateConfig(c); err != nil {
		return Config{}, fmt.Errorf("parsing gsxui.json: %w", err)
	}
	return c, nil
}

func (c Config) Save(dir string) error {
	data, err := c.canonicalJSON()
	if err != nil {
		return fmt.Errorf("saving gsxui.json: %w", err)
	}
	configPath, err := artifactPath(dir, "gsxui.json")
	if err != nil {
		return fmt.Errorf("saving gsxui.json: %w", err)
	}
	return os.WriteFile(configPath, data, 0o644)
}

func (c Config) canonicalJSON() ([]byte, error) {
	if err := validateConfig(c); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func validateConfig(c Config) error {
	for _, field := range []struct {
		name  string
		value string
	}{
		{name: "ui", value: c.UI},
		{name: "js", value: c.JS},
		{name: "css", value: c.CSS},
	} {
		if err := validateManagedPath(field.value); err != nil {
			return fmt.Errorf("%s path %q: %w", field.name, field.value, err)
		}
	}
	return validateManaged(c.Managed)
}
