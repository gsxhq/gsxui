// Menubar behavior on the native popover API: top layer, light dismiss and
// Esc are free. Adapted from dropdown.js: same role="menu" reads, arrow-key
// roving focus WITHIN an open content, close-on-select, and toggle-driven
// state/aria sync — reused (duplicated, not imported — see ui/menubar.gsx's
// own HOOK NAMESPACING paragraph and the "no shared JS module" constraint,
// registry.Deps is derived by go/parser over the generated .x.go, so a
// JS-only shared module would be invisible to vendoring) for the item/
// submenu machinery Task 1 already established. Their
// `data-gsxui-slot-menubar-*` hooks are namespaced to THIS component, never a
// prefix shared with dropdown.js's own `data-gsxui-slot-dropdown-menu-*` or
// context-menu.js's own `data-gsxui-slot-context-menu-*` — see ui/menubar.gsx's
// own HOOK NAMESPACING paragraph for the incident this fixes.
//
// WHAT'S NEW (the only two things this task adds, per its own brief):
//   a. Roving tabindex across MenubarTrigger — the whole bar is one tab
//      stop; ArrowLeft/ArrowRight move it. Follows ui/toggle-group.js's own
//      model verbatim: server renders every trigger with NO tabindex
//      attribute at all (graceful no-JS fallback — every trigger is its own
//      tab stop until this module loads), then a one-time init scan
//      collapses each bar to exactly one tabIndex=0 trigger, maintained on
//      every arrow move and click.
//   b. Open-follows-hover once one menu is open — hovering a SIBLING
//      trigger switches to its menu with no click, but hovering a trigger
//      while NO menu in this bar is open does nothing (a menubar, unlike a
//      plain dropdown, requires an explicit click/Enter/Space/ArrowDown to
//      open the first menu). Sibling menus are NOT nested in one another
//      (only a menu and its own submenus are, per the SUBMENUS mechanism
//      below), so switching is a close-then-open, not a reparent.
// Both are DERIVED-NOT-READ from the public Radix Menubar contract (source
// map `## menubar` §2/§3) — no @radix-ui/react-menubar dist exists in the
// reference checkout, so this is built to the documented public behavior,
// not verified against Radix's own runtime.
//
// SUBMENUS (MenubarSub/SubTrigger/SubContent) and checkbox/radio items: see
// ui/menubar.gsx's file header for the DOM-nesting-not-portalling rationale
// and the server-rendered-checked MECHANISM — same mechanism as dropdown.js's
// own. A submenu's trigger/content pair reuses hover-card.js's
// scheduleShow/scheduleHide timer shape for its pointer-leave grace period
// (SUB_CLOSE_DELAY, same value as dropdown.js's own) — OPEN is immediate on
// pointerover/ArrowRight/click, unlike hover-card's own delayed open.
//
// toggle doesn't bubble — capture.
import { on, emit, position } from "./gsxui.js";

// contentOf resolves a trigger (or anything inside its own MenubarMenu) to
// THAT menu's own content — scoped by data-gsxui-slot-menubar-menu, the
// per-menu pairing wrapper, NOT data-gsxui-slot-menubar (the whole bar) — the
// same "closest content, not closest bar" proximity shape dropdown.js's own
// contentOf uses for its single trigger/content pair.
const contentOf = (el) =>
  el
    .closest("[data-gsxui-slot-menubar-menu]")
    ?.querySelector("[data-gsxui-slot-menubar-content]");

// Any popover=auto menu surface — the top-level content or a nested
// submenu's own content — wherever a keydown/toggle handler must resolve
// "which content level is this happening at," picking the NEAREST one via
// closest().
const CONTENT_SELECTOR =
  "[data-gsxui-slot-menubar-content],[data-gsxui-slot-menubar-sub-content]";
const ITEM_SELECTOR =
  "[data-gsxui-slot-menubar-item]:not([aria-disabled]),[data-gsxui-slot-menubar-checkbox-item]:not([aria-disabled]),[data-gsxui-slot-menubar-radio-item]:not([aria-disabled]),[data-gsxui-slot-menubar-sub-trigger]:not([aria-disabled])";

