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
// ARCHITECTURE (no portal, no DOM-node-moving either — see
// ui/navigation-menu.gsx's own file header for the full rationale): every
// NavigationMenuContent is independently popover="manual", DOM-nested
// inside its own NavigationMenuItem (never moved or portalled). When a
// shared NavigationMenuViewport exists (data-viewport!="false"), it is a
// SEPARATE, second popover="manual" panel shown/positioned/sized in
// lockstep with whichever Content is currently active — coincident, not
// containing: both open at the identical fixed left/top rect (computed off
// the List's own boundingClientRect), Viewport supplies the chrome and is
// resized to the active Content's own measured width/height, Content stays
// unchromed in this mode.
//
// SIZING (source map ## navigation-menu §3, decisive): one measurement on
// activation (a getBoundingClientRect read off the newly-active Content)
// plus a ResizeObserver on that same Content for late-arriving resizes —
// explicitly NOT a per-frame tween between panel sizes, since no CSS in any
// of the eight nova-family stylesheets transitions width/height on the
// viewport, only a scale transform gated on open/close.
//
// data-motion (new-york-v4's own direction-aware "which side did this panel
// enter from" slide) is NOT reproduced — GAP, ledgered in
// ui/navigation-menu.gsx's own NavigationMenuContent doc comment: the
// task's own binding decision reuses the popover family's discrete-
// transition block byte-identically instead.
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
const viewportOf = (menu) => menu?.querySelector("[data-gsxui-navigation-menu-viewport]");

// Radix's own real delayDuration is unread (no dist in the reference
// checkout) — reusing hover-card.js's own OPEN_DELAY/CLOSE_DELAY (100ms
// each) is the house default for a hover-delay component, per the brief's
// own "match the shipped hover-card.js precedent" instruction.
const OPEN_DELAY = 100;
const CLOSE_DELAY = 100;
const timers = new WeakMap(); // trigger -> pending open/close setTimeout id
const resizeObservers = new WeakMap(); // viewport -> ResizeObserver

function clearTimer(trigger) {
  clearTimeout(timers.get(trigger));
  timers.delete(trigger);
}

function isAnyOpen(menu) {
  return [...menu.querySelectorAll("[data-gsxui-navigation-menu-trigger]")].some((t) =>
    contentOf(t)?.matches(":popover-open"),
  );
}

function sizeViewport(viewport, content) {
  const r = content.getBoundingClientRect();
  viewport.style.width = `${r.width}px`;
  viewport.style.height = `${r.height}px`;
}

// One measurement on activation, plus a ResizeObserver for content that
// changes size AFTER it's already active (an image finishing loading, etc) —
// the source map's own decisive "discrete state, not a per-frame tween" scope.
function observeViewport(viewport, content) {
  let ro = resizeObservers.get(viewport);
  if (ro) ro.disconnect();
  else {
    ro = new ResizeObserver(() => sizeViewport(viewport, content));
    resizeObservers.set(viewport, ro);
  }
  ro.observe(content);
}

function positionAt(el, rect, offset) {
  el.style.position = "fixed";
  el.style.inset = "auto";
  el.style.left = `${rect.left}px`;
  el.style.top = `${rect.bottom + offset}px`;
}

function stillWithin(trigger, related) {
  if (!(related instanceof Element)) return false;
  const item = itemOf(trigger);
  if (item?.contains(related)) return true;
  const viewport = viewportOf(menuOf(trigger));
  return !!viewport?.contains(related);
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
  const list = listOf(menu) ?? menu;
  const rect = list.getBoundingClientRect();

  // Close any OTHER open panel in this menu first — sibling panels are
  // never DOM-nested in one another, so switching is close-then-open, the
  // same shape as menubar.js's own openMenu.
  for (const other of menu.querySelectorAll("[data-gsxui-navigation-menu-content]")) {
    if (other !== content && other.matches(":popover-open")) close(triggerOf(other));
  }

  positionAt(content, rect, 6);
  // Stamp open BEFORE showing — same flash-avoidance rule as every other
  // popover in this codebase (a queued toggle event can otherwise leave one
  // frame painted in the stale closed state).
  content.dataset.state = "open";
  content.showPopover();
  trigger.dataset.state = "open";
  trigger.setAttribute("aria-expanded", "true");
  emit(content, "gsxui:open");

  const viewport = viewportOf(menu);
  if (viewport) {
    positionAt(viewport, rect, 6);
    sizeViewport(viewport, content); // AFTER showPopover — a hidden popover has no layout box
    viewport.dataset.state = "open";
    if (!viewport.matches(":popover-open")) viewport.showPopover();
    observeViewport(viewport, content);
  }
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
  if (menu && !isAnyOpen(menu)) {
    const viewport = viewportOf(menu);
    if (viewport?.matches(":popover-open")) {
      viewport.hidePopover();
      viewport.dataset.state = "closed";
    }
    hideIndicator(menu);
  }
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
