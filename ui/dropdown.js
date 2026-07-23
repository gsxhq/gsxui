// Dropdown behavior on the native popover API: top layer, light dismiss and
// Esc are free. JS adds anchored positioning, state/aria sync, arrow-key
// roving focus, and close-on-select. toggle doesn't bubble — capture.
import { on, emit } from "./gsxui.js";

const contentOf = (el) =>
  el.closest("[data-gsxui-dropdown]")?.querySelector("[data-gsxui-dropdown-content]");

// A pointerdown on the trigger records whether the menu was open at that
// instant: popover="auto" light-dismisses on outside pointerdown (the
// trigger is outside the content), so by click time the popover may already
// be closed and a bare toggle would wrongly reopen it.
on("pointerdown", "[data-gsxui-dropdown-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (content) trigger.dataset.gsxuiWasOpen = content.matches(":popover-open") ? "true" : "false";
});

on("click", "[data-gsxui-dropdown-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (!content) return;
  const wasOpen = trigger.dataset.gsxuiWasOpen === "true";
  delete trigger.dataset.gsxuiWasOpen;
  if (wasOpen) {
    // If light dismiss didn't fire (e.g. a caller overrode popover="manual"),
    // converge on the actual state instead of assuming it closed.
    if (content.matches(":popover-open")) content.hidePopover();
    return;
  }
  if (content.matches(":popover-open")) {     // keyboard activation close path
    content.hidePopover();
    return;
  }
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.inset = "auto";
  content.style.left = `${r.left}px`;
  content.style.top = `${r.bottom + 4}px`;
  content.showPopover();
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
    if (open) {
      // clear only on open — clearing on close races the trigger-click task that needs to read the flag
      delete trigger?.dataset.gsxuiWasOpen;
      content.querySelector('[role="menuitem"]:not([aria-disabled])')?.focus();
    }
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
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  const content = item.closest("[data-gsxui-dropdown-content]");
  emit(item, "gsxui:select");
  content?.hidePopover();
});

// Hover highlight IS focus (Radix's roving-focus-follows-pointer): the
// shadcn item classes style focus: only, so pointer hover must move focus
// onto the item for the highlight to appear.
on("pointerover", "[data-gsxui-dropdown-item]", (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  item.focus();
});

// Leaving the menu entirely clears the item highlight by parking focus on
// the content (tabindex="-1") — not body, so arrow keys keep working.
on("pointerout", "[data-gsxui-dropdown-content]", (e, content) => {
  if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget)) return;
  if (content.contains(document.activeElement)) content.focus();
});
