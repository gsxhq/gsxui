// gsxui delegation core. One document-level listener per (event type, phase);
// behaviors register selector→handler pairs — elements added later (DOM
// swaps, innerHTML) just work with no re-scan. init() (below) covers the
// complementary case: per-instance init WORK that on()'s delegation model
// doesn't reach for free.
//
// Non-bubbling events (toggle, close, focus, blur, …) must be registered with
// { capture: true } — the document-level listener only sees them during the
// capture descent.

const registry = new Map(); // "type:capture" -> [{ selector, handler }]

export function on(type, selector, handler, { capture = false } = {}) {
  const key = `${type}:${capture}`;
  let handlers = registry.get(key);
  if (!handlers) {
    handlers = [];
    registry.set(key, handlers);
    document.addEventListener(
      type,
      (event) => dispatch(handlers, event),
      capture,
    );
  }
  handlers.push({ selector, handler });
}

function dispatch(handlers, event) {
  const target = event.target instanceof Element ? event.target : null;
  if (!target) return;
  for (const { selector, handler } of handlers) {
    const el = target.closest(selector);
    if (el) handler(event, el);
  }
}

export function emit(el, name, detail) {
  return el.dispatchEvent(
    new CustomEvent(name, { bubbles: true, cancelable: true, detail }),
  );
}

// --- dynamic init ---------------------------------------------------------
//
// on() needs no re-scan for late DOM — document-level delegation already
// covers it. Init WORK (reflecting server-rendered state, per-instance
// wiring) does not get that for free: a module's load-time scan runs once.
// init(selector, fn) closes the gap: fn runs for current matches at
// registration, for matches added later (DOM swaps, morphs, innerHTML
// writes), and again when a match's subtree is morphed back to server
// state.
//
// Contract: fn is IDEMPOTENT. The core deliberately re-runs it on mutated
// roots — a component needing true once-per-element wiring guards inside
// its own fn (WeakSet / data flag). Mutations caused by an init pass
// itself are discarded — see takeRecords() in flush() below, the sole
// loop-prevention mechanism; a concurrent external mutation landing in
// that exact window is missed until the next mutation re-heals it.
// Double-registering the same (selector, fn) runs fn twice per match —
// callers are expected to register once at module load. By convention,
// the registered function is named `initRoot` and lives in each module's
// trailing `--- init ---` section; the shared observer is instantiated by
// whichever barrel import calls init() first (avatar today).
//
// WARNING: fn fires on ANY attribute, text, or childList mutation at or
// under a match — not just an external DOM swap/morph. A live interaction
// inside the same subtree (e.g. a roving-tabindex module writing
// tabindex="0"/"-1" on click/arrow-key, a toggle handler flipping
// data-state) is itself such a mutation, and schedules fn again one
// microtask later. Two obligations follow: fn must be CHEAP (it may run
// once per interaction, not once per page load), and it must not FIGHT
// live interaction state a handler just set (an unconditional "reset to
// entry state" fn will stomp its own component's just-changed tab
// stop/open state right back — see ui/menubar.js's and
// ui/toggle-group.js's own normalize() for the "bail when a live value
// already exists" guard this forces).

const inits = []; // [{ selector, fn }]
let observer = null;
let pending = null; // Map<fn, Set<Element>> while a flush is queued

export function init(selector, fn) {
  inits.push({ selector, fn });
  for (const el of document.querySelectorAll(selector)) fn(el);
  if (!observer) {
    observer = new MutationObserver(collect);
    observer.observe(document.documentElement, {
      childList: true,
      subtree: true,
      attributes: true,
      characterData: true,
    });
  }
}

function schedule(fn, el) {
  if (!pending) {
    pending = new Map();
    queueMicrotask(flush);
  }
  let set = pending.get(fn);
  if (!set) pending.set(fn, (set = new Set()));
  set.add(el);
}

function collect(records) {
  for (const record of records) {
    if (record.type === "childList") {
      for (const node of record.addedNodes) {
        if (!(node instanceof Element)) continue;
        for (const { selector, fn } of inits) {
          if (node.matches(selector)) schedule(fn, node);
          for (const el of node.querySelectorAll(selector)) schedule(fn, el);
        }
      }
    }
    const target =
      record.target instanceof Element
        ? record.target
        : record.target.parentElement;
    if (!target) continue;
    for (const { selector, fn } of inits) {
      const root = target.closest(selector);
      if (root) schedule(fn, root);
    }
  }
}

function flush() {
  const batch = pending;
  pending = null;
  if (!batch) return;
  try {
    for (const [fn, els] of batch) {
      for (const el of els) if (el.isConnected) fn(el);
    }
  } finally {
    observer.takeRecords(); // discard mutations our own inits produced
  }
}

// --- shared init helpers --------------------------------------------------

// One counter for every generated id in the library. Uniqueness within the
// document is the entire contract — which module the number came from never
// matters, so modules share it instead of each keeping a private counter.
let uidCounter = 0;
export function uid(prefix) {
  return `${prefix}-${++uidCounter}`;
}

