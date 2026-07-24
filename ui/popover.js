// Popover behavior on the native popover API: top layer, light dismiss and
// Esc are free. Adapted from dropdown.js: same anchored-positioning/state-
// aria-sync/trigger-click-close-guard shape, WITHOUT menu semantics — no
// role="menu" reads, no arrow-key roving focus, no click-on-item close.
// A popover holds arbitrary content (a form, free text), not a list of
// selectable items, so none of dropdown.js's item-focused machinery
// applies. The one positioning change: centered below the trigger, not
// left-aligned (Radix's own Popover default is side=bottom align=center;
// DropdownMenuContent's is align=start). toggle doesn't bubble — capture.
import { on, emit } from "./gsxui.js";

const contentOf = (el) =>
  el.closest("[data-gsxui-popover]")?.querySelector("[data-gsxui-popover-content]");

// A pointerdown on the trigger records whether the popover was open at that
// instant: popover="auto" light-dismisses on outside pointerdown (the
// trigger is outside the content), so by click time the popover may already
// be closed and a bare toggle would wrongly reopen it. Same guard as
// dropdown.js's own (see its ledger MECHANISM entry).
on("pointerdown", "[data-gsxui-popover-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (content) trigger.dataset.gsxuiWasOpen = content.matches(":popover-open") ? "true" : "false";
});

on("click", "[data-gsxui-popover-trigger]", (_e, trigger) => {
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
  content.showPopover();
  // Position numerically AFTER showing (hidden popovers have no box) and
  // never via transform: the animate-in enter keyframes animate transform,
  // so a positioning translate would be overridden for the animation's
  // duration — the popover would enter at the untranslated spot and snap
  // (same rationale as tooltip.js's own comment). offsetWidth is a layout
  // size, unaffected by the in-flight enter scale.
  //
  // Centered below the trigger (Radix's own Popover default is side=bottom
  // align=center — dropdown.js's left-aligned `r.left` is NOT reused here):
  // left is the trigger's horizontal midpoint minus half the content's own
  // width, so the content straddles the trigger's center the way Radix's
  // Floating-UI align=center placement would. No viewport-edge clamping is
  // done, consistent with dropdown.js/tooltip.js's own documented NOTE
  // (hand-rolled position:fixed anchoring, not Floating-UI's collision-
  // aware placement — a stopgap until CSS anchor positioning is Baseline).
  content.style.left = `${r.left + r.width / 2 - content.offsetWidth / 2}px`;
  content.style.top = `${r.bottom + 4}px`;
});

on(
  "toggle",
  "[data-gsxui-popover-content]",
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const trigger = content
      .closest("[data-gsxui-popover]")
      ?.querySelector("[data-gsxui-popover-trigger]");
    trigger?.setAttribute("aria-expanded", open ? "true" : "false");
    // clear only on open — clearing on close races the trigger-click task
    // that needs to read the flag (same ordering rationale as dropdown.js).
    if (open) delete trigger?.dataset.gsxuiWasOpen;
    emit(content, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);
