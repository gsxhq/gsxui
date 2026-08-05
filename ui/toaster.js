// Toaster (toasts) — the behavior module for ui/toaster.gsx. Unlike every
// other gsxui behavior module (which attaches delegated behavior to
// server-rendered markup), the toast card markup is still authored
// server-side — as the separately registered ui.Toast component — but shipped as
// inert per-type <template>s inside ui.Toaster. This module CLONES the
// matching type's template on each toast() call (never builds card DOM from
// JS string concatenation), then owns the whole lifecycle: mount → stack →
// timer → dismiss.
//
// Because the card is server markup, the same lifecycle is applied to rows
// the SERVER inserts, not just JS-triggered ones. A MutationObserver on the
// <ol> adopts any inserted `li[data-gsxui-slot-toast]` — so a full-page-load
// flash rendered inline, an HTMX out-of-band swap
// (`hx-swap-oob="beforeend:#gsxui-toaster"`), an HTMX partial append, or an
// SSE-driven insert all animate/stack/auto-dismiss with ZERO HTMX-specific
// code here. This is the one-viewport-per-page server-flash model
// (docs/jsx-parity.md ## sonner): the server is the single source of toast
// markup and the observer is the single adoption path.
//
// Public imperative API (re-exported through the ui/index.js barrel):
//   import { toast } from "gsxui";
//   toast(msg, opts); toast.success/.info/.warning/.error/.loading(msg, opts);
//   toast.promise(promiseOrFn, { loading, success, error }); toast.dismiss(id?);
// opts: { description, duration, action: { label, onClick }, cancel: { label, onClick } }.
// Also reachable as window.gsxui.toast for inline <script> demo pages that
// cannot import the barrel.
import { on, emit } from "./gsxui.js";

// --- Stacking / timing constants -------------------------------------------
const DEFAULT_DURATION = 4000; // ms; a loading toast overrides this to Infinity
const MAX_VISIBLE = 3; // sonner's VISIBLE_TOASTS_AMOUNT — 4th+ is queued
const EXPAND_GAP = 14; // px between toasts when the stack is hover-expanded
const COLLAPSE_PEEK = 16; // px each stacked toast peeks above the front when collapsed
const SCALE_STEP = 0.05; // scale reduction per stack level when collapsed
const REMOVE_CAP = 600; // ms fallback if transitionend never fires
const HOVER_LEAVE_MS = 80; // debounce so crossing a gap between toasts doesn't collapse
// Closed (enter/exit) visual state — the discrete-transition start/end point
// the CSS transition on the <li> animates to/from (docs/jsx-parity.md
// ## animations, adapted to a template-cloned node: set closed state, force a
// frame, flip to open; exit reverses it). Bottom position → slide up on
// enter, back down on exit.
const CLOSED_TRANSFORM = "translateY(20px) scale(0.9)";

// --- State -----------------------------------------------------------------
// A plain array of toast records (oldest first, newest last = the front),
// NOT sonner's CSS-custom-property machine — we recompute the
// scale/translate stack per toast via inline style.
// The record map is the ownership boundary for every wired card, including a
// card in its exit transition. Region replacement explicitly disposes these
// records, which is what lets timers, listeners, and detached targets be
// reconciled instead of leaking behind a WeakSet marker.
const toasts = []; // { id, el, type, duration, remaining, timer, startedAt, onAction, onCancel }
const records = new Map();
let uid = 0;
let expanded = false; // hover-expands the stack AND pauses every timer
let leaveTimer = null;
let regionEl = null;
let fallbackSectionEl = null;
let fallbackRegionEl = null;

const regionObserver =
  typeof MutationObserver !== "undefined"
    ? new MutationObserver(handleRegionMutations)
    : null;
const documentObserver =
  typeof MutationObserver !== "undefined"
    ? new MutationObserver(handleDocumentMutations)
    : null;

