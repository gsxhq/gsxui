// Custom-listbox Select behavior. Built on the SAME native-popover machinery
// as dropdown.js (popover="auto" top layer, light dismiss + free Esc, the
// closest()-proximity trigger/content wiring, the sync-data-state-before-
// showPopover flash fix, the wasOpen pointerdown/click guard, and the
// hover-is-focus / arrow-key .focus() roving-focus idioms). What this module
// adds on top of dropdown's machinery — all genuinely new, none ported:
//   - a VALUE MODEL (one checked item per root; trigger text follows the
//     selected item; data-placeholder cleared on first selection),
//   - aria-selected recomputed as (isValue AND isFocused) on every focus
//     change — an item that IS the value but isn't the highlighted one
//     reports aria-selected="false"; data-state="checked"|"unchecked" is the
//     separate attribute that tracks the value alone and drives the checkmark,
//   - a bespoke 1000ms prefix TYPEAHEAD (startsWith, wrap-from-current,
//     same-character-repeat cycling) that runs on the CLOSED trigger too
//     (typing there changes the value in place); Space is swallowed as a
//     search char only while a search is already in progress,
//   - a hidden native <select> FORM BRIDGE, populated from the DOM items at
//     init and kept in sync (value assignment + a bubbling change event).
// See docs/jsx-parity.md ## select and ui/select.gsx's header.
import { on, emit } from "./gsxui.js";

const rootOf = (el) => el.closest("[data-gsxui-select]");
const contentOf = (el) => rootOf(el)?.querySelector("[data-gsxui-select-content]");
const triggerOf = (el) => rootOf(el)?.querySelector("[data-gsxui-select-trigger]");
const bridgeOf = (root) => root.querySelector("[data-gsxui-select-bridge]");
const itemsOf = (content) => [...content.querySelectorAll("[data-gsxui-select-item]")];
const isDisabled = (item) =>
  item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset;
const enabledItems = (content) => itemsOf(content).filter((i) => !isDisabled(i));
const labelOf = (item) =>
  (item.querySelector('[data-slot="select-item-text"]') ?? item).textContent.trim();
const focusedItem = (content) => {
  const a = document.activeElement;
  return a instanceof Element && content.contains(a) && a.matches("[data-gsxui-select-item]")
    ? a
    : null;
};

const OPEN_KEYS = [" ", "Enter", "ArrowUp", "ArrowDown"];

let uid = 0;

// --- value model -----------------------------------------------------------

// applyValue stamps exactly one item checked, follows the trigger's displayed
// text, clears the placeholder, and syncs the hidden bridge. silent skips the
// public gsxui:select event (used by the init pass reflecting a server-checked
// item, which is not a user action).
function applyValue(root, item, { silent = false } = {}) {
  const content = root.querySelector("[data-gsxui-select-content]");
  const trigger = root.querySelector("[data-gsxui-select-trigger]");
  if (content) {
    for (const other of itemsOf(content)) {
      other.dataset.state = other === item ? "checked" : "unchecked";
    }
  }
  const value = item.dataset.value ?? "";
  const valueEl = trigger?.querySelector('[data-slot="select-value"]');
  if (valueEl) valueEl.textContent = labelOf(item);
  trigger?.removeAttribute("data-placeholder");
  const bridge = bridgeOf(root);
  if (bridge) {
    bridge.value = value;
    bridge.dispatchEvent(new Event("change", { bubbles: true }));
  }
  if (!silent) emit(root, "gsxui:select", { value });
}

// selectItem is the user-facing "pick this option": apply the value, close.
function selectItem(item) {
  if (isDisabled(item)) return;
  applyValue(rootOf(item), item);
  contentOf(item)?.hidePopover();
}

// --- focus + aria-selected recompute --------------------------------------

// focusItem moves real DOM focus onto item and recomputes aria-selected on the
// whole listbox: aria-selected = (isValue AND isFocused), so only the newly-
// focused item can be true, and only when it is also the checked value.
function focusItem(item, { scroll = false } = {}) {
  const content = contentOf(item);
  if (content) {
    for (const other of itemsOf(content)) {
      if (other !== item) other.setAttribute("aria-selected", "false");
    }
  }
  item.setAttribute("aria-selected", item.dataset.state === "checked" ? "true" : "false");
  item.focus({ preventScroll: true });
  if (scroll) item.scrollIntoView({ block: "nearest" });
}

