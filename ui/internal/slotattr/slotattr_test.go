package slotattr

import (
	"reflect"
	"testing"

	gsx "github.com/gsxhq/gsx"
)

func TestWithPrependsDeduplicatesAndPreservesAttrs(t *testing.T) {
	attrs := gsx.Attrs{
		{Key: "data-gsxui-slot", Value: []string{"spinner", "icon"}},
		{Key: "class", Value: "size-6"},
		{Key: "id", Value: "loading"},
	}
	wantInput := append(gsx.Attrs(nil), attrs...)

	got := With("icon", attrs)

	if value, _ := got.Get("data-gsxui-slot"); value != "icon spinner" {
		t.Fatalf("slot tokens = %q", value)
	}
	if got.Class() != "size-6" {
		t.Fatalf("class = %q", got.Class())
	}
	if value, _ := got.Get("id"); value != "loading" {
		t.Fatalf("id = %q", value)
	}
	if !reflect.DeepEqual(attrs, wantInput) {
		t.Fatalf("input attrs mutated: got %#v, want %#v", attrs, wantInput)
	}
}

func TestWithIgnoresUnrenderableExistingSlot(t *testing.T) {
	got := With("icon", gsx.Attrs{
		{Key: "data-gsxui-slot", Value: struct{}{}},
		{Key: "id", Value: "loading"},
	})
	if value, _ := got.Get("data-gsxui-slot"); value != "icon" {
		t.Fatalf("slot tokens = %q", value)
	}
	if value, _ := got.Get("id"); value != "loading" {
		t.Fatalf("id = %q", value)
	}
}

func TestWithPanicsForEmptyName(t *testing.T) {
	defer func() {
		if got := recover(); got != "gsxui: empty slot name" {
			t.Fatalf("panic = %#v", got)
		}
	}()
	With("", nil)
}
