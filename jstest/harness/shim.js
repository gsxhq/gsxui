// Recording replacement for ui/gsxui.js, served at /shim/gsxui.js.
//
// It has gsxui.js's exact export surface but registers nothing: on() records
// what a module WOULD have bound, so the Playwright suite can check that no
// two modules can claim the same element for the same (type, capture) pair.
// That is the mechanical form of the hook-prefix collision that shipped in
// Tier 4 Batch B, where dropdown.js and context-menu.js both matched
// data-gsxui-menu-* and one click ran both handlers.
//
// emit() is the real implementation, unchanged — nothing depends on it here,
// but a module that calls emit at import time must not throw.
//
// SCOPE — read "the registry" as "the DELEGATED registry", not "everything
// the library binds." Only what a module routes through gsxui.js's on() is
// recorded here, so only that is visible to the checks built on this data
// (selector disjointness, selector coverage, selector parseability).
//
// A module that calls addEventListener directly is invisible to all of
// them. That is correct code where it appears, not an oversight, but it is
// the hole in this file's picture. Two shapes exist in ui/ today:
//
//   1. Document/window-level signals with no selector to delegate on —
//      ui/input-otp.js's document "selectionchange" (the case the review
//      called out), ui/command.js's ⌘K "keydown", ui/sidebar.js's shortcut
//      "keydown", ui/avatar.js's window "load" sweep, ui/sonner.js's
//      "DOMContentLoaded" init.
//   2. Listeners bound to a specific element the module already holds a
//      reference to — ui/carousel.js's per-root pointer/focus pause,
//      ui/sonner.js's per-toast buttons, ui/context-menu.js's one-shot
//      "pointerup".
//
// Neither shape can collide the way delegated selectors can, which is why
// their absence does not weaken the disjointness invariant. It does mean
// selector coverage cannot speak for them: if a .gsx renamed an attribute
// only a direct listener reads, nothing here would notice.

const registrations = [];
window.__gsxuiRegistrations = registrations;

export function on(type, selector, handler, { capture = false } = {}) {
  registrations.push({ type, capture, selector, module: callerModule() });
}

export function emit(el, name, detail) {
  return el.dispatchEvent(
    new CustomEvent(name, { bubbles: true, cancelable: true, detail }),
  );
}

// callerModule walks the V8 stack to the first frame outside this file. Frame
// 0 is "Error", frame 1 is on() itself, frame 2 is the module that called it.
// Returns the bare filename ("dropdown.js") so failure output is readable.
function callerModule() {
  const frames = (new Error().stack || "").split("\n");
  for (const frame of frames) {
    const match = frame.match(/\/shim\/([A-Za-z0-9._-]+\.js)/);
    if (match && match[1] !== "gsxui.js") return match[1];
  }
  return "unknown";
}
