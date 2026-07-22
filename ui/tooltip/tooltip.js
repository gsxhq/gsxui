// Tooltip: pointerover/out + focusin/out delegation (these bubble), 300ms
// open delay, manual popover so hover can't light-dismiss it.
import { on, emit } from "../core/gsxui.js";

let timer;
const contentOf = (el) =>
  el.closest("[data-gsxui-tooltip]")?.querySelector("[data-gsxui-tooltip-content]");

function show(trigger) {
  const content = contentOf(trigger);
  if (!content || content.matches(":popover-open")) return;
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.left = `${r.left + r.width / 2}px`;
  content.style.top = `${r.top - 6}px`;
  content.style.transform = "translate(-50%, -100%)";
  content.showPopover();
  content.dataset.state = "open";
  emit(content, "gsxui:open");
}

function hide(trigger) {
  clearTimeout(timer);
  const content = contentOf(trigger);
  if (!content || !content.matches(":popover-open")) return;
  content.hidePopover();
  content.dataset.state = "closed";
  emit(content, "gsxui:close");
}

on("pointerover", "[data-gsxui-tooltip-trigger]", (_e, t) => {
  clearTimeout(timer);
  timer = setTimeout(() => show(t), 300);
});
on("pointerout", "[data-gsxui-tooltip-trigger]", (_e, t) => hide(t));
on("focusin", "[data-gsxui-tooltip-trigger]", (_e, t) => show(t));
on("focusout", "[data-gsxui-tooltip-trigger]", (_e, t) => hide(t));
