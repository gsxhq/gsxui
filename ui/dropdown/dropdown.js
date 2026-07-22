// Dropdown behavior on the native popover API: top layer, light dismiss and
// Esc are free. JS adds anchored positioning, state/aria sync, arrow-key
// roving focus, and close-on-select. toggle doesn't bubble — capture.
import { on, emit } from "../core/gsxui.js";

const contentOf = (el) =>
  el.closest("[data-gsxui-dropdown]")?.querySelector("[data-gsxui-dropdown-content]");

on("click", "[data-gsxui-dropdown-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (!content) return;
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.left = `${r.left}px`;
  content.style.top = `${r.bottom + 4}px`;
  content.togglePopover();
});

on(
  "toggle",
  "[data-gsxui-dropdown-content]",
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const trigger = content
      .closest("[data-gsxui-dropdown]")
      ?.querySelector("[data-gsxui-dropdown-trigger]");
    trigger?.setAttribute("aria-expanded", open ? "true" : "false");
    if (open) content.querySelector('[role="menuitem"]:not([aria-disabled])')?.focus();
    emit(content, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);

on("keydown", "[data-gsxui-dropdown-content]", (e, content) => {
  const dir = { ArrowDown: 1, ArrowUp: -1 }[e.key];
  if (!dir) return;
  const items = [...content.querySelectorAll('[role="menuitem"]:not([aria-disabled])')];
  const i = items.indexOf(document.activeElement);
  items[(i + dir + items.length) % items.length]?.focus();
  e.preventDefault();
});

on("click", "[data-gsxui-dropdown-item]", (_e, item) => {
  const content = item.closest("[data-gsxui-dropdown-content]");
  emit(item, "gsxui:select");
  content?.hidePopover();
});
