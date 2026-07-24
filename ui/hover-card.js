// Hover card behavior on the native popover API: top layer, no light
// dismiss (hover/focus drive it, not outside clicks or Esc) — same
// popover="manual" mechanism as tooltip.js. Adapted from tooltip.js:
//   - anchored BELOW the trigger, not above (Radix HoverCard's own default
//     side is bottom; Tooltip's is top) — the top calc is flipped, and the
//     left calc (already centered in tooltip.js) is reused as-is.
//   - Radix HoverCard's own openDelay/closeDelay (700ms/300ms) replace
//     tooltip's flat 300ms-open/immediate-close. Because HoverCardContent
//     can hold real interactive content (a link, a bio — see hover-card.tsx's
//     own @nextjs demo), a closeDelay only does its job if hovering onto the
//     content itself also cancels the pending close: two extra listeners on
//     [data-gsxui-hover-card-content] do exactly that (cancel on enter,
//     reschedule on leave) — tooltip.js needs neither, since its content is
//     never interactive.
//   - no arrow (hover-card has none, unlike tooltip's diamond span).
import { on, emit } from "./gsxui.js";

const OPEN_DELAY = 700; // Radix HoverCard's default openDelay
const CLOSE_DELAY = 300; // Radix HoverCard's default closeDelay

const timers = new WeakMap(); // trigger -> pending open/close setTimeout id
const contentOf = (el) =>
  el.closest("[data-gsxui-hovercard]")?.querySelector("[data-gsxui-hovercard-content]");
const triggerOf = (el) =>
  el.closest("[data-gsxui-hovercard]")?.querySelector("[data-gsxui-hovercard-trigger]");

function clearTimer(trigger) {
  clearTimeout(timers.get(trigger));
  timers.delete(trigger);
}

function show(trigger) {
  clearTimer(trigger);
  const content = contentOf(trigger);
  if (!content || content.matches(":popover-open")) return;
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.inset = "auto";
  content.showPopover();
  // Position numerically AFTER showing (hidden popovers have no box) and
  // never via transform — same rationale as tooltip.js's own comment.
  // Centered below the trigger (Radix HoverCard's own align=center,
  // side=bottom default): left is the trigger's horizontal midpoint minus
  // half the content's own width (identical calc to tooltip.js's); top is
  // the trigger's bottom edge plus Radix's own 4px sideOffset (hover-
  // card.tsx's default, same value as popover.tsx's — no arrow to clear
  // room for, unlike tooltip's 6px-plus-arrow gap above the trigger).
  content.style.left = `${r.left + r.width / 2 - content.offsetWidth / 2}px`;
  content.style.top = `${r.bottom + 4}px`;
  content.dataset.state = "open";
  emit(content, "gsxui:open");
}

function hide(trigger) {
  clearTimer(trigger);
  const content = contentOf(trigger);
  if (!content || !content.matches(":popover-open")) return;
  content.hidePopover();
  content.dataset.state = "closed";
  emit(content, "gsxui:close");
}

function scheduleShow(trigger) {
  clearTimer(trigger);
  timers.set(trigger, setTimeout(() => show(trigger), OPEN_DELAY));
}

function scheduleHide(trigger) {
  clearTimer(trigger);
  timers.set(trigger, setTimeout(() => hide(trigger), CLOSE_DELAY));
}

on("pointerover", "[data-gsxui-hovercard-trigger]", (_e, t) => scheduleShow(t));
on("pointerout", "[data-gsxui-hovercard-trigger]", (_e, t) => scheduleHide(t));
on("focusin", "[data-gsxui-hovercard-trigger]", (_e, t) => show(t));
on("focusout", "[data-gsxui-hovercard-trigger]", (_e, t) => hide(t));

// The content itself can hold real interactive children — hovering onto it
// must cancel any pending close the trigger's own pointerout just
// scheduled (the pointer is still within the hover-card's own hit area, it
// simply moved from the trigger to the content sitting 4px below it), and
// leaving it schedules the same delayed close the trigger uses. Both
// helpers are keyed by trigger throughout (contentOf/triggerOf), so a
// trigger's pending timer is the single source of truth regardless of
// which of the two elements the pointer is currently over.
on("pointerover", "[data-gsxui-hovercard-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) clearTimer(trigger);
});
on("pointerout", "[data-gsxui-hovercard-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) scheduleHide(trigger);
});
