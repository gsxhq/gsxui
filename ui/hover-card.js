// Hover card behavior on the native popover API: top layer, no light
// dismiss (hover/focus drive it, not outside clicks or Esc) — same
// popover="manual" mechanism as tooltip.js. Adapted from tooltip.js:
//   - anchored BELOW the trigger, not above (Radix HoverCard's own default
//     side is bottom; Tooltip's is top) — the top calc is flipped, and the
//     left calc (already centered in tooltip.js) is reused as-is.
//   - Radix HoverCard's own openDelay/closeDelay (700ms/300ms) replace
//     tooltip's flat 300ms-open/immediate-close, and — unlike tooltip.js —
//     BOTH the pointer and the keyboard-focus leave paths ride the same
//     scheduleHide/closeDelay grace period, not just pointer. Radix treats
//     hover and focus as one unified "is the user still interacting with
//     this trigger-or-content pair" model, not two independent show/hide
//     controls the way tooltip.js's flat immediate-focusout-hide does; a
//     keyboard user tabbing from the trigger into a focusable child of
//     HoverCardContent (a link, a bio — see hover-card.tsx's own @nextjs
//     demo) must not have the popover display:none'd out from under them
//     mid-tab. Content-side listeners mirror the trigger-side pair for both
//     input modalities: pointerover/focusin on
//     [data-gsxui-hovercard-content] cancel the pending close scheduled by
//     leaving the trigger, pointerout/focusout on it reschedule the same
//     delayed close — tooltip.js needs none of this, since its content is
//     never interactive and never receives focus or a hover of its own.
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
// Leaving the trigger by keyboard rides the same closeDelay grace as
// leaving it by pointer (scheduleHide, not an immediate hide) — a Tab press
// that lands on a focusable child of HoverCardContent must not race the
// popover closing under it. contentOf/triggerOf both key off the trigger,
// so the content-side focusin below cancels this exact timer.
on("focusout", "[data-gsxui-hovercard-trigger]", (_e, t) => scheduleHide(t));

// The content itself can hold real interactive children — entering it (by
// pointer OR by keyboard focus) must cancel any pending close the
// trigger's own pointerout/focusout just scheduled (the pointer/focus is
// still within the hover-card's own hit area, it simply moved from the
// trigger to the content sitting 4px below it), and leaving it (by either
// modality) schedules the same delayed close the trigger uses — including
// the case where focus/pointer leaves the content for somewhere outside
// both the trigger and the content entirely, which is exactly the plain
// scheduleHide → nothing-cancels-it → hide()-fires-after-closeDelay path.
// Both helpers are keyed by trigger throughout (contentOf/triggerOf), so a
// trigger's pending timer is the single source of truth regardless of
// which of the two elements — or which input modality — is currently
// active.
on("pointerover", "[data-gsxui-hovercard-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) clearTimer(trigger);
});
on("pointerout", "[data-gsxui-hovercard-content]", (e, content) => {
  // Moving between two child elements still inside the content fires
  // pointerout/pointerover on each boundary crossed even though the
  // pointer never left the content's own hit area — same guard as
  // dropdown.js's own content pointerout handler, otherwise every internal
  // move would needlessly churn a schedule/clear pair.
  if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget)) return;
  const trigger = triggerOf(content);
  if (trigger) scheduleHide(trigger);
});
on("focusin", "[data-gsxui-hovercard-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) clearTimer(trigger);
});
on("focusout", "[data-gsxui-hovercard-content]", (_e, content) => {
  const trigger = triggerOf(content);
  if (trigger) scheduleHide(trigger);
});
