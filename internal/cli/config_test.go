package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	c := DefaultConfig()
	c.Managed["ui/button.gsx"] = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if err := c.Save(dir); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, c) {
		t.Fatalf("got %+v want %+v", got, c)
	}
}

func TestConfigSaveCanonicalizesManagedOrdering(t *testing.T) {
	dir := t.TempDir()
	c := DefaultConfig()
	c.Managed["web/z.js"] = strings.Repeat("b", 64)
	c.Managed["ui/a.gsx"] = strings.Repeat("a", 64)

	if err := c.Save(dir); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		t.Fatal(err)
	}
	want := `{
  "ui": "ui",
  "js": "web/gsxui",
  "css": "web/gsxui/index.css",
  "managed": {
    "ui/a.gsx": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    "web/z.js": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  }
}
`
	if string(got) != want {
		t.Fatalf("config bytes:\n%s\nwant:\n%s", got, want)
	}
}

func TestLoadConfigLegacyInitializesManagedWithoutWriting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "gsxui.json")
	legacy := []byte(`{"ui":"ui","js":"web/gsxui","css":"web/gsxui/index.css"}`)
	if err := os.WriteFile(path, legacy, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Managed == nil || len(got.Managed) != 0 {
		t.Fatalf("Managed = %#v, want initialized empty map", got.Managed)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(legacy) {
		t.Fatalf("LoadConfig rewrote legacy config:\n%s", after)
	}
}

func TestConfigRejectsInvalidManagedEntries(t *testing.T) {
	validHash := strings.Repeat("a", 64)
	tests := []struct {
		name string
		path string
		hash string
	}{
		{name: "empty path", path: "", hash: validHash},
		{name: "absolute path", path: "/tmp/button.gsx", hash: validHash},
		{name: "escape", path: "../button.gsx", hash: validHash},
		{name: "unclean alias", path: "ui/../ui/button.gsx", hash: validHash},
		{name: "backslash alias", path: `ui\button.gsx`, hash: validHash},
		{name: "NUL", path: "ui/\x00button.gsx", hash: validHash},
		{name: "invalid UTF-8", path: "ui/\xffbutton.gsx", hash: validHash},
		{name: "config self hash", path: "gsxui.json", hash: validHash},
		{name: "short hash", path: "ui/button.gsx", hash: "abc"},
		{name: "uppercase hash", path: "ui/button.gsx", hash: strings.Repeat("A", 64)},
		{name: "non hex hash", path: "ui/button.gsx", hash: strings.Repeat("g", 64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.Managed[tt.path] = tt.hash
			if err := cfg.Save(t.TempDir()); err == nil {
				t.Fatal("Save accepted invalid managed entry")
			}
		})
	}
}

func TestConfigRejectsInvalidProjectPathsOnSaveAndLoad(t *testing.T) {
	tests := []struct {
		name  string
		field string
		value string
	}{
		{name: "empty UI", field: "ui", value: ""},
		{name: "unclean UI", field: "ui", value: "ui/../components"},
		{name: "escaping JS", field: "js", value: "../web"},
		{name: "absolute CSS", field: "css", value: "/tmp/index.css"},
		{name: "backslash CSS", field: "css", value: `web\index.css`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			switch tt.field {
			case "ui":
				cfg.UI = tt.value
			case "js":
				cfg.JS = tt.value
			case "css":
				cfg.CSS = tt.value
			default:
				t.Fatalf("unknown field %q", tt.field)
			}
			if err := cfg.Save(t.TempDir()); err == nil {
				t.Fatal("Save accepted invalid project path")
			}

			dir := t.TempDir()
			raw, err := json.Marshal(cfg)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "gsxui.json"), raw, 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadConfig(dir); err == nil {
				t.Fatal("LoadConfig accepted invalid project path")
			}
		})
	}
}

func TestConfigRejectsSymlinkEscapeOnLoadAndSave(t *testing.T) {
	for _, operation := range []string{"load", "save"} {
		t.Run(operation, func(t *testing.T) {
			dir := t.TempDir()
			outside := filepath.Join(t.TempDir(), "gsxui.json")
			if err := os.WriteFile(outside, []byte(`{"ui":"ui","js":"web/gsxui","css":"web/gsxui/index.css"}`), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, filepath.Join(dir, "gsxui.json")); err != nil {
				t.Fatal(err)
			}

			var err error
			if operation == "load" {
				_, err = LoadConfig(dir)
			} else {
				err = DefaultConfig().Save(dir)
			}
			if err == nil || !strings.Contains(err.Error(), "escapes module root") {
				t.Fatalf("%s error = %v, want symlink escape", operation, err)
			}
		})
	}
}

func TestLoadConfigMissing(t *testing.T) {
	_, err := LoadConfig(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "gsxui init") {
		t.Fatalf("want actionable error, got %v", err)
	}
}
