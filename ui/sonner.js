// Sonner (toasts). Unlike every other gsxui behavior module (which attaches
// delegated behavior to server-rendered markup), sonner's card markup is
// still authored server-side — as the ui.Toast component — but shipped as
// inert per-type <template>s inside ui.Toaster. This module CLONES the
// matching type's template on each toast() call (never builds card DOM from
// JS string concatenation), then owns the whole lifecycle: mount → stack →
// timer → dismiss.
//
// Because the card is server markup, the same lifecycle is applied to rows
// the SERVER inserts, not just JS-triggered ones. A MutationObserver on the
// <ol> adopts any inserted `li[data-slot="toast"]` — so a full-page-load
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
const REMOVE_CAP = 600; // ms fallback if transitionend never fires (dialog.js's cap idea)
const HOVER_LEAVE_MS = 80; // debounce so crossing a gap between toasts doesn't collapse
// Closed (enter/exit) visual state — the discrete-transition start/end point
// the CSS transition on the <li> animates to/from (docs/jsx-parity.md
// ## animations, adapted to a template-cloned node: set closed state, force a
// frame, flip to open; exit reverses it). Bottom position → slide up on
// enter, back down on exit.
const CLOSED_TRANSFORM = "translateY(20px) scale(0.9)";

// --- State -----------------------------------------------------------------
// A plain array of toast records (oldest first, newest last = the front),
// NOT sonner's CSS-custom-property machine — we ship fixed Tailwind classes
// and recompute the scale/translate stack per toast via inline style. A
// WeakSet marks every <li> already owned by a record, so the observer never
// double-adopts a row the imperative API just inserted.
const toasts = []; // { id, el, type, duration, remaining, timer, startedAt, onAction, onCancel }
const registered = new WeakSet();
let uid = 0;
let expanded = false; // hover-expands the stack AND pauses every timer
let leaveTimer = null;

// --- Toaster region --------------------------------------------------------
// Uses ui.Toaster's server-rendered region if present; otherwise builds a
// fallback so the region always exists (the imperative API still needs the
// per-type <template>s ui.Toaster ships — a page firing toast() must mount
// <ui.Toaster/>).
let olEl = null;
function ol() {
  if (olEl && olEl.isConnected) return olEl;
  olEl = document.querySelector("[data-gsxui-toaster]");
  if (!olEl) {
    const section = document.createElement("section");
    section.setAttribute("aria-label", "Notifications");
    section.tabIndex = -1;
    olEl = document.createElement("ol");
    olEl.dataset.slot = "toaster";
    olEl.setAttribute("data-gsxui-toaster", "");
    olEl.id = "gsxui-toaster";
    olEl.className =
      "pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0";
    section.appendChild(olEl);
    document.body.appendChild(section);
  }
  return olEl;
}

// --- Template clone --------------------------------------------------------
// The server ships one inert <template data-gsxui-toast-template="TYPE"> per
// type inside ui.Toaster; each wraps a pre-rendered ui.Toast card. Cloning
// the matching template is how a card is created — the card markup (classes,
// icons, aria) is authored ONCE, in the Go ui.Toast component.
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
    ? template.content.querySelector("[data-icon]")
    : null;
  const slot = el.querySelector("[data-icon]");
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

  const title = el.querySelector("[data-title]");
  if (title) title.textContent = opts.message ?? "";

  const desc = el.querySelector("[data-description]");
  if (desc) {
    if (opts.description) desc.textContent = opts.description;
    else desc.remove();
  }

  const actionBtn = el.querySelector("[data-action]");
  if (actionBtn) {
    if (opts.action && opts.action.label) {
      actionBtn.textContent = opts.action.label;
      rec.onAction = opts.action.onClick;
    } else {
      actionBtn.remove();
    }
  }

  const cancelBtn = el.querySelector("[data-cancel]");
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
  const close = el.querySelector("[data-close-button]");
  if (close) close.addEventListener("click", () => dismiss(rec.id));

  const actionBtn = el.querySelector("[data-action]");
  if (actionBtn) {
    actionBtn.addEventListener("click", () => {
      emit(el, "gsxui:toast-action", { id: rec.id });
      if (typeof rec.onAction === "function") rec.onAction();
      dismiss(rec.id);
    });
  }

  const cancelBtn = el.querySelector("[data-cancel]");
  if (cancelBtn) {
    cancelBtn.addEventListener("click", () => {
      if (typeof rec.onCancel === "function") rec.onCancel();
      dismiss(rec.id);
    });
  }

  el.addEventListener("pointerenter", () => setExpanded(true));
  el.addEventListener("pointerleave", () => setExpanded(false));
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
  rec.timer = setTimeout(() => dismiss(rec.id), rec.remaining);
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
  };
  rec.remaining = rec.duration;
  rec.el = build(rec, opts);
  if (!rec.el) return rec.id; // ui.Toaster (and its templates) not mounted
  registered.add(rec.el); // mark BEFORE insertion so the observer skips it
  toasts.push(rec);
  wire(rec.el, rec);
  ol().appendChild(rec.el);
  enter(rec.el);
  return rec.id;
}

// Adopt a server-inserted (or otherwise externally-appended) toast row into
// the same lifecycle as an imperative one: assign an id, read data-type and
// optional data-duration (loading defaults to no auto-dismiss), wire the
// interactive parts, run the enter animation, and register it. Idempotent —
// a row already owned by a record is skipped (the imperative path marks its
// own rows before insertion).
function adopt(el) {
  if (registered.has(el)) return;
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
  };
  registered.add(el);
  toasts.push(rec);
  wire(el, rec);
  enter(el);
  return rec.id;
}

function finalize(el) {
  if (el.parentNode) el.remove();
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
  // race-against-a-hard-cap idea as dialog.js's requestClose.
  const cap = setTimeout(() => finalize(el), REMOVE_CAP);
  el.addEventListener(
    "transitionend",
    (e) => {
      if (e.target !== el) return;
      clearTimeout(cap);
      finalize(el);
    },
    { once: true },
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
  const title = rec.el.querySelector("[data-title]");
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
// reading the same fields the imperative API takes — mirrors the
// data-gsxui-dialog-trigger idiom so a docs page needs no page-specific
// <script>. The cloned toast <li> also carries data-gsxui-toast (an empty
// slot marker, per the Toast card markup), so the empty-value guard is what
// stops a click INSIDE a toast from spawning a blank one.
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
// A MutationObserver on the <ol> adopts any externally-inserted toast row
// (server full-load flash, HTMX OOB/partial append, SSE insert). Rows the
// imperative API inserts are pre-marked in `registered`, so they are never
// double-adopted. Rows already present at module init (full-page-load
// flashes drained server-side into the <ol>) are adopted once at startup.
const observer =
  typeof MutationObserver !== "undefined"
    ? new MutationObserver((mutations) => {
        for (const m of mutations) {
          for (const node of m.addedNodes) {
            if (
              node.nodeType === 1 &&
              node.matches &&
              node.matches('li[data-slot="toast"]') &&
              !registered.has(node)
            ) {
              adopt(node);
            }
          }
        }
      })
    : null;

function init() {
  const region = ol();
  if (observer) observer.observe(region, { childList: true });
  region.querySelectorAll('li[data-slot="toast"]').forEach((el) => {
    if (!registered.has(el)) adopt(el);
  });
}

if (typeof document !== "undefined") {
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
}

// Barrel re-export makes `import { toast } from "gsxui"` work; the window
// global covers inline <script> demo pages that cannot import the barrel.
if (typeof window !== "undefined") {
  window.gsxui = Object.assign(window.gsxui ?? {}, { toast });
}

export { toast };
