package ui

import (
	"reflect"
	"testing"

	gsx "github.com/gsxhq/gsx"
)

func TestWithSlotPrependsAndDeduplicates(t *testing.T) {
	attrs := gsx.Attrs{
		{Key: "data-gsxui-slot", Value: "alert-dialog-action button"},
		{Key: "class", Value: "h-12"},
		{Key: "id", Value: "save"},
	}

	got := withSlot("button", attrs)
	if v, _ := got.Get("data-gsxui-slot"); v != "button alert-dialog-action" {
		t.Fatalf("slot tokens = %q", v)
	}
	if got.Class() != "h-12" {
		t.Fatalf("class = %q", got.Class())
	}
	if v, _ := got.Get("id"); v != "save" {
		t.Fatalf("id = %q", v)
	}
}

func TestWithSlotAddsTokenToNilAttrs(t *testing.T) {
	got := withSlot("button", nil)
	if v, _ := got.Get("data-gsxui-slot"); v != "button" {
		t.Fatalf("slot tokens = %q", v)
	}
}

func TestWithSlotPreservesNestedInnerToOuterOrder(t *testing.T) {
	got := withSlot("button", withSlot("pagination-link", withSlot("pagination-previous", nil)))
	if v, _ := got.Get("data-gsxui-slot"); v != "button pagination-link pagination-previous" {
		t.Fatalf("slot tokens = %q", v)
	}
}

func TestWithSlotUsesAttrsRenderingSemantics(t *testing.T) {
	attrs := gsx.Attrs{{Key: "data-gsxui-slot", Value: []string{"alert-dialog-action", "button"}}}

	got := withSlot("button", attrs)
	if v, _ := got.Get("data-gsxui-slot"); v != "button alert-dialog-action" {
		t.Fatalf("slot tokens = %q", v)
	}
}

func TestWithSlotDoesNotMutateInput(t *testing.T) {
	attrs := gsx.Attrs{
		{Key: "data-gsxui-slot", Value: "alert-dialog-action button"},
		{Key: "id", Value: "save"},
	}
	want := append(gsx.Attrs(nil), attrs...)

	_ = withSlot("button", attrs)
	if !reflect.DeepEqual(attrs, want) {
		t.Fatalf("input attrs mutated: got %#v, want %#v", attrs, want)
	}
}

func TestWithSlotPanicsForEmptyName(t *testing.T) {
	defer func() {
		if got := recover(); got != "gsxui: empty slot name" {
			t.Fatalf("panic = %#v", got)
		}
	}()

	withSlot("", nil)
}
