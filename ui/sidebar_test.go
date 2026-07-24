package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSidebarProviderStampsServerState(t *testing.T) {
	// The whole point of the design: state arrives as a parameter and is
	// rendered, so there is no flash and no hydration step.
	open := render(t, ui.SidebarProvider(true, gsx.Raw("x"), nil))
	if !strings.Contains(open, `data-state="expanded"`) {
		t.Errorf("want expanded\nin: %s", open)
	}
	closed := render(t, ui.SidebarProvider(false, gsx.Raw("x"), nil))
	if !strings.Contains(closed, `data-state="collapsed"`) {
		t.Errorf("want collapsed\nin: %s", closed)
	}
}

func TestSidebarProviderCarriesWidthVars(t *testing.T) {
	got := render(t, ui.SidebarProvider(true, gsx.Raw("x"), nil))
	for _, want := range []string{"--sidebar-width", "--sidebar-width-icon", `data-slot="sidebar-wrapper"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarRendersBothTrees(t *testing.T) {
	// Mobile is CSS-gated, not JS-swapped: the Sheet tree and the desktop
	// tree both exist in the DOM, gated md:hidden / hidden md:block.
	got := render(t, ui.Sidebar("", "", "", gsx.Raw("CONTENT"), nil))
	if strings.Count(got, "CONTENT") != 2 {
		t.Errorf("want children rendered in both trees, got %d\nin: %s", strings.Count(got, "CONTENT"), got)
	}
	for _, want := range []string{`data-mobile="true"`, "md:hidden", "md:block"} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarCollapsibleNoneIsFlat(t *testing.T) {
	// The reference short-circuits to one plain div; no gap, no container,
	// no Sheet, children rendered exactly once.
	got := render(t, ui.Sidebar("", "", "none", gsx.Raw("CONTENT"), nil))
	if strings.Count(got, "CONTENT") != 1 {
		t.Errorf("collapsible=none renders children once, got %d\nin: %s", strings.Count(got, "CONTENT"), got)
	}
	if strings.Contains(got, `data-mobile="true"`) {
		t.Errorf("collapsible=none must not render the Sheet tree\nin: %s", got)
	}
}

func TestSidebarStampsVariantSideCollapsible(t *testing.T) {
	got := render(t, ui.Sidebar("right", "floating", "icon", gsx.Raw("x"), nil))
	for _, want := range []string{`data-side="right"`, `data-variant="floating"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarMenuButtonActiveAndTooltip(t *testing.T) {
	got := render(t, ui.SidebarMenuButton(true, "", "", "Inbox", gsx.Raw("Inbox"), nil))
	for _, want := range []string{
		`data-slot="sidebar-menu-button"`,
		`data-active="true"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarMenuSkeletonShowIcon(t *testing.T) {
	with := render(t, ui.SidebarMenuSkeleton(true, nil))
	without := render(t, ui.SidebarMenuSkeleton(false, nil))
	if strings.Count(with, `data-slot="skeleton"`) <= strings.Count(without, `data-slot="skeleton"`) {
		t.Errorf("showIcon must add a skeleton\n with: %s\nwithout: %s", with, without)
	}
}
