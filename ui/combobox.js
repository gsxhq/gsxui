// Filtered single-select listbox on ComboboxInput's text box. Open/close and
// the value model/form bridge are ui/select.js's own machinery restated for
// an input instead of a button trigger — popover="auto" (top layer, light
// dismiss, free Esc), position:fixed anchoring below the input, the
// data-state="open"-BEFORE-showPopover() flash fix, and a hidden bridge kept
// in sync with a bubbling "change". What's genuinely new, modeled on
// command.js's interaction loop (see its own header comment — FOCUS STAYS IN
// THE INPUT the whole time; items are never tab stops, never receive real
// DOM focus): a data-highlighted + aria-activedescendant highlight cursor,
// and an `input`-event filter.
//
// FIX (review round 1, CRITICAL, live-browser-traced): the popover used to
// open on the input's "focus" event. `popover="auto"` light-dismisses on
// the pointer gesture that is ALREADY IN FLIGHT when focus fires mid-click
// (pointerdown -> focus -> the browser's own light-dismiss evaluation ->
// mouseup -> click, all on the SAME gesture) — opening from focus opens the
// popover mid-gesture, so the completion of that same gesture immediately
// closes it again: click the input, watch the list flash open and vanish.
// The trigger chevron never had this bug because it opens on "click", which
// lands AFTER light-dismiss has already been evaluated for that gesture —
// the fix is to open the input the same way: on "click", not "focus". This
// also matches Base UI's own real behavior (opens on click/typing, not
// focus) and stops Tab-through from popping the list open on every combobox
// it passes.
//
// FILTER (ADAPT, web-verified, not read from @base-ui/react's own source —
// see ui/combobox.gsx's package doc comment): Base UI's default Combobox
// filter is `useFilter().contains`, an Intl.Collator-backed BOOLEAN match at
// sensitivity: "base" — case- and accent-insensitive, no ranking, no
// reordering. This is a real divergence from ui/command.js's commandScore
// fuzzy-ranking engine (## command): items here are only ever hidden/shown,
// never reordered in the DOM. Duplicated here rather than imported from
// command.js — a JS-only import would be invisible to
// internal/registry.Deps' go/parser scan over the generated .x.go, silently
// breaking CLI vendoring (see the .gsx file's own GAP entry).
import { on, emit } from "./gsxui.js";

const rootOf = (el) => el.closest("[data-gsxui-combobox]");
const inputOf = (root) => root?.querySelector("[data-gsxui-combobox-input]") ?? null;
const contentOf = (root) => root?.querySelector("[data-gsxui-combobox-content]") ?? null;
const listOf = (root) => root?.querySelector("[data-gsxui-combobox-list]") ?? null;
const bridgeOf = (root) => root?.querySelector("[data-gsxui-combobox-bridge]") ?? null;
const itemsOf = (root) => (root ? [...root.querySelectorAll("[data-gsxui-combobox-item]")] : []);
const isDisabled = (item) =>
  item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset;
const labelOf = (item) => item.textContent.trim();
const visibleItems = (root) => itemsOf(root).filter((i) => !i.hidden && !isDisabled(i));
const highlightedOf = (root) =>
  root?.querySelector('[data-gsxui-combobox-item][data-highlighted="true"]') ?? null;

let uid = 0;

// --- filter: Intl.Collator-backed boolean `contains` (Base UI's useFilter
// default) — the exact algorithm traced from Base UI's docs/issue tracker,
// mirroring React Aria's own implementation. See the file header ADAPT.
const collator = new Intl.Collator(undefined, { usage: "search", sensitivity: "base" });

function contains(string, substring) {
  if (substring.length === 0) return true;
  const s = string.normalize("NFC");
  const q = substring.normalize("NFC");
  for (let i = 0; i <= s.length - q.length; i++) {
    if (collator.compare(s.slice(i, i + q.length), q) === 0) return true;
  }
  return false;
}

// --- highlight cursor (NOT real DOM focus — command.js's model) ----------

// FIX (review round 1, CRITICAL): aria-activedescendant must never survive
// on a CLOSED combobox — clearValue()'s own filter() pass used to stamp it
// unconditionally, leaving the input pointing at an option inside a hidden
// popover while aria-expanded="false". highlight() now only stamps
// data-highlighted/aria-activedescendant while the popover is actually
// open; closed, it clears both regardless of what item was asked for.
function highlight(root, item) {
  const input = inputOf(root);
  const open = contentOf(root)?.matches(":popover-open") ?? false;
  for (const other of itemsOf(root)) {
    if (other !== item) delete other.dataset.highlighted;
  }
  if (!item || !open) {
    if (item) delete item.dataset.highlighted;
    input?.removeAttribute("aria-activedescendant");
    return;
  }
  item.dataset.highlighted = "true";
  if (!item.id) item.id = `gsxui-combobox-item-${++uid}`;
  input?.setAttribute("aria-activedescendant", item.id);
  item.scrollIntoView({ block: "nearest" });
}

