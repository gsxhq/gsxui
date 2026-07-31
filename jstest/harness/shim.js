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
//      "keydown", ui/avatar.js's window "load" sweep, ui/toaster.js's
//      "DOMContentLoaded" init.
//   2. Listeners bound to a specific element the module already holds a
//      reference to — ui/carousel.js's per-root pointer/focus pause,
//      ui/toaster.js's per-toast buttons, ui/context-menu.js's one-shot
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

// position/release are the real implementations, copied verbatim from
// ui/gsxui.js — positioning is behavior under test, not delegation, so the
// recording layer must not distort it. Keep in lockstep with ui/gsxui.js
// ("exact export surface" above): a module importing a name the shim lacks
// fails to evaluate, and with it every registration this file exists to
// record (the 127-test incident). See ui/gsxui.js's own doc comments for
// the semantics; none are repeated here.

const CAP_PADDING = 8;

const tracked = new Map();

export function position(content, anchor, opts = {}) {
  release(content);
  const {
    side = "bottom",
    align = "start",
    sideOffset = 0,
    alignOffset = 0,
    flip = true,
    cap = false,
  } = opts;
  const anchorRect =
    anchor instanceof Element ? () => anchor.getBoundingClientRect() : () => anchor;
  const update = () =>
    applyPlacement(content, anchorRect(), {
      side,
      align,
      sideOffset,
      alignOffset,
      flip,
      cap,
    });
  const onToggle = (e) => {
    if (e.newState === "closed") release(content);
  };
  content.addEventListener("toggle", onToggle);
  window.addEventListener("scroll", update, { capture: true, passive: true });
  window.addEventListener("resize", update);
  tracked.set(content, {
    update,
    onToggle,
    hadSide: content.hasAttribute("data-side"),
    authoredSide: content.getAttribute("data-side"),
  });
  update();
}

export function release(content) {
  const entry = tracked.get(content);
  if (!entry) return;
  tracked.delete(content);
  content.removeEventListener("toggle", entry.onToggle);
  window.removeEventListener("scroll", entry.update, { capture: true });
  window.removeEventListener("resize", entry.update);
  if (entry.hadSide) content.setAttribute("data-side", entry.authoredSide);
  else content.removeAttribute("data-side");
}

function applyPlacement(content, a, o) {
  content.style.position = "fixed";
  content.style.inset = "auto";
  if (o.cap) content.style.removeProperty("--gsxui-available-height");
  const vw = document.documentElement.clientWidth;
  const vh = document.documentElement.clientHeight;
  const w = content.offsetWidth;
  let h = content.offsetHeight;
  const vertical = o.side === "top" || o.side === "bottom";
  let side = o.side;
  if (o.flip) {
    const room = {
      top: a.top - o.sideOffset,
      bottom: vh - a.bottom - o.sideOffset,
      left: a.left - o.sideOffset,
      right: vw - a.right - o.sideOffset,
    };
    const opposite = { top: "bottom", bottom: "top", left: "right", right: "left" }[side];
    const need = vertical ? h : w;
    if (need > room[side] && room[opposite] > room[side]) side = opposite;
  }
  if (o.cap && vertical) {
    const avail =
      (side === "bottom" ? vh - a.bottom - o.sideOffset : a.top - o.sideOffset) -
      CAP_PADDING;
    content.style.setProperty(
      "--gsxui-available-height",
      `${Math.max(0, Math.floor(avail))}px`,
    );
    h = content.offsetHeight;
  }
  let left, top;
  if (vertical) {
    top = side === "bottom" ? a.bottom + o.sideOffset : a.top - o.sideOffset - h;
    left = alignedCoord(o.align, a.left, a.width, w) + o.alignOffset;
    left = Math.max(0, Math.min(left, vw - w));
  } else {
    left = side === "right" ? a.right + o.sideOffset : a.left - o.sideOffset - w;
    top = alignedCoord(o.align, a.top, a.height, h) + o.alignOffset;
    top = Math.max(0, Math.min(top, vh - h));
  }
  left = Math.max(0, Math.min(left, vw - w));
  top = Math.max(0, Math.min(top, vh - h));
  if (o.flip) content.dataset.side = side;
  content.style.left = `${left}px`;
  content.style.top = `${top}px`;
}

function alignedCoord(align, start, anchorSize, contentSize) {
  if (align === "center") return start + anchorSize / 2 - contentSize / 2;
  if (align === "end") return start + anchorSize - contentSize;
  return start;
}
