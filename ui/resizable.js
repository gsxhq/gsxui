// Resizable behavior — the plain-JS half of the react-resizable-panels
// substitution (see ui/resizable.gsx's own package doc comment): the
// library itself is absent from the reference checkout, so nothing here is
// ported from its dist build. Two adjacent panels around a dragged/keyed
// handle resize by shifting a single shared boundary — sibling panels
// further away are untouched, matching a real splitter, not a
// redistribute-everything model. Percentages, not pixels, are the unit of
// record (data-min-size/data-max-size and the emitted gsxui:change sizes
// are all percentages of the group's own PANELS' combined px size — see
// FIX below, not the group's own box), so a later viewport resize doesn't
// invalidate a size a caller persisted from a gsxui:change payload.
//
// FIX (2026-07-24 review round 2, CRITICAL): panels are laid out via
// `flex: <n> 1 0px` (see ui/resizable.gsx's own doc comment) — a LENGTH
// zero basis, not a percentage one, because a percentage flex-basis can't
// resolve against an indefinite main-axis container size (e.g. a vertical
// group sized only by `min-h-[...]`), which made every non-`max-w`-pinned
// example simply not resizable at all: a written percentage flex-basis
// silently failed to apply, panels fell back to content size, and nothing
// moved. With a `0px` basis, grow distributes 100% of the FREE space
// (content-box size minus every panel's own 0px basis, i.e. minus nothing
// panel-side but STILL minus every handle's own rendered width/height) —
// so applyDeltaPct below writes `style.flexGrow`, not `style.flexBasis`,
// and every percentage computed FROM pixels in this file (drag start,
// pointermove, keyboard step/Home/End, aria sync) is computed against the
// PANELS' OWN summed px size (panelsTotalPx), not the group's full
// clientWidth/clientHeight — using the full box would introduce a scaling
// error proportional to total handle thickness, since grow only ever
// distributes what's left AFTER the handles' own space is already spoken
// for. currentSizes() already worked this way (round 1's own fix); this
// round makes the drag/keyboard math consistent with it.
import { on, emit, isRTL, init } from "./gsxui.js";

const HANDLE_SELECTOR = "[data-gsxui-slot-resizable-handle]";

const rootOf = (el) => el.closest("[data-gsxui-slot-resizable-panel-group]");

const isPanel = (el) => !!el && el.hasAttribute("data-gsxui-slot-resizable-panel");

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

function panelsOf(root) {
  return [...root.children].filter(isPanel);
}

function panelSize(panel, vertical) {
  const rect = panel.getBoundingClientRect();
  return vertical ? rect.height : rect.width;
}

// Summed px size of root's own direct-child PANELS (handles excluded) —
// the one denominator used everywhere a pixel needs to become a percentage
// in this file. See the FIX note in the module header comment for why this
// replaces the group's own clientWidth/clientHeight.
function panelsTotalPx(root, vertical) {
  return panelsOf(root).reduce((sum, p) => sum + panelSize(p, vertical), 0);
}

function readPct(el, key, fallback) {
  const raw = el.dataset[key];
  const n = raw === undefined ? NaN : parseFloat(raw);
  return Number.isFinite(n) ? n : fallback;
}

// One entry per direct-child panel of root, in DOM order — the shape
// gsxui:change's own `sizes` payload documents. Normalized against the
// PANELS' OWN total (excluding handle widths) and rounded to 2 decimal
// places, so a caller persisting this payload gets clean numbers that
// actually sum to 100 — reporting raw panel-size/group-size fractions
// instead would leak every handle's few px as missing space (e.g. two
// panels either side of one handle reporting ~49.95/49.95) and carry
// binary-float noise forever.
function currentSizes(root) {
  const vertical = isVertical(root);
  const rawSizes = panelsOf(root).map((panel) => panelSize(panel, vertical));
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
  const total = panelsTotalPx(root, vertical);
  const prevSizePct = total ? (panelSize(neighbours.prev, vertical) / total) * 100 : 0;
  syncAria(handle, prevSizePct, readPct(neighbours.prev, "minSize", 0), readPct(neighbours.prev, "maxSize", 100));
}

