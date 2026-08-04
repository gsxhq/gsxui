package cli

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
)

// TestAnimateAssetMatchesInstalledPackage pins assets/css/animate.css to the
// npm package it was copied from, so a tw-animate-css upgrade fails loudly
// until `make generate-animate` refreshes the embedded copy.
func TestAnimateAssetMatchesInstalledPackage(t *testing.T) {
	embedded, err := fs.ReadFile(gsxui.Files, "assets/css/animate.css")
	if err != nil {
		t.Fatal(err)
	}
	upstream, err := os.ReadFile(filepath.Join("..", "..", "node_modules", "tw-animate-css", "dist", "tw-animate.css"))
	if err != nil {
		t.Fatal(err)
	}
	header, body, found := bytes.Cut(embedded, []byte("*/\n"))
	if !found || !bytes.Contains(header, []byte("tw-animate-css")) {
		t.Fatalf("assets/css/animate.css must start with a tw-animate-css attribution comment, got %q", embedded[:min(len(embedded), 80)])
	}
	if !bytes.Equal(bytes.TrimSpace(body), bytes.TrimSpace(upstream)) {
		t.Fatal("assets/css/animate.css drifted from node_modules/tw-animate-css/dist/tw-animate.css — run `make generate-animate`")
	}
}