// Items belonging to THIS content alone. content.querySelectorAll would also
// recurse into a DOM-nested submenu's own items — that nesting is exactly
// what makes submenus work at all (see menubar.gsx's file header) — so a
// naive query would let a submenu's items leak into its PARENT's roving
// arrow-key list even while the submenu is closed. The closest() check
// excludes any item whose nearest content ancestor isn't this one.
function ownItems(content) {
  return [...content.querySelectorAll(ITEM_SELECTOR)].filter(
    (item) => item.closest(CONTENT_SELECTOR) === content,
  );
}

const subRootOf = (el) => el.closest("[data-gsxui-slot-menubar-sub]");
const subContentOf = (trigger) =>
  subRootOf(trigger)?.querySelector("[data-gsxui-slot-menubar-sub-content]");
const subTriggerOf = (content) =>
  subRootOf(content)?.querySelector("[data-gsxui-slot-menubar-sub-trigger]");

// --- (a) roving tabindex across MenubarTrigger — ui/toggle-group.js's own model ---

const triggersOf = (bar) =>
  [...bar.querySelectorAll("[data-gsxui-slot-menubar-trigger]")].filter(
    (t) => t.closest("[data-gsxui-slot-menubar]") === bar,
  );

// Enabled-only view, same "who can become the tab stop / arrow-key
// destination" filter toggle-group.js's own normalize()/keydown handler
// apply via `items.filter(i => !i.disabled)`. Without this, a bar whose
// FIRST trigger is disabled would hand it tabIndex=0 in normalize() below —
// an unfocusable disabled button, leaving the whole bar with no reachable
// tab stop at all — and the arrow walk could land focus on a disabled
// trigger mid-bar.
const enabledTriggersOf = (bar) => triggersOf(bar).filter((t) => !t.disabled);

// Exactly one trigger (the current tab stop) gets tabindex="0"; every other
// trigger gets "-1" — same invariant toggle-group.js's own setActiveTabStop
// maintains, restated here rather than imported (no shared JS module, see
// the file header). Iterates the FULL (not enabled-only) list: a disabled
// trigger still needs its own tabindex explicitly set to -1, same as every
// other non-active trigger.
function setActiveTrigger(bar, trigger) {
  for (const t of triggersOf(bar)) t.tabIndex = t === trigger ? 0 : -1;
}

// Entry-tabstop assignment for bars rendered without JS having run yet
// (server renders every trigger as a plain tab stop) — one-time init scan,
// same shape as toggle-group.js's own normalize(). Enabled-only: the first
// trigger might be disabled.
function normalize(bar) {
  const triggers = enabledTriggersOf(bar);
  if (triggers.length) setActiveTrigger(bar, triggers[0]);
}
for (const bar of document.querySelectorAll("[data-gsxui-slot-menubar]"))
  normalize(bar);

// isAnyOpen gates (b) open-follows-hover: a menubar requires an explicit
// open of the FIRST menu (click/Enter/Space/ArrowDown) — only once one
// menu in this bar is open does hovering a sibling trigger switch menus
// with no click.
function isAnyOpen(bar) {
  return triggersOf(bar).some((t) => contentOf(t)?.matches(":popover-open"));
}

// openMenu shows trigger's own content, closing any OTHER open content in
// the SAME bar first — sibling menus are never DOM-nested in one another
// (only a menu and its own submenus are, see the SUBMENUS section below),
// so switching which one is open is a close-then-open, not a reparent.
// Idempotent: re-opening an already-open trigger's own menu is a no-op, so
// callers (click, hover, arrow-key-with-a-menu-open) can call it
// unconditionally without checking state themselves first.
function openMenu(trigger) {
  const content = contentOf(trigger);
  if (!content) return;
  const bar = trigger.closest("[data-gsxui-slot-menubar]");
  if (bar) {
    for (const t of triggersOf(bar)) {
      if (t === trigger) continue;
      const other = contentOf(t);
      if (other?.matches(":popover-open")) other.hidePopover();
    }
  }
  if (content.matches(":popover-open")) return;
  // Stamp open BEFORE showing: the toggle event that also stamps it is
  // queued as a separate task, and a paint can land in the gap — see
  // dropdown.js's own comment on the identical rule.
  content.dataset.state = "open";
  content.showPopover();
  // Below the trigger, left-aligned, 4px sideOffset. Flip/shift/clamp +
  // scroll/resize tracking: see gsxui.js.
  position(content, trigger, { side: "bottom", sideOffset: 4 });
}

