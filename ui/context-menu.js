// Context-menu behavior on the native popover API: top layer, light dismiss
// and Esc are free. Adapted from dropdown.js: same role="menu" reads,
// arrow-key roving focus, close-on-select, and toggle-driven state/aria
// sync — reused (mostly verbatim, see the shared block below) for the
// shared-items parts (CheckboxItem/RadioItem/Sub/SubTrigger/SubContent)
// this task adds. Their `data-gsxui-slot-context-menu-*` hooks are namespaced to
// THIS component, not shared with dropdown.js's own `data-gsxui-slot-dropdown-menu-*`
// equivalents — see this file's own MECHANISM note further down and
// docs/jsx-parity.md's ## dropdown ledger entry for why a shared
// `data-gsxui-menu-*` prefix across modules was tried and reverted (fix
// round 1: gsxui.js's delegation registry dispatches to every handler whose
// selector matches, regardless of which module registered it — an
// unnamespaced shared selector fires BOTH modules' handlers on one click).
// What's different is how the menu OPENS:
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
//   - CLAMPING (context-menu was first; the siblings adopted the same
//     positioning engine in gsxui.js (position(), which clamps as its last resort) after the theme
//     editor rendered a picker popover flush offscreen-left). A context menu
//     has no such fixed anchor — it opens wherever the cursor was, which
//     can be arbitrarily close to any edge — so an unclamped menu could
//     render partially or fully offscreen on an ordinary right-click near
//     the right or bottom edge, not just an unusual one. Numeric clamping
//     against document.documentElement.clientWidth/clientHeight (the
//     scrollbar-free client area) using the content's own
//     offsetWidth/offsetHeight (read AFTER showPopover() — a hidden popover
//     has no layout box) closes that gap; see docs/jsx-parity.md's ##
//     context-menu ledger entry for the full ADAPT writeup.
//   - ContextMenuTrigger carries no aria-haspopup/aria-expanded of its own
//     (it's a passive AREA, not an interactive control, unlike
//     DropdownMenuTrigger) — the toggle handler below correspondingly has
//     no top-level trigger aria to sync, only a SubTrigger's.
//
// SUBMENUS (ContextMenuSub/SubTrigger/SubContent) and checkbox/radio items:
// see ui/context-menu.gsx's file header for the DOM-nesting-not-portalling
// rationale (spec §1) and the server-rendered-checked MECHANISM — same
// mechanism as dropdown.js's own, not re-derived here. A submenu's
// trigger/content pair reuses hover-card.js's scheduleShow/scheduleHide
// timer shape for its pointer-leave grace period (SUB_CLOSE_DELAY, same
// value as dropdown.js's own) — OPEN is immediate on pointerover/
// ArrowRight/click, unlike hover-card's own delayed open.
//
// No static data-side is ever stamped on the top-level content — see
// context-menu.gsx's own doc comment. toggle doesn't bubble — capture.
import { on, emit, position, isRTL } from "./gsxui.js";

const contentOf = (el) =>
  el
    .closest("[data-gsxui-slot-context-menu]")
    ?.querySelector("[data-gsxui-slot-context-menu-content]");

// Any popover=auto menu surface — the top-level content or a nested
// submenu's own content — wherever a keydown/toggle handler must resolve
// "which content level is this happening at," picking the NEAREST one via
// closest().
const CONTENT_SELECTOR =
  "[data-gsxui-slot-context-menu-content],[data-gsxui-slot-context-menu-sub-content]";
const ITEM_SELECTOR =
  "[data-gsxui-slot-context-menu-item]:not([aria-disabled]),[data-gsxui-slot-context-menu-checkbox-item]:not([aria-disabled]),[data-gsxui-slot-context-menu-radio-item]:not([aria-disabled]),[data-gsxui-slot-context-menu-sub-trigger]:not([aria-disabled])";