// --- Toaster region --------------------------------------------------------
// One document observer discovers region replacement. One region observer is
// rebound to exactly one current viewport and owns row adoption. A fallback
// is tracked by object identity, so a later server viewport can replace it
// without adding a private DOM marker or guessing from markup.
function createFallbackRegion() {
  fallbackSectionEl = document.createElement("section");
  fallbackSectionEl.setAttribute("aria-label", "Notifications");
  fallbackSectionEl.tabIndex = -1;
  fallbackRegionEl = document.createElement("ol");
  fallbackRegionEl.setAttribute("data-gsxui-slot-toaster", "");
  fallbackRegionEl.id = "gsxui-toaster";
  fallbackSectionEl.appendChild(fallbackRegionEl);
  document.body.appendChild(fallbackSectionEl);
  return fallbackRegionEl;
}

function desiredRegion() {
  if (regionEl && regionEl.isConnected) {
    if (regionEl !== fallbackRegionEl) return regionEl;
    const serverRegion = [
      ...document.querySelectorAll("[data-gsxui-slot-toaster]"),
    ].find((candidate) => candidate !== fallbackRegionEl);
    return serverRegion || regionEl;
  }

  const serverRegion = [
    ...document.querySelectorAll("[data-gsxui-slot-toaster]"),
  ].find((candidate) => candidate !== fallbackRegionEl);
  if (serverRegion) return serverRegion;
  if (fallbackRegionEl && fallbackRegionEl.isConnected) {
    return fallbackRegionEl;
  }
  return createFallbackRegion();
}

function ensureRegion() {
  const next = desiredRegion();
  if (next !== regionEl) bindRegion(next);
  return regionEl;
}

function ol() {
  return ensureRegion();
}

function disposeRecord(rec, removeTarget) {
  pauseTimer(rec);
  if (rec.removeTimer) {
    clearTimeout(rec.removeTimer);
    rec.removeTimer = null;
  }
  rec.controller.abort();
  const index = toasts.indexOf(rec);
  if (index >= 0) toasts.splice(index, 1);
  records.delete(rec.el);
  if (removeTarget && rec.el.parentNode) rec.el.remove();
}

function reconcileRows(root) {
  if (root.nodeType !== 1) return;
  if (root.matches("li[data-gsxui-slot-toast]")) adopt(root);
  root.querySelectorAll("li[data-gsxui-slot-toast]").forEach((row) => adopt(row));
}

function bindRegion(next) {
  const previous = regionEl;
  regionObserver?.disconnect();

  if (leaveTimer) {
    clearTimeout(leaveTimer);
    leaveTimer = null;
  }
  expanded = false;

  // All records belonged to the previous viewport. If a row was physically
  // moved into the replacement before this callback, keep its DOM but release
  // its old ownership; reconciliation below will wire it once for the new
  // viewport. Every other stale target is removed with its record.
  for (const rec of [...records.values()]) {
    disposeRecord(rec, !next.contains(rec.el));
  }

  regionEl = next;
  if (previous === fallbackRegionEl && next !== fallbackRegionEl) {
    const fallbackSection = fallbackSectionEl;
    fallbackSectionEl = null;
    fallbackRegionEl = null;
    fallbackSection?.remove();
  }

  regionObserver?.observe(regionEl, { childList: true, subtree: true });
  reconcileRows(regionEl);
  refresh();
}

function handleRegionMutations(mutations) {
  for (const mutation of mutations) {
    for (const node of mutation.addedNodes) reconcileRows(node);
  }

  // A server may remove or move a row without dismissing it. Mutation records
  // are delivered after the whole DOM operation, so a row moved elsewhere
  // inside the current viewport remains owned while a genuinely detached row
  // is disposed.
  let disposed = false;
  for (const rec of [...records.values()]) {
    if (!regionEl?.contains(rec.el)) {
      disposeRecord(rec, false);
      disposed = true;
    }
  }
  if (disposed) refresh();
}

