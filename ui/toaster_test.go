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

// novaToasterRecipe loads the default style's Toaster recipe once, the same
// pattern novaToastRecipe uses in toast_test.go.
var novaToasterRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "toaster.css")
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

// canonicalToasterClass renders the class attribute ui.Toaster's single slot
// emits, with any caller classes merged in after the recipe's utilities.
func canonicalToasterClass(caller ...string) string {
	rule, ok := novaToasterRecipe().Lookup("gsxui-recipe-toaster")
	if !ok {
		panic("default style declares no recipe gsxui-recipe-toaster")
	}
	classes := append(append([]string(nil), rule.Utilities...), caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestToasterContract(t *testing.T) {
	got := render(t, ui.Toaster(nil))
	requireMarkup(t, got,
		`<section aria-label="Notifications" tabindex="-1">`,
		`<ol id="gsxui-toaster" `+canonicalToasterClass()+` data-gsxui-slot-toaster></ol>`,
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
	// The region's own <ol> is pinned above with no class attribute; the
	// nested <template> cards are ui.Toast, which is on the slot axis and
	// legitimately renders its recipe classes.
	forbidMarkup(t, got, `data-slot=`)
}

func TestToasterAttrsMergeAndCallerClass(t *testing.T) {
	got := render(t, ui.Toaster(gsx.Attrs{
		{Key: "id", Value: "my-toaster"},
		{Key: "class", Value: "caller-region"},
		{Key: "data-gsxui-slot-caller-token", Value: true},
	}))
	requireMarkup(t, got,
		`id="my-toaster"`,
		`data-gsxui-slot-toaster`,
		canonicalToasterClass("caller-region"),
		`data-gsxui-slot-caller-token data-gsxui-slot-toaster`,
	)
	forbidMarkup(t, got, `id="gsxui-toaster"`)
}