// --- filter pass: hide/show only, no reordering, no ranking --------------

// Groups/separators hide (or become inert) while filtered; data-empty is
// stamped on BOTH the content (gates ComboboxEmpty via its own
// group-data-empty/combobox-content: selector) and the list (its own
// data-empty:p-0 selector, an independent target — see ui/combobox.gsx's
// ComboboxList doc comment).
//
// queryOverride lets openContent() force an EMPTY query on reopen when the
// input's current text came from a commit/init/reset rather than the user
// typing — see the gsxuiCommitted flag below (FIX, review round 1,
// CRITICAL: reopening used to always filter against input.value, so after
// picking "Remix" every other option stayed hidden on the next open —
// reopening a combobox that already has a value must show every option,
// not just the one whose label happens to be sitting in the input).
function filter(root, queryOverride) {
  const input = inputOf(root);
  const query = queryOverride !== undefined ? queryOverride : (input?.value ?? "").trim();
  let any = false;
  for (const item of itemsOf(root)) {
    const match = contains(labelOf(item), query);
    item.hidden = !match;
    if (match) any = true;
  }
  for (const group of root.querySelectorAll('[data-slot="combobox-group"]')) {
    const items = [...group.querySelectorAll("[data-gsxui-combobox-item]")];
    group.hidden = items.length > 0 && items.every((i) => i.hidden);
  }
  for (const sep of root.querySelectorAll('[data-slot="combobox-separator"]')) {
    sep.hidden = query !== "";
  }
  contentOf(root)?.toggleAttribute("data-empty", !any);
  listOf(root)?.toggleAttribute("data-empty", !any);
  // Prefer the checked (picked) item if it's still visible, else the first
  // visible one — the same priority ## select's own open-time focus uses.
  const visible = visibleItems(root);
  const checked = visible.find((i) => i.dataset.state === "checked");
  highlight(root, checked ?? visible[0] ?? null);
}

// --- value model: commit (pick) / clear -----------------------------------

// FIX (review round 1, IMPORTANT): mouse commit used to leave real DOM
// focus wherever the click landed (usually document.body, since items are
// never tab stops) — subsequent typing/arrowing was dead. Restore focus to
// the input on every commit, the same as ui/command.js:255's own activate()
// does for its selection.
function commit(root, item) {
  if (!item || isDisabled(item)) return;
  for (const other of itemsOf(root)) {
    const isIt = other === item;
    other.dataset.state = isIt ? "checked" : "unchecked";
    other.setAttribute("aria-selected", isIt ? "true" : "false");
  }
  const value = item.dataset.value ?? "";
  const input = inputOf(root);
  if (input) {
    input.value = labelOf(item);
    // This text came from US, not the user typing — the next reopen must
    // show every option, not just this one (see filter()'s own header).
    input.dataset.gsxuiCommitted = "true";
  }
  const bridge = bridgeOf(root);
  if (bridge) {
    bridge.value = value;
    bridge.dispatchEvent(new Event("change", { bubbles: true }));
  }
  emit(root, "gsxui:select", { value });
  contentOf(root)?.hidePopover();
  input?.focus();
}

function clearValue(root) {
  for (const item of itemsOf(root)) {
    item.dataset.state = "unchecked";
    item.setAttribute("aria-selected", "false");
  }
  const input = inputOf(root);
  if (input) {
    input.value = "";
    delete input.dataset.gsxuiCommitted;
  }
  const bridge = bridgeOf(root);
  if (bridge) {
    bridge.value = "";
    bridge.dispatchEvent(new Event("change", { bubbles: true }));
  }
  emit(root, "gsxui:select", { value: "" });
  filter(root);
  input?.focus();
}