// switchBarMenu moves the roving tab stop to the adjacent ENABLED trigger
// (wrapping) and opens its menu unconditionally — the shared step both (a)'s
// own trigger-level ArrowLeft/ArrowRight and the real "arrow-switch between
// open menus" behavior (wired at the CONTENT level below — see that
// handler's own comment for why the trigger level alone can't carry it)
// call into.
function switchBarMenu(trigger, dir) {
  const bar = trigger.closest("[data-gsxui-slot-menubar]");
  if (!bar) return;
  const triggers = enabledTriggersOf(bar);
  const i = triggers.indexOf(trigger);
  if (i === -1 || !triggers.length) return;
  const next = triggers[(i + dir + triggers.length) % triggers.length];
  setActiveTrigger(bar, next);
  next.focus();
  openMenu(next);
}

// A pointerdown on the trigger records whether ITS OWN menu was open at
// that instant: popover="auto" light-dismisses on outside pointerdown (the
// trigger is outside the content), so by click time the popover may already
// be closed and a bare toggle would wrongly reopen it — same MECHANISM as
// dropdown.js's own trigger click guard.
on("pointerdown", "[data-gsxui-slot-menubar-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (content)
    trigger.dataset.gsxuiWasOpen = content.matches(":popover-open")
      ? "true"
      : "false";
});

on("click", "[data-gsxui-slot-menubar-trigger]", (_e, trigger) => {
  const content = contentOf(trigger);
  if (!content) return;
  const wasOpen = trigger.dataset.gsxuiWasOpen === "true";
  delete trigger.dataset.gsxuiWasOpen;
  const bar = trigger.closest("[data-gsxui-slot-menubar]");
  if (bar) setActiveTrigger(bar, trigger);
  if (wasOpen) {
    if (content.matches(":popover-open")) content.hidePopover();
    return;
  }
  if (content.matches(":popover-open")) {
    // keyboard activation close path
    content.hidePopover();
    return;
  }
  openMenu(trigger);
});

// (b) open-follows-hover: only once some menu in this bar is already open.
on("pointerover", "[data-gsxui-slot-menubar-trigger]", (_e, trigger) => {
  const bar = trigger.closest("[data-gsxui-slot-menubar]");
  if (!bar || !isAnyOpen(bar)) return;
  const content = contentOf(trigger);
  if (content?.matches(":popover-open")) return; // already this one — no-op
  setActiveTrigger(bar, trigger);
  openMenu(trigger);
});