// Group→label aria wiring shared by listbox-shaped components (select,
// combobox): every group under root that lacks aria-labelledby gets pointed
// at its own label, generating the label's id on demand. Skip-if-present
// per group keeps it idempotent under re-init.
export function wireGroupLabels(root, groupSelector, labelSelector, idPrefix) {
  for (const group of root.querySelectorAll(groupSelector)) {
    if (group.getAttribute("aria-labelledby")) continue;
    const label = group.querySelector(labelSelector);
    if (!label) continue;
    if (!label.id) label.id = uid(idPrefix);
    group.setAttribute("aria-labelledby", label.id);
  }
}

// Initializers re-run on ANY mutation under a match (see init() above), so
// they split into two kinds of work: REFLECTION (recompute state from the
// DOM — must re-run unguarded, that is the whole point of re-init) and
// RESOURCE BINDS (timers, listeners, per-element observers — must happen
// once per element or they stack). once(fn) is the guard for the second
// kind: the returned function runs fn at most once per element, forever.
// It is NOT for semantic "a live value already exists" bails — those must
// re-arm after a morph strips the value (see hasLiveTabStop below).
export function once(fn) {
  const ran = new WeakSet();
  return (el) => {
    if (ran.has(el)) return;
    ran.add(el); // Marks before invoking: a throwing fn stays marked (no retry) — callers' binds must not throw.
    fn(el);
  };
}

// Roving-tabindex bail shared by menubar and toggle-group: "has normalize
// already picked an entry tab stop, or has the user moved it?" Checks the
// tabindex ATTRIBUTE, never the IDL property — a plain <button> reports
// .tabIndex === 0 even with no attribute at all, which would make this
// always-true on server-fresh DOM and prevent the initial roving collapse.
// Server markup renders no tabindex attributes; interaction and normalize
// write them; a morph back to server state strips them — exactly the reset
// that must re-trigger normalize, which is why this is a semantic check
// and not a once() guard.
export function hasLiveTabStop(els) {
  return [...els].some((el) => el.getAttribute("tabindex") === "0");
}

// --- anchored positioning --------------------------------------------------
//
// position() is the shared placement engine every floating surface in ui/
// uses (popover, hover-card, tooltip, dropdown-menu + sub, select, menubar
// + sub, combobox, navigation-menu, context-menu + sub). It reproduces the
// slice of Radix's popper (Floating UI middleware) those components rely
// on, in Radix's middleware order:
//
//   offset (sideOffset/alignOffset) → flip (main axis) → shift (cross
//   axis) → hard clamp (both axes, last resort) → size (the
//   --gsxui-available-height cap) → autoUpdate (window scroll capture +
//   resize while open).
//
// DELTAS from Radix popper, all deliberate — this is a small parity layer,
// not a Floating-UI port:
//   - flip is MAIN AXIS ONLY: preferred side vs its opposite. No
//     cross-axis fallback placements, no fallbackAxisSideDirection.
//   - no anchor-detached "hide" middleware (content is never auto-hidden
//     when the trigger scrolls out of view).
//   - autoUpdate is window scroll (capture, passive) + resize only — no
//     MutationObserver / element-resize tracking. gsxui's DOM is
//     server-rendered and static while a popover is open; scroll and
//     resize are the movements that actually occur.
//   - no arrow support, no collisionBoundary/collisionPadding options:
//     the boundary is always the viewport (documentElement client box —
//     innerWidth/Height include classic-scrollbar gutters and clamping
//     against them tucks an edge under the scrollbar).
//
// Call AFTER showPopover() — a hidden popover has no layout box. The
// engine sets position:fixed / inset:auto / left / top itself, stamps
// data-side with the PLACED side (the enter/exit animations key on
// data-[side=...], so a flipped placement must animate from the flipped
// direction, exactly as Radix recomputes placement), and restores the
// authored data-side on close so the next open starts from the preference.
// Cleanup is automatic: the content's own "toggle" (newState "closed")
// releases the listeners, so a module only ever calls position().

const CAP_PADDING = 8; // px kept clear of the viewport edge by the size cap

// tracked: content element -> { update, onToggle, hadSide, authoredSide }.
// One entry per open surface; position() on an already-tracked content
// releases first, so listeners attach exactly once per open.
const tracked = new Map();

// position places `content` against `anchor` (an Element, re-measured on
// every update so the content tracks it through scroll/resize, or a static
// DOMRect-like object — context-menu's pointer anchor) and keeps it placed
// while the popover stays open.
//   side:        "top" | "bottom" | "left" | "right" (preferred side)
//   align:       "start" | "center" | "end" (cross-axis alignment)
//   sideOffset:  gap between anchor and content on the main axis
//   alignOffset: shift along the cross axis from the aligned position
//   flip:        set false to keep the preferred side unconditionally AND
//                leave data-side untouched (context-menu's top-level content
//                deliberately has none — see context-menu.gsx)
//   cap:         set --gsxui-available-height (px from the placed edge to
//                the viewport edge minus CAP_PADDING) for recipe CSS to
//                consume as a max-height arm (select, combobox)
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
    anchor instanceof Element
      ? () => anchor.getBoundingClientRect()
      : () => anchor;
  const update = () => {
    if (!content.isConnected) return release(content);
    applyPlacement(content, anchorRect(), {
      side,
      align,
      sideOffset,
      alignOffset,
      flip,
      cap,
    });
  };
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

