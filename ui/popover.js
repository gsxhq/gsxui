// Popover behavior on the native popover API: top layer, light dismiss and
// Esc are free. Adapted from dropdown.js: same anchored-positioning/state-
// aria-sync/trigger-click-close-guard shape, WITHOUT menu semantics — no
// role="menu" reads, no arrow-key roving focus, no click-on-item close.
// A popover holds arbitrary content (a form, free text), not a list of
// selectable items, so none of dropdown.js's item-focused machinery
// applies. The one positioning change: centered below the trigger, not
// left-aligned (Radix's own Popover default is side=bottom align=center;
// DropdownMenuContent's is align=start). toggle doesn't bubble — capture.
import { on, emit, position } from "./gsxui.js";

const contentOf = (el) =>
  el
    .closest("[data-gsxui-slot-popover]")
    ?.querySelector("[data-gsxui-slot-popover-content]");

// A pointerdown on the trigger records whether the popover was open at that
// instant: popover="auto" light-dismisses on outside pointerdown (the
// trigger is outside the content), so by click time the popover may already
// be closed and a bare toggle would wrongly reopen it. Same guard as
// dropdown.js's own (see its ledger MECHANISM entry).
// A popover opens from PopoverTrigger, identified by its slot marker — the
// slot attribute carries the behavior contract along with the styling role.
const TRIGGER_SELECTOR = "[data-gsxui-slot-popover-trigger]";

on("pointerdown", TRIGGER_SELECTOR, (_e, trigger) => {
  const content = contentOf(trigger);
  if (content)
    trigger.dataset.gsxuiWasOpen = content.matches(":popover-open")
      ? "true"
      : "false";
});

on("click", TRIGGER_SELECTOR, (_e, trigger) => {
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
  if (content.matches(":popover-open")) {
    // keyboard activation close path
    content.hidePopover();
    return;
  }
  // Stamp open BEFORE showing — the toggle event that also stamps it is a
  // queued task, and a paint in the gap flashes one closed-state frame
  // before the enter animation restarts (see dropdown.js's comment).
  content.dataset.state = "open";
  content.showPopover();
  // Position AFTER showing (hidden popovers have no box) via the shared
  // engine — never via translate/transform: the discrete-transition
  // enter/exit animates the `translate` and `scale` properties (see
  // popover.gsx's ADAPT comment), so a positioning translate would be
  // fought by the transition in both directions.
  //
  // Centered below the trigger (Radix's own Popover default is side=bottom
  // align=center — dropdown.js's align=start is NOT reused here), 4px
  // sideOffset. Flip/shift/clamp + scroll/resize tracking: see gsxui.js.
  position(content, trigger, { side: "bottom", align: "center", sideOffset: 4 });
});

// Focus management, the Radix FocusScope trio popover.js was missing
// (found in real-keyboard verification — with none of it, Tab from an open
// popover walked the page underneath):
//   - on open, focus the first tabbable descendant, falling back to the
//     content itself (tabindex="-1" in popover.gsx exists for this) —
//     FocusScope's own mount behavior, and what makes a form popover
//     immediately typeable;
//   - while open, Tab/Shift+Tab wrap within the content (the keydown
//     handler below);
//   - on close, return focus to the trigger — but only when focus is still
//     inside the content (Esc, or a wrapped Tab). An outside click has
//     already moved focus to the clicked element by the time light dismiss
//     hides the popover; stealing it back would fight the user's click.
//     This runs on beforetoggle, not toggle: toggle is queued as a task, by
//     which point the UA has already parked focus on <body> and "was focus
//     inside?" is unanswerable.
const TABBABLE =
  'a[href], button:not([disabled]), input:not([disabled]):not([type="hidden"]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

on(
  "beforetoggle",
  "[data-gsxui-slot-popover-content]",
  (e, content) => {
    if (e.newState !== "closed" || !content.contains(document.activeElement))
      return;
    content
      .closest("[data-gsxui-slot-popover]")
      ?.querySelector(TRIGGER_SELECTOR)
      ?.focus();
  },
  { capture: true },
);

on("keydown", "[data-gsxui-slot-popover-content]", (e, content) => {
  if (e.key !== "Tab") return;
  const items = [...content.querySelectorAll(TABBABLE)];
  if (!items.length) return;
  const edge = e.shiftKey ? items[0] : items[items.length - 1];
  if (document.activeElement === edge || document.activeElement === content) {
    e.preventDefault();
    (e.shiftKey ? items[items.length - 1] : items[0]).focus();
  }
});

on(
  "toggle",
  "[data-gsxui-slot-popover-content]",
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const trigger = content
      .closest("[data-gsxui-slot-popover]")
      ?.querySelector(TRIGGER_SELECTOR);
    trigger?.setAttribute("aria-expanded", open ? "true" : "false");
    // clear only on open — clearing on close races the trigger-click task
    // that needs to read the flag (same ordering rationale as dropdown.js).
    if (open) {
      delete trigger?.dataset.gsxuiWasOpen;
      (content.querySelector(TABBABLE) ?? content).focus();
    }
    emit(content, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);