function handleDocumentMutations(mutations) {
  if (!regionEl?.isConnected) {
    ensureRegion();
    return;
  }
  if (regionEl !== fallbackRegionEl) return;

  for (const mutation of mutations) {
    for (const node of mutation.addedNodes) {
      if (
        node.nodeType === 1 &&
        (node.matches("[data-gsxui-slot-toaster]") ||
          node.querySelector("[data-gsxui-slot-toaster]"))
      ) {
        ensureRegion();
        return;
      }
    }
  }
}

// --- Template clone --------------------------------------------------------
// The server ships one inert <template data-gsxui-toast-template="TYPE"> per
// type inside ui.Toaster; each wraps a pre-rendered ui.Toast card. Cloning
// the matching template is how a card is created — the card markup (slots,
// behavior hooks, icons, aria) is authored ONCE, in the Go ui.Toast component.
function tpl(type) {
  return document.querySelector(
    `template[data-gsxui-toast-template="${type}"]`,
  );
}

// Replace el's icon slot with the target type's template icon (used by
// promise-morph). The loading→success/error cards have different glyphs; a
// default card has none.
function setIconFromTemplate(el, type) {
  const template = tpl(type);
  const srcIcon = template
    ? template.content.querySelector("[data-gsxui-slot-toast-icon]")
    : null;
  const slot = el.querySelector("[data-gsxui-slot-toast-icon]");
  if (!srcIcon) {
    if (slot) slot.remove();
    return;
  }
  const fresh = srcIcon.cloneNode(true);
  if (slot) slot.replaceWith(fresh);
  else el.insertBefore(fresh, el.firstChild);
}

// --- Toast construction (clone + fill) -------------------------------------
// Clone the matching type template's card, then fill or remove the
// title/description/action/cancel parts to match opts. Returns null when the
// template is missing (ui.Toaster not mounted) — show() then no-ops.
function build(rec, opts) {
  const template = tpl(rec.type) || tpl("default");
  if (!template) return null;
  const el = template.content.firstElementChild.cloneNode(true);
  el.dataset.type = rec.type;

  const title = el.querySelector("[data-gsxui-slot-toast-title]");
  if (title) title.textContent = opts.message ?? "";

  const desc = el.querySelector("[data-gsxui-slot-toast-description]");
  if (desc) {
    if (opts.description) desc.textContent = opts.description;
    else desc.remove();
  }

  const actionBtn = el.querySelector("[data-gsxui-slot-toast-action]");
  if (actionBtn) {
    if (opts.action && opts.action.label) {
      actionBtn.textContent = opts.action.label;
      rec.onAction = opts.action.onClick;
    } else {
      actionBtn.remove();
    }
  }

  const cancelBtn = el.querySelector("[data-gsxui-slot-toast-cancel]");
  if (cancelBtn) {
    if (opts.cancel && opts.cancel.label) {
      cancelBtn.textContent = opts.cancel.label;
      rec.onCancel = opts.cancel.onClick;
    } else {
      cancelBtn.remove();
    }
  }

  return el;
}

// Wire a card's interactive parts to a record. Applied to BOTH imperative
// cards (build) and server-adopted rows — the only difference is that
// server rows carry no JS onClick callbacks (action/cancel still emit +
// dismiss). Hover any toast → expand the whole stack + pause every timer;
// leave → collapse + resume (debounced so crossing a gap doesn't flicker).
function wire(el, rec) {
  const listenerOptions = { signal: rec.controller.signal };
  const close = el.querySelector("[data-gsxui-slot-toast-close]");
  if (close) {
    close.addEventListener("click", () => dismiss(rec.id), listenerOptions);
  }

  const actionBtn = el.querySelector("[data-gsxui-slot-toast-action]");
  if (actionBtn) {
    actionBtn.addEventListener(
      "click",
      () => {
        emit(el, "gsxui:toast-action", { id: rec.id });
        if (typeof rec.onAction === "function") rec.onAction();
        dismiss(rec.id);
      },
      listenerOptions,
    );
  }

  const cancelBtn = el.querySelector("[data-gsxui-slot-toast-cancel]");
  if (cancelBtn) {
    cancelBtn.addEventListener(
      "click",
      () => {
        if (typeof rec.onCancel === "function") rec.onCancel();
        dismiss(rec.id);
      },
      listenerOptions,
    );
  }

  el.addEventListener(
    "pointerenter",
    () => setExpanded(true),
    listenerOptions,
  );
  el.addEventListener(
    "pointerleave",
    () => setExpanded(false),
    listenerOptions,
  );
}

