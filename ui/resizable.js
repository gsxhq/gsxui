// Resizable behavior — the plain-JS half of the react-resizable-panels
// substitution (see ui/resizable.gsx's own package doc comment): the
// library itself is absent from the reference checkout, so nothing here is
// ported from its dist build. Two adjacent panels around a dragged/keyed
// handle resize by shifting a single shared boundary — sibling panels
// further away are untouched, matching a real splitter, not a
// redistribute-everything model. Percentages, not pixels, are the unit of
// record (data-min-size/data-max-size and the emitted gsxui:change sizes
// are all percentages of the group's own content-box size along its main
// axis), so a later viewport resize doesn't invalidate a size a caller
// persisted from a gsxui:change payload.
import { on, emit } from "./gsxui.js";

const rootOf = (el) => el.closest("[data-gsxui-resizable]");

const isPanel = (el) => !!el && el.dataset.slot === "resizable-panel";

// The handle's two neighbours are its immediate DOM siblings — true
// regardless of nesting, since a nested ResizablePanelGroup always sits
// INSIDE a ResizablePanel (one level deeper), never as a direct sibling of
// an outer handle.
function neighboursOf(handle) {
  const prev = handle.previousElementSibling;
  const next = handle.nextElementSibling;
  if (!isPanel(prev) || !isPanel(next)) return null;
  return { prev, next };
}

function isVertical(root) {
  return root.getAttribute("aria-orientation") === "vertical";
}

// Main-axis content-box size of the group itself — clientWidth/Height
// (excludes border, includes padding; the group is expected to carry no
// padding of its own, so this is effectively the content box).
function groupSize(root, vertical) {
  return vertical ? root.clientHeight : root.clientWidth;
}

function panelSize(panel, vertical) {
  const rect = panel.getBoundingClientRect();
  return vertical ? rect.height : rect.width;
}

function readPct(el, key, fallback) {
  const raw = el.dataset[key];
  const n = raw === undefined ? NaN : parseFloat(raw);
  return Number.isFinite(n) ? n : fallback;
}

// One entry per direct-child panel of root, in DOM order — the shape
// gsxui:change's own `sizes` payload documents. Normalized against the
// PANELS' OWN total (excluding handle widths, unlike the drag/keyboard
// math above which deliberately works in fractions of the group's full
// content box) and rounded to 2 decimal places, so a caller persisting
// this payload gets clean numbers that actually sum to 100 — reporting
// raw panel-size/group-size fractions instead would leak every handle's
// few px as missing space (e.g. two panels either side of one handle
// reporting ~49.95/49.95) and carry binary-float noise forever.
function currentSizes(root) {
  const vertical = isVertical(root);
  const rawSizes = [...root.children].filter(isPanel).map((panel) => panelSize(panel, vertical));
  const total = rawSizes.reduce((sum, s) => sum + s, 0);
  if (!total) return rawSizes.map(() => 0);
  return rawSizes.map((s) => Math.round((s / total) * 100 * 100) / 100);
}

function syncAria(handle, prevSizePct, prevMin, prevMax) {
  handle.setAttribute("aria-valuenow", String(Math.round(prevSizePct)));
  handle.setAttribute("aria-valuemin", String(Math.round(prevMin)));
  handle.setAttribute("aria-valuemax", String(Math.round(prevMax)));
}

// Recomputes one handle's own aria-value* from its current DOM geometry —
// the shared body behind both the module-init scan and syncGroupAria below.
function syncHandleAria(handle) {
  const neighbours = neighboursOf(handle);
  const root = rootOf(handle);
  if (!neighbours || !root) return;
  const vertical = isVertical(root);
  const size = groupSize(root, vertical);
  const prevSizePct = size ? (panelSize(neighbours.prev, vertical) / size) * 100 : 0;
  syncAria(handle, prevSizePct, readPct(neighbours.prev, "minSize", 0), readPct(neighbours.prev, "maxSize", 100));
}

// A commit (drag end or a keyboard step/Home/End) can move a PANEL that is
// itself another handle's own neighbour — e.g. in a 3-panel group, handle 1
// resizes panel 2, which is also handle 2's `prev` — so every handle
// sharing this root needs its aria-value* re-derived after a commit, not
// just the one the user actually touched.
function syncGroupAria(root) {
  for (const el of root.children) {
    if (el.dataset && el.dataset.slot === "resizable-handle") syncHandleAria(el);
  }
}

// Applies a percentage-point delta to the boundary between prev/next,
// clamped so BOTH panels stay within their own min/max (defaulting to
// 0%/100% when unset — an unconstrained panel can go anywhere the OTHER
// side's constraint still allows). prevStartPct/nextStartPct are the sizes
// the delta is relative to (the drag's recorded start, or the current
// live size for a discrete keyboard step).
function applyDeltaPct(handle, prev, next, prevStartPct, nextStartPct, deltaPct) {
  const prevMin = readPct(prev, "minSize", 0);
  const prevMax = readPct(prev, "maxSize", 100);
  const nextMin = readPct(next, "minSize", 0);
  const nextMax = readPct(next, "maxSize", 100);
  const dMin = Math.max(prevMin - prevStartPct, nextStartPct - nextMax);
  const dMax = Math.min(prevMax - prevStartPct, nextStartPct - nextMin);
  const d = Math.min(dMax, Math.max(dMin, deltaPct));
  const prevSize = prevStartPct + d;
  const nextSize = nextStartPct - d;
  prev.style.flexBasis = `${prevSize}%`;
  next.style.flexBasis = `${nextSize}%`;
  syncAria(handle, prevSize, prevMin, prevMax);
}

