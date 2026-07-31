// Dropdown behavior on the native popover API: top layer, light dismiss and
// Esc are free. JS adds anchored positioning, state/aria sync, arrow-key
// roving focus, and close-on-select. toggle doesn't bubble — capture.
//
// SUBMENUS (DropdownMenuSub/SubTrigger/SubContent) and checkbox/radio items:
// see ui/dropdown.gsx's file header for the DOM-nesting-not-portalling
// rationale (spec §1) and the server-rendered-checked MECHANISM. A submenu's
// trigger/content pair reuses hover-card.js's scheduleShow/scheduleHide
// timer shape for its pointer-leave grace period (CLOSE_DELAY, same value,
// per the task brief's "match the shipped hover-card.js precedent") — but
// OPEN is immediate on pointerover/ArrowRight/click, unlike hover-card's own
// delayed open, per the brief's own step-by-step spec. Escape and outside-
// pointerdown light dismiss are free here too: closing an ancestor auto
// popover closes its nested (DOM-descendant) auto popovers as part of the
// same platform stack-unwind, so a submenu never needs to be closed by hand
// when its parent menu closes.
import { on, emit } from "./gsxui.js";

const contentOf = (el) =>
  el.closest("[data-gsxui-slot-dropdown-menu]")?.querySelector("[data-gsxui-slot-dropdown-menu-content]");

// Any popover=auto menu surface — the top-level content or a nested
// submenu's own content — wherever a keydown/toggle handler must resolve
// "which content level is this happening at," picking the NEAREST one via
// closest().
const CONTENT_SELECTOR = "[data-gsxui-slot-dropdown-menu-content],[data-gsxui-slot-dropdown-menu-sub-content]";
const ITEM_SELECTOR =
  "[data-gsxui-slot-dropdown-menu-item]:not([aria-disabled]),[data-gsxui-slot-dropdown-menu-checkbox-item]:not([aria-disabled]),[data-gsxui-slot-dropdown-menu-radio-item]:not([aria-disabled]),[data-gsxui-slot-dropdown-menu-sub-trigger]:not([aria-disabled])";

// Items belonging to THIS content alone. content.querySelectorAll would also
// recurse into a DOM-nested submenu's own items — that nesting is exactly
// what makes submenus work at all (see dropdown.gsx's file header) — so a
// naive query would let a submenu's items leak into its PARENT's roving
// arrow-key list even while the submenu is closed. The closest() check
// excludes any item whose nearest content ancestor isn't this one.
function ownItems(content) {
  return [...content.querySelectorAll(ITEM_SELECTOR)].filter(
    (item) => item.closest(CONTENT_SELECTOR) === content,
  );
}

const subRootOf = (el) => el.closest("[data-gsxui-slot-dropdown-menu-sub]");
const subContentOf = (trigger) => subRootOf(trigger)?.querySelector("[data-gsxui-slot-dropdown-menu-sub-content]");
const subTriggerOf = (content) => subRootOf(content)?.querySelector("[data-gsxui-slot-dropdown-menu-sub-trigger]");

// A pointerdown on the trigger records whether the menu was open at that
// instant: popover="auto" light-dismisses on outside pointerdown (the
// trigger is outside the content), so by click time the popover may already
// be closed and a bare toggle would wrongly reopen it.
on("pointerdown", "[data-gsxui-slot-dropdown-menu-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (content) trigger.dataset.gsxuiWasOpen = content.matches(":popover-open") ? "true" : "false";
});

on("click", "[data-gsxui-slot-dropdown-menu-trigger]", (_e, trigger) => {
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
  // Stamp open BEFORE showing: the toggle event that also stamps it is
  // queued as a separate task (spec: "queue a popover toggle event task"),
  // and a paint can land in the gap — one frame of the menu fully visible
  // in its closed state, then the enter animation restarting from opacity
  // 0 reads as a flash (routinely visible when the window is inactive and
  // task scheduling is throttled, occasionally on an active window too).
  // The toggle handler's own stamp stays as reconciliation for closes.
  content.dataset.state = "open";
  content.showPopover();
});