// Items belonging to THIS content alone. content.querySelectorAll would also
// recurse into a DOM-nested submenu's own items — that nesting is exactly
// what makes submenus work at all (see context-menu.gsx's file header) — so
// a naive query would let a submenu's items leak into its PARENT's roving
// arrow-key list even while the submenu is closed. The closest() check
// excludes any item whose nearest content ancestor isn't this one.
function ownItems(content) {
  return [...content.querySelectorAll(ITEM_SELECTOR)].filter(
    (item) => item.closest(CONTENT_SELECTOR) === content,
  );
}

const subRootOf = (el) => el.closest("[data-gsxui-slot-context-menu-sub]");
const subContentOf = (trigger) =>
  subRootOf(trigger)?.querySelector(
    "[data-gsxui-slot-context-menu-sub-content]",
  );
const subTriggerOf = (content) =>
  subRootOf(content)?.querySelector(
    "[data-gsxui-slot-context-menu-sub-trigger]",
  );

on("contextmenu", "[data-gsxui-slot-context-menu-trigger]", (e, trigger) => {
  const content = contentOf(trigger);
  if (!content) return;
  e.preventDefault();
  // Defensive: normally already closed by the right-click's own pointerdown
  // (light dismiss) before this event fires — see the header comment above.
  if (content.matches(":popover-open")) content.hidePopover();
  const openAt = () => {
    // Stamp open BEFORE showing — the toggle event that also stamps it is a
    // queued task, and a paint in the gap flashes one closed-state frame
    // before the enter animation restarts (see dropdown.js's comment).
    content.dataset.state = "open";
    content.showPopover();
    // Position AFTER showing (hidden popovers have no box) and never via
    // translate/transform: the discrete-transition enter/exit animates the
    // `translate` and `scale` properties (see popover.gsx's ADAPT comment),
    // so a positioning translate would be fought by the transition in both
    // directions.
    //
    // Pointer-anchored: a zero-size static rect at the cursor (Radix's own
    // virtual anchor), flip:false — a context menu keeps its top-left at
    // the cursor and only shift-clamps against the viewport edges (the
    // ADAPT from the siblings' no-clamp precedent, see the header comment
    // above), matching the prior inline clamp. flip:false also leaves
    // data-side alone — no static data-side is ever stamped on this
    // content (see context-menu.gsx). The engine still supplies the
    // scroll/resize reposition tracking; the anchor rect is viewport-fixed
    // client coordinates, so on scroll the menu re-clamps in place rather
    // than tracking page content — Radix's virtual-anchor behaviour.
    const rect = {
      left: e.clientX,
      top: e.clientY,
      right: e.clientX,
      bottom: e.clientY,
      width: 0,
      height: 0,
    };
    position(content, rect, { flip: false });
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

// Everything below is dropdown.js's own menu-semantics block, adapted for
// this component's own `data-gsxui-slot-context-menu-*` selectors throughout —
// including the shared-items parts (CheckboxItem/RadioItem/Sub-family),
// which use `data-gsxui-slot-context-menu-checkbox-item` etc, NOT a prefix shared
// with dropdown.js's own `data-gsxui-slot-dropdown-menu-checkbox-item` etc (fix round
// 1: gsxui.js's single delegation registry, keyed only by `${type}:
// ${capture}`, dispatches to EVERY handler whose selector matches the
// target regardless of which module registered it — ui/index.js imports
// both dropdown.js and context-menu.js, so an identical selector in both
// files double-fires on one event; measured live: a checkbox item's
// data-state/aria-checked net-zeroed back to its starting value and
// gsxui:change fired twice). Duplicated here, not imported, per the task
// brief's own "no shared JS module" constraint (registry.Deps is derived by
// go/parser over the generated .x.go, so a JS-only shared module would be
// invisible to
// vendoring).

on(
  "toggle",
  CONTENT_SELECTOR,
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const isSub = content.matches("[data-gsxui-slot-context-menu-sub-content]");
    // Only a SubTrigger has aria-expanded/data-state to sync — the
    // top-level ContextMenuTrigger is a passive AREA with neither (see this
    // file's own header comment).
    if (isSub) {
      const trigger = subTriggerOf(content);
      trigger?.setAttribute("aria-expanded", open ? "true" : "false");
      if (trigger) trigger.dataset.state = open ? "open" : "closed";
    }
    if (open) {
      // A hover-opened submenu must NOT steal keyboard focus (Radix's own
      // convention, and this port's own "hover highlight IS focus" idiom
      // already covers the mouse case via pointerover on whatever's
      // hovered) — only the top-level content auto-focuses its first item
      // on open. A keyboard-opened submenu (ArrowRight) focuses its own
      // first item explicitly, from the ArrowRight handler itself, not here.
      if (!isSub) ownItems(content)[0]?.focus();
    }
    emit(content, open ? "gsxui:open" : "gsxui:close");
  },
  { capture: true },
);

const ACTIVATABLE_SELECTOR =
  "[data-gsxui-slot-context-menu-item],[data-gsxui-slot-context-menu-checkbox-item],[data-gsxui-slot-context-menu-radio-item],[data-gsxui-slot-context-menu-sub-trigger]";

on("keydown", CONTENT_SELECTOR, (e, content) => {
  // Items are <div role="menuitem...">, not buttons — Enter/Space produce
  // no native click, so menu-pattern activation is synthesized here. For a
  // sub-trigger, Enter/Space opens its submenu and moves focus in, same as
  // ArrowRight below (WAI-ARIA menu convention: Enter on a submenu trigger
  // both activates AND opens).
  if (e.key === "Enter" || e.key === " ") {
    const item = e.target.closest(ACTIVATABLE_SELECTOR);
    if (!item) return;
    e.preventDefault();
    if (item.matches("[data-gsxui-slot-context-menu-sub-trigger]")) {
      openSubAndFocusFirst(item);
      return;
    }
    item.click();
    return;
  }
  // RTL flips "into the menu" / "out of the menu" — mirror image of LTR's
  // ArrowRight-opens/ArrowLeft-closes (WAI-ARIA menu convention read against
  // resolved direction, per isRTL — see gsxui.js).
  const openKey = isRTL(content) ? "ArrowLeft" : "ArrowRight";
  const closeKey = isRTL(content) ? "ArrowRight" : "ArrowLeft";
  if (e.key === openKey) {
    const trigger = e.target.closest(
      "[data-gsxui-slot-context-menu-sub-trigger]",
    );
    if (!trigger) return;
    e.preventDefault();
    openSubAndFocusFirst(trigger);
    return;
  }
  if (
    e.key === closeKey &&
    content.matches("[data-gsxui-slot-context-menu-sub-content]")
  ) {
    e.preventDefault();
    const trigger = subTriggerOf(content);
    content.hidePopover();
    trigger?.focus();
    return;
  }
  const dir = { ArrowDown: 1, ArrowUp: -1 }[e.key];
  if (!dir) return;
  const items = ownItems(content);
  const i = items.indexOf(document.activeElement);
  // When nothing in this content is currently focused (i === -1, e.g. focus
  // is parked on the content itself), (i + dir) % length lands one short of
  // the intended end (ArrowUp: -1 + -1 = -2, wrapping to the SECOND-to-last
  // item, not the last) — special-case the no-selection start explicitly
  // instead of relying on the wraparound arithmetic to do it by accident.
  const next =
    i === -1
      ? dir === 1
        ? items[0]
        : items[items.length - 1]
      : items[(i + dir + items.length) % items.length];
  if (next) {
    // Moving keyboard focus to a different item at this same level closes
    // any OTHER open submenu among this level's own sub-triggers — without
    // this, a submenu opened earlier by hover (mouse stationary, so no
    // pointerout ever fires) would stay visually open while focus moves
    // elsewhere in the list. Mouse-driven moves don't need this: leaving a
    // sub-trigger by an actual pointer movement already fires pointerout on
    // it, which schedules the same close via the grace-period timer below.
    for (const item of items) {
      if (
        item !== next &&
        item.matches("[data-gsxui-slot-context-menu-sub-trigger]") &&
        item.dataset.state === "open"
      ) {
        closeSub(item);
      }
    }
    next.focus();
  }
  e.preventDefault();
});

on("click", "[data-gsxui-slot-context-menu-item]", (_e, item) => {
  if (
    item.getAttribute("aria-disabled") === "true" ||
    "disabled" in item.dataset
  )
    return;
  // Always hide the ROOT content, not item.closest(CONTENT_SELECTOR) — from
  // inside a submenu, that resolves to the SUB-content, closing only the
  // submenu and leaving the root open with the sub-trigger still
  // highlighted. One call on the root is enough: the measured popover
  // stack cascade (see context-menu.gsx's file header) closes any
  // nested-open submenu along with it.
  const content = contentOf(item);
  emit(item, "gsxui:select");
  content?.hidePopover();
});

// A checkbox item flips in place and does NOT close its menu — a deliberate
// ADAPT per the task brief, not a Radix default (Radix's actual CheckboxItem
// onSelect closes unless preventDefault is called; this port simply never
// closes on a checkbox select).
on("click", "[data-gsxui-slot-context-menu-checkbox-item]", (_e, item) => {
  if (
    item.getAttribute("aria-disabled") === "true" ||
    "disabled" in item.dataset
  )
    return;
  const checked = item.dataset.state !== "checked";
  item.dataset.state = checked ? "checked" : "unchecked";
  item.setAttribute("aria-checked", checked ? "true" : "false");
  emit(item, "gsxui:change", { checked, value: item.dataset.value ?? "" });
});

// A radio item sets itself checked, clears every OTHER item in the SAME
// radio group (data-gsxui-slot-context-menu-radio-group scopes the sibling walk — a page
// may have more than one group), and closes the menu like a plain Item.
on("click", "[data-gsxui-slot-context-menu-radio-item]", (_e, item) => {
  if (
    item.getAttribute("aria-disabled") === "true" ||
    "disabled" in item.dataset
  )
    return;
  const group = item.closest("[data-gsxui-slot-context-menu-radio-group]");
  const siblings = group
    ? [...group.querySelectorAll("[data-gsxui-slot-context-menu-radio-item]")]
    : [item];
  for (const sibling of siblings) {
    const isThis = sibling === item;
    sibling.dataset.state = isThis ? "checked" : "unchecked";
    sibling.setAttribute("aria-checked", isThis ? "true" : "false");
  }
  const value = item.dataset.value ?? "";
  if (group) group.dataset.value = value;
  emit(group ?? item, "gsxui:change", { value });
  // Same root-not-nearest-content fix as the plain item click handler above.
  const content = contentOf(item);
  content?.hidePopover();
});

// Hover highlight IS focus (Radix's roving-focus-follows-pointer): the
// shadcn item classes style focus: only, so pointer hover must move focus
// onto the item for the highlight to appear. Covers every item type that
// participates in roving focus, including a sub-trigger (its OWN open-on-
// hover is wired separately below).
const HOVERABLE_SELECTOR =
  "[data-gsxui-slot-context-menu-item],[data-gsxui-slot-context-menu-checkbox-item],[data-gsxui-slot-context-menu-radio-item],[data-gsxui-slot-context-menu-sub-trigger]";
on("pointerover", HOVERABLE_SELECTOR, (_e, item) => {
  if (
    item.getAttribute("aria-disabled") === "true" ||
    "disabled" in item.dataset
  )
    return;
  item.focus();
});

// Leaving a content entirely clears the item highlight by parking focus on
// that content (tabindex="-1") — not body, so arrow keys keep working.
// Moving into a DOM-nested (open) submenu doesn't count as leaving — it's
// still `content.contains(relatedTarget)` true, same guard as hover-card.js.
on("pointerout", CONTENT_SELECTOR, (e, content) => {
  if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget))
    return;
  if (content.contains(document.activeElement)) content.focus();
});

