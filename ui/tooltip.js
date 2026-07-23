// Tooltip: pointerover/out + focusin/out delegation (these bubble), 300ms
// open delay, manual popover so hover can't light-dismiss it.
import { on, emit } from "./gsxui.js";

const timers = new WeakMap();
const contentOf = (el) =>
  el.closest("[data-gsxui-tooltip]")?.querySelector("[data-gsxui-tooltip-content]");

function show(trigger) {
  const content = contentOf(trigger);
  if (!content || content.matches(":popover-open")) return;
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.inset = "auto";
  content.showPopover();
  // Position numerically AFTER showing (hidden popovers have no box) and
  // never via transform: the animate-in enter keyframes animate transform,
  // so a positioning translate would be overridden for the animation's
  // duration — the tooltip would enter at the untranslated spot and snap.
  // offsetWidth/Height are layout sizes, unaffected by the in-flight
  // enter scale.
  content.style.left = `${r.left + r.width / 2 - content.offsetWidth / 2}px`;
  content.style.top = `${r.top - 6 - content.offsetHeight}px`;
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

on("pointerover", "[data-gsxui-tooltip-trigger]", (_e, t) => {
  clearTimeout(timers.get(t));
  timers.set(t, setTimeout(() => show(t), 300));
});
on("pointerout", "[data-gsxui-tooltip-trigger]", (_e, t) => hide(t));
on("focusin", "[data-gsxui-tooltip-trigger]", (_e, t) => show(t));
on("focusout", "[data-gsxui-tooltip-trigger]", (_e, t) => hide(t));
