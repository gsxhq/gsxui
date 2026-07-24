// InputOTP behavior. The architecture is the entire mechanism (see
// ui/input-otp.gsx's own doc comment, byte-faithful to the 2026-07-24
// controls source map's `## input-otp`): ONE real, visually-hidden <input>
// owns focus/caret/paste; native text editing already supplies
// focus-advance, backspace-across-slots, and arrow navigation for free — no
// per-slot keyboard JS exists here at all. This module only:
//   (a) stamps each root's slots with a positional data-index once, at
//       mount or on first use (DOM-order stamping — InputOTPSlot itself
//       takes no index param, see the .gsx file's ADAPT doc comment; same
//       "stamp source order once" idiom command.js's data-gsxui-index
//       already establishes),
//   (b) recomputes every slot's rendered char / data-active / fake-caret
//       from the real input's live value + selection on every
//       input/selectionchange/focus/blur,
//   (c) filters keystrokes (and pastes — the ordinary `input` event
//       recompute covers paste with no special-case code path) against an
//       optional per-character data-gsxui-input-otp-pattern, and
//   (d) routes a slot click to input.focus() + setSelectionRange(index,
//       index) — a fallback on top of the native click-through the real
//       input's own opacity-0/z-10 (not pointer-events-none) already gives
//       for free.
import { on } from "./gsxui.js";

const rootOf = (el) => el.closest("[data-gsxui-input-otp]");
const inputOf = (root) => root.querySelector("[data-gsxui-input-otp-input]");
const slotsOf = (root) => [...root.querySelectorAll('[data-slot="input-otp-slot"]')];

// One-time positional data-index stamp, DOM order, idempotent (mirrors
// command.js's `if (!("gsxuiIndex" in item.dataset))` guard).
function stamp(root) {
  slotsOf(root).forEach((slot, i) => {
    if (!("index" in slot.dataset)) slot.dataset.index = String(i);
  });
}

// The fake-caret markup, byte-ported from the map: a centered overlay
// containing a single blinking line, shown only inside the active slot
// while that slot's own character is empty (i.e. where the next typed
// character will land).
function renderSlot(slot, char, active) {
  slot.dataset.active = active ? "true" : "false";
  slot.replaceChildren();
  if (char) {
    slot.append(char);
    return;
  }
  if (!active) return;
  const overlay = document.createElement("div");
  overlay.className = "pointer-events-none absolute inset-0 flex items-center justify-center";
  const caret = document.createElement("div");
  caret.className = "h-4 w-px animate-caret-blink bg-foreground duration-1000";
  overlay.appendChild(caret);
  slot.appendChild(overlay);
}

function recompute(root) {
  const input = inputOf(root);
  if (!input) return;
  const focused = document.activeElement === input;
  const collapsed = input.selectionStart === input.selectionEnd;
  const activeIndex = focused && collapsed ? input.selectionStart : -1;
  slotsOf(root).forEach((slot, i) => {
    const index = Number(slot.dataset.index ?? i);
    renderSlot(slot, input.value[index] ?? "", index === activeIndex);
  });
}

// Per-character filter RegExp, constructed once per real input (not once
// per keystroke) and cached — an invalid pattern source is caught here so
// it degrades to "no filtering" instead of throwing on every input event.
const patternCache = new WeakMap();

function patternOf(input) {
  const src = input.dataset.gsxuiInputOtpPattern;
  if (!src) return null;
  if (patternCache.has(input)) return patternCache.get(input);
  let re = null;
  try {
    re = new RegExp(src);
  } catch {
    re = null;
  }
  patternCache.set(input, re);
  return re;
}

// Strips characters that don't match the per-character pattern (typed or
// pasted alike — paste needs no special path, this handler runs on the
// same `input` event either way) before recomputing slot display.
function filterPattern(input) {
  const re = patternOf(input);
  if (!re) return;
  if ([...input.value].every((c) => re.test(c))) return;
  const pos = input.selectionStart;
  const filtered = [...input.value].filter((c) => re.test(c)).join("");
  const removed = input.value.length - filtered.length;
  input.value = filtered;
  const newPos = Math.max(0, pos - removed);
  input.setSelectionRange(newPos, newPos);
}

on("input", "[data-gsxui-input-otp-input]", (_e, input) => {
  const root = rootOf(input);
  if (!root) return;
  filterPattern(input);
  stamp(root);
  recompute(root);
});

on(
  "focus",
  "[data-gsxui-input-otp-input]",
  (_e, input) => {
    const root = rootOf(input);
    if (root) recompute(root);
  },
  { capture: true },
);

on(
  "blur",
  "[data-gsxui-input-otp-input]",
  (_e, input) => {
    const root = rootOf(input);
    if (root) recompute(root);
  },
  { capture: true },
);

// selectionchange fires on `document`, never on the input itself — the
// delegated `on()` helper dispatches by walking event.target.closest(),
// which can't work here (event.target IS document), so this is a plain
// document-level listener gated on document.activeElement instead (same
// shape as command.js's own raw ⌘K addEventListener).
document.addEventListener("selectionchange", () => {
  const input = document.activeElement;
  if (!(input instanceof Element) || !input.matches("[data-gsxui-input-otp-input]")) return;
  const root = rootOf(input);
  if (root) recompute(root);
});

on("click", '[data-slot="input-otp-slot"]', (_e, slot) => {
  const root = rootOf(slot);
  const input = root && inputOf(root);
  if (!input) return;
  stamp(root);
  const index = Number(slot.dataset.index);
  input.focus();
  input.setSelectionRange(index, index);
});

// Initial paint: populate every root present at parse time immediately
// (e.g. a server-rendered initial value) rather than waiting for the first
// interaction — same one-time init-scan shape as toggle-group.js's
// normalize() / command.js's filter() loops.
for (const root of document.querySelectorAll("[data-gsxui-input-otp]")) {
  stamp(root);
  recompute(root);
}