// --- typeahead (bespoke: prefix + 1000ms buffer + same-char cycle) --------

let buffer = "";
let bufferTimer;

function pushSearch(key) {
  buffer += key;
  clearTimeout(bufferTimer);
  bufferTimer = setTimeout(() => {
    buffer = "";
  }, 1000);
}

// findMatch does prefix startsWith matching, wrapping the candidate list to
// start at the current item. When the whole buffer is one repeated character
// (e.g. "bbb") it is normalized to a single char and the current item is
// excluded, so repeated presses of the same key step through matches one at a
// time — the exact findNextItem contract traced from Radix.
function findMatch(items, search, current) {
  const norm = search.toLowerCase();
  const allSame = [...norm].every((c) => c === norm[0]);
  const query = allSame ? norm[0] : norm;
  const start = Math.max(0, current ? items.indexOf(current) : 0);
  let wrapped = items.slice(start).concat(items.slice(0, start));
  if (query.length === 1 && current) wrapped = wrapped.filter((it) => it !== current);
  return wrapped.find((it) => labelOf(it).toLowerCase().startsWith(query)) ?? null;
}

// --- init: bridge population, aria wiring, server-checked reflection -------

function populateBridge(content, bridge) {
  bridge.textContent = "";
  const empty = document.createElement("option");
  empty.value = "";
  bridge.appendChild(empty);
  for (const item of itemsOf(content)) {
    const opt = document.createElement("option");
    opt.value = item.dataset.value ?? "";
    opt.textContent = labelOf(item);
    bridge.appendChild(opt);
  }
}

function init(root) {
  const content = root.querySelector("[data-gsxui-select-content]");
  const trigger = root.querySelector("[data-gsxui-select-trigger]");
  if (content) {
    // Wire each group's aria-labelledby to its own label's generated id.
    for (const group of content.querySelectorAll("[data-gsxui-select-group]")) {
      if (group.getAttribute("aria-labelledby")) continue;
      const label = group.querySelector('[data-slot="select-label"]');
      if (!label) continue;
      if (!label.id) label.id = `gsxui-select-label-${++uid}`;
      group.setAttribute("aria-labelledby", label.id);
    }
  }
  const bridge = bridgeOf(root);
  if (bridge && content) {
    populateBridge(content, bridge);
    // aria-required lives on the combobox trigger; the bridge carries the real
    // required attribute (there is no context to pass it directly to the trigger).
    if (trigger && bridge.required) trigger.setAttribute("aria-required", "true");
  }
  // Reflect a server-rendered checked item into the trigger text + bridge value.
  const checked = content?.querySelector(
    '[data-gsxui-select-item][data-state="checked"]',
  );
  if (checked) applyValue(root, checked, { silent: true });
}

for (const root of document.querySelectorAll("[data-gsxui-select]")) init(root);

// --- open / close (ported dropdown.js machinery) --------------------------

function openContent(trigger, content) {
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.inset = "auto";
  content.style.left = `${r.left}px`;
  content.style.top = `${r.bottom + 4}px`;
  // Popper-equivalent width: never narrower than the trigger (Radix's
  // --radix-select-trigger-width). min-w-36 from the class still floors it.
  content.style.minWidth = `${r.width}px`;
  // Stamp open BEFORE showing — the toggle event that also stamps it is a
  // separate queued task; a paint can land in the gap and flash the closed
  // state (same fix dropdown.js documents).
  content.dataset.state = "open";
  content.showPopover();
}

on("pointerdown", "[data-gsxui-select-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (content) {
    trigger.dataset.gsxuiWasOpen = content.matches(":popover-open") ? "true" : "false";
  }
});

on("click", "[data-gsxui-select-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (!content) return;
  const wasOpen = trigger.dataset.gsxuiWasOpen === "true";
  delete trigger.dataset.gsxuiWasOpen;
  if (wasOpen) {
    // Light dismiss on the outside-pointerdown may already have closed it;
    // converge on the real state rather than assuming.
    if (content.matches(":popover-open")) content.hidePopover();
    return;
  }
  if (content.matches(":popover-open")) {
    content.hidePopover();
    return;
  }
  openContent(trigger, content);
});