// Enter: stamp the closed visual state, force one frame so the transition has
// a start point, then let refresh() set the open transform — the <li>'s CSS
// transition animates the slide-up/fade-in (and shifts the rest of the stack
// in the same frame).
function enter(el) {
  el.dataset.state = "closed";
  el.style.opacity = "0";
  el.style.transform = CLOSED_TRANSFORM;
  void el.offsetHeight;
  el.dataset.state = "open";
  refresh();
}

// --- Stack layout ----------------------------------------------------------
// Recompute every visible toast's inline transform from its stack position.
// Collapsed: only the last MAX_VISIBLE show; the front (pos 0) sits at full
// scale, each older one peeks COLLAPSE_PEEK px above and shrinks SCALE_STEP
// per level (origin-bottom keeps bottoms aligned), z-index descending so the
// front paints on top. Expanded: the same toasts separate upward, each
// lifted by the running sum of the ones below it plus EXPAND_GAP, measured
// from live offsetHeight (heights vary with description/action). Toasts past
// MAX_VISIBLE are hidden (queued) until a dismiss promotes them.
function layout() {
  const n = toasts.length;
  let lift = 0; // running px offset for the expanded stack, from the front up
  for (let i = n - 1, pos = 0; i >= 0; i--, pos++) {
    const el = toasts[i].el;
    const z = 100 - pos;
    el.style.zIndex = String(z);
    if (pos >= MAX_VISIBLE) {
      el.dataset.visible = "false";
      el.style.opacity = "0";
      el.style.pointerEvents = "none";
      el.style.transform = `translateY(${-(MAX_VISIBLE - 1) * COLLAPSE_PEEK}px) scale(${1 - MAX_VISIBLE * SCALE_STEP})`;
      continue;
    }
    el.dataset.visible = "true";
    el.style.opacity = "1";
    el.style.pointerEvents = "auto";
    if (expanded) {
      el.style.transform = `translateY(${-lift}px) scale(1)`;
      lift += el.offsetHeight + EXPAND_GAP;
    } else {
      el.style.transform = `translateY(${-pos * COLLAPSE_PEEK}px) scale(${1 - pos * SCALE_STEP})`;
    }
  }
}

// --- Timers ----------------------------------------------------------------
// Only visible, non-loading toasts count down, and never while the stack is
// hover-expanded. remaining tracks time left across pause/resume.
function startTimer(rec) {
  if (rec.timer || rec.duration === Infinity || rec.remaining <= 0) return;
  rec.startedAt = Date.now();
  rec.timer = setTimeout(() => {
    rec.timer = null;
    dismiss(rec.id);
  }, rec.remaining);
}
function pauseTimer(rec) {
  if (!rec.timer) return;
  clearTimeout(rec.timer);
  rec.timer = null;
  rec.remaining -= Date.now() - rec.startedAt;
}
function syncTimers() {
  const n = toasts.length;
  toasts.forEach((rec, i) => {
    const visible = n - 1 - i < MAX_VISIBLE;
    if (expanded || !visible) pauseTimer(rec);
    else startTimer(rec);
  });
}

function refresh() {
  layout();
  syncTimers();
}