on(
  "toggle",
  CONTENT_SELECTOR,
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const isSub = content.matches("[data-gsxui-slot-dropdown-menu-sub-content]");
    const trigger = isSub
      ? subTriggerOf(content)
      : content.closest("[data-gsxui-slot-dropdown-menu]")?.querySelector("[data-gsxui-slot-dropdown-menu-trigger]");
    trigger?.setAttribute("aria-expanded", open ? "true" : "false");
    // Sub-trigger's own class keys :open highlighting off data-state (same
    // selector shape DropdownMenuItem's data-variant uses) — the top-level
    // DropdownMenuTrigger has no such selector, so only subs get this.
    if (isSub && trigger) trigger.dataset.state = open ? "open" : "closed";
    if (open) {
      // clear only on open — clearing on close races the trigger-click task that needs to read the flag
      delete trigger?.dataset.gsxuiWasOpen;
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
  "[data-gsxui-slot-dropdown-menu-item],[data-gsxui-slot-dropdown-menu-checkbox-item],[data-gsxui-slot-dropdown-menu-radio-item],[data-gsxui-slot-dropdown-menu-sub-trigger]";

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
    if (item.matches("[data-gsxui-slot-dropdown-menu-sub-trigger]")) {
      openSubAndFocusFirst(item);
      return;
    }
    item.click();
    return;
  }
  if (e.key === "ArrowRight") {
    const trigger = e.target.closest("[data-gsxui-slot-dropdown-menu-sub-trigger]");
    if (!trigger) return;
    e.preventDefault();
    openSubAndFocusFirst(trigger);
    return;
  }
  if (e.key === "ArrowLeft" && content.matches("[data-gsxui-slot-dropdown-menu-sub-content]")) {
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
  const next = i === -1 ? (dir === 1 ? items[0] : items[items.length - 1]) : items[(i + dir + items.length) % items.length];
  if (next) {
    // Moving keyboard focus to a different item at this same level closes
    // any OTHER open submenu among this level's own sub-triggers — without
    // this, a submenu opened earlier by hover (mouse stationary, so no
    // pointerout ever fires) would stay visually open while focus moves
    // elsewhere in the list. Mouse-driven moves don't need this: leaving a
    // sub-trigger by an actual pointer movement already fires pointerout on
    // it, which schedules the same close via the grace-period timer below.
    for (const item of items) {
      if (item !== next && item.matches("[data-gsxui-slot-dropdown-menu-sub-trigger]") && item.dataset.state === "open") {
        closeSub(item);
      }
    }
    next.focus();
  }
  e.preventDefault();
});

on("click", "[data-gsxui-slot-dropdown-menu-item]", (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  // Always hide the ROOT content, not item.closest(CONTENT_SELECTOR) — from
  // inside a submenu, that resolves to the SUB-content, closing only the
  // submenu and leaving the root open with the sub-trigger still
  // highlighted. One call on the root is enough: the measured popover
  // stack cascade (see dropdown.gsx's file header) closes any nested-open
  // submenu along with it.
  const content = contentOf(item);
  emit(item, "gsxui:select");
  content?.hidePopover();
});

// A checkbox item flips in place and does NOT close its menu — a deliberate
// ADAPT per the task brief, not a Radix default (Radix's actual CheckboxItem
// onSelect closes unless preventDefault is called; this port simply never
// closes on a checkbox select).
on("click", "[data-gsxui-slot-dropdown-menu-checkbox-item]", (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  const checked = item.dataset.state !== "checked";
  item.dataset.state = checked ? "checked" : "unchecked";
  item.setAttribute("aria-checked", checked ? "true" : "false");
  emit(item, "gsxui:change", { checked, value: item.dataset.value ?? "" });
});

