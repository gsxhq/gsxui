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

// Light dismiss: a click on the <dialog> itself (not its children) is a
// click on the backdrop area.
on("click", "dialog[data-gsxui-dialog-content]", (event, dialog) => {
  if (event.target === dialog) dialog.close();
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