// --- submenus: open/close + the hover-card-shaped pointer-leave grace period

const SUB_CLOSE_DELAY = 100; // hover-card.js's own CLOSE_DELAY — brief: "match the shipped hover-card.js grace-period precedent"
const subTimers = new WeakMap(); // sub-trigger -> pending close setTimeout id

function clearSubTimer(trigger) {
  clearTimeout(subTimers.get(trigger));
  subTimers.delete(trigger);
}

// openSub is idempotent (hover re-entering an already-open sub, or pressing
// ArrowRight twice, must not reposition/re-show it) and returns the content
// element either way, so callers can chain a focus-first-item step.
function openSub(trigger) {
  clearSubTimer(trigger);
  const content = subContentOf(trigger);
  if (!content) return null;
  if (content.matches(":popover-open")) return content;
  // Stamp open BEFORE showing — same flash-avoidance rule as the top-level
  // opener's own contextmenu handler (a queued toggle event can otherwise
  // leave one frame painted in the stale closed state).
  trigger.dataset.state = "open";
  trigger.setAttribute("aria-expanded", "true");
  content.dataset.state = "open";
  content.showPopover();
  // To the right of its trigger; alignOffset -4 offsets the content's own
  // p-1 so the first item's text roughly aligns with the trigger's. Unlike
  // the pointer-anchored top-level content, a SubContent is an ordinary
  // trigger-anchored surface with an authored data-side="right", so it
  // gets the full flip/shift/clamp treatment like dropdown-menu's own sub.
  position(content, trigger, { side: "right", alignOffset: -4 });
  return content;
}

