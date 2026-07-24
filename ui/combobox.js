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

function highlight(root, item) {
  const input = inputOf(root);
  for (const other of itemsOf(root)) {
    if (other !== item) delete other.dataset.highlighted;
  }
  if (!item) {
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
function filter(root) {
  const input = inputOf(root);
  const query = (input?.value ?? "").trim();
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

function commit(root, item) {
  if (!item || isDisabled(item)) return;
  for (const other of itemsOf(root)) {
    const isIt = other === item;
    other.dataset.state = isIt ? "checked" : "unchecked";
    other.setAttribute("aria-selected", isIt ? "true" : "false");
  }
  const value = item.dataset.value ?? "";
  const input = inputOf(root);
  if (input) input.value = labelOf(item);
  const bridge = bridgeOf(root);
  if (bridge) {
    bridge.value = value;
    bridge.dispatchEvent(new Event("change", { bubbles: true }));
  }
  emit(root, "gsxui:select", { value });
  contentOf(root)?.hidePopover();
}

function clearValue(root) {
  for (const item of itemsOf(root)) {
    item.dataset.state = "unchecked";
    item.setAttribute("aria-selected", "false");
  }
  const input = inputOf(root);
  if (input) input.value = "";
  const bridge = bridgeOf(root);
  if (bridge) {
    bridge.value = "";
    bridge.dispatchEvent(new Event("change", { bubbles: true }));
  }
  emit(root, "gsxui:select", { value: "" });
  filter(root);
  input?.focus();
}

// --- init: group aria-labelledby wiring, reflect a server-checked item ---

function init(root) {
  for (const group of root.querySelectorAll('[data-slot="combobox-group"]')) {
    if (group.getAttribute("aria-labelledby")) continue;
    const label = group.querySelector('[data-slot="combobox-label"]');
    if (!label) continue;
    if (!label.id) label.id = `gsxui-combobox-label-${++uid}`;
    group.setAttribute("aria-labelledby", label.id);
  }
  const checked = root.querySelector('[data-gsxui-combobox-item][data-state="checked"]');
  const input = inputOf(root);
  if (checked && input && !input.value) input.value = labelOf(checked);
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
  filter(root);
}

function closeContent(root) {
  contentOf(root)?.hidePopover();
}

on("toggle", "[data-gsxui-combobox-content]", (e, content) => {
  const open = e.newState === "open";
  content.dataset.state = open ? "open" : "closed";
  const root = rootOf(content);
  const input = inputOf(root);
  input?.setAttribute("aria-expanded", open ? "true" : "false");
  if (open) {
    // aria-controls points at the listbox (role="listbox" lives on
    // ComboboxList, not on this popover surface), falling back to the
    // content itself if a caller omitted ComboboxList.
    const list = listOf(root) ?? content;
    if (!list.id) list.id = `gsxui-combobox-list-${++uid}`;
    input?.setAttribute("aria-controls", list.id);
    emit(content, "gsxui:open");
  } else {
    input?.removeAttribute("aria-controls");
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

// --- the input: opens on focus, filters on input, keyboard nav -----------

on(
  "focus",
  "[data-gsxui-combobox-input]",
  (_e, input) => {
    openContent(rootOf(input));
  },
  { capture: true },
);

on("input", "[data-gsxui-combobox-input]", (_e, input) => {
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
