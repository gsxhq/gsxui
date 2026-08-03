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

const DIALOG = "dialog[data-gsxui-slot-dialog-content]";

function generatedId(prefix) {
  for (;;) {
    const words = crypto.getRandomValues(new Uint32Array(4));
    const suffix = [...words]
      .map((word) => word.toString(16).padStart(8, "0"))
      .join("");
    const id = `gsxui-${prefix}-${suffix}`;
    if (!document.getElementById(id)) return id;
  }
}

function ensureId(element, prefix) {
  if (!element.id) {
    element.id = generatedId(prefix);
    if (element.matches(DIALOG)) element.setAttribute("data-gsxui-dialog-generated-id", "");
  }
  return element.id;
}

const rootOf = (element) => element.closest("[data-gsxui-slot-dialog]");

function owned(root, selector) {
  return [...root.querySelectorAll(selector)].filter(
    (element) => rootOf(element) === root,
  );
}

function dialogOf(root) {
  return owned(root, DIALOG)[0];
}

function authoredInvokers(dialog) {
  if (!dialog.id || dialog.hasAttribute("data-gsxui-dialog-generated-id")) return [];
  const target = CSS.escape(dialog.id);
  return [...document.querySelectorAll(`[commandfor="${target}"][command="show-modal"]`)];
}

// A dialog opens from its family Trigger components, identified by their
// slot markers — the slot attribute carries the behavior contract along
// with the styling role.
const TRIGGER_SELECTOR =
  "[data-gsxui-slot-dialog-trigger],[data-gsxui-slot-alert-dialog-trigger],[data-gsxui-slot-drawer-trigger],[data-gsxui-slot-sheet-trigger]";

// Idempotent: name/describe the dialog and point triggers at it.
function wireA11y(root, dialog) {
  const title = owned(root, "[data-gsxui-slot-dialog-title]")[0];
  const desc = owned(root, "[data-gsxui-slot-dialog-description]")[0];
  if (title && !dialog.hasAttribute("aria-labelledby"))
    dialog.setAttribute("aria-labelledby", ensureId(title, "title"));
  if (desc && !dialog.hasAttribute("aria-describedby"))
    dialog.setAttribute("aria-describedby", ensureId(desc, "desc"));
  for (const t of owned(root, TRIGGER_SELECTOR))
    if (!t.hasAttribute("aria-controls"))
      t.setAttribute("aria-controls", ensureId(dialog, "dialog"));
}

function performOpen(dialog) {
  const root = rootOf(dialog);
  if (root) wireA11y(root, dialog);
  // Stamping open also aborts an in-flight close before its animation settles.
  dialog.dataset.state = "open";
  if (!dialog.open) dialog.showModal();
}

async function performClose(dialog) {
  if (!dialog.open || dialog.dataset.state === "closed") return;
  dialog.dataset.state = "closed";
  while (dialog.open && dialog.dataset.state === "closed") {
    const animations = dialog.getAnimations().filter((animation) => {
      const endTime = animation.effect?.getComputedTiming().endTime;
      return (
        animation.playState !== "finished" &&
        typeof endTime === "number" &&
        Number.isFinite(endTime)
      );
    });
    if (animations.length === 0) break;
    await Promise.allSettled(animations.map((animation) => animation.finished));
  }
  if (dialog.open && dialog.dataset.state === "closed") dialog.close();
}

function request(dialog, type, detail) {
  return emit(dialog, type, detail);
}

on("gsxui:request-open", DIALOG, (event, dialog) => {
  queueMicrotask(() => {
    if (!event.defaultPrevented) performOpen(dialog);
  });
});

on("gsxui:request-close", DIALOG, (event, dialog) => {
  queueMicrotask(() => {
    if (!event.defaultPrevented) void performClose(dialog);
  });
});

on("click", TRIGGER_SELECTOR, (_event, trigger) => {
  const root = rootOf(trigger);
  const dialog = root && dialogOf(root);
  if (dialog) request(dialog, "gsxui:request-open", { reason: "trigger" });
});

on("click", "[data-gsxui-dialog-close]", (_event, closer) => {
  const root = rootOf(closer);
  const dialog = root && dialogOf(root);
  if (dialog) request(dialog, "gsxui:request-close", { reason: "close-button" });
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
on("click", "dialog[data-gsxui-slot-dialog-content]", (event, dialog) => {
  if (event.target !== dialog) return;
  if (dialog.hasAttribute("data-gsxui-dialog-static")) return;
  const r = dialog.getBoundingClientRect();
  const inBox =
    event.clientX >= r.left && event.clientX <= r.right &&
    event.clientY >= r.top && event.clientY <= r.bottom;
  if (!inBox) request(dialog, "gsxui:request-close", { reason: "backdrop" });
});

// Esc: intercept the native cancel so the exit animation can run; the
// The request-close path ends in close(), which still fires toggle below.
on(
  "cancel",
  "dialog[data-gsxui-slot-dialog-content]",
  (event, dialog) => {
    event.preventDefault();
    request(dialog, "gsxui:request-close", { reason: "cancel" });
  },
  { capture: true },
);

// Native dialog methods and Invoker Commands transition asynchronously: stamp
// the impending state before the browser can paint, then leave `toggle` to
// confirm it and notify observers.
on(
  "beforetoggle",
  "dialog[data-gsxui-slot-dialog-content]",
  (event, dialog) => {
    const open = event.newState === "open";
    const root = rootOf(dialog);
    if (open && root) wireA11y(root, dialog);
    dialog.dataset.state = open ? "open" : "closed";
  },
  { capture: true },
);

// Single source of truth for state, events, and aria-expanded sync — covers
// programmatic showModal()/close() too.
on(
  "toggle",
  "dialog[data-gsxui-slot-dialog-content]",
  (event, dialog) => {
    const open = event.newState === "open";
    dialog.dataset.state = open ? "open" : "closed";
    const root = rootOf(dialog);
    if (root) {
      for (const t of owned(root, TRIGGER_SELECTOR))
        t.setAttribute("aria-expanded", open ? "true" : "false");
      if (open) wireA11y(root, dialog);
    }
    for (const invoker of authoredInvokers(dialog))
      invoker.setAttribute("aria-expanded", open ? "true" : "false");
    emit(dialog, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);