function closeSub(trigger) {
  clearSubTimer(trigger);
  const content = subContentOf(trigger);
  if (content?.matches(":popover-open")) content.hidePopover();
}

function scheduleCloseSub(trigger) {
  clearSubTimer(trigger);
  subTimers.set(
    trigger,
    setTimeout(() => closeSub(trigger), SUB_CLOSE_DELAY),
  );
}

// The keyboard-only open path (ArrowRight / Enter on a focused sub-trigger):
// explicit focus-management, deliberately NOT delegated to native popover
// focus restoration (which would race this port's own hover-follows-pointer
// mechanism — see context-menu.gsx's file header and this file's own toggle
// handler comment).
function openSubAndFocusFirst(trigger) {
  const content = openSub(trigger);
  if (content) ownItems(content)[0]?.focus();
}

on(
  "pointerover",
  "[data-gsxui-slot-context-menu-sub-trigger]",
  (_e, trigger) => {
    if (
      trigger.getAttribute("aria-disabled") === "true" ||
      "disabled" in trigger.dataset
    )
      return;
    openSub(trigger);
  },
);
on("pointerout", "[data-gsxui-slot-context-menu-sub-trigger]", (e, trigger) => {
  const root = subRootOf(trigger);
  // Moving from the trigger into its OWN content (the gap between them, or
  // the content itself) is not a leave — same "moved within" guard as
  // hover-card.js's own trigger/content pair.
  if (
    root &&
    e.relatedTarget instanceof Element &&
    root.contains(e.relatedTarget)
  )
    return;
  scheduleCloseSub(trigger);
});
on(
  "pointerover",
  "[data-gsxui-slot-context-menu-sub-content]",
  (_e, content) => {
    const trigger = subTriggerOf(content);
    if (trigger) clearSubTimer(trigger);
  },
);
on("pointerout", "[data-gsxui-slot-context-menu-sub-content]", (e, content) => {
  const root = subRootOf(content);
  if (
    root &&
    e.relatedTarget instanceof Element &&
    root.contains(e.relatedTarget)
  )
    return;
  const trigger = subTriggerOf(content);
  if (trigger) scheduleCloseSub(trigger);
});

// Click toggles: closed -> open (idempotent-safe, moving focus in, same as
// ArrowRight — covers touch/no-hover pointer types, which never fire the
// pointerover open above); already-open -> closed. Without the open branch
// checking current state first, clicking an ALREADY-open sub-trigger (e.g.
// one opened by hover) would re-focus its first item and steal the hover
// highlight instead of acting as a close toggle.
on("click", "[data-gsxui-slot-context-menu-sub-trigger]", (_e, trigger) => {
  if (
    trigger.getAttribute("aria-disabled") === "true" ||
    "disabled" in trigger.dataset
  )
    return;
  if (trigger.dataset.state === "open") {
    closeSub(trigger);
    return;
  }
  openSubAndFocusFirst(trigger);
});
