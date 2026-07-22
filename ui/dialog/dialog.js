// Dialog behavior. Trigger→content wiring by proximity (closest root), state
// stamped on the <dialog>, gsxui:open/gsxui:close CustomEvents as the API.
import { on, emit } from "../core/gsxui.js";

on("click", "[data-gsxui-dialog-trigger]", (_event, trigger) => {
  const dialog = trigger
    .closest("[data-gsxui-dialog]")
    ?.querySelector("dialog[data-gsxui-dialog-content]");
  if (!dialog || dialog.open) return;
  dialog.showModal();
  dialog.dataset.state = "open";
  emit(dialog, "gsxui:open");
});

on("click", "[data-gsxui-dialog-close]", (_event, closer) => {
  closer.closest("dialog[data-gsxui-dialog-content]")?.close();
});

// Light dismiss: a click on the <dialog> itself (not its children) is a
// click on the backdrop area.
on("click", "dialog[data-gsxui-dialog-content]", (event, dialog) => {
  if (event.target === dialog) dialog.close();
});

// Native close event (covers Esc and every close() path above). It does not
// bubble — delegate in the capture phase.
on(
  "close",
  "dialog[data-gsxui-dialog-content]",
  (_event, dialog) => {
    dialog.dataset.state = "closed";
    emit(dialog, "gsxui:close");
  },
  { capture: true },
);
