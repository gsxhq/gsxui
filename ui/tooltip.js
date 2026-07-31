// Tooltip: pointerover/out + focusin/out delegation (these bubble), 300ms
// open delay, manual popover so hover can't light-dismiss it.
import { on, emit, position } from "./gsxui.js";

const timers = new WeakMap();
const contentOf = (el) =>
  el
    .closest("[data-gsxui-slot-tooltip]")
    ?.querySelector("[data-gsxui-slot-tooltip-content]");

function show(trigger) {
  const content = contentOf(trigger);
  if (!content || content.matches(":popover-open")) return;
  content.showPopover();
  // Position AFTER showing (hidden popovers have no box) and never via
  // transform: the animate-in enter keyframes animate transform, so a
  // positioning translate would be overridden for the animation's duration
  // — the tooltip would enter at the untranslated spot and snap. Above the
  // trigger, centered (Radix Tooltip's side=top default), 6px sideOffset
  // (arrow clearance). Flip/shift/clamp + scroll/resize tracking: see
  // gsxui.js.
  position(content, trigger, { side: "top", align: "center", sideOffset: 6 });
  content.dataset.state = "open";
  emit(content, "gsxui:open");
}

function hide(trigger) {
  clearTimeout(timers.get(trigger));
  timers.delete(trigger);
  const content = contentOf(trigger);
  if (!content || !content.matches(":popover-open")) return;
  content.hidePopover();
  content.dataset.state = "closed";
  emit(content, "gsxui:close");
}

// A tooltip attaches to TooltipTrigger — identified by its slot marker — or
// to any element opted in with data-gsxui-tooltip-trigger, the idiom for
// arbitrary triggers (Sidebar's menu button uses it). Component markup does
// not stamp the role hook: for the family trigger the role is implied by identity.
const TRIGGER_SELECTOR =
  "[data-gsxui-tooltip-trigger],[data-gsxui-slot-tooltip-trigger]";

on("pointerover", TRIGGER_SELECTOR, (_e, t) => {
  clearTimeout(timers.get(t));
  timers.set(
    t,
    setTimeout(() => show(t), 300),
  );
});
on("pointerout", TRIGGER_SELECTOR, (_e, t) => hide(t));
on("focusin", TRIGGER_SELECTOR, (_e, t) => show(t));
on("focusout", TRIGGER_SELECTOR, (_e, t) => hide(t));
