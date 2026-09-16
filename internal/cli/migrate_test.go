package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateRTLTransformsManagedFilesInPlaceAndIsIdempotent(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "native-select", "table"}); err != nil {
		t.Fatal(err)
	}
	// The consumer edited a vendored file; migrate must keep that edit.
	tablePath := filepath.Join(dir, "ui", "table.gsx")
	table, _ := os.ReadFile(tablePath)
	edited := strings.Replace(string(table), "package ui\n", "package ui\n\n// consumer note\n", 1)
	if err := os.WriteFile(tablePath, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"migrate", "rtl"}); err == nil || !strings.Contains(err.Error(), `"rtl": true`) {
		t.Fatalf("migrate rtl without the flag must refuse and say how to set it, got %v", err)
	}
	cfg, _ := LoadConfig(dir)
	cfg.RTL = true
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"migrate", "rtl"}); err != nil {
		t.Fatal(err)
	}
	ns := readFile(t, dir, "ui/native-select.gsx")
	if strings.Contains(ns, "right-2.5") || !strings.Contains(ns, "end-2.5") {
		t.Fatalf("native-select not migrated:\n%s", ns)
	}
	got := readFile(t, dir, "ui/table.gsx")
	if !strings.Contains(got, "// consumer note") {
		t.Fatal("migrate discarded the consumer's edit")
	}
	if strings.Contains(got, "text-left") || !strings.Contains(got, "text-start") {
		t.Fatal("table not migrated")
	}
	style := readFile(t, dir, "web/gsxui/style.css")
	if strings.Contains(style, "@apply pl-1.5") || !strings.Contains(style, "ps-1.5") {
		t.Fatal("style.css not migrated")
	}
	// Idempotent: a second run changes nothing and reports nothing to do.
	snapshot := map[string]string{}
	for _, rel := range []string{"ui/native-select.gsx", "ui/table.gsx", "web/gsxui/style.css", "gsxui.json"} {
		snapshot[rel] = readFile(t, dir, rel)
	}
	if err := Run([]string{"migrate", "rtl"}); err != nil {
		t.Fatal(err)
	}
	for rel, before := range snapshot {
		if readFile(t, dir, rel) != before {
			t.Errorf("second migrate changed %s", rel)
		}
	}
	// The managed hash now matches the migrated content: add --overwrite is a no-op on native-select.
	if err := Run([]string{"add", "native-select"}); err != nil {
		t.Fatal(err)
	}
	if readFile(t, dir, "ui/native-select.gsx") != snapshot["ui/native-select.gsx"] {
		t.Fatal("add after migrate re-vendored a file whose hash should already match")
	}
}

func TestMigrateListsMigrations(t *testing.T) {
	err := Run([]string{"migrate"})
	if err == nil || !strings.Contains(err.Error(), "rtl") {
		t.Fatalf("bare migrate must name the available migrations, got %v", err)
	}
	if err := Run([]string{"migrate", "nope"}); err == nil {
		t.Fatal("unknown migration accepted")
	}
}
