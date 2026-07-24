// ToggleGroup behavior: roving tabindex (server renders every item as a
// plain tab stop — graceful no-JS fallback — this module collapses that to
// exactly one tab stop at init and on every interaction), arrow-key
// navigation, and click activation with the type="single"/"multiple" split
// baked into each item's own role (role="radio" => single, ported straight
// off Radix's own ARIA split rather than a redundant data-type stamp —
// ToggleGroup/ToggleGroupItem already only ever pair a "radio" role with
// type="single").
import { on, emit } from "./gsxui.js";

const itemsOf = (root) =>
  [...root.querySelectorAll("[data-gsxui-toggle-group-item]")].filter(
    (i) => i.closest("[data-gsxui-toggle-group]") === root,
  );
const isSingle = (item) => item.getAttribute("role") === "radio";

function setPressed(item, pressed) {
  item.dataset.state = pressed ? "on" : "off";
  if (isSingle(item)) {
    item.setAttribute("aria-checked", pressed ? "true" : "false");
  } else {
    item.setAttribute("aria-pressed", pressed ? "true" : "false");
  }
}

// Exactly one item (the current tab stop) gets tabindex="0"; every other
// item gets "-1". Because this is the ONLY tab stop the group ever exposes
// once normalized, a plain Shift+Tab from it already lands on whatever
// precedes the group in the page — "exits the group in one press" falls out
// of roving tabindex itself, no onItemShiftTab-style keydown interception
// needed (that Radix mechanism earns its keep inside their FocusScope
// machinery, which this plain-DOM port doesn't have or need).
function setActiveTabStop(root, item) {
  for (const i of itemsOf(root)) i.tabIndex = i === item ? 0 : -1;
}

// Entry-tabstop priority at init, mirroring RovingFocusGroup's own
// fallback chain minus "last-focused" (nothing to remember across a fresh
// page load): the pressed item wins if one exists and isn't disabled,
// else the first non-disabled item.
function normalize(root) {
  const items = itemsOf(root);
  const enabled = items.filter((i) => !i.disabled);
  if (!enabled.length) return;
  const pressed = items.find((i) => i.dataset.state === "on" && !i.disabled);
  setActiveTabStop(root, pressed ?? enabled[0]);
}

on("click", "[data-gsxui-toggle-group-item]", (_e, item) => {
  if (item.disabled) return;
  const root = item.closest("[data-gsxui-toggle-group]");
  if (!root) return;
  const single = isSingle(item);
  // Single-type replace-on-activate: activating a new item just sets a new
  // value (every OTHER item in the group is force-cleared below), and
  // re-activating the already-pressed item toggles it off — Radix allows
  // an empty single value unless a caller opts otherwise (ToggleGroupImplSingle's
  // onItemActivate === setValue, no "who wins" race to resolve here since
  // there's no shared value state to update, just per-item data-state).
  const pressed = single ? item.dataset.state !== "on" : item.getAttribute("aria-pressed") !== "true";
  setPressed(item, pressed);
  if (single && pressed) {
    for (const other of itemsOf(root)) {
      if (other !== item) setPressed(other, false);
    }
  }
  setActiveTabStop(root, item);
  const value = single
    ? pressed ? item.dataset.value : ""
    : itemsOf(root).filter((i) => i.dataset.state === "on").map((i) => i.dataset.value);
  emit(root, "gsxui:change", { value });
});

on("keydown", "[data-gsxui-toggle-group-item]", (e, item) => {
  // A held modifier suppresses the focus-move entirely (Radix: "if
  // (event.metaKey || ...) return") — arrow+modifier does nothing, not even
  // scroll, so this handler steps aside for the browser's own default.
  if (e.metaKey || e.ctrlKey || e.altKey || e.shiftKey) return;
  const root = item.closest("[data-gsxui-toggle-group]");
  if (!root) return;
  const items = itemsOf(root).filter((i) => !i.disabled);
  if (!items.length) return;
  const dir = { ArrowLeft: -1, ArrowUp: -1, ArrowRight: 1, ArrowDown: 1 }[e.key];
  let next;
  if (dir) {
    const i = items.indexOf(item);
    next = items[(i + dir + items.length) % items.length];
  } else if (e.key === "Home" || e.key === "PageUp") {
    next = items[0];
  } else if (e.key === "End" || e.key === "PageDown") {
    next = items[items.length - 1];
  } else {
    return;
  }
  e.preventDefault();
  setActiveTabStop(root, next);
  next.focus();
});

// Initial tab-stop assignment for groups rendered without JS having run yet
// (server renders every item as a plain tab stop — see the package doc
// comment on ui/toggle-group.gsx) — same one-time init-scan shape as
// command.js's own initial filter() pass.
for (const root of document.querySelectorAll("[data-gsxui-toggle-group]")) normalize(root);
