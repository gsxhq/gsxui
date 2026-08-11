package port

import (
	"reflect"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		token     string
		kind      Kind
		dimension string
		value     string
		slot      string
		child     bool
		rest      string
	}{
		{"p-4", KindPlain, "", "", "", false, "p-4"},
		{"data-[size=sm]:max-w-xs", KindDimension, "size", "sm", "", false, "max-w-xs"},
		{"data-[size=default]:sm:max-w-md", KindDimension, "size", "default", "", false, "sm:max-w-md"},
		{"**:data-[slot=accordion-trigger-icon]:size-4", KindSlot, "", "", "accordion-trigger-icon", false, "size-4"},
		{"*:data-[slot=alert-description]:text-destructive/90", KindSlot, "", "", "alert-description", true, "text-destructive/90"},
		{"data-open:bg-muted/50", KindPlain, "", "", "", false, "data-open:bg-muted/50"},
		{"data-closed:animate-out", KindPlain, "", "", "", false, "data-closed:animate-out"},
		{"group-data-[size=default]/alert-dialog-content:place-items-start", KindPlain, "", "", "", false, "group-data-[size=default]/alert-dialog-content:place-items-start"},
		{"group-has-data-[slot=item-description]/item:translate-y-0.5", KindPlain, "", "", "", false, "group-has-data-[slot=item-description]/item:translate-y-0.5"},
		{"in-data-[slot=dropdown-menu-content]:p-0", KindPlain, "", "", "", false, "in-data-[slot=dropdown-menu-content]:p-0"},
		{"has-data-[icon=inline-end]:pr-1.5", KindPlain, "", "", "", false, "has-data-[icon=inline-end]:pr-1.5"},
		{"dark:aria-invalid:ring-destructive/40", KindPlain, "", "", "", false, "dark:aria-invalid:ring-destructive/40"},
		{"*:[svg:not([class*='size-'])]:size-8", KindPlain, "", "", "", false, "*:[svg:not([class*='size-'])]:size-8"},
		{"not-last:border-b", KindPlain, "", "", "", false, "not-last:border-b"},
		{"supports-backdrop-filter:backdrop-blur-xs", KindPlain, "", "", "", false, "supports-backdrop-filter:backdrop-blur-xs"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			got := Classify(tt.token)

			want := Classified{
				Kind:      tt.kind,
				Dimension: tt.dimension,
				Value:     tt.value,
				Slot:      tt.slot,
				Child:     tt.child,
				Rest:      tt.rest,
				Raw:       tt.token,
			}
			if got != want {
				t.Errorf("Classify(%q) = %+v, want %+v", tt.token, got, want)
			}
		})
	}
}

func TestSplitVariants(t *testing.T) {
	tests := []struct {
		token string
		want  []string
	}{
		{"p-4", []string{"p-4"}},
		{"data-[size=sm]:max-w-xs", []string{"data-[size=sm]", "max-w-xs"}},
		{"data-[size=default]:sm:max-w-md", []string{"data-[size=default]", "sm", "max-w-md"}},
		{"**:data-[slot=accordion-trigger-icon]:size-4", []string{"**", "data-[slot=accordion-trigger-icon]", "size-4"}},
		{"*:data-[slot=alert-description]:text-destructive/90", []string{"*", "data-[slot=alert-description]", "text-destructive/90"}},
		// The critical hazard: a `:` inside the arbitrary selector's own
		// nested pseudo-class call (`not(...)`) and a nested `[...]` inside
		// that call must not be treated as split points. The colon right
		// after the bare `*` combinator IS still a split point (bracket
		// depth is 0 there), giving three segments — Classify is what
		// later decides this doesn't match the `*`+`data-[slot=x]`
		// compound pattern and falls back to KindPlain.
		{"*:[svg:not([class*='size-'])]:size-8", []string{"*", "[svg:not([class*='size-'])]", "size-8"}},
		{"group-data-[size=default]/alert-dialog-content:place-items-start", []string{"group-data-[size=default]/alert-dialog-content", "place-items-start"}},
		{"dark:aria-invalid:ring-destructive/40", []string{"dark", "aria-invalid", "ring-destructive/40"}},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			got := splitVariants(tt.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitVariants(%q) = %#v, want %#v", tt.token, got, tt.want)
			}
		})
	}
}