on(
  "toggle",
  "[data-gsxui-select-content]",
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const trigger = triggerOf(content);
    trigger?.setAttribute("aria-expanded", open ? "true" : "false");
    if (open) {
      if (!content.id) content.id = `gsxui-select-content-${++uid}`;
      trigger?.setAttribute("aria-controls", content.id);
      delete trigger?.dataset.gsxuiWasOpen;
      // Focus the checked item if any, else the first enabled one.
      const target =
        content.querySelector(
          '[data-gsxui-select-item][data-state="checked"]:not([aria-disabled="true"])',
        ) ?? enabledItems(content)[0];
      if (target) focusItem(target, { scroll: true });
      emit(content, "gsxui:open");
    } else {
      trigger?.removeAttribute("aria-controls");
      emit(content, "gsxui:close");
    }
  },
  { capture: true },
);

// --- keyboard on the open listbox -----------------------------------------

on("keydown", "[data-gsxui-select-content]", (e, content) => {
  // Tab never leaves an open listbox (matches a native <select>'s open dropdown).
  if (e.key === "Tab") {
    e.preventDefault();
    return;
  }
  // Printable char → typeahead. Space is a search char only mid-search;
  // otherwise it falls through to selection below.
  if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey && !(e.key === " " && buffer === "")) {
    e.preventDefault();
    pushSearch(e.key);
    const match = findMatch(enabledItems(content), buffer, focusedItem(content));
    if (match) focusItem(match, { scroll: true });
    return;
  }
  if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    const item = focusedItem(content);
    if (item) selectItem(item);
    return;
  }
  const items = enabledItems(content);
  if (!items.length) return;
  const cur = items.indexOf(focusedItem(content));
  let next;
  switch (e.key) {
    case "ArrowDown":
      next = items[Math.min(items.length - 1, cur + 1)] ?? items[0];
      break;
    case "ArrowUp":
      next = items[Math.max(0, cur - 1)] ?? items[items.length - 1];
      break;
    case "Home":
      next = items[0];
      break;
    case "End":
      next = items[items.length - 1];
      break;
    default:
      return;
  }
  e.preventDefault();
  focusItem(next, { scroll: true });
});

// contextmenu inside the listbox is suppressed (Radix does the same).
on("contextmenu", "[data-gsxui-select-content]", (e) => e.preventDefault());

// --- keyboard on the CLOSED trigger ---------------------------------------

on("keydown", "[data-gsxui-select-trigger]", (e, trigger) => {
  const content = contentOf(trigger);
  if (!content || content.matches(":popover-open")) return;
  // Typeahead on the closed trigger selects the value in place (no open),
  // exactly like a native <select>. Space with no active search is an OPEN
  // request, not a search char.
  if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey && !(e.key === " " && buffer === "")) {
    e.preventDefault();
    pushSearch(e.key);
    const current = content.querySelector(
      '[data-gsxui-select-item][data-state="checked"]',
    );
    const match = findMatch(enabledItems(content), buffer, current);
    if (match) applyValue(rootOf(trigger), match);
    return;
  }
  if (OPEN_KEYS.includes(e.key)) {
    e.preventDefault();
    openContent(trigger, content);
  }
});

// --- pointer: hover-is-focus + pointer-type-aware selection ---------------

// Hovering an item with the MOUSE focuses it (Radix gates this on
// pointerType==="mouse"; touch/pen do not focus-on-hover).
on("pointermove", "[data-gsxui-select-item]", (e, item) => {
  if (e.pointerType !== "mouse" || isDisabled(item)) return;
  if (document.activeElement === item) return;
  focusItem(item);
});

// Mouse selects on pointerup (Radix's snappier path); a suppression flag stops
// the trailing click from double-firing. Touch/pen fall through to click.
let suppressClick = false;
on("pointerup", "[data-gsxui-select-item]", (e, item) => {
  if (e.pointerType !== "mouse" || isDisabled(item)) return;
  suppressClick = true;
  selectItem(item);
  setTimeout(() => {
    suppressClick = false;
  }, 0);
});

on("click", "[data-gsxui-select-item]", (_e, item) => {
  if (suppressClick) {
    suppressClick = false;
    return;
  }
  selectItem(item);
});

// Leaving the listbox parks focus on the content (tabindex="-1", so arrow keys
// keep working) and clears every item's aria-selected — same shape as
// dropdown.js's content-level pointerout handler.
on("pointerout", "[data-gsxui-select-content]", (e, content) => {
  if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget)) return;
  if (!content.contains(document.activeElement)) return;
  for (const item of itemsOf(content)) item.setAttribute("aria-selected", "false");
  content.focus();
});
