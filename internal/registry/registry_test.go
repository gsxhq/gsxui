package registry_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/registry"
)

func TestComponents(t *testing.T) {
	got, err := registry.Components()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alert", "avatar", "badge", "button", "card", "checkbox", "dialog", "icon", "input", "label", "radio", "selectbox", "separator", "skeleton", "switchctl", "table", "textarea"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestDepsDerivedFromImports(t *testing.T) {
	deps, err := registry.Deps("dialog")
	if err != nil {
		t.Fatal(err)
	}
	// dialog.gsx imports ui/button (DialogFooter's Close button).
	if !reflect.DeepEqual(deps, []string{"button"}) {
		t.Fatalf("dialog deps = %v, want [button]", deps)
	}
	deps, err = registry.Deps("badge")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("badge deps = %v, want none", deps)
	}
	deps, err = registry.Deps("selectbox")
	if err != nil {
		t.Fatal(err)
	}
	// select.gsx imports ui/icon (the chevron).
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("selectbox deps = %v, want [icon]", deps)
	}
}

func TestHasJS(t *testing.T) {
	if !registry.HasJS("dialog") {
		t.Error("dialog should have JS")
	}
	if registry.HasJS("badge") {
		t.Error("badge should not have JS")
	}
}

func TestResolveTransitive(t *testing.T) {
	got, err := registry.Resolve([]string{"dialog"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"button", "dialog"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestResolveUnknown(t *testing.T) {
	if _, err := registry.Resolve([]string{"nope"}); err == nil {
		t.Fatal("want error for unknown component")
	}
}

func TestResolveRejectsNonComponentDir(t *testing.T) {
	// "core" is a real directory under ui/ but isn't an installable
	// component — it must be rejected the same as any unknown name.
	_, err := registry.Resolve([]string{"core"})
	if err == nil {
		t.Fatal("want error for ui/core, which is not a component")
	}
	if !strings.Contains(err.Error(), `unknown component "core"`) {
		t.Fatalf("got %v, want unknown-component error", err)
	}
}

func TestDepsRejectsNonComponentDir(t *testing.T) {
	if _, err := registry.Deps("core"); err == nil {
		t.Fatal("want error for ui/core, which is not a component")
	}
}
