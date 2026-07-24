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
import { on, emit } from "./gsxui.js";

const rootOf = (el) => el.closest("[data-gsxui-carousel]");
const viewportOf = (root) => root.querySelector('[data-slot="carousel-content"]');
const itemsOf = (root) => [...root.querySelectorAll('[data-slot="carousel-item"]')];
const isVertical = (root) => root.dataset.orientation === "vertical";

// Sub-pixel rounding epsilon for the scrollLeft/scrollTop-vs-bounds
// comparisons below (fractional scroll offsets are routine on high-DPI/
// zoomed viewports) — same rationale the map's own behavior-contract section
// gives for embla's canScrollPrev/canScrollNext internals.
const EPS = 1;

// Scroll amount for one prev/next press: the FIRST item's own measured
// width (border-box, so its own pl-4/pt-4 spacing padding is already
// included) — not one full viewport width, which would visibly skip slides
// whenever more than one slide is visible per view (the size/orientation
// demos). Because CarouselContent's flex track lays items out with no gap
// property (spacing comes from each item's own padding plus the track's
// compensating negative margin), one item's rendered width IS the distance
// from its own start to the next item's start — embla's default
// slidesToScroll: 1 behavior, reproduced without needing embla's own
// per-slide bookkeeping.
function itemExtent(root) {
  const items = itemsOf(root);
  if (!items.length) return 0;
  const rect = items[0].getBoundingClientRect();
  return isVertical(root) ? rect.height : rect.width;
}

function scrollByItems(root, dir) {
  const viewport = viewportOf(root);
  if (!viewport) return;
  const amount = itemExtent(root) * dir;
  if (!amount) return;
  viewport.scrollBy(isVertical(root) ? { top: amount, behavior: "smooth" } : { left: amount, behavior: "smooth" });
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
  const pos = vertical ? viewport.scrollTop : viewport.scrollLeft;
  const max = vertical
    ? viewport.scrollHeight - viewport.clientHeight
    : viewport.scrollWidth - viewport.clientWidth;
  const prev = root.querySelector("[data-gsxui-carousel-prev]");
  const next = root.querySelector("[data-gsxui-carousel-next]");
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
  const viewport = viewportOf(root);
  const items = itemsOf(root);
  if (!viewport || !items.length) return;
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
function initAutoplay(root) {
  const ms = Number(root.dataset.gsxuiCarouselAutoplay);
  if (!ms) return;
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
      const pos = vertical ? viewport.scrollTop : viewport.scrollLeft;
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
}

on("click", "[data-gsxui-carousel-prev]", (_e, btn) => {
  const root = rootOf(btn);
  if (root) scrollByItems(root, -1);
});

on("click", "[data-gsxui-carousel-next]", (_e, btn) => {
  const root = rootOf(btn);
  if (root) scrollByItems(root, 1);
});

// ArrowLeft/ArrowRight always map to prev/next, even for orientation=
// "vertical" carousels — mirroring shadcn's own source exactly (Carousel's
// onKeyDownCapture hard-codes ArrowLeft => scrollPrev()/ArrowRight =>
// scrollNext() unconditionally, never ArrowUp/ArrowDown, regardless of
// axis).
on("keydown", "[data-gsxui-carousel]", (e, root) => {
  const dir = { ArrowLeft: -1, ArrowRight: 1 }[e.key];
  if (!dir) return;
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
  '[data-slot="carousel-content"]',
  (_e, viewport) => {
    const root = rootOf(viewport);
    if (!root || scheduled.has(root)) return;
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

// Module-init scan, same one-time shape as toggle-group.js's normalize()
// loop and command.js's filter() loop — late-added carousels (an HTMX swap
// after this module has already run) are not picked up, the same accepted
// limitation those two modules' own init loops carry.
for (const root of document.querySelectorAll("[data-gsxui-carousel]")) {
  root.gsxuiCarousel = {
    scrollTo: (index) => scrollToItem(root, index),
    next: () => scrollByItems(root, 1),
    prev: () => scrollByItems(root, -1),
  };
  recompute(root);
  initAutoplay(root);
  const viewport = viewportOf(root);
  if (viewport) resizeObserver.observe(viewport);
}
