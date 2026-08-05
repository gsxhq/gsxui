// Carousel behavior — the plain-JS half of the native-scroll substitution
// for embla (see ui/carousel.gsx's own package doc comment): prev/next
// scroll-by-one-item, disabled-state/current-index bookkeeping recomputed
// off the viewport's own scrollLeft/scrollTop (not an embla API — there is
// no embla here), ArrowLeft/ArrowRight keyboard, and an imperative
// `el.gsxuiCarousel` handle standing in for embla's much larger CarouselApi
// (API-surface reduction — see the package doc comment / docs/jsx-parity.md
// `## carousel` for the full ledger: embla's scrollTo/scrollPrev/scrollNext/
// canScrollPrev/canScrollNext/selectedScrollSnap/scrollSnapList/on/off plus
// its whole plugin-extension mechanism collapse to {scrollTo,next,prev} plus
// the gsxui:carousel-select CustomEvent).
import { on, emit, isRTL, init, once } from "./gsxui.js";

const rootOf = (el) => el.closest("[data-gsxui-slot-carousel]");
const viewportOf = (root) => root.querySelector("[data-gsxui-slot-carousel-content]");
const itemsOf = (root) => [...root.querySelectorAll("[data-gsxui-slot-carousel-item]")];
const isVertical = (root) => root.dataset.orientation === "vertical";

// Sub-pixel rounding epsilon for the scrollLeft/scrollTop-vs-bounds
// comparisons below (fractional scroll offsets are routine on high-DPI/
// zoomed viewports) — same rationale the map's own behavior-contract section
// gives for embla's canScrollPrev/canScrollNext internals.
const EPS = 1;

// Rapid prev/next presses must compound by INTENT, not by mid-flight
// position: a relative scroll issued while the previous smooth scroll is
// still animating reads the interpolated position, and mandatory
// scroll-snap then rounds the compounded target back onto the snap point
// already in flight — the second press visibly does nothing. embla avoids
// this by tracking its own target index internally; same scheme here.
// `pending` holds the in-flight target index per root until scrolling
// settles (SETTLE_MS after the last scroll event — `scrollend` support is
// not yet universal), so every press advances from where the carousel is
// GOING, not where it happens to be mid-animation. Manual wheel/drag
// scrolling also settles and clears it, so stale intent can't outlive a
// user taking over.
const pending = new WeakMap();
const settleTimers = new WeakMap();
const SETTLE_MS = 150;

function bumpSettle(root) {
  clearTimeout(settleTimers.get(root));
  settleTimers.set(
    root,
    setTimeout(() => pending.delete(root), SETTLE_MS),
  );
}

// The item whose leading edge sits nearest the viewport's own leading edge
// — the item scroll-snap is resting on (or nearest mid-flight).
function nearestIndex(root) {
  const viewport = viewportOf(root);
  const items = itemsOf(root);
  if (!viewport || !items.length) return 0;
  const vertical = isVertical(root);
  const viewportEdge = vertical
    ? viewport.getBoundingClientRect().top
    : viewport.getBoundingClientRect().left;
  let index = 0;
  let nearest = Infinity;
  items.forEach((item, i) => {
    const rect = item.getBoundingClientRect();
    const edge = vertical ? rect.top : rect.left;
    const distance = Math.abs(edge - viewportEdge);
    if (distance < nearest) {
      nearest = distance;
      index = i;
    }
  });
  return index;
}

function scrollToIndex(root, index) {
  const items = itemsOf(root);
  if (!items.length) return;
  const target = Math.max(0, Math.min(items.length - 1, index));
  pending.set(root, target);
  bumpSettle(root);
  scrollToItem(root, target);
}

// One prev/next press = one item (embla's default slidesToScroll: 1),
// expressed as an absolute index move so presses stay deterministic under
// rapid clicking (see the pending-intent comment above).
function scrollByItems(root, dir) {
  const base = pending.has(root) ? pending.get(root) : nearestIndex(root);
  scrollToIndex(root, base + dir);
}

