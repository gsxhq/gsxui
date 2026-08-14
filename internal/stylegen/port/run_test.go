package port

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

const upstreamRootPath = "/Users/jackieli/personal/shadcn-ui"

// repoRootPath resolves this package's containing repository root
// (internal/stylegen/port -> internal/stylegen -> internal -> root), the
// same three-levels-up relationship emit_test.go's golden test already
// relies on for locating registry/styles/nova.
func repoRootPath(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return abs
}

func fileChecksum(t *testing.T, path string) [32]byte {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return sha256.Sum256(src)
}

// TestRunDryRunMaia is the plan's required end-to-end smoke test: porting
// maia for real (dryRun=false, done once by hand as part of Task 6's real
// port, see the task report) must produce Written == 54 with Unmapped == 0
// once the mapping table is complete — this test pins that same claim in
// dry-run form, which additionally must touch nothing on disk. "chart"
// (registry/canonical/shapes/chart.go) joined the canonical set after that
// survey, so Written is now 55, still with Unmapped == 0 — a real "MARK:
// Chart" section maps it (mapping.go's ignoredSections/styleInvariantComponents
// comments have the history of Chart's own section moving from ignored to
// mapped as its shape grew a real tooltip slot).
func TestRunDryRunMaia(t *testing.T) {
	if _, err := os.Stat(upstreamRootPath); err != nil {
		t.Skipf("upstream shadcn-ui checkout not present at %s", upstreamRootPath)
	}
	root := repoRootPath(t)

	buttonPath := filepath.Join(root, "registry", "styles", "maia", "button.css")
	before := fileChecksum(t, buttonPath)

	report, err := Run(root, upstreamRootPath, []string{"maia"}, true)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if report.Written != 55 {
		t.Errorf("Written = %d, want 55", report.Written)
	}
	if report.Unmapped != 0 {
		t.Errorf("Unmapped = %d, want 0: %v", report.Unmapped, report.UnmappedDetail)
	}

	after := fileChecksum(t, buttonPath)
	if before != after {
		t.Error("dry run modified registry/styles/maia/button.css on disk")
	}
}

// TestRunRejectsUnknownStyle exercises Run's own validation independent of
// the upstream checkout (no t.Skip needed, since Run validates every
// requested style name before it ever touches upstreamRoot): an
// unrecognized style name must error rather than silently reading a
// nonexistent upstream file.
func TestRunRejectsUnknownStyle(t *testing.T) {
	root := repoRootPath(t)
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	_, err := Run(root, upstreamRootPath, []string{"not-a-real-style"}, true)
	if err == nil {
		t.Fatal("Run accepted an unknown style, want error")
	}
}
