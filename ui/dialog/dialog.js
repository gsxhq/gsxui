// Dialog behavior. Trigger→content wiring by proximity (closest root), state
// stamped on the <dialog>, gsxui:open/gsxui:close CustomEvents as the API.
//
// State sync and events ride the ToggleEvent (`toggle`, capture-delegated,
// Baseline 2024) rather than the native `close` event: `toggle` fires on
// every open/close path — trigger, Escape, light dismiss, and programmatic
// showModal()/close() from user code — while `close` proved unreliable in
// current Chrome (never delivered, verified empirically; `cancel` alone
// fires on Esc).
import { on, emit } from "../core/gsxui.js";

on("click", "[data-gsxui-dialog-trigger]", (_event, trigger) => {
  const dialog = trigger
    .closest("[data-gsxui-dialog]")
    ?.querySelector("dialog[data-gsxui-dialog-content]");
  if (!dialog || dialog.open) return;
  dialog.showModal();
});

on("click", "[data-gsxui-dialog-close]", (_event, closer) => {
  closer.closest("dialog[data-gsxui-dialog-content]")?.close();
});

// Light dismiss: only a click outside the dialog's own box — i.e. on the
// ::backdrop — dismisses. Clicks in the panel's padding and grid gaps also
// target the <dialog> element itself, so target identity alone is not enough.
on("click", "dialog[data-gsxui-dialog-content]", (event, dialog) => {
  if (event.target !== dialog) return;
  const r = dialog.getBoundingClientRect();
  const inBox =
    event.clientX >= r.left && event.clientX <= r.right &&
    event.clientY >= r.top && event.clientY <= r.bottom;
  if (!inBox) dialog.close();
});

// Single source of truth for state + events, all open/close paths included.
on(
  "toggle",
  "dialog[data-gsxui-dialog-content]",
  (event, dialog) => {
    const open = event.newState === "open";
    dialog.dataset.state = open ? "open" : "closed";
    emit(dialog, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);