// scrollTo(index): scrolls so item `index`'s leading edge aligns with the
// viewport's own leading edge, computed off getBoundingClientRect() deltas
// rather than offsetLeft/offsetTop (whose offsetParent is not guaranteed to
// be the viewport div) — works uniformly for both single- and
// multi-item-per-view layouts.
function scrollToItem(root, index) {
  const items = itemsOf(root);
  const item = items[index];
  const viewport = viewportOf(root);
  if (!item || !viewport) return;
  const itemRect = item.getBoundingClientRect();
  const viewportRect = viewport.getBoundingClientRect();
  const vertical = isVertical(root);
  const delta = vertical ? itemRect.top - viewportRect.top : itemRect.left - viewportRect.left;
  viewport.scrollBy(vertical ? { top: delta, behavior: "smooth" } : { left: delta, behavior: "smooth" });
}

function updateDisabled(root) {
  const viewport = viewportOf(root);
  if (!viewport) return;
  const vertical = isVertical(root);
  const pos = vertical ? viewport.scrollTop : Math.abs(viewport.scrollLeft);
  const max = vertical
    ? viewport.scrollHeight - viewport.clientHeight
    : viewport.scrollWidth - viewport.clientWidth;
  const prev = root.querySelector("[data-gsxui-slot-carousel-previous]");
  const next = root.querySelector("[data-gsxui-slot-carousel-next]");
  if (prev) prev.disabled = pos <= EPS;
  if (next) next.disabled = pos >= max - EPS;
}

// The current index is the item whose leading edge is nearest the
// viewport's own leading edge — the item scroll-snap is currently resting
// on. Stamps data-current-index on the root (CSS-only dot-indicator hook,
// e.g. a caller-authored `[data-index="N"]` dot list needs no JS of its
// own) and emits gsxui:carousel-select with {index, count} — both 0-based,
// matching gsxuiCarousel.scrollTo(i)'s own indexing — only when the index
// actually changes, so a caller's "Slide X of Y" listener isn't re-run on
// every rAF tick while mid-scroll.
function updateIndex(root) {
  const items = itemsOf(root);
  if (!viewportOf(root) || !items.length) return;
  const index = nearestIndex(root);
  if (root.dataset.currentIndex === String(index)) return;
  root.dataset.currentIndex = String(index);
  emit(root, "gsxui:carousel-select", { index, count: items.length });
}

function recompute(root) {
  updateDisabled(root);
  updateIndex(root);
}

// Optional bespoke autoplay: data-gsxui-carousel-autoplay="<ms>" on the
// root. Stands in for embla-carousel-autoplay (the one plugin the docs
// demos actually use — carousel-plugin.tsx) without porting embla's whole
// plugin system; reproduces that demo's ACTUAL behavior (explicit
// onMouseEnter={plugin.stop}/onMouseLeave={plugin.reset} hover pause/
// resume), not embla Autoplay's own stopOnInteraction semantics (which
// trigger on drag/click, not hover). No loop mode (see the package doc
// comment's GAP) — autoplay simply stops advancing once it reaches the
// last slide rather than wrapping back to the first.
// GUARD (required — audited): the bindAutoplay() function binds four
// per-call pointerenter/focusin/pointerleave/focusout listeners on `root`
// plus a `timer` variable captured in ITS OWN closure. Under init()'s
// self-healing re-init (a re-init pass over an already-inited carousel,
// e.g. this root's subtree morphed back to server state), a second,
// unguarded call would add a SECOND set of listeners with their own
// independent `timer` closure — start() in the second call sees its own
// `timer` as null even while the first call's interval is still running,
// so a single pointerleave/focusout would start TWO intervals advancing
// the same carousel concurrently (and stopping one on hover wouldn't stop
// the other). Confirmed by reading the function body: it is genuinely
// stateful (setInterval handle + 4 addEventListener calls), not pure
// reflection, so it needs the same one-time-bind guard toggle-group.js's
// own per-item state and command.js's data-gsxui-index stamp use elsewhere
// in this codebase — once() (ui/gsxui.js) is exactly that guard, keyed on
// the root, since "has autoplay listeners ever been bound for this root"
// is exactly the once-only fact being guarded. The interval-configured
// check MUST be OUTSIDE once() so a later attribute mutation can still
// trigger the bind.
const bindAutoplay = once(function bindAutoplayOnce(root) {
  const ms = Number(root.dataset.gsxuiCarouselAutoplay);
  let timer = null;
  const stop = () => {
    if (timer) clearInterval(timer);
    timer = null;
  };
  const start = () => {
    if (timer) return;
    timer = setInterval(() => {
      const viewport = viewportOf(root);
      if (!viewport) return stop();
      const vertical = isVertical(root);
      const pos = vertical ? viewport.scrollTop : Math.abs(viewport.scrollLeft);
      const max = vertical
        ? viewport.scrollHeight - viewport.clientHeight
        : viewport.scrollWidth - viewport.clientWidth;
      if (pos >= max - EPS) return stop();
      scrollByItems(root, 1);
    }, ms);
  };
  root.addEventListener("pointerenter", stop);
  root.addEventListener("focusin", stop);
  root.addEventListener("pointerleave", start);
  root.addEventListener("focusout", start);
  start();
});