// FIX (review round 1, IMPORTANT): a native <form>.reset() reverts the
// hidden bridge's .value to its content-attribute default (Combobox's own
// server-rendered `value` param) automatically — but the browser has no
// idea the VISIBLE input and every item's data-state/aria-selected are
// separate elements that need to follow it. Re-derive them from the
// bridge's own (already-reset, by the time this runs) value. Chosen over
// keeping a second copy of the value in the bridge's `value` ATTRIBUTE in
// parallel with the property: this rides the browser's own native reset
// semantics instead of duplicating them.
function reflectFromBridge(root) {
  const bridge = bridgeOf(root);
  const input = inputOf(root);
  const value = bridge?.value ?? "";
  let matched = null;
  for (const item of itemsOf(root)) {
    const isIt = value !== "" && (item.dataset.value ?? "") === value;
    item.dataset.state = isIt ? "checked" : "unchecked";
    item.setAttribute("aria-selected", isIt ? "true" : "false");
    if (isIt) matched = item;
  }
  if (input) {
    input.value = matched ? labelOf(matched) : "";
    if (matched) input.dataset.gsxuiCommitted = "true";
    else delete input.dataset.gsxuiCommitted;
  }
}

// Scoped to a form that actually contains a combobox — a bare "form" would
// otherwise claim every <form> on every page for "reset:false", colliding
// with any OTHER module that also wants a reset hook on its own component
// (ui/calendar.js's own reset handler, added in Task 5, hit exactly this:
// jstest/specs/invariants.spec.ts's Invariant 4 flagged calendar/form.gsx's
// plain <form> — no combobox in sight — as "claimed by two modules for one
// event" the moment both modules used unscoped "form"). Each handler is a
// no-op on a form that doesn't contain its own component regardless (both
// scope internally via `form.querySelectorAll(...)`), so the collision was
// harmless in practice, but scoping here — the same fix calendar.js's own
// reset handler already uses — is what keeps every future form-bridge
// module from having to relitigate the same overlap.
on("reset", "form:has([data-gsxui-combobox])", (_e, form) => {
  for (const root of form.querySelectorAll("[data-gsxui-combobox]")) {
    reflectFromBridge(root);
  }
});

// --- init: group aria-labelledby wiring, reflect a server-checked item,
// wire the permanent aria-controls (APG expects it present regardless of
// open/closed state — see ui/combobox.gsx's ComboboxInput doc comment) ----

function init(root) {
  for (const group of root.querySelectorAll('[data-slot="combobox-group"]')) {
    if (group.getAttribute("aria-labelledby")) continue;
    const label = group.querySelector('[data-slot="combobox-label"]');
    if (!label) continue;
    if (!label.id) label.id = `gsxui-combobox-label-${++uid}`;
    group.setAttribute("aria-labelledby", label.id);
  }
  const input = inputOf(root);
  const list = listOf(root);
  if (input && list) {
    if (!list.id) list.id = `gsxui-combobox-list-${++uid}`;
    input.setAttribute("aria-controls", list.id);
  }
  // FIX (review round 2, IMPORTANT): a caller is now expected to
  // server-render the checked item's own label as ComboboxInput's `value`
  // directly (docs/superpowers/specs/2026-07-24-tier4-batch-a-design.md
  // §4 — state is a server-rendered parameter reflected in the DOM, not
  // JS-seeded) — this used to be the ONLY place that label ever reached
  // the input, so it wrote it unconditionally whenever the input was
  // still empty. It stays as a fallback for a caller who doesn't supply
  // `value` (input still empty), but either way the gsxuiCommitted flag
  // must be set once the input's text equals the checked item's label —
  // server-rendered or JS-seeded, an unflagged first reopen re-triggers
  // the exact "reopening filtered to the committed label" bug this flag
  // exists to prevent (see docs/jsx-parity.md `## combobox`'s own FIX
  // entry): every option except this one would stay hidden.
  const checked = root.querySelector('[data-gsxui-combobox-item][data-state="checked"]');
  if (checked && input) {
    const label = labelOf(checked);
    if (!input.value) input.value = label;
    if (input.value === label) input.dataset.gsxuiCommitted = "true";
  }
}

for (const root of document.querySelectorAll("[data-gsxui-combobox]")) init(root);

// --- open / close (ported dropdown.js/select.js machinery) ---------------

function openContent(root) {
  const content = contentOf(root);
  const input = inputOf(root);
  if (!content || content.matches(":popover-open")) return;
  const anchor = input?.closest('[data-slot="input-group"]') ?? input;
  if (anchor) {
    const r = anchor.getBoundingClientRect();
    content.style.position = "fixed";
    content.style.inset = "auto";
    content.style.left = `${r.left}px`;
    content.style.top = `${r.bottom + 4}px`;
    content.style.minWidth = `${r.width}px`;
  }
  // Stamp open BEFORE showing — the toggle event that also stamps it is a
  // separate queued task; a paint can land in the gap and flash the closed
  // state (same fix select.js/dropdown.js document).
  content.dataset.state = "open";
  content.showPopover();
  // Reopening after a commit/init/reset shows every option, not just the
  // one whose label happens to be sitting in the input — see filter()'s and
  // commit()'s own headers (FIX, review round 1, CRITICAL).
  filter(root, input?.dataset.gsxuiCommitted ? "" : undefined);
}