// --- Hover expand ----------------------------------------------------------
function setExpanded(on) {
  if (on) {
    if (leaveTimer) {
      clearTimeout(leaveTimer);
      leaveTimer = null;
    }
    if (expanded) return;
    expanded = true;
    ol().dataset.expanded = "true";
    refresh();
  } else {
    if (leaveTimer) return;
    leaveTimer = setTimeout(() => {
      leaveTimer = null;
      expanded = false;
      ol().dataset.expanded = "false";
      refresh();
    }, HOVER_LEAVE_MS);
  }
}

// --- Lifecycle -------------------------------------------------------------
function byId(id) {
  return toasts.find((r) => r.id === id);
}

function show(opts) {
  // Region binding reconciles and disposes ownership from the previous
  // viewport, so it must finish before this toast creates its record.
  const region = ol();
  const type = opts.type || "default";
  const rec = {
    id: opts.id ?? ++uid,
    el: null,
    type,
    duration:
      opts.duration ?? (type === "loading" ? Infinity : DEFAULT_DURATION),
    remaining: 0,
    timer: null,
    startedAt: 0,
    onAction: undefined,
    onCancel: undefined,
    controller: new AbortController(),
    removeTimer: null,
  };
  rec.remaining = rec.duration;
  rec.el = build(rec, opts);
  if (!rec.el) return rec.id; // ui.Toaster (and its templates) not mounted
  records.set(rec.el, rec); // own BEFORE insertion so the observer skips it
  toasts.push(rec);
  wire(rec.el, rec);
  region.appendChild(rec.el);
  enter(rec.el);
  return rec.id;
}

// Adopt a server-inserted (or otherwise externally-appended) toast row into
// the same lifecycle as an imperative one: assign an id, read data-type and
// optional data-duration (loading defaults to no auto-dismiss), wire the
// interactive parts, run the enter animation, and register it. Idempotent —
// a row already owned by a record is skipped (the imperative path records
// its own rows before insertion).
function adopt(el) {
  if (records.has(el)) return;
  const type = el.dataset.type || "default";
  let duration;
  const durAttr = el.dataset.duration;
  if (durAttr != null && durAttr !== "") {
    duration = Number(durAttr);
    if (Number.isNaN(duration)) duration = DEFAULT_DURATION;
  } else {
    duration = type === "loading" ? Infinity : DEFAULT_DURATION;
  }
  const rec = {
    id: ++uid,
    el,
    type,
    duration,
    remaining: duration,
    timer: null,
    startedAt: 0,
    onAction: undefined,
    onCancel: undefined,
    controller: new AbortController(),
    removeTimer: null,
  };
  records.set(el, rec);
  toasts.push(rec);
  wire(el, rec);
  enter(el);
  return rec.id;
}

function finalize(rec) {
  disposeRecord(rec, true);
  refresh();
}

function dismiss(id) {
  const idx = toasts.findIndex((r) => r.id === id);
  if (idx < 0) return;
  const rec = toasts[idx];
  pauseTimer(rec);
  toasts.splice(idx, 1);
  const el = rec.el;
  el.dataset.state = "closed";
  el.style.pointerEvents = "none";
  el.style.opacity = "0";
  el.style.transform = CLOSED_TRANSFORM;
  refresh(); // reflow the remaining stack (a queued toast may promote in)

  // Remove after the exit transition, capped so a backgrounded tab (frozen
  // transition clock) or a missing transitionend still cleans up — the same
  // hard-cap strategy, unlike dialog.js's finite-own-animation completion.
  rec.removeTimer = setTimeout(() => finalize(rec), REMOVE_CAP);
  el.addEventListener(
    "transitionend",
    (e) => {
      if (e.target !== el) return;
      finalize(rec);
    },
    { once: true, signal: rec.controller.signal },
  );
}

