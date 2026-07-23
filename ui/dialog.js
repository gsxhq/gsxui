// Dialog behavior. Trigger→content wiring by proximity (closest root), state
// stamped on the <dialog>, gsxui:open/gsxui:close CustomEvents as the API.
//
// State sync and events ride the ToggleEvent (`toggle`, capture-delegated,
// Baseline 2024) rather than the native `close` event: `toggle` fires on
// every open/close path — trigger, Escape, light dismiss, and programmatic
// showModal()/close() from user code — while `close` proved unreliable in
// current Chrome (never delivered, verified empirically; `cancel` alone
// fires on Esc).
//
// Accessibility: DialogTrigger server-renders aria-haspopup and the initial
// aria-expanded; ids for aria-labelledby/-describedby/-controls are ensured
// lazily here (authored ids and aria-* attributes are always respected).
import { on, emit } from "./gsxui.js";

let uid = 0;
const ensureId = (el, prefix) => el.id || (el.id = `gsxui-${prefix}-${++uid}`);

// Idempotent: name/describe the dialog and point triggers at it.
function wireA11y(root, dialog) {
  const title = dialog.querySelector('[data-slot="dialog-title"]');
  const desc = dialog.querySelector('[data-slot="dialog-description"]');
  if (title && !dialog.hasAttribute("aria-labelledby"))
    dialog.setAttribute("aria-labelledby", ensureId(title, "title"));
  if (desc && !dialog.hasAttribute("aria-describedby"))
    dialog.setAttribute("aria-describedby", ensureId(desc, "desc"));
  for (const t of root.querySelectorAll("[data-gsxui-dialog-trigger]"))
    if (!t.hasAttribute("aria-controls"))
      t.setAttribute("aria-controls", ensureId(dialog, "dialog"));
}

// Animated close: stamp the closed state first so data-[state=closed] exit
// animations start, then release the top layer when they finish. With no
// active animations (buildless page, prefers-reduced-motion) this closes
// immediately.
function requestClose(dialog) {
  if (!dialog.open) return;
  dialog.dataset.state = "closed";
  const anims = dialog.getAnimations({ subtree: true });
  if (!anims.length) {
    dialog.close();
    return;
  }
  Promise.allSettled(anims.map((a) => a.finished)).then(() => dialog.close());
}

const rootOf = (el) => el.closest("[data-gsxui-dialog]");

on("click", "[data-gsxui-dialog-trigger]", (_event, trigger) => {
  const root = rootOf(trigger);
  const dialog = root?.querySelector("dialog[data-gsxui-dialog-content]");
  if (!dialog || dialog.open) return;
  wireA11y(root, dialog);
  dialog.showModal();
});

on("click", "[data-gsxui-dialog-close]", (_event, closer) => {
  const dialog = closer.closest("dialog[data-gsxui-dialog-content]");
  if (dialog) requestClose(dialog);
});

// Light dismiss: only a click outside the dialog's own box — i.e. on the
// ::backdrop — dismisses. Clicks in the panel's padding and grid gaps also
// target the <dialog> element itself, so target identity alone is not enough.
//
// data-gsxui-dialog-static (ui/alert-dialog.gsx's AlertDialogContent) opts
// a content element out of this path entirely: Radix's own AlertDialog
// ignores outside clicks while Esc still closes it, and this early return
// reproduces exactly that — Esc/cancel below and the close-button/toggle
// handlers are untouched, only the backdrop-click path is skipped.
on("click", "dialog[data-gsxui-dialog-content]", (event, dialog) => {
  if (event.target !== dialog) return;
  if (dialog.hasAttribute("data-gsxui-dialog-static")) return;
  const r = dialog.getBoundingClientRect();
  const inBox =
    event.clientX >= r.left && event.clientX <= r.right &&
    event.clientY >= r.top && event.clientY <= r.bottom;
  if (!inBox) requestClose(dialog);
});

// Esc: intercept the native cancel so the exit animation can run; the
// requestClose path ends in close(), which still fires toggle below.
on(
  "cancel",
  "dialog[data-gsxui-dialog-content]",
  (event, dialog) => {
    event.preventDefault();
    requestClose(dialog);
  },
  { capture: true },
);

// Single source of truth for state, events, and aria-expanded sync — covers
// programmatic showModal()/close() too.
on(
  "toggle",
  "dialog[data-gsxui-dialog-content]",
  (event, dialog) => {
    const open = event.newState === "open";
    dialog.dataset.state = open ? "open" : "closed";
    const root = rootOf(dialog);
    if (root) {
      for (const t of root.querySelectorAll("[data-gsxui-dialog-trigger]"))
        t.setAttribute("aria-expanded", open ? "true" : "false");
      if (open) wireA11y(root, dialog);
    }
    emit(dialog, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);
