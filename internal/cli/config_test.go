package cli

import (
	"strings"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	c := DefaultConfig()
	if err := c.Save(dir); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != c {
		t.Fatalf("got %+v want %+v", got, c)
	}
}

func TestLoadConfigMissing(t *testing.T) {
	_, err := LoadConfig(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "gsxui init") {
		t.Fatalf("want actionable error, got %v", err)
	}
}