// release detaches the reposition listeners and restores the authored
// data-side. Idempotent; runs automatically on the content's close toggle.
// A content element REMOVED from the DOM while open (e.g. a boosted-swap
// navigation mid-open) never fires that toggle; update() covers this by
// self-releasing on the first scroll/resize that finds the content
// disconnected — exactly the events the leaked entry is subscribed to.
export function release(content) {
  const entry = tracked.get(content);
  if (!entry) return;
  tracked.delete(content);
  content.removeEventListener("toggle", entry.onToggle);
  window.removeEventListener("scroll", entry.update, { capture: true });
  window.removeEventListener("resize", entry.update);
  // --gsxui-available-height is deliberately NOT removed here: release runs
  // on the close toggle, while the 150ms discrete-transition exit is still
  // playing, and un-capping max-height mid-exit makes the closing box grow.
  // applyPlacement clears it before measuring on the next open instead.
  // Restoring data-side now is safe — only the @starting-style ENTER arms
  // key on data-[side=...]; the exit is a side-agnostic fade/scale.
  if (entry.hadSide) content.setAttribute("data-side", entry.authoredSide);
  else content.removeAttribute("data-side");
}

// isRTL reports the resolved direction at el — computed style, so it honors
// both dir attributes and CSS `direction`.
export function isRTL(el) {
  return getComputedStyle(el).direction === "rtl";
}

function applyPlacement(content, a, o) {
  content.style.position = "fixed";
  content.style.inset = "auto";
  // Clear a previous update's cap before measuring — a stale (smaller)
  // cap would shrink the box and corrupt the flip decision's natural size.
  if (o.cap) content.style.removeProperty("--gsxui-available-height");
  const vw = document.documentElement.clientWidth;
  const vh = document.documentElement.clientHeight;
  const w = content.offsetWidth;
  let h = content.offsetHeight;

  // Resolve logical options (side/align/alignOffset) to physical ones
  // before the rest of the math — anchor rects have no element when a
  // virtual rect is passed (context-menu's pointer anchor), so direction
  // is read from `content`, which lives in the same direction context.
  const vertical = o.side === "top" || o.side === "bottom";
  let side = o.side;
  let align = o.align;
  let alignOffset = o.alignOffset;
  if (isRTL(content)) {
    if (vertical) {
      // horizontal cross-axis: mirror alignment and its offset
      if (align === "start") align = "end";
      else if (align === "end") align = "start";
      alignOffset = -alignOffset;
    } else {
      // horizontal main axis: mirror the preferred side
      side = side === "left" ? "right" : "left";
    }
  }

  // Main-axis flip: if the preferred side can't hold the content and the
  // opposite side has more room, flip (Radix flip middleware, main axis).
  if (o.flip) {
    const room = {
      top: a.top - o.sideOffset,
      bottom: vh - a.bottom - o.sideOffset,
      left: a.left - o.sideOffset,
      right: vw - a.right - o.sideOffset,
    };
    const opposite = {
      top: "bottom",
      bottom: "top",
      left: "right",
      right: "left",
    }[side];
    const need = vertical ? h : w;
    if (need > room[side] && room[opposite] > room[side]) side = opposite;
  }
  // Size cap (Radix size middleware, height only): expose the room on the
  // PLACED side; recipe CSS min()s it into its max-height, so re-measure.
  if (o.cap && vertical) {
    const avail =
      (side === "bottom"
        ? vh - a.bottom - o.sideOffset
        : a.top - o.sideOffset) - CAP_PADDING;
    content.style.setProperty(
      "--gsxui-available-height",
      `${Math.max(0, Math.floor(avail))}px`,
    );
    h = content.offsetHeight;
  }
  let left, top;
  if (vertical) {
    top =
      side === "bottom" ? a.bottom + o.sideOffset : a.top - o.sideOffset - h;
    left = alignedCoord(align, a.left, a.width, w) + alignOffset;
    left = Math.max(0, Math.min(left, vw - w)); // cross-axis shift
  } else {
    left =
      side === "right" ? a.right + o.sideOffset : a.left - o.sideOffset - w;
    top = alignedCoord(align, a.top, a.height, h) + alignOffset;
    top = Math.max(0, Math.min(top, vh - h)); // cross-axis shift
  }
  // Last-resort clamp on BOTH axes (the pre-flip clampToViewport
  // behaviour, preserved): even a flipped side that still overflows stays
  // inside the viewport.
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