// Promise: render a loading toast, then MORPH THE SAME NODE in place on
// settle — swap data-type, replace the icon slot with the target type's
// template icon, swap the title, restart the dismiss timer. No re-animation
// and no stack reflow (a dismiss-old/spawn-new approach would visibly jump).
function morph(id, type, message) {
  const rec = byId(id);
  if (!rec) return;
  rec.type = type;
  rec.el.dataset.type = type;
  setIconFromTemplate(rec.el, type);
  const title = rec.el.querySelector("[data-gsxui-slot-toast-title]");
  if (title && message != null) title.textContent = message;
  clearTimeout(rec.timer);
  rec.timer = null;
  rec.duration = DEFAULT_DURATION;
  rec.remaining = DEFAULT_DURATION;
  syncTimers();
}

function resolveMsg(m, value) {
  return typeof m === "function" ? m(value) : m;
}

// --- Public API ------------------------------------------------------------
function toast(message, opts = {}) {
  return show({ ...opts, message, type: opts.type || "default" });
}
toast.success = (message, opts = {}) =>
  show({ ...opts, message, type: "success" });
toast.info = (message, opts = {}) => show({ ...opts, message, type: "info" });
toast.warning = (message, opts = {}) =>
  show({ ...opts, message, type: "warning" });
toast.error = (message, opts = {}) => show({ ...opts, message, type: "error" });
toast.loading = (message, opts = {}) =>
  show({
    ...opts,
    message,
    type: "loading",
    duration: opts.duration ?? Infinity,
  });
toast.promise = (promiseOrFn, msgs = {}) => {
  const id = show({
    message: resolveMsg(msgs.loading),
    type: "loading",
    duration: Infinity,
  });
  const p = typeof promiseOrFn === "function" ? promiseOrFn() : promiseOrFn;
  Promise.resolve(p).then(
    (value) => morph(id, "success", resolveMsg(msgs.success, value)),
    (error) => morph(id, "error", resolveMsg(msgs.error, error)),
  );
  return id;
};
toast.dismiss = (id) => {
  if (id == null) {
    for (const rec of [...toasts]) dismiss(rec.id);
  } else {
    dismiss(id);
  }
};

// --- Declarative trigger (zero-JS demo/doc pages) --------------------------
// Any element with a NON-EMPTY data-gsxui-toast fires a toast on click,
// reading the same fields the imperative API takes, so a docs page needs
// no page-specific <script>. Toast cards themselves no longer carry data-gsxui-toast (their
// identity is data-gsxui-slot-toast), but the empty-value guard stays as a
// defensive filter against stray valueless markers.
on("click", "[data-gsxui-toast]", (_event, el) => {
  if (!el.dataset.gsxuiToast) return;
  const label = el.dataset.gsxuiToastAction;
  show({
    message: el.dataset.gsxuiToast,
    description: el.dataset.gsxuiToastDescription,
    type: el.dataset.gsxuiToastType || "default",
    action: label ? { label } : undefined,
  });
});

// --- Server-flash adoption -------------------------------------------------
// The document observer detects viewport replacement without polling. Binding
// the viewport reconciles preexisting rows, while the subtree region observer
// adopts direct rows, rows inside inserted wrappers/fragments, and descendants
// added to an existing wrapper later. Explicit record ownership makes both
// paths idempotent.
// bootstrap, not the core init(): toast cards are cloned from <template>s,
// not server DOM matching a stable selector, so the selector-based
// self-healing model does not fit.
function bootstrap() {
  documentObserver?.observe(document.documentElement, {
    childList: true,
    subtree: true,
  });
  ensureRegion();
}

if (typeof document !== "undefined") {
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", bootstrap);
  } else {
    bootstrap();
  }
}

// Barrel re-export makes `import { toast } from "gsxui"` work; the window
// global covers inline <script> demo pages that cannot import the barrel.
if (typeof window !== "undefined") {
  window.gsxui = Object.assign(window.gsxui ?? {}, { toast });
}

export { toast };