function closeContent(root) {
  contentOf(root)?.hidePopover();
}

// aria-controls is now a PERMANENT attribute wired once at init() (APG
// expects it present regardless of open/closed state — see
// ui/combobox.gsx's own ComboboxInput doc comment); this handler only
// tracks open/closed state and clears the highlight cursor on close.
on("toggle", "[data-gsxui-combobox-content]", (e, content) => {
  const open = e.newState === "open";
  content.dataset.state = open ? "open" : "closed";
  const root = rootOf(content);
  const input = inputOf(root);
  input?.setAttribute("aria-expanded", open ? "true" : "false");
  if (open) {
    emit(content, "gsxui:open");
  } else {
    input?.removeAttribute("aria-activedescendant");
    for (const item of itemsOf(root)) delete item.dataset.highlighted;
    emit(content, "gsxui:close");
  }
}, { capture: true });

// contextmenu inside the listbox is suppressed (dropdown.js/select.js do the
// same).
on("contextmenu", "[data-gsxui-combobox-content]", (e) => e.preventDefault());

// --- trigger button --------------------------------------------------------

on("pointerdown", "[data-gsxui-combobox-trigger]", (_e, trigger) => {
  const content = contentOf(rootOf(trigger));
  if (content) {
    trigger.dataset.gsxuiWasOpen = content.matches(":popover-open") ? "true" : "false";
  }
});

on("click", "[data-gsxui-combobox-trigger]", (_e, trigger) => {
  const root = rootOf(trigger);
  const content = contentOf(root);
  if (!content) return;
  const wasOpen = trigger.dataset.gsxuiWasOpen === "true";
  delete trigger.dataset.gsxuiWasOpen;
  if (wasOpen || content.matches(":popover-open")) {
    if (content.matches(":popover-open")) closeContent(root);
  } else {
    openContent(root);
  }
  inputOf(root)?.focus();
});

// --- clear button ------------------------------------------------------

on("click", "[data-gsxui-combobox-clear]", (_e, clear) => {
  clearValue(rootOf(clear));
});

// --- the input: opens on click, filters on input, keyboard nav -----------

// FIX (review round 1, CRITICAL — see the file header for the full trace):
// opening from "focus" opened the popover MID-GESTURE, so the pointer
// gesture that focused the input also light-dismissed the popover it just
// opened — click the input, watch the list flash open and vanish. "click"
// fires after light-dismiss has already been evaluated for that gesture,
// the same reason the trigger chevron (which always opened on "click")
// never had this bug.
on("click", "[data-gsxui-combobox-input]", (_e, input) => {
  openContent(rootOf(input));
});

on("input", "[data-gsxui-combobox-input]", (_e, input) => {
  // The user is typing now — any earlier commit/init/reset-seeded text no
  // longer applies; a real, query-driven filter pass is always correct
  // from here (see filter()'s and openContent()'s own headers).
  delete input.dataset.gsxuiCommitted;
  const root = rootOf(input);
  const content = contentOf(root);
  if (content && !content.matches(":popover-open")) {
    openContent(root);
  } else {
    filter(root);
  }
});

on("keydown", "[data-gsxui-combobox-input]", (e, input) => {
  const root = rootOf(input);
  const content = contentOf(root);
  const isOpen = content?.matches(":popover-open") ?? false;

  if (e.key === "Escape") {
    if (isOpen) {
      e.preventDefault();
      closeContent(root);
    }
    return;
  }
  if (e.key === "Enter") {
    if (isOpen) {
      e.preventDefault();
      commit(root, highlightedOf(root));
    }
    return;
  }
  if (e.key === "ArrowDown" || e.key === "ArrowUp") {
    e.preventDefault();
    if (!isOpen) {
      openContent(root);
      return;
    }
    const items = visibleItems(root);
    if (!items.length) return;
    const dir = e.key === "ArrowDown" ? 1 : -1;
    const cur = items.indexOf(highlightedOf(root));
    const next =
      cur === -1
        ? dir === 1
          ? items[0]
          : items[items.length - 1]
        : items[Math.min(items.length - 1, Math.max(0, cur + dir))];
    highlight(root, next);
  }
});

// --- pointer on items: hover highlights, click commits (command.js model) -

on("pointerover", "[data-gsxui-combobox-item]", (_e, item) => {
  if (isDisabled(item)) return;
  highlight(rootOf(item), item);
});

on("click", "[data-gsxui-combobox-item]", (_e, item) => {
  commit(rootOf(item), item);
});
