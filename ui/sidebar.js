// Sidebar behavior. SidebarProvider's data-state is the server-rendered
// source of truth (no cookie, no localStorage — see ui/sidebar.gsx's own
// package doc comment for why); this module only flips it and its direct
// DOM consequences, and tells the page via gsxui:change so a caller's own
// Go handler/Alpine/htmx owns actually persisting it — the
// site/examples/sidebar/persisted.gsx recipe demonstrates the cookie
// round-trip shadcn's own React version bakes into the component instead.
//
// Trigger/rail both resolve "the provider" via
// closest('[data-slot="sidebar-wrapper"]') (SidebarProvider's own root) —
// no imperative handle, the DOM is the API surface, matching dialog.js's
// own proximity-wiring convention.
//
// Mobile vs desktop is resolved PER CLICK, not cached, via getComputedStyle
// on the desktop tree's own root: its `hidden md:block` class is what CSS
// mobile-gates in the first place (see ui/sidebar.gsx's package doc
// comment MOBILE section), so reading its computed display instead of
// re-implementing the md breakpoint as a JS constant means a caller's own
// custom md: override (or an entirely different container-query scheme)
// is still honored.
import { on, emit } from "./gsxui.js";

const desktopRootOf = (wrapper) => wrapper.querySelector('[data-slot="sidebar"]:not([data-mobile])');
const mobileDialogOf = (wrapper) => wrapper.querySelector('dialog[data-mobile="true"]');

function isMobile(desktopRoot) {
  return !desktopRoot || getComputedStyle(desktopRoot).display === "none";
}

// Desktop toggle: flips data-state on the wrapper AND on the desktop root —
// the group/peer element every group-data-*/peer-data-* selector in this
// component keys off (see the source map's own data-* consumption table).
// data-collapsible is restored from data-gsxui-sidebar-collapsible, the
// constant configured mode Sidebar stamps once at render time and never
// itself toggles — this mirrors the reference's own
// `state === "collapsed" ? collapsible : ""` ternary, computed here instead
// of from React state.
function toggleDesktop(wrapper, desktopRoot) {
  const open = wrapper.dataset.state !== "expanded";
  wrapper.dataset.state = open ? "expanded" : "collapsed";
  desktopRoot.dataset.state = wrapper.dataset.state;
  desktopRoot.dataset.collapsible = open ? "" : desktopRoot.dataset.gsxuiSidebarCollapsible || "offcanvas";
  emit(wrapper, "gsxui:change", { open });
}

// Mobile toggle: opens the Sheet tree exactly like a SheetTrigger click
// would — dialog.js's own delegated "toggle" listener picks up the
// resulting native `toggle` event to sync aria/state/gsxui:open (see
// ui/dialog.js). No gsxui:change here: mobile openMobile and desktop open
// are independent booleans in the reference too (see the source map's own
// Sidebar §3 note), and dialog.js already emits its own gsxui:open/
// gsxui:close on the sheet's own <dialog> for anyone tracking mobile state.
function openMobile(wrapper) {
  const dialog = mobileDialogOf(wrapper);
  if (!dialog || dialog.open) return;
  dialog.dataset.state = "open";
  dialog.showModal();
}

function toggle(wrapper) {
  const desktopRoot = desktopRootOf(wrapper);
  if (isMobile(desktopRoot)) {
    openMobile(wrapper);
  } else {
    toggleDesktop(wrapper, desktopRoot);
  }
}

on("click", '[data-slot="sidebar-trigger"]', (_e, trigger) => {
  const wrapper = trigger.closest('[data-slot="sidebar-wrapper"]');
  if (wrapper) toggle(wrapper);
});

on("click", '[data-slot="sidebar-rail"]', (_e, rail) => {
  const wrapper = rail.closest('[data-slot="sidebar-wrapper"]');
  if (wrapper) toggle(wrapper);
});

// Cmd/Ctrl+B toggles the FIRST provider on the page (SIDEBAR_KEYBOARD_
// SHORTCUT — see ui/sidebar.gsx's own package doc comment) — module scope,
// same shape as command.js's own ⌘K handler, plus the repeat-fire guard:
// holding the chord down must not spam toggle on every OS key-repeat.
addEventListener("keydown", (e) => {
  if (e.key.toLowerCase() !== "b" || !(e.metaKey || e.ctrlKey) || e.repeat) return;
  const wrapper = document.querySelector('[data-slot="sidebar-wrapper"]');
  if (!wrapper) return;
  e.preventDefault();
  toggle(wrapper);
});
