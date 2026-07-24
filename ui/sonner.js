// Sonner (toasts) behavior — the codebase's first CLIENT-CONSTRUCTED-DOM
// module. Every other gsxui behavior (dialog/dropdown/command/…) attaches
// delegated behavior to server-rendered markup; a toast has no server
// markup to attach to (a toast is definitionally a client-triggered
// response to some JS event), so this module BUILDS each toast <li> from
// scratch and appends it into ui.Toaster's one static <ol data-gsxui-toaster>
// region, then owns its whole lifecycle: mount → stack → timer → dismiss.
//
// There is no Tailwind source upstream to port (sonner ships a non-Tailwind
// stylesheet), so the toast card's classes below are the synthesized spec
// from docs/superpowers/plans/2026-07-24-tier3-source-map-wrapped.md
// `## sonner`, reconstructed to match our popover/card surfaces.
//
// Public imperative API (first in this codebase to be re-exported through
// the ui/index.js barrel for page authors, not just sibling modules):
//   import { toast } from "gsxui";
//   toast(msg, opts); toast.success/.info/.warning/.error/.loading(msg, opts);
//   toast.promise(promiseOrFn, { loading, success, error }); toast.dismiss(id?);
// opts: { description, duration, action: { label, onClick }, cancel: { label, onClick } }.
// Also reachable as window.gsxui.toast for inline <script> demo pages that
// cannot import the barrel.
import { on, emit } from "./gsxui.js";

// --- Icon glyphs -----------------------------------------------------------
// Hand-copied SVG path data from ui/icon/icon_data.go (Lucide, ISC) — the
// toast <li> is built in JS, so icon.CircleCheck (a server-side Go call) is
// unreachable here; these strings are the same "static data ported into a
// JS module" precedent as command.js's commandScore port. MAINTENANCE SEAM:
// if ui/icon's glyphs are ever regenerated from a newer Lucide, these copies
// do NOT update automatically (ledgered in docs/jsx-parity.md ## sonner).
// Provenance (icon_data.go keys): circle-check / info / triangle-alert /
// octagon-x / loader-circle / x.
const GLYPHS = {
  success: '<circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/>',
  info: '<circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>',
  warning:
    '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/><path d="M12 9v4"/><path d="M12 17h.01"/>',
  error:
    '<path d="m15 9-6 6"/><path d="M2.586 16.726A2 2 0 0 1 2 15.312V8.688a2 2 0 0 1 .586-1.414l4.688-4.688A2 2 0 0 1 8.688 2h6.624a2 2 0 0 1 1.414.586l4.688 4.688A2 2 0 0 1 22 8.688v6.624a2 2 0 0 1-.586 1.414l-4.688 4.688a2 2 0 0 1-1.414.586H8.688a2 2 0 0 1-1.414-.586z"/><path d="m9 9 6 6"/>',
  loading: '<path d="M21 12a9 9 0 1 1-6.219-8.56"/>',
};
const X_GLYPH = '<path d="M18 6 6 18"/><path d="m6 6 12 12"/>';

// Lucide's shared <svg> wrapper attributes (matches ui/icon's own output and
// the theme-toggle glyph in site/pages/layout.gsx).
function svg(paths, cls) {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"${cls ? ` class="${cls}"` : ""}>${paths}</svg>`;
}

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
// ## animations, adapted to a JS-constructed node: set closed state, force a
// frame, flip to open; exit reverses it). Bottom position → slide up on
// enter, back down on exit.
const CLOSED_TRANSFORM = "translateY(20px) scale(0.9)";

// --- State -----------------------------------------------------------------
// A plain array of toast records (oldest first, newest last = the front),
// NOT sonner's CSS-custom-property machine — we ship fixed Tailwind classes
// and recompute the scale/translate stack per toast via inline style.
const toasts = []; // { id, el, type, duration, remaining, timer, startedAt }
let uid = 0;
let expanded = false; // hover-expands the stack AND pauses every timer
let leaveTimer = null;

// --- Toaster region --------------------------------------------------------
// Uses ui.Toaster's server-rendered region if present; otherwise builds a
// fallback so toast() works on any page (matches Toaster's classes).
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
    olEl.className =
      "pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0";
    section.appendChild(olEl);
    document.body.appendChild(section);
  }
  return olEl;
}

// --- Toast construction ----------------------------------------------------
const TOAST_CLASS =
  // Card surface — synthesized spec (see header). rounded-2xl / bg-popover /
  // w-[356px] verbatim from the map; `relative` → `absolute bottom-6 right-6`
  // + a transform transition are gsxui's stacking additions (the map's card
  // string has no positioning of its own — sonner's own stylesheet owned it).
  "pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out " +
  "data-[type=success]:[&>[data-icon]]:text-emerald-500 data-[type=info]:[&>[data-icon]]:text-sky-500 data-[type=warning]:[&>[data-icon]]:text-amber-500 data-[type=error]:[&>[data-icon]]:text-destructive";

