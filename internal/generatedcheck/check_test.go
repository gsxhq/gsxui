package generatedcheck

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckGeneratedCleanPass(t *testing.T) {
	root := generatedFixture(t, "generated")

	if err := Check(root, func() error { return nil }); err != nil {
		t.Fatalf("Check(clean): %v", err)
	}
}

func TestCheckGeneratedIgnoresGitMetadata(t *testing.T) {
	root := generatedFixture(t, "generated")
	gitDir := filepath.Join(root, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "metadata.gsx"), []byte("not repository source"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "metadata.x.go"), []byte("not generated output"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "orphan.x.go"), []byte("repository metadata"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Check(root, func() error { return nil }); err != nil {
		t.Fatalf("Check(.git metadata): %v", err)
	}
}

func TestCheckGeneratedRejectsOutputRewrittenByGenerator(t *testing.T) {
	root := generatedFixture(t, "stale")
	output := filepath.Join(root, "component.x.go")

	err := Check(root, func() error {
		return os.WriteFile(output, []byte("generated"), 0o644)
	})
	assertCheckError(t, err, "generated output inventory changed")
}

func TestCheckGeneratedRejectsGeneratorFailureWithoutInventoryChange(t *testing.T) {
	root := generatedFixture(t, "generated")

	err := Check(root, func() error { return errors.New("exit status 23") })
	assertCheckError(t, err, "generator failed: exit status 23")
}

func TestCheckGeneratedRejectsOrphanOutput(t *testing.T) {
	root := generatedFixture(t, "generated")
	if err := os.WriteFile(filepath.Join(root, "orphan.x.go"), []byte("orphan"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Check(root, func() error { return nil })
	assertCheckError(t, err, `unexpected generated output "orphan.x.go"`)
}

func TestCheckGeneratedRejectsMissingOutputEvenWhenGeneratorRecreatesIt(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "component.gsx")
	output := filepath.Join(root, "component.x.go")
	if err := os.WriteFile(source, []byte("package fixture"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Check(root, func() error {
		return os.WriteFile(output, []byte("generated"), 0o644)
	})
	assertCheckError(t, err, `missing generated output "component.x.go"`)
	assertCheckError(t, err, "generated output inventory changed")
}

func generatedFixture(t *testing.T, output string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "component.gsx"), []byte("package fixture"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "component.x.go"), []byte(output), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func assertCheckError(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("Check() = nil, want error containing %q", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("Check() error = %q, want substring %q", err, want)
	}
}