// (a) roving tabindex, focused ON A TRIGGER: ArrowLeft/ArrowRight move the
// tab stop (opening the destination's menu too, IF one was already open in
// this bar — the trigger-level echo of (b)'s "connected bar" feel). This
// branch is reachable while focus is actually on a trigger element — the
// common "arrow through an OPEN menu's own items" case is NOT this: opening
// a menu moves focus onto its first ITEM (the toggle handler below), so a
// subsequent arrow keydown targets that item, not the trigger — see the
// CONTENT-level keydown handler's own ArrowLeft/ArrowRight branch, where
// that case is actually handled. This handler still matters on its own for
// a menu with zero items (focus never left the trigger) and for a trigger
// refocused after Shift+Tab. ArrowDown/Enter/Space open the focused
// trigger's own menu — the explicit-open path a menubar requires before
// hover-follow (isAnyOpen above) ever activates.
on("keydown", "[data-gsxui-slot-menubar-trigger]", (e, trigger) => {
  // A held modifier suppresses the whole handler (Radix's own guard,
  // ported verbatim from ui/toggle-group.js's own keydown handler) — e.g.
  // Cmd+ArrowRight is a browser/OS shortcut, not a request to rove.
  if (e.metaKey || e.ctrlKey || e.altKey || e.shiftKey) return;
  const bar = trigger.closest("[data-gsxui-slot-menubar]");
  if (!bar) return;
  const dir = { ArrowLeft: -1, ArrowRight: 1 }[e.key];
  if (dir) {
    e.preventDefault();
    const triggers = enabledTriggersOf(bar);
    const i = triggers.indexOf(trigger);
    if (i === -1 || !triggers.length) return;
    const next = triggers[(i + dir + triggers.length) % triggers.length];
    setActiveTrigger(bar, next);
    next.focus();
    // Only switch which menu is SHOWING if one already was — a bare rove
    // must not auto-open the first menu (a menubar requires an explicit
    // click/Enter/Space/ArrowDown, see (b)'s own isAnyOpen gate above).
    if (isAnyOpen(bar)) openMenu(next);
    return;
  }
  if (e.key === "ArrowDown" || e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    setActiveTrigger(bar, trigger);
    openMenu(trigger);
  }
});

on(
  "toggle",
  CONTENT_SELECTOR,
  (e, content) => {
    const open = e.newState === "open";
    content.dataset.state = open ? "open" : "closed";
    const isSub = content.matches("[data-gsxui-slot-menubar-sub-content]");
    const trigger = isSub
      ? subTriggerOf(content)
      : content
          .closest("[data-gsxui-slot-menubar-menu]")
          ?.querySelector("[data-gsxui-slot-menubar-trigger]");
    trigger?.setAttribute("aria-expanded", open ? "true" : "false");
    // Unlike dropdown.js (where only a SubTrigger's own class keys :open
    // highlighting off data-state), MenubarTrigger's own class ALSO carries
    // data-[state=open]: — sync it for both levels here.
    if (trigger) trigger.dataset.state = open ? "open" : "closed";
    if (open) {
      // clear only on open — clearing on close races the trigger-click task
      // that needs to read the flag.
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
  "[data-gsxui-slot-menubar-item],[data-gsxui-slot-menubar-checkbox-item],[data-gsxui-slot-menubar-radio-item],[data-gsxui-slot-menubar-sub-trigger]";

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
    if (item.matches("[data-gsxui-slot-menubar-sub-trigger]")) {
      openSubAndFocusFirst(item);
      return;
    }
    item.click();
    return;
  }
  // ArrowRight: opening a highlighted sub-trigger's own submenu takes
  // priority (WAI-ARIA submenu convention) — only once that doesn't apply
  // AND we're not already inside a submenu does ArrowRight mean "switch to
  // the NEXT top-level menu in the bar," the actual reachable path for the
  // arrow-switch-between-open-menus behavior (a held modifier is not
  // checked here, unlike the trigger-level handler above: these are
  // WAI-ARIA menu navigation keys operating on already-open menu content,
  // not the bar's own roving-tabindex entry point a browser/OS shortcut
  // could plausibly intercept).
  if (e.key === "ArrowRight") {
    const subTrigger = e.target.closest(
      "[data-gsxui-slot-menubar-sub-trigger]",
    );
    if (subTrigger) {
      e.preventDefault();
      openSubAndFocusFirst(subTrigger);
      return;
    }
    if (content.matches("[data-gsxui-slot-menubar-sub-content]")) return; // no meaning for a plain item inside a submenu
    e.preventDefault();
    const trigger = content
      .closest("[data-gsxui-slot-menubar-menu]")
      ?.querySelector("[data-gsxui-slot-menubar-trigger]");
    if (trigger) switchBarMenu(trigger, 1);
    return;
  }
  // ArrowLeft: closing the CURRENT submenu and returning focus to its own
  // sub-trigger takes priority when content IS a sub-content (unchanged
  // meaning) — otherwise (top-level content), it's the mirror of
  // ArrowRight above: switch to the PREVIOUS top-level menu in the bar.
  if (e.key === "ArrowLeft") {
    if (content.matches("[data-gsxui-slot-menubar-sub-content]")) {
      e.preventDefault();
      const trigger = subTriggerOf(content);
      content.hidePopover();
      trigger?.focus();
      return;
    }
    e.preventDefault();
    const trigger = content
      .closest("[data-gsxui-slot-menubar-menu]")
      ?.querySelector("[data-gsxui-slot-menubar-trigger]");
    if (trigger) switchBarMenu(trigger, -1);
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
        item.matches("[data-gsxui-slot-menubar-sub-trigger]") &&
        item.dataset.state === "open"
      ) {
        closeSub(item);
      }
    }
    next.focus();
  }
  e.preventDefault();
});

