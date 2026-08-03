package showcase_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/examples/showcase"
)

// render executes a showcase card and returns its HTML.
func render(t *testing.T, node interface {
	Render(ctx context.Context, w io.Writer) error
}) string {
	t.Helper()
	var buf bytes.Buffer
	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func TestSignInCard(t *testing.T) {
	html := render(t, showcase.SignInCard())
	for _, want := range []string{
		"home-showcase-email",
		"home-showcase-password",
		"home-showcase-remember",
		"Sign in",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("SignInCard output missing %q", want)
		}
	}
}

func TestSettingsCard(t *testing.T) {
	html := render(t, showcase.SettingsCard())
	for _, want := range []string{
		"home-showcase-notifications",
		"home-showcase-autosave",
		"home-showcase-theme",
		"home-showcase-density",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("SettingsCard output missing %q", want)
		}
	}
}

func TestStatsCard(t *testing.T) {
	html := render(t, showcase.StatsCard())
	for _, want := range []string{
		"Storage",
		"Bandwidth",
		"Ada Lovelace",
		"Active",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("StatsCard output missing %q", want)
		}
	}
}