function iconHTML(type) {
  const glyph = GLYPHS[type];
  if (!glyph) return ""; // default type has no icon (sonner's own default)
  return svg(glyph, type === "loading" ? "animate-spin" : "");
}

function setIcon(el, type) {
  let slot = el.querySelector("[data-icon]");
  const html = iconHTML(type);
  if (!html) {
    if (slot) slot.remove();
    return;
  }
  if (!slot) {
    slot = document.createElement("div");
    slot.setAttribute("data-icon", "");
    slot.className = "mt-0.5 shrink-0 [&>svg]:size-4";
    el.insertBefore(slot, el.firstChild);
  }
  slot.innerHTML = html;
}

function build(rec, opts) {
  const el = document.createElement("li");
  el.className = TOAST_CLASS;
  el.dataset.slot = "toast";
  el.setAttribute("data-gsxui-toast", "");
  el.dataset.type = rec.type;
  el.setAttribute("role", "status");
  el.setAttribute("aria-live", rec.type === "error" ? "assertive" : "polite");
  el.setAttribute("aria-atomic", "true");

  setIcon(el, rec.type);

  const content = document.createElement("div");
  content.setAttribute("data-content", "");
  content.className = "flex flex-1 flex-col gap-1";
  const title = document.createElement("div");
  title.setAttribute("data-title", "");
  title.className = "font-medium text-foreground";
  title.textContent = opts.message ?? "";
  content.appendChild(title);
  if (opts.description) {
    const desc = document.createElement("div");
    desc.setAttribute("data-description", "");
    desc.className = "text-muted-foreground";
    desc.textContent = opts.description;
    content.appendChild(desc);
  }
  el.appendChild(content);

  if (opts.action && opts.action.label) {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.setAttribute("data-action", "");
    btn.className =
      "shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline";
    btn.textContent = opts.action.label;
    btn.addEventListener("click", () => {
      emit(el, "gsxui:toast-action", { id: rec.id });
      if (typeof opts.action.onClick === "function") opts.action.onClick();
      dismiss(rec.id);
    });
    el.appendChild(btn);
  }
  if (opts.cancel && opts.cancel.label) {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.setAttribute("data-cancel", "");
    btn.className =
      "shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline";
    btn.textContent = opts.cancel.label;
    btn.addEventListener("click", () => {
      if (typeof opts.cancel.onClick === "function") opts.cancel.onClick();
      dismiss(rec.id);
    });
    el.appendChild(btn);
  }

  const close = document.createElement("button");
  close.type = "button";
  close.setAttribute("data-close-button", "");
  close.setAttribute("aria-label", "Close");
  close.className =
    "absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm";
  close.innerHTML = svg(X_GLYPH, "size-3");
  close.addEventListener("click", () => dismiss(rec.id));
  el.appendChild(close);

  // Hover the toast → expand the whole stack + pause every timer; leave →
  // collapse + resume (debounced so crossing a gap doesn't flicker).
  el.addEventListener("pointerenter", () => setExpanded(true));
  el.addEventListener("pointerleave", () => setExpanded(false));

  return el;
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
  };
  rec.remaining = rec.duration;
  rec.el = build(rec, opts);
  toasts.push(rec);
  ol().appendChild(rec.el);

  // Enter: stamp the closed visual state, force one frame so the transition
  // has a start point, then let refresh() set the open transform — the <li>'s
  // CSS transition animates the slide-up/fade-in (and shifts the rest of the
  // stack in the same frame).
  rec.el.dataset.state = "closed";
  rec.el.style.opacity = "0";
  rec.el.style.transform = CLOSED_TRANSFORM;
  void rec.el.offsetHeight;
  rec.el.dataset.state = "open";
  refresh();
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
// settle — swap icon/type/title, restart the dismiss timer, no re-animation
// and no stack reflow (a dismiss-old/spawn-new approach would visibly jump).
function morph(id, type, message) {
  const rec = byId(id);
  if (!rec) return;
  rec.type = type;
  rec.el.dataset.type = type;
  setIcon(rec.el, type);
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
  show({ ...opts, message, type: "loading", duration: opts.duration ?? Infinity });
toast.promise = (promiseOrFn, msgs = {}) => {
  const id = show({ message: resolveMsg(msgs.loading), type: "loading", duration: Infinity });
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
// <script>. The constructed toast <li> also carries data-gsxui-toast (an
// empty slot marker, per the map's markup), so the empty-value guard is what
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

// Barrel re-export makes `import { toast } from "gsxui"` work; the window
// global covers inline <script> demo pages that cannot import the barrel.
if (typeof window !== "undefined") {
  window.gsxui = Object.assign(window.gsxui ?? {}, { toast });
}

export { toast };
