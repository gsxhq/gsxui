package ui_test

import (
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// novaToastRecipe loads the default style's Toast recipe once. ui.Toast is
// generated output now, so the classes it renders are the default style's
// concrete utilities — same pattern as novaButtonRecipe in button_test.go.
var novaToastRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "toast.css")
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	style, err := recipe.ParseStyle(path, src)
	if err != nil {
		panic(err)
	}
	return style
})

func toastRecipeUtilities(class string) []string {
	rule, ok := novaToastRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalToastClass renders the class attribute ui.Toast's slot emits,
// with any caller classes merged in after the recipe's own utilities.
func canonicalToastClass(slot string, caller ...string) string {
	name := "gsxui-recipe-toast"
	if slot != "" {
		name += "-" + slot
	}
	classes := append([]string(nil), toastRecipeUtilities(name)...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func requireMarkup(t *testing.T, got string, fragments ...string) {
	t.Helper()
	for _, fragment := range fragments {
		if !strings.Contains(got, fragment) {
			t.Errorf("missing %q\nin: %s", fragment, got)
		}
	}
}

func forbidMarkup(t *testing.T, got string, fragments ...string) {
	t.Helper()
	for _, fragment := range fragments {
		if strings.Contains(got, fragment) {
			t.Errorf("unexpected %q\nin: %s", fragment, got)
		}
	}
}

func TestToastContract(t *testing.T) {
	got := render(t, ui.Toast("error", "Failed", "Try again", "Retry", "Dismiss", nil))
	requireMarkup(t, got,
		`<li data-type="error" role="status" aria-live="assertive" aria-atomic="true" `+canonicalToastClass("")+` data-gsxui-slot-toast>`,
		`data-gsxui-slot-toast-icon`,
		`data-gsxui-slot-toast-icon data-gsxui-slot-icon`,
		`<div `+canonicalToastClass("content")+` data-gsxui-slot-toast-content>`,
		`<div `+canonicalToastClass("title")+` data-gsxui-slot-toast-title>Failed</div>`,
		`<div `+canonicalToastClass("description")+` data-gsxui-slot-toast-description>Try again</div>`,
		`<button type="button" `+canonicalToastClass("action")+` data-gsxui-slot-toast-action>Retry</button>`,
		`<button type="button" `+canonicalToastClass("cancel")+` data-gsxui-slot-toast-cancel>Dismiss</button>`,
		`<button type="button" `+canonicalToastClass("close")+` aria-label="Close" data-gsxui-slot-toast-close>`,
		`data-gsxui-slot-toast-close-icon data-gsxui-slot-icon`,
	)
	forbidMarkup(t, got,
		`data-slot=`,
		`data-icon`,
		`data-content`,
		`data-title`,
		`data-description`,
		`data-action`,
		`data-cancel`,
		`data-close-button`,
	)
}

func TestToastCallerClassAndSlotComposition(t *testing.T) {
	got := render(t, ui.Toast("default", "Hello", "", "", "", gsx.Attrs{
		{Key: "class", Value: "caller-toast"},
		{Key: "data-gsxui-slot-caller-token", Value: true},
	}))
	requireMarkup(t, got,
		canonicalToastClass("", "caller-toast"),
		`data-gsxui-slot-caller-token data-gsxui-slot-toast`,
	)
	if strings.Count(got, `caller-toast`) != 1 {
		t.Errorf("caller class merged %d times, want once\nin: %s", strings.Count(got, `caller-toast`), got)
	}
}

func TestToastTypeIconsAndAria(t *testing.T) {
	cases := []struct {
		typ      string
		ariaLive string
		glyph    string
	}{
		{"default", "polite", ""},
		{"success", "polite", `<path d="m9 12 2 2 4-4"/>`},
		{"info", "polite", `<path d="M12 16v-4"/>`},
		{"warning", "polite", `d="m21.73 18-8-14`},
		{"error", "assertive", `<path d="m15 9-6 6"/>`},
		{"loading", "polite", `d="M21 12a9 9 0 1 1-6.219-8.56"`},
	}
	for _, c := range cases {
		t.Run(c.typ, func(t *testing.T) {
			got := render(t, ui.Toast(c.typ, "Title", "", "", "", nil))
			requireMarkup(t, got,
				`data-type="`+c.typ+`"`,
				`aria-live="`+c.ariaLive+`"`,
			)
			if c.typ == "default" {
				forbidMarkup(t, got, `data-gsxui-slot-toast-icon`)
			} else {
				requireMarkup(t, got, `data-gsxui-slot-toast-icon`, c.glyph)
			}
		})
	}
}

func TestToastOptionalPartsAndDuration(t *testing.T) {
	full := render(t, ui.Toast("info", "Title", "Detail", "Retry", "Dismiss",
		gsx.Attrs{{Key: "data-duration", Value: "8000"}}))
	requireMarkup(t, full,
		`data-gsxui-slot-toast-description`,
		`data-gsxui-slot-toast-action`,
		`data-gsxui-slot-toast-cancel`,
		`data-duration="8000"`,
	)

	minimal := render(t, ui.Toast("", "Hello", "", "", "", nil))
	requireMarkup(t, minimal, `data-type="default"`, `aria-live="polite"`)
	forbidMarkup(t, minimal,
		`data-gsxui-slot-toast-icon`,
		`data-gsxui-slot-toast-description`,
		`data-gsxui-slot-toast-action`,
		`data-gsxui-slot-toast-cancel`,
	)
}
