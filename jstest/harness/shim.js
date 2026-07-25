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
