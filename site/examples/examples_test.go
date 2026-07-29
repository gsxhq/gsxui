package examples_test

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/examples"
)

func TestRegisterRejectsDuplicateExampleNames(t *testing.T) {
	const childEnvironment = "GSXUI_TEST_DUPLICATE_EXAMPLE_REGISTRATION"
	if os.Getenv(childEnvironment) == "1" {
		examples.Register("__duplicate-registration-test__", examples.Example{Name: "same-name"})
		examples.Register("__duplicate-registration-test__", examples.Example{Name: "same-name"})
		return
	}

	command := exec.Command(os.Args[0], "-test.run=^TestRegisterRejectsDuplicateExampleNames$")
	command.Env = append(os.Environ(), childEnvironment+"=1")
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("duplicate component/example registration succeeded; output:\n%s", output)
	}
	const want = `examples: duplicate registration for component "__duplicate-registration-test__" example "same-name"`
	if !strings.Contains(string(output), want) {
		t.Fatalf("duplicate registration failure = %v\noutput:\n%s\nwant panic containing:\n%s", err, output, want)
	}
}

func TestFindResolvesOnlyExactRegisteredExamples(t *testing.T) {
	title, node, ok := examples.Find("button", "basic", nil)
	if !ok || title != "Basic" || node == nil {
		t.Fatalf("Find(button, basic) = %q, %#v, %t", title, node, ok)
	}

	for _, key := range [][2]string{
		{"button", "missing"},
		{"missing", "basic"},
		{"Button", "basic"},
	} {
		if _, _, ok := examples.Find(key[0], key[1], nil); ok {
			t.Errorf("Find(%q, %q) accepted an unregistered key", key[0], key[1])
		}
	}
}

func TestFindResolvesOnlyRegisteredPreviewCases(t *testing.T) {
	title, node, ok := examples.Find(
		"sidebar",
		"variants",
		url.Values{examples.PreviewQueryKey: []string{"floating"}},
	)
	if !ok || title != "Variants · variant=floating" || node == nil {
		t.Fatalf("Find(sidebar, variants, floating) = %q, %#v, %t", title, node, ok)
	}

	for _, values := range []url.Values{
		{examples.PreviewQueryKey: []string{"missing"}},
		{examples.PreviewQueryKey: []string{""}},
		{examples.PreviewQueryKey: []string{"floating", "inset"}},
	} {
		if _, _, ok := examples.Find("sidebar", "variants", values); ok {
			t.Fatalf("Find(sidebar, variants, %v) accepted an inexact preview selector", values)
		}
	}
}

func TestFindUsesTheRegisteredQueryRenderer(t *testing.T) {
	_, node, ok := examples.Find(
		"calendar",
		"basic",
		url.Values{"month": []string{"2026-02"}},
	)
	if !ok {
		t.Fatal("Find(calendar, basic) = not found")
	}

	var output bytes.Buffer
	if err := node.Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "February 2026") {
		t.Fatalf("query-rendered calendar does not contain February 2026:\n%s", output.String())
	}
}

// TestExampleSourceMatchesFile is the drift guard for
// source-shown-is-source-run: for every registered example, the bytes
// examples.Source returns (read from the embedded FS) must be byte-
// identical to the on-disk .gsx file. If they ever diverge — e.g. someone
// edits the embedded copy without touching examples.go's //go:embed
// pattern, or vice versa — the component page would display source that
// isn't what's actually running.
func TestExampleSourceMatchesFile(t *testing.T) {
	for _, component := range examples.Components() {
		for _, ex := range examples.For(component) {
			t.Run(component+"/"+ex.Name, func(t *testing.T) {
				got, err := examples.Source(ex.SourcePath)
				if err != nil {
					t.Fatalf("examples.Source(%q): %v", ex.SourcePath, err)
				}

				want, err := os.ReadFile(filepath.FromSlash(ex.SourcePath))
				if err != nil {
					t.Fatalf("reading on-disk file %q: %v", ex.SourcePath, err)
				}

				if got != string(want) {
					t.Errorf("embedded Source(%q, %q) != on-disk %q\n--- embedded ---\n%s\n--- on-disk ---\n%s",
						component, ex.Name, ex.SourcePath, got, string(want))
				}
			})
		}
	}
}

// TestExamplesRender is a smoke test for rendering all registered examples.
// For every component × example pair, it calls ex.Node.Render(context.Background(), &buf),
// fails on error, and fails if the output contains gsx's blocked-URL sentinel "about:invalid#gsx".
func TestExamplesRender(t *testing.T) {
	for _, component := range examples.Components() {
		for _, ex := range examples.For(component) {
			t.Run(component+"/"+ex.Name, func(t *testing.T) {
				nodes := []gsx.Node{ex.Node}
				for _, preview := range ex.Previews {
					nodes = append(nodes, preview.Node)
				}
				for index, node := range nodes {
					var buf bytes.Buffer
					if err := node.Render(context.Background(), &buf); err != nil {
						t.Errorf("node %d Render: %v", index, err)
					}
					if strings.Contains(buf.String(), "about:invalid#gsx") {
						t.Errorf("node %d output contains blocked URL sentinel about:invalid#gsx", index)
					}
				}
			})
		}
	}
}
