package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

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

func TestSonnerToasterContract(t *testing.T) {
	got := render(t, ui.Toaster(nil))
	requireMarkup(t, got,
		`<section aria-label="Notifications" tabindex="-1">`,
		`<ol id="gsxui-toaster" data-gsxui-toaster data-gsxui-slot="toaster"></ol>`,
		`<template data-gsxui-toast-template="default">`,
		`<template data-gsxui-toast-template="success">`,
		`<template data-gsxui-toast-template="info">`,
		`<template data-gsxui-toast-template="warning">`,
		`<template data-gsxui-toast-template="error">`,
		`<template data-gsxui-toast-template="loading">`,
	)
	if gotCount := strings.Count(got, `data-gsxui-toast-template=`); gotCount != 6 {
		t.Errorf("template count = %d, want 6\nin: %s", gotCount, got)
	}
	forbidMarkup(t, got, `data-slot=`, ` class="`)
}

func TestSonnerToasterAttrsMergeAndCallerClass(t *testing.T) {
	got := render(t, ui.Toaster(gsx.Attrs{
		{Key: "id", Value: "my-toaster"},
		{Key: "class", Value: "caller-region"},
		{Key: "data-gsxui-slot", Value: "caller-token"},
	}))
	requireMarkup(t, got,
		`id="my-toaster"`,
		`data-gsxui-toaster`,
		`class="caller-region"`,
		`data-gsxui-slot="toaster caller-token"`,
	)
	forbidMarkup(t, got, `id="gsxui-toaster"`)
}

func TestSonnerToastContract(t *testing.T) {
	got := render(t, ui.Toast("error", "Failed", "Try again", "Retry", "Dismiss", nil))
	requireMarkup(t, got,
		`<li data-gsxui-toast data-type="error" role="status" aria-live="assertive" aria-atomic="true" data-gsxui-slot="toast">`,
		`data-gsxui-toast-icon`,
		`data-gsxui-slot="icon toast-icon"`,
		`<div data-gsxui-slot="toast-content">`,
		`<div data-gsxui-toast-title data-gsxui-slot="toast-title">Failed</div>`,
		`<div data-gsxui-toast-description data-gsxui-slot="toast-description">Try again</div>`,
		`<button type="button" data-gsxui-toast-action data-gsxui-slot="toast-action">Retry</button>`,
		`<button type="button" data-gsxui-toast-cancel data-gsxui-slot="toast-cancel">Dismiss</button>`,
		`<button type="button" data-gsxui-toast-close aria-label="Close" data-gsxui-slot="toast-close">`,
		`data-gsxui-slot="icon toast-close-icon"`,
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
		` class="`,
	)
}

func TestSonnerToastCallerClassAndSlotComposition(t *testing.T) {
	got := render(t, ui.Toast("default", "Hello", "", "", "", gsx.Attrs{
		{Key: "class", Value: "caller-toast"},
		{Key: "data-gsxui-slot", Value: "caller-token"},
	}))
	requireMarkup(t, got,
		`class="caller-toast"`,
		`data-gsxui-slot="toast caller-token"`,
	)
}

func TestSonnerToastTypeIconsAndAria(t *testing.T) {
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
				forbidMarkup(t, got, `data-gsxui-toast-icon`)
			} else {
				requireMarkup(t, got, `data-gsxui-toast-icon`, c.glyph)
			}
		})
	}
}

func TestSonnerToastOptionalPartsAndDuration(t *testing.T) {
	full := render(t, ui.Toast("info", "Title", "Detail", "Retry", "Dismiss",
		gsx.Attrs{{Key: "data-duration", Value: "8000"}}))
	requireMarkup(t, full,
		`data-gsxui-toast-description`,
		`data-gsxui-toast-action`,
		`data-gsxui-toast-cancel`,
		`data-duration="8000"`,
	)

	minimal := render(t, ui.Toast("", "Hello", "", "", "", nil))
	requireMarkup(t, minimal, `data-type="default"`, `aria-live="polite"`)
	forbidMarkup(t, minimal,
		`data-gsxui-toast-icon`,
		`data-gsxui-toast-description`,
		`data-gsxui-toast-action`,
		`data-gsxui-toast-cancel`,
	)
}
