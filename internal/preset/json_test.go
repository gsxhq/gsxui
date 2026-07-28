package preset

import (
	"bytes"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestParseJSONRejectsDuplicateKeysAtEveryObjectLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		path string
	}{
		{
			name: "top level",
			src:  `{"$schema":"a","$schema":"b"}`,
			path: "$schema",
		},
		{
			name: "theme",
			src:  `{"theme":{"light":{},"light":{}}}`,
			path: "theme.light",
		},
		{
			name: "dark mode",
			src:  `{"theme":{"dark":{},"dark":{}}}`,
			path: "theme.dark",
		},
		{
			name: "light token",
			src:  `{"theme":{"light":{"primary":"red","primary":"blue"}}}`,
			path: "theme.light.primary",
		},
		{
			name: "dark token",
			src:  `{"theme":{"dark":{"primary":"red","primary":"blue"}}}`,
			path: "theme.dark.primary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseJSON([]byte(tt.src))
			if err == nil {
				t.Fatal("ParseJSON returned nil error")
			}
			if !strings.Contains(err.Error(), tt.path) ||
				!strings.Contains(err.Error(), "duplicate") {
				t.Fatalf("error = %q, want duplicate path %q", err, tt.path)
			}
		})
	}

	t.Run("theme property", func(t *testing.T) {
		golden := readTestFile(t, "testdata/default-nova.json")
		closing := bytes.LastIndex(golden, []byte("\n}"))
		if closing == -1 {
			t.Fatal("golden fixture has no top-level object end")
		}
		src := slices.Concat(
			golden[:closing],
			[]byte(`, "theme": {}`),
			golden[closing:],
		)
		_, err := ParseJSON(src)
		if err == nil ||
			!strings.Contains(err.Error(), "theme") ||
			!strings.Contains(err.Error(), "duplicate") {
			t.Fatalf("error = %v, want duplicate theme", err)
		}
	})
}

func TestParseJSONRejectsUnknownKeysAtEveryObjectLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		path string
	}{
		{name: "top level", src: `{"future":true}`, path: "future"},
		{name: "theme", src: `{"theme":{"future":{}}}`, path: "theme.future"},
		{name: "light token", src: `{"theme":{"light":{"future":"red"}}}`, path: "theme.light.future"},
		{name: "dark token", src: `{"theme":{"dark":{"future":"red"}}}`, path: "theme.dark.future"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseJSON([]byte(tt.src))
			if err == nil {
				t.Fatal("ParseJSON returned nil error")
			}
			if !strings.Contains(err.Error(), tt.path) ||
				!strings.Contains(err.Error(), "unknown") {
				t.Fatalf("error = %q, want unknown path %q", err, tt.path)
			}
		})
	}
}

func TestParseJSONRejectsMissingRequiredFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		path string
	}{
		{name: "schema", src: `{}`, path: "$schema"},
		{
			name: "schema version",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json"}`,
			path: "schemaVersion",
		},
		{
			name: "style",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1}`,
			path: "style",
		},
		{
			name: "radius",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova"}`,
			path: "radius",
		},
		{
			name: "theme",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova","radius":"0"}`,
			path: "theme",
		},
		{
			name: "light mode",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova","radius":"0","theme":{"dark":{}}}`,
			path: "theme.light",
		},
		{
			name: "dark mode",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova","radius":"0","theme":{"light":{}}}`,
			path: "theme.dark",
		},
		{
			name: "light token",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova","radius":"0","theme":{"light":{},"dark":{}}}`,
			path: "theme.light.background",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseJSON([]byte(tt.src))
			if err == nil {
				t.Fatal("ParseJSON returned nil error")
			}
			if !strings.Contains(err.Error(), tt.path) ||
				!strings.Contains(err.Error(), "missing") {
				t.Fatalf("error = %q, want missing path %q", err, tt.path)
			}
		})
	}
}

func TestParseJSONRejectsTypeMismatchesAndFutureVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		path string
	}{
		{name: "schema number", src: `{"$schema":1}`, path: "$schema"},
		{
			name: "version string",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":"1"}`,
			path: "schemaVersion",
		},
		{
			name: "style number",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":1}`,
			path: "style",
		},
		{
			name: "radius number",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova","radius":1}`,
			path: "radius",
		},
		{
			name: "theme array",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":1,"style":"nova","radius":"0","theme":[]}`,
			path: "theme",
		},
		{
			name: "token number",
			src:  `{"theme":{"light":{"background":1}}}`,
			path: "theme.light.background",
		},
		{
			name: "future version",
			src:  `{"$schema":"https://ui.gsxhq.dev/schemas/preset-v1.json","schemaVersion":2}`,
			path: "schemaVersion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseJSON([]byte(tt.src))
			if err == nil {
				t.Fatal("ParseJSON returned nil error")
			}
			if !strings.Contains(err.Error(), tt.path) {
				t.Fatalf("error = %q, want path %q", err, tt.path)
			}
		})
	}
}

func TestParseJSONRejectsTrailingJSONValue(t *testing.T) {
	t.Parallel()

	golden := readTestFile(t, "testdata/default-nova.json")
	_, err := ParseJSON(append(golden, []byte(` {}`)...))
	if err == nil || !strings.Contains(err.Error(), "trailing") {
		t.Fatalf("error = %v, want trailing JSON value", err)
	}
}

func TestCanonicalJSONMatchesGoldenAndRoundTripsByteIdentically(t *testing.T) {
	t.Parallel()

	golden := readTestFile(t, "testdata/default-nova.json")
	got, err := CanonicalJSON(Default(StyleNova))
	if err != nil {
		t.Fatalf("CanonicalJSON: %v", err)
	}
	if !bytes.Equal(got, golden) {
		t.Fatalf("CanonicalJSON mismatch\n--- got ---\n%s\n--- want ---\n%s", got, golden)
	}
	if len(got) == 0 || got[len(got)-1] != '\n' ||
		(len(got) > 1 && got[len(got)-2] == '\n') {
		t.Fatal("CanonicalJSON must end in exactly one newline")
	}

	parsed, err := ParseJSON(got)
	if err != nil {
		t.Fatalf("ParseJSON(canonical): %v", err)
	}
	again, err := CanonicalJSON(parsed)
	if err != nil {
		t.Fatalf("CanonicalJSON(parsed): %v", err)
	}
	if !bytes.Equal(again, golden) {
		t.Fatalf("round trip mismatch\n--- got ---\n%s\n--- want ---\n%s", again, golden)
	}
}

func TestParseJSONReturnsIndependentThemeMaps(t *testing.T) {
	t.Parallel()

	golden := readTestFile(t, "testdata/default-nova.json")
	first, err := ParseJSON(golden)
	if err != nil {
		t.Fatalf("ParseJSON(first): %v", err)
	}
	first.Theme.Light["primary"] = "red"
	first.Theme.Dark["primary"] = "blue"

	again, err := ParseJSON(golden)
	if err != nil {
		t.Fatalf("ParseJSON(again): %v", err)
	}
	if again.Theme.Light["primary"] != "oklch(0% 0 0)" ||
		again.Theme.Dark["primary"] != "oklch(0.922 0 0)" {
		t.Fatal("ParseJSON returned shared theme map storage")
	}
}

func readTestFile(t *testing.T, name string) []byte {
	t.Helper()
	src, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return src
}
