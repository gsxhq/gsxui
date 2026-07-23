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
	want := []string{"accordion", "alert", "aspect-ratio", "avatar", "badge", "breadcrumb", "button", "card", "checkbox", "dialog", "dropdown", "icon", "input", "kbd", "label", "progress", "radio", "select", "separator", "skeleton", "spinner", "switch", "table", "tabs", "textarea", "tooltip"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for _, unwanted := range []string{"core", "gsxui", "index", "selectbox", "switchctl"} {
		if slicesContains(got, unwanted) {
			t.Fatalf("Components() = %v, should not contain %q", got, unwanted)
		}
	}
}

func slicesContains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

func TestDeps(t *testing.T) {
	// dialog.x.go references Button (DialogFooter's Close button) — an
	// intra-package edge with no import to scan, resolved via declIndex.
	deps, err := registry.Deps("dialog")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button"}) {
		t.Fatalf("dialog deps = %v, want [button]", deps)
	}

	// accordion.gsx imports ui/icon (AccordionTrigger's chevron).
	deps, err = registry.Deps("accordion")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("accordion deps = %v, want [icon]", deps)
	}

	// select.gsx imports ui/icon (the chevron).
	deps, err = registry.Deps("select")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("select deps = %v, want [icon]", deps)
	}

	// spinner.gsx imports ui/icon (icon.LoaderCircle).
	deps, err = registry.Deps("spinner")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("spinner deps = %v, want [icon]", deps)
	}

	// breadcrumb.gsx imports ui/icon (BreadcrumbSeparator's default
	// ChevronRight, BreadcrumbEllipsis's Ellipsis/MoreHorizontal).
	deps, err = registry.Deps("breadcrumb")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("breadcrumb deps = %v, want [icon]", deps)
	}

	// kbd.gsx, aspect-ratio.gsx, and progress.gsx have no icon import and no
	// intra-package reference to another component.
	deps, err = registry.Deps("kbd")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("kbd deps = %v, want none", deps)
	}

	deps, err = registry.Deps("aspect-ratio")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("aspect-ratio deps = %v, want none", deps)
	}

	deps, err = registry.Deps("progress")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("progress deps = %v, want none", deps)
	}

	deps, err = registry.Deps("button")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("button deps = %v, want none", deps)
	}

	deps, err = registry.Deps("icon")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("icon deps = %v, want none", deps)
	}

	if _, err := registry.Deps("nosuch"); err == nil || !strings.Contains(err.Error(), "gsxui list") {
		t.Fatalf("Deps(nosuch) err = %v, want error mentioning 'gsxui list'", err)
	}
	if _, err := registry.Deps("core"); err == nil || !strings.Contains(err.Error(), "gsxui list") {
		t.Fatalf("Deps(core) err = %v, want error mentioning 'gsxui list'", err)
	}
}

func TestHasJS(t *testing.T) {
	if !registry.HasJS("dropdown") {
		t.Error("dropdown should have JS")
	}
	if registry.HasJS("button") {
		t.Error("button should not have JS")
	}
	if registry.HasJS("gsxui") {
		t.Error("gsxui should not have JS")
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

	got, err = registry.Resolve([]string{"select"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"icon", "select"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestResolveUnknown(t *testing.T) {
	if _, err := registry.Resolve([]string{"nope"}); err == nil {
		t.Fatal("want error for unknown component")
	}
}

func TestResolveRejectsNonComponentFile(t *testing.T) {
	// "core" is not a real path under the flat ui/ anymore, but must still
	// be rejected the same as any other unknown name.
	_, err := registry.Resolve([]string{"core"})
	if err == nil {
		t.Fatal("want error for core, which is not a component")
	}
	if !strings.Contains(err.Error(), `unknown component "core"`) {
		t.Fatalf("got %v, want unknown-component error", err)
	}
}

func TestDepsRejectsNonComponentFile(t *testing.T) {
	if _, err := registry.Deps("core"); err == nil {
		t.Fatal("want error for core, which is not a component")
	}
}