on("click", "[data-gsxui-slot-menubar-item]", (_e, item) => {
  if (
    item.getAttribute("aria-disabled") === "true" ||
    "disabled" in item.dataset
  )
    return;
  // Always hide the ROOT content, not item.closest(CONTENT_SELECTOR) — from
  // inside a submenu, that resolves to the SUB-content, closing only the
  // submenu and leaving the root open with the sub-trigger still
  // highlighted. contentOf() resolves via data-gsxui-slot-menubar-menu (this
  // item's own top-level menu), not the whole bar, so a select inside one
  // MenubarMenu never touches a sibling's own content. One call on the root
  // is enough: the measured popover stack cascade (see menubar.gsx's file
  // header) closes any nested-open submenu along with it.
  const content = contentOf(item);
  emit(item, "gsxui:select");
  content?.hidePopover();
});

// A checkbox item flips in place and does NOT close its menu — a deliberate
// ADAPT per the task brief, not a Radix default (Radix's actual CheckboxItem
// onSelect closes unless preventDefault is called; this port simply never
// closes on a checkbox select).
on("click", "[data-gsxui-slot-menubar-checkbox-item]", (_e, item) => {
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
// radio group (data-gsxui-slot-menubar-radio-group scopes the sibling walk — a
// page may have more than one group), and closes the menu like a plain
// Item.
on("click", "[data-gsxui-slot-menubar-radio-item]", (_e, item) => {
  if (
    item.getAttribute("aria-disabled") === "true" ||
    "disabled" in item.dataset
  )
    return;
  const group = item.closest("[data-gsxui-slot-menubar-radio-group]");
  const siblings = group
    ? [...group.querySelectorAll("[data-gsxui-slot-menubar-radio-item]")]
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
  "[data-gsxui-slot-menubar-item],[data-gsxui-slot-menubar-checkbox-item],[data-gsxui-slot-menubar-radio-item],[data-gsxui-slot-menubar-sub-trigger]";
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
  // trigger's own openMenu.
  trigger.dataset.state = "open";
  trigger.setAttribute("aria-expanded", "true");
  content.dataset.state = "open";
  content.showPopover();
  // To the right of its trigger; alignOffset -4 offsets the content's own
  // p-1 so the first item's text roughly aligns with the trigger's.
  // Flip/shift/clamp + scroll/resize tracking: see gsxui.js.
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
// mechanism — see menubar.gsx's file header and this file's own toggle
// handler comment).
function openSubAndFocusFirst(trigger) {
  const content = openSub(trigger);
  if (content) ownItems(content)[0]?.focus();
}

on("pointerover", "[data-gsxui-slot-menubar-sub-trigger]", (_e, trigger) => {
  if (
    trigger.getAttribute("aria-disabled") === "true" ||
    "disabled" in trigger.dataset
  )
    return;
  openSub(trigger);
});
on("pointerout", "[data-gsxui-slot-menubar-sub-trigger]", (e, trigger) => {
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
on("pointerover", "[data-gsxui-slot-menubar-sub-content]", (_e, content) => {
  const trigger = subTriggerOf(content);
  if (trigger) clearSubTimer(trigger);
});
on("pointerout", "[data-gsxui-slot-menubar-sub-content]", (e, content) => {
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
on("click", "[data-gsxui-slot-menubar-sub-trigger]", (_e, trigger) => {
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