// A commit (drag end or a keyboard step/Home/End) can move a PANEL that is
// itself another handle's own neighbour — e.g. in a 3-panel group, handle 1
// resizes panel 2, which is also handle 2's `prev` — so every handle
// sharing this root needs its aria-value* re-derived after a commit, not
// just the one the user actually touched.
function syncGroupAria(root) {
  for (const el of root.children) {
    if (el.matches?.(HANDLE_SELECTOR)) syncHandleAria(el);
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
  // flex-GROW, not flex-basis (review round 2 FIX — see the module header
  // comment): every panel's basis is a fixed `0px` (server-rendered), so
  // the boundary moves by reweighting how the two panels split the group's
  // free space, not by resizing either panel's own basis length.
  prev.style.flexGrow = String(prevSize);
  next.style.flexGrow = String(nextSize);
  syncAria(handle, prevSize, prevMin, prevMax);
}

// --- Pointer drag ---------------------------------------------------------

let drag = null;

on("pointerdown", HANDLE_SELECTOR, (e, handle) => {
  const neighbours = neighboursOf(handle);
  const root = rootOf(handle);
  if (!neighbours || !root) return;
  const vertical = isVertical(root);
  const total = panelsTotalPx(root, vertical);
  if (!total) return;
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
    panelsTotalPx: total,
    ...neighbours,
    prevStartPct: (panelSize(neighbours.prev, vertical) / total) * 100,
    nextStartPct: (panelSize(neighbours.next, vertical) / total) * 100,
    pointerStart: vertical ? e.clientY : e.clientX,
  };
});

// pointermove/up/cancel: once a handle has captured the pointer, every
// subsequent pointer event for that pointerId targets the capturing
// element itself (per the Pointer Events spec) regardless of where the
// pointer physically is — closest() below always resolves back to the same
// handle.
on("pointermove", HANDLE_SELECTOR, (e, handle) => {
  if (!drag || drag.handle !== handle || drag.pointerId !== e.pointerId) return;
  const pos = drag.vertical ? e.clientY : e.clientX;
  // Horizontal only: under RTL, the inline-start (first) panel sits
  // visually on the right, so a pointer moving toward the viewport-left
  // (decreasing clientX) must still GROW it — negate the raw pixel delta.
  const rtlFlip = !drag.vertical && isRTL(drag.root) ? -1 : 1;
  const deltaPct = (((pos - drag.pointerStart) * rtlFlip) / drag.panelsTotalPx) * 100;
  applyDeltaPct(drag.handle, drag.prev, drag.next, drag.prevStartPct, drag.nextStartPct, deltaPct);
});

function endDrag(_e, handle) {
  if (!drag || drag.handle !== handle) return;
  const root = drag.root;
  drag = null;
  syncGroupAria(root);
  emit(root, "gsxui:change", { sizes: currentSizes(root) });
}

on("pointerup", HANDLE_SELECTOR, endDrag);
on("pointercancel", HANDLE_SELECTOR, endDrag);
on("lostpointercapture", HANDLE_SELECTOR, endDrag);

// --- Keyboard --------------------------------------------------------------

// Step size: derived-not-read (the map's own `## resizable` §5 — the
// library's real keyboard step was never read, react-resizable-panels
// being absent from the checkout entirely) — 10 percentage points per
// press, a gsxui-authored value, not ported from the library.
const STEP = 10;

on("keydown", HANDLE_SELECTOR, (e, handle) => {
  const neighbours = neighboursOf(handle);
  const root = rootOf(handle);
  if (!neighbours || !root) return;
  const { prev, next } = neighbours;
  const vertical = isVertical(root);
  const total = panelsTotalPx(root, vertical);
  if (!total) return;

  // Horizontal only: native-direction convention — keys mirror under RTL,
  // so ArrowLeft still GROWS the first (inline-start) panel, same as
  // ArrowRight does under LTR.
  const rtlFlip = !vertical && isRTL(root) ? -1 : 1;
  const stepKeys = vertical
    ? { ArrowUp: -STEP, ArrowDown: STEP }
    : { ArrowLeft: -STEP * rtlFlip, ArrowRight: STEP * rtlFlip };
  const prevStartPct = (panelSize(prev, vertical) / total) * 100;
  const nextStartPct = (panelSize(next, vertical) / total) * 100;

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

// --- init --------------------------------------------------------------
// Sync aria-valuenow/-min/-max from the server-rendered geometry (no context
// needed — every handle's neighbours are its own DOM siblings). Self-healing
// via init() (ui/gsxui.js): current matches, later-added matches (e.g. a DOM
// swap), and any match morphed back to server state all get this
// initializer re-run. Both writes are idempotent style/aria writes with no
// per-call resource bound — touchAction is a plain style property
// (re-setting it to the same value is a no-op) and syncHandleAria() is pure
// reflection off the handle's current DOM geometry — so re-running the pair
// is safe.
function initRoot(handle) {
  // Without this, touch input's default pan/scroll gesture wins the
  // pointerdown-to-pointermove race — the browser claims the gesture and
  // fires pointercancel mid-drag, so resizing never works on touch at all.
  // The pinned class string itself stays untouched; this is a runtime
  // style, not a class-string change.
  handle.style.touchAction = "none";
  syncHandleAria(handle);
}
init(HANDLE_SELECTOR, initRoot);
