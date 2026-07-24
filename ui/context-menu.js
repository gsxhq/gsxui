// Context-menu behavior on the native popover API: top layer, light dismiss
// and Esc are free. Adapted from dropdown.js: same role="menu" reads,
// arrow-key roving focus, close-on-select, and toggle-driven state/aria
// sync — reused verbatim (see the shared block below, unchanged from
// dropdown.js except for the data-gsxui-contextmenu-* selectors). What's
// different is how the menu OPENS:
//   - dropdown.js opens on a `click` on the trigger BUTTON; context-menu.js
//     opens on a `contextmenu` (right-click) event anywhere inside the
//     trigger AREA (event delegation via closest(), same as every other
//     selector here) and calls preventDefault() to suppress the native
//     browser context menu.
//   - dropdown.js anchors to the trigger's own getBoundingClientRect();
//     context-menu.js positions at the cursor (event.clientX/clientY) —
//     there is no single "trigger rect" to anchor to, the whole area is
//     clickable anywhere inside it.
//   - dropdown.js needs the pointerdown/click wasOpen guard (## dropdown's
//     MECHANISM in docs/jsx-parity.md) because a left-click on the trigger
//     BUTTON is itself an outside pointerdown relative to the content,
//     racing popover="auto"'s own light dismiss. context-menu.js has no
//     equivalent guard: a right-click's pointerdown ALSO counts as an
//     outside pointerdown and already light-dismisses an open menu before
//     the contextmenu event fires, so by the time this handler runs the
//     popover has normally already closed on its own — the defensive
//     hidePopover() below only matters for a contextmenu event dispatched
//     without a preceding pointerdown (e.g. the keyboard Menu key), which
//     still needs to reposition to the new (keyboard-relative) coordinates.
//   - CLAMPING (the one deviation from dropdown.js/tooltip.js/popover.js's
//     own documented no-clamp NOTE): those siblings anchor to a FIXED side
//     of a known trigger element and accept imprecision near viewport edges
//     as a stopgap until CSS anchor positioning is Baseline. A context menu
//     has no such fixed anchor — it opens wherever the cursor was, which
//     can be arbitrarily close to any edge — so an unclamped menu could
//     render partially or fully offscreen on an ordinary right-click near
//     the right or bottom edge, not just an unusual one. Numeric clamping
//     against document.documentElement.clientWidth/clientHeight (the
//     scrollbar-free client area) using the content's own
//     offsetWidth/offsetHeight (read AFTER showPopover() — a hidden popover
//     has no layout box) closes that gap; see docs/jsx-parity.md's ##
//     context-menu ledger entry for the full ADAPT writeup.
// No static data-side is ever stamped — see context-menu.gsx's own doc
// comment. toggle doesn't bubble — capture.
import { on, emit } from "./gsxui.js";

const contentOf = (el) =>
  el.closest("[data-gsxui-contextmenu]")?.querySelector("[data-gsxui-contextmenu-content]");

on("contextmenu", "[data-gsxui-contextmenu-trigger]", (e, trigger) => {
  const content = contentOf(trigger);
  if (!content) return;
  e.preventDefault();
  // Defensive: normally already closed by the right-click's own pointerdown
  // (light dismiss) before this event fires — see the header comment above.
  if (content.matches(":popover-open")) content.hidePopover();
  const openAt = () => {
    content.style.position = "fixed";
    content.style.inset = "auto";
    content.showPopover();
    // Position numerically AFTER showing (hidden popovers have no box) and
    // never via transform: the animate-in enter keyframes animate transform,
    // so a positioning translate would be overridden for the animation's
    // duration (same rationale as dropdown.js/popover.js's own comment).
    // Clamp to the viewport edges (the ADAPT from the siblings' no-clamp
    // precedent, see the header comment above) so a right-click near the
    // right/bottom edge doesn't spawn an offscreen menu. clientWidth/Height,
    // not innerWidth/Height: the inner metrics include classic-scrollbar
    // gutters, and clamping against them tucks the menu's edge under the
    // scrollbar (found in the Task 7 browser pass on a real-scrollbar
    // window).
    const maxLeft = document.documentElement.clientWidth - content.offsetWidth;
    const maxTop = document.documentElement.clientHeight - content.offsetHeight;
    content.style.left = `${Math.max(0, Math.min(e.clientX, maxLeft))}px`;
    content.style.top = `${Math.max(0, Math.min(e.clientY, maxTop))}px`;
  };
  // On macOS the contextmenu event fires at mousedown time, so a button is
  // still held here (e.buttons != 0) and the gesture's own pointerup is
  // still coming. popover="auto"'s light dismiss pairs that pointerup with
  // its pointerdown — both outside a popover that didn't exist yet — and
  // hides whatever we show in between: showing now makes the menu flash
  // once and vanish (found in real-mouse verification; synthetic occluded-
  // tab passes can't reproduce it, UA light dismiss is disabled there).
  // Defer past the gesture: its pointerup, then a task so the UA's light-
  // dismiss processing for that event sees nothing open. On Windows/Linux
  // contextmenu fires after the release (e.buttons == 0), and the keyboard
  // Menu key has no pointer gesture at all — both show immediately.
  if (e.buttons) {
    addEventListener("pointerup", () => setTimeout(openAt), { once: true });
  } else {
    openAt();
  }
});

// Everything below is dropdown.js's own menu-semantics block, unchanged
// except for the data-gsxui-contextmenu-* selectors — see dropdown.js for
// the full rationale of each handler.

on(
  "toggle",
  "[data-gsxui-contextmenu-content]",
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    if (open) {
      content.querySelector('[role="menuitem"]:not([aria-disabled])')?.focus();
    }
    emit(content, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);

on("keydown", "[data-gsxui-contextmenu-content]", (e, content) => {
  // Items are <div role="menuitem">, not buttons — Enter/Space produce no
  // native click, so menu-pattern activation is synthesized here.
  if (e.key === "Enter" || e.key === " ") {
    const item = e.target.closest("[data-gsxui-contextmenu-item]");
    if (item) {
      e.preventDefault();
      item.click();
    }
    return;
  }
  const dir = { ArrowDown: 1, ArrowUp: -1 }[e.key];
  if (!dir) return;
  const items = [...content.querySelectorAll('[role="menuitem"]:not([aria-disabled])')];
  const i = items.indexOf(document.activeElement);
  items[(i + dir + items.length) % items.length]?.focus();
  e.preventDefault();
});

on("click", "[data-gsxui-contextmenu-item]", (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  const content = item.closest("[data-gsxui-contextmenu-content]");
  emit(item, "gsxui:select");
  content?.hidePopover();
});

// Hover highlight IS focus (Radix's roving-focus-follows-pointer): the
// shadcn item classes style focus: only, so pointer hover must move focus
// onto the item for the highlight to appear.
on("pointerover", "[data-gsxui-contextmenu-item]", (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  item.focus();
});

// Leaving the menu entirely clears the item highlight by parking focus on
// the content (tabindex="-1") — not body, so arrow keys keep working.
on("pointerout", "[data-gsxui-contextmenu-content]", (e, content) => {
  if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget)) return;
  if (content.contains(document.activeElement)) content.focus();
});
