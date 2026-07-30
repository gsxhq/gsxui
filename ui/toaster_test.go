package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestToasterContract(t *testing.T) {
	got := render(t, ui.Toaster(nil))
	requireMarkup(t, got,
		`<section aria-label="Notifications" tabindex="-1">`,
		`<ol id="gsxui-toaster" data-gsxui-toaster data-gsxui-slot-toaster></ol>`,
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
		`data-gsxui-toaster`,
		`class="caller-region"`,
		`data-gsxui-slot-caller-token data-gsxui-slot-toaster`,
	)
	forbidMarkup(t, got, `id="gsxui-toaster"`)
}
