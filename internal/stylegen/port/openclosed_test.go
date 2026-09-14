package port

import "testing"

// TestRewriteOpenClosedMarkerPerSlot pins the translations of upstream's bare
// data-open/data-closed variants: gsxui's JS-stamped data-[state=…] form for
// components with a behavior module; the native open:/not-open: variant on
// the one slot of a <details>-backed component that IS the <details>; and a
// refusal (reported unmapped by Transform) on that component's other slots,
// where neither form can match correctly.
func TestRewriteOpenClosedMarkerPerSlot(t *testing.T) {
	cases := []struct {
		component, slot, token string
		want                   string
		ok                     bool
	}{
		{"dialog", "content", "data-open:animate-in", "data-[state=open]:animate-in", true},
		{"dialog", "content", "data-closed:animate-out", "data-[state=closed]:animate-out", true},
		{"accordion", "item", "data-open:bg-muted/50", "open:bg-muted/50", true},
		{"accordion", "item", "data-closed:opacity-50", "not-open:opacity-50", true},
		{"accordion", "item", "not-last:border-b", "not-last:border-b", true},
		{"accordion", "content", "data-open:animate-accordion-down", "data-open:animate-accordion-down", false},
		{"accordion", "content", "px-4", "px-4", true},
		{"collapsible", "", "data-open:hover:bg-muted", "open:hover:bg-muted", true},
		{"collapsible", "trigger", "data-open:underline", "data-open:underline", false},
	}
	for _, tc := range cases {
		got, ok := rewriteOpenClosedMarker(tc.component, tc.slot, tc.token)
		if got != tc.want || ok != tc.ok {
			t.Errorf("rewriteOpenClosedMarker(%q, %q, %q) = %q, %v; want %q, %v",
				tc.component, tc.slot, tc.token, got, ok, tc.want, tc.ok)
		}
	}
	for _, component := range []string{"accordion", "collapsible"} {
		if !NativeDetails(component) {
			t.Errorf("NativeDetails(%q) = false, want true", component)
		}
	}
	if NativeDetails("dialog") {
		t.Error("NativeDetails(\"dialog\") = true, want false: dialog.js stamps data-state")
	}
}