// --- Pointer drag ---------------------------------------------------------

let drag = null;

on("pointerdown", '[data-slot="resizable-handle"]', (e, handle) => {
  const neighbours = neighboursOf(handle);
  const root = rootOf(handle);
  if (!neighbours || !root) return;
  const vertical = isVertical(root);
  const size = groupSize(root, vertical);
  if (!size) return;
  handle.setPointerCapture(e.pointerId);
  // Capture succeeded: this gesture belongs to the handle now, not to
  // whatever native gesture the browser would otherwise start (text
  // selection over the panels' own content) or to whatever element
  // previously had focus (the keyboard path needs the handle focused).
  e.preventDefault();
  handle.focus();
  drag = {
    pointerId: e.pointerId,
    handle,
    root,
    vertical,
    groupSize: size,
    ...neighbours,
    prevStartPct: (panelSize(neighbours.prev, vertical) / size) * 100,
    nextStartPct: (panelSize(neighbours.next, vertical) / size) * 100,
    pointerStart: vertical ? e.clientY : e.clientX,
  };
});

// pointermove/up/cancel: once a handle has captured the pointer, every
// subsequent pointer event for that pointerId targets the capturing
// element itself (per the Pointer Events spec) regardless of where the
// pointer physically is — closest() below always resolves back to the same
// handle.
on("pointermove", '[data-slot="resizable-handle"]', (e, handle) => {
  if (!drag || drag.handle !== handle || drag.pointerId !== e.pointerId) return;
  const pos = drag.vertical ? e.clientY : e.clientX;
  const deltaPct = ((pos - drag.pointerStart) / drag.groupSize) * 100;
  applyDeltaPct(drag.handle, drag.prev, drag.next, drag.prevStartPct, drag.nextStartPct, deltaPct);
});

function endDrag(_e, handle) {
  if (!drag || drag.handle !== handle) return;
  const root = drag.root;
  drag = null;
  syncGroupAria(root);
  emit(root, "gsxui:change", { sizes: currentSizes(root) });
}

on("pointerup", '[data-slot="resizable-handle"]', endDrag);
on("pointercancel", '[data-slot="resizable-handle"]', endDrag);
on("lostpointercapture", '[data-slot="resizable-handle"]', endDrag);

// --- Keyboard --------------------------------------------------------------

// Step size: derived-not-read (the map's own `## resizable` §5 — the
// library's real keyboard step was never read, react-resizable-panels
// being absent from the checkout entirely) — 10 percentage points per
// press, a gsxui-authored value, not ported from the library.
const STEP = 10;

on("keydown", '[data-slot="resizable-handle"]', (e, handle) => {
  const neighbours = neighboursOf(handle);
  const root = rootOf(handle);
  if (!neighbours || !root) return;
  const { prev, next } = neighbours;
  const vertical = isVertical(root);
  const size = groupSize(root, vertical);
  if (!size) return;

  const stepKeys = vertical ? { ArrowUp: -STEP, ArrowDown: STEP } : { ArrowLeft: -STEP, ArrowRight: STEP };
  const prevStartPct = (panelSize(prev, vertical) / size) * 100;
  const nextStartPct = (panelSize(next, vertical) / size) * 100;

  let deltaPct;
  if (e.key in stepKeys) {
    deltaPct = stepKeys[e.key];
  } else if (e.key === "Home") {
    deltaPct = readPct(prev, "minSize", 0) - prevStartPct;
  } else if (e.key === "End") {
    deltaPct = readPct(prev, "maxSize", 100) - prevStartPct;
  } else {
    return;
  }

  e.preventDefault();
  applyDeltaPct(handle, prev, next, prevStartPct, nextStartPct, deltaPct);
  syncGroupAria(root);
  emit(root, "gsxui:change", { sizes: currentSizes(root) });
});

// --- Init: sync aria-valuenow/-min/-max from the server-rendered geometry --
// (no context needed — every handle's neighbours are its own DOM siblings).
// Same one-time module-load scan shape as toggle-group.js's normalize()
// loop and carousel.js's own init loop: a handle added later via an HTMX
// swap is not picked up, the same accepted limitation those modules carry.
for (const handle of document.querySelectorAll('[data-slot="resizable-handle"]')) {
  // Without this, touch input's default pan/scroll gesture wins the
  // pointerdown-to-pointermove race — the browser claims the gesture and
  // fires pointercancel mid-drag, so resizing never works on touch at all.
  // The pinned class string itself stays untouched; this is a runtime
  // style, not a class-string change.
  handle.style.touchAction = "none";
  syncHandleAria(handle);
}