// A radio item sets itself checked, clears every OTHER item in the SAME
// radio group (data-gsxui-slot-dropdown-menu-radio-group scopes the sibling walk — a page
// may have more than one group), and closes the menu like a plain Item.
on("click", "[data-gsxui-slot-dropdown-menu-radio-item]", (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  const group = item.closest("[data-gsxui-slot-dropdown-menu-radio-group]");
  const siblings = group
    ? [...group.querySelectorAll("[data-gsxui-slot-dropdown-menu-radio-item]")]
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
  "[data-gsxui-slot-dropdown-menu-item],[data-gsxui-slot-dropdown-menu-checkbox-item],[data-gsxui-slot-dropdown-menu-radio-item],[data-gsxui-slot-dropdown-menu-sub-trigger]";
on("pointerover", HOVERABLE_SELECTOR, (_e, item) => {
  if (item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset) return;
  item.focus();
});

// Leaving a content entirely clears the item highlight by parking focus on
// that content (tabindex="-1") — not body, so arrow keys keep working.
// Moving into a DOM-nested (open) submenu doesn't count as leaving — it's
// still `content.contains(relatedTarget)` true, same guard as hover-card.js.
on("pointerout", CONTENT_SELECTOR, (e, content) => {
  if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget)) return;
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
  const r = trigger.getBoundingClientRect();
  content.style.position = "fixed";
  content.style.inset = "auto";
  content.style.left = `${r.right}px`;
  content.style.top = `${r.top - 4}px`; // -4 offsets the content's own p-1 so the first item's text roughly aligns with the trigger's
  // Stamp open BEFORE showing — same flash-avoidance rule as the top-level
  // trigger's own click handler (a queued toggle event can otherwise leave
  // one frame painted in the stale closed state).
  trigger.dataset.state = "open";
  trigger.setAttribute("aria-expanded", "true");
  content.dataset.state = "open";
  content.showPopover();
  return content;
}

function closeSub(trigger) {
  clearSubTimer(trigger);
  const content = subContentOf(trigger);
  if (content?.matches(":popover-open")) content.hidePopover();
}

function scheduleCloseSub(trigger) {
  clearSubTimer(trigger);
  subTimers.set(trigger, setTimeout(() => closeSub(trigger), SUB_CLOSE_DELAY));
}

// The keyboard-only open path (ArrowRight / Enter on a focused sub-trigger):
// explicit focus-management, deliberately NOT delegated to native popover
// focus restoration (which would race this port's own hover-follows-pointer
// mechanism — see dropdown.gsx's file header and this file's own toggle
// handler comment).
function openSubAndFocusFirst(trigger) {
  const content = openSub(trigger);
  if (content) ownItems(content)[0]?.focus();
}

on("pointerover", "[data-gsxui-slot-dropdown-menu-sub-trigger]", (_e, trigger) => {
  if (trigger.getAttribute("aria-disabled") === "true" || "disabled" in trigger.dataset) return;
  openSub(trigger);
});
on("pointerout", "[data-gsxui-slot-dropdown-menu-sub-trigger]", (e, trigger) => {
  const root = subRootOf(trigger);
  // Moving from the trigger into its OWN content (the gap between them, or
  // the content itself) is not a leave — same "moved within" guard as
  // hover-card.js's own trigger/content pair.
  if (root && e.relatedTarget instanceof Element && root.contains(e.relatedTarget)) return;
  scheduleCloseSub(trigger);
});
on("pointerover", "[data-gsxui-slot-dropdown-menu-sub-content]", (_e, content) => {
  const trigger = subTriggerOf(content);
  if (trigger) clearSubTimer(trigger);
});
on("pointerout", "[data-gsxui-slot-dropdown-menu-sub-content]", (e, content) => {
  const root = subRootOf(content);
  if (root && e.relatedTarget instanceof Element && root.contains(e.relatedTarget)) return;
  const trigger = subTriggerOf(content);
  if (trigger) scheduleCloseSub(trigger);
});

// Click toggles: closed -> open (idempotent-safe, moving focus in, same as
// ArrowRight — covers touch/no-hover pointer types, which never fire the
// pointerover open above); already-open -> closed. Without the open branch
// checking current state first, clicking an ALREADY-open sub-trigger (e.g.
// one opened by hover) would re-focus its first item and steal the hover
// highlight instead of acting as a close toggle.
on("click", "[data-gsxui-slot-dropdown-menu-sub-trigger]", (_e, trigger) => {
  if (trigger.getAttribute("aria-disabled") === "true" || "disabled" in trigger.dataset) return;
  if (trigger.dataset.state === "open") {
    closeSub(trigger);
    return;
  }
  openSubAndFocusFirst(trigger);
});