function initAutoplay(root) {
  const ms = Number(root.dataset.gsxuiCarouselAutoplay);
  if (!ms) return;
  bindAutoplay(root);
}

on("click", "[data-gsxui-slot-carousel-previous]", (_e, btn) => {
  const root = rootOf(btn);
  if (root) scrollByItems(root, -1);
});

on("click", "[data-gsxui-slot-carousel-next]", (_e, btn) => {
  const root = rootOf(btn);
  if (root) scrollByItems(root, 1);
});

// ArrowLeft/ArrowRight always map to prev/next, even for orientation=
// "vertical" carousels — mirroring shadcn's own source exactly (Carousel's
// onKeyDownCapture hard-codes ArrowLeft => scrollPrev()/ArrowRight =>
// scrollNext() unconditionally, never ArrowUp/ArrowDown, regardless of
// axis).
on("keydown", "[data-gsxui-slot-carousel]", (e, root) => {
  let dir = { ArrowLeft: -1, ArrowRight: 1 }[e.key];
  if (!dir) return;
  if (isRTL(root)) dir = -dir;
  e.preventDefault();
  scrollByItems(root, dir);
});

// scroll doesn't bubble — delegated via capture, same pattern ui/gsxui.js's
// own header comment documents for toggle/close/focus/blur. rAF-throttled:
// scroll fires far faster than layout needs to be re-measured for disabled-
// state/current-index bookkeeping.
const scheduled = new WeakSet();
on(
  "scroll",
  "[data-gsxui-slot-carousel-content]",
  (_e, viewport) => {
    const root = rootOf(viewport);
    if (!root) return;
    bumpSettle(root);
    if (scheduled.has(root)) return;
    scheduled.add(root);
    requestAnimationFrame(() => {
      scheduled.delete(root);
      recompute(root);
    });
  },
  { capture: true },
);

// ResizeObserver on every viewport: layout changes (e.g. a responsive
// basis-* breakpoint changing how many slides fit) can flip disabled state
// or the resting index without a scroll event ever firing.
const resizeObserver = new ResizeObserver((entries) => {
  for (const entry of entries) {
    const root = rootOf(entry.target);
    if (root) recompute(root);
  }
});

// --- init --------------------------------------------------------------

// Self-healing via init() (ui/gsxui.js): current matches, later-added
// matches (e.g. a DOM swap), and any match morphed back to server state
// all get initRoot() re-run. Audited per-step, since this is the one
// STATEFUL initializer in this batch (see the task-3 report for the full
// writeup):
//   - root.gsxuiCarousel reassign: a fresh object with the same shape every
//     call — safe to overwrite.
//   - recompute(root): reflection off the viewport's own scroll geometry —
//     safe to re-run.
//   - initAutoplay(root): binds per-call resources (a setInterval handle
//     plus 4 event listeners) — internally guarded above via once() so a
//     re-init is a no-op past the first bind.
//   - resizeObserver.observe(viewport): per MDN, calling observe() again
//     with a target already in the observer's own [[ObservationTargets]]
//     is a no-op (the spec: "If target is in this.[[ObservationTargets]],
//     then return") — no duplicate ResizeObserverEntry per resize, so this
//     needs no guard of its own.
function initRoot(root) {
  root.gsxuiCarousel = {
    scrollTo: (index) => scrollToIndex(root, index),
    next: () => scrollByItems(root, 1),
    prev: () => scrollByItems(root, -1),
  };
  recompute(root);
  initAutoplay(root);
  const viewport = viewportOf(root);
  if (viewport) resizeObserver.observe(viewport);
}
init("[data-gsxui-slot-carousel]", initRoot);
