// Navigation menu behavior on the native popover API: a hover mega-menu,
// not a menu of items — no role="menu"/role="menuitem" machinery here (see
// ui/dropdown.js for that shape; nothing in this file needs it). Modeled on
// ui/hover-card.js's own hover/focus-open + grace-period-close precedent —
// the task brief's own explicit pointer — NOT on dropdown.js's own
// click+toggle-event model: every popover here is popover="manual" (top
// layer, no light dismiss/Esc; hover and focus alone drive open/close,
// matching hover-card's own deliberate choice, restated in
// ui/navigation-menu.gsx's own file header and docs/jsx-parity.md
// ## navigation-menu).
//
// FIX ROUND 1 (shared viewport removed, not reworked): a v1 of this file
// positioned/sized a SEPARATE NavigationMenuViewport popover in lockstep
// with the active Content, both shown at the same rect. That does not
// work: two coincident top-layer popovers stack, and whichever is shown
// second paints an OPAQUE surface directly over the other — confirmed live
// (the panel rendered as an empty box, `elementFromPoint` returned the
// viewport at every sampled point inside it, the content's own links were
// completely unreachable). This file now drives shadcn's own
// `viewport={false}` configuration instead (see ui/navigation-menu.gsx's
// own file header GAP paragraph for the full rationale): each
// NavigationMenuContent is independently popover="manual", DOM-nested in
// its own NavigationMenuItem, carries its own full chrome, and is
// positioned under ITS OWN trigger's rect — no shared viewport element
// exists to position, size, or occlude anything.
//
// No Escape/outside-pointerdown light dismiss — same deliberate choice as
// hover-card.js's own (hover/focus drive it, not outside clicks or Esc).
import { on, emit } from "./gsxui.js";

const menuOf = (el) => el.closest("[data-gsxui-navigation-menu]");
const listOf = (menu) => menu?.querySelector("[data-gsxui-navigation-menu-list]");
const itemOf = (el) => el.closest("[data-gsxui-navigation-menu-item]");
const contentOf = (trigger) =>
  itemOf(trigger)?.querySelector(":scope > [data-gsxui-navigation-menu-content]");
const triggerOf = (content) =>
  itemOf(content)?.querySelector(":scope > [data-gsxui-navigation-menu-trigger]");

// Radix's own real delayDuration is unread (no dist in the reference
// checkout) — reusing hover-card.js's own OPEN_DELAY/CLOSE_DELAY (100ms
// each) is the house default for a hover-delay component, per the brief's
// own "match the shipped hover-card.js precedent" instruction.
const OPEN_DELAY = 100;
const CLOSE_DELAY = 100;
const timers = new WeakMap(); // trigger -> pending open/close setTimeout id

function clearTimer(trigger) {
  clearTimeout(timers.get(trigger));
  timers.delete(trigger);
}

function isAnyOpen(menu) {
  return [...menu.querySelectorAll("[data-gsxui-navigation-menu-trigger]")].some((t) =>
    contentOf(t)?.matches(":popover-open"),
  );
}

function positionAt(el, rect, offset) {
  el.style.position = "fixed";
  el.style.inset = "auto";
  el.style.left = `${rect.left}px`;
  el.style.top = `${rect.bottom + offset}px`;
}

function stillWithin(trigger, related) {
  if (!(related instanceof Element)) return false;
  return !!itemOf(trigger)?.contains(related);
}

function positionIndicator(menu, trigger) {
  const list = listOf(menu);
  const indicator = list?.querySelector(":scope > [data-gsxui-navigation-menu-indicator]");
  if (!list || !indicator) return;
  // Radix's own runtime sets position:relative on whatever hosts the
  // indicator's own translate-relative-to inline, at runtime, rather than
  // baking it into a class — same idiom every fixed-positioned popover in
  // this codebase already uses for its own position:fixed.
  if (!list.style.position) list.style.position = "relative";
  indicator.style.position = "absolute";
  const listRect = list.getBoundingClientRect();
  const triggerRect = trigger.getBoundingClientRect();
  indicator.style.left = `${triggerRect.left - listRect.left}px`;
  indicator.style.width = `${triggerRect.width}px`;
  indicator.dataset.state = "visible";
}

function hideIndicator(menu) {
  const indicator = listOf(menu)?.querySelector(":scope > [data-gsxui-navigation-menu-indicator]");
  if (indicator) indicator.dataset.state = "hidden";
}

function open(trigger) {
  clearTimer(trigger);
  const content = contentOf(trigger);
  if (!content || content.matches(":popover-open")) return; // idempotent
  const menu = menuOf(trigger);
  if (!menu) return;

  // Close any OTHER open panel in this menu first — sibling panels are
  // never DOM-nested in one another, so switching is close-then-open, the
  // same shape as menubar.js's own openMenu.
  for (const other of menu.querySelectorAll("[data-gsxui-navigation-menu-content]")) {
    if (other !== content && other.matches(":popover-open")) close(triggerOf(other));
  }

  // Positioned under THIS trigger's own rect (viewport={false} semantics —
  // each panel is its own self-contained floating panel, not a shared,
  // list-anchored one).
  positionAt(content, trigger.getBoundingClientRect(), 6);
  // Stamp open BEFORE showing — same flash-avoidance rule as every other
  // popover in this codebase (a queued toggle event can otherwise leave one
  // frame painted in the stale closed state).
  content.dataset.state = "open";
  content.showPopover();
  trigger.dataset.state = "open";
  trigger.setAttribute("aria-expanded", "true");
  emit(content, "gsxui:open");

  positionIndicator(menu, trigger);
}

function close(trigger) {
  if (!trigger) return;
  clearTimer(trigger);
  const content = contentOf(trigger);
  if (!content || !content.matches(":popover-open")) return;
  content.hidePopover();
  content.dataset.state = "closed";
  trigger.dataset.state = "closed";
  trigger.setAttribute("aria-expanded", "false");
  emit(content, "gsxui:close");

  const menu = menuOf(trigger);
  if (menu && !isAnyOpen(menu)) hideIndicator(menu);
}

// Switching between two ALREADY-open triggers in the same menu is
// immediate, no grace period — the "connected bar" feel menubar.js's own
// open-follows-hover already established for its own sibling switch.
// Opening the FIRST panel in a menu still rides the grace period, avoiding
// an accidental trigger while the pointer merely passes over the bar.
function scheduleOpen(trigger) {
  clearTimer(trigger);
  const menu = menuOf(trigger);
  if (menu && isAnyOpen(menu)) {
    open(trigger);
    return;
  }
  timers.set(trigger, setTimeout(() => open(trigger), OPEN_DELAY));
}

function scheduleClose(trigger) {
  clearTimer(trigger);
  timers.set(trigger, setTimeout(() => close(trigger), CLOSE_DELAY));
}

on("pointerover", "[data-gsxui-navigation-menu-trigger]", (_e, t) => scheduleOpen(t));
on("pointerout", "[data-gsxui-navigation-menu-trigger]", (e, t) => {
  if (stillWithin(t, e.relatedTarget)) return;
  scheduleClose(t);
});
on("focusin", "[data-gsxui-navigation-menu-trigger]", (_e, t) => open(t));
on("focusout", "[data-gsxui-navigation-menu-trigger]", (e, t) => {
  if (stillWithin(t, e.relatedTarget)) return;
  scheduleClose(t);
});
// Click toggles — covers touch/no-hover pointer types, which never fire the
// pointerover open above, same rationale as every submenu's own click
// handler in this codebase.
on("click", "[data-gsxui-navigation-menu-trigger]", (_e, t) => {
  const content = contentOf(t);
  if (content?.matches(":popover-open")) close(t);
  else open(t);
});

on("pointerover", "[data-gsxui-navigation-menu-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) clearTimer(trigger);
});
on("pointerout", "[data-gsxui-navigation-menu-content]", (e, content) => {
  const trigger = triggerOf(content);
  if (!trigger || stillWithin(trigger, e.relatedTarget)) return;
  scheduleClose(trigger);
});
on("focusin", "[data-gsxui-navigation-menu-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) clearTimer(trigger);
});
on("focusout", "[data-gsxui-navigation-menu-content]", (e, content) => {
  const trigger = triggerOf(content);
  if (!trigger || stillWithin(trigger, e.relatedTarget)) return;
  scheduleClose(trigger);
});

// Selecting a Link inside an open Content closes that panel. A Link with no
// Content ancestor (a plain top-level nav link with no dropdown, styled via
// NavigationMenuTriggerStyle()) has nothing to close — the click just
// navigates.
on("click", "[data-gsxui-navigation-menu-link]", (_e, link) => {
  const content = link.closest("[data-gsxui-navigation-menu-content]");
  if (!content) return;
  emit(link, "gsxui:select");
  close(triggerOf(content));
});
