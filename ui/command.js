// Command palette behavior — the cmdk primitive's contract without React:
// score-ranked filtering, DOM reordering by score, group/separator hiding,
// empty state, and a roving selection that keeps FOCUS in the input the
// whole time (data-selected + aria-activedescendant, cmdk's own model —
// items are never tab stops). Activation: Enter or click emits gsxui:select
// on the item; an item carrying data-href then navigates. ⌘K/Ctrl-K toggles
// the first dialog marked data-gsxui-command-dialog (see command.gsx),
// riding dialog.js's machinery for open state and animated close.
import { on, emit } from "./gsxui.js";
import { requestClose } from "./dialog.js";

// ---------------------------------------------------------------------------
// commandScore — verbatim port of command-score 0.1.2 (MIT, (c) Superhuman
// Labs; the exact ranking cmdk uses), reformatted for this module. See
// https://github.com/superhuman/command-score. Kept byte-faithful in
// constants and control flow so filtering feels identical to shadcn's.
const SCORE_CONTINUE_MATCH = 1,
  SCORE_WORD_JUMP = 0.9,
  SCORE_CHARACTER_JUMP = 0.3,
  SCORE_TRANSPOSITION = 0.1,
  SCORE_LONG_JUMP = 0,
  PENALTY_SKIPPED = 0.999,
  PENALTY_CASE_MISMATCH = 0.9999,
  PENALTY_NOT_COMPLETE = 0.99;

const IS_GAP_REGEXP = /[\\/\-_+.# \t"@[({&]/,
  COUNT_GAPS_REGEXP = /[\\/\-_+.# \t"@[({&]/g;

function commandScoreInner(string, abbreviation, lowerString, lowerAbbreviation, stringIndex, abbreviationIndex) {
  if (abbreviationIndex === abbreviation.length) {
    if (stringIndex === string.length) return SCORE_CONTINUE_MATCH;
    return PENALTY_NOT_COMPLETE;
  }

  const abbreviationChar = lowerAbbreviation.charAt(abbreviationIndex);
  let index = lowerString.indexOf(abbreviationChar, stringIndex);
  let highScore = 0;
  let score, transposedScore, wordBreaks;

  while (index >= 0) {
    score = commandScoreInner(string, abbreviation, lowerString, lowerAbbreviation, index + 1, abbreviationIndex + 1);
    if (score > highScore) {
      if (index === stringIndex) {
        score *= SCORE_CONTINUE_MATCH;
      } else if (IS_GAP_REGEXP.test(string.charAt(index - 1))) {
        score *= SCORE_WORD_JUMP;
        wordBreaks = string.slice(stringIndex, index - 1).match(COUNT_GAPS_REGEXP);
        if (wordBreaks && stringIndex > 0) {
          score *= Math.pow(PENALTY_SKIPPED, wordBreaks.length);
        }
      } else if (IS_GAP_REGEXP.test(string.slice(stringIndex, index - 1))) {
        score *= SCORE_LONG_JUMP;
        if (stringIndex > 0) {
          score *= Math.pow(PENALTY_SKIPPED, index - stringIndex);
        }
      } else {
        score *= SCORE_CHARACTER_JUMP;
        if (stringIndex > 0) {
          score *= Math.pow(PENALTY_SKIPPED, index - stringIndex);
        }
      }
      if (string.charAt(index) !== abbreviation.charAt(abbreviationIndex)) {
        score *= PENALTY_CASE_MISMATCH;
      }
    }

    if (
      score < SCORE_TRANSPOSITION &&
      lowerString.charAt(index - 1) === lowerAbbreviation.charAt(abbreviationIndex + 1) &&
      lowerString.charAt(index - 1) !== lowerAbbreviation.charAt(abbreviationIndex)
    ) {
      transposedScore = commandScoreInner(string, abbreviation, lowerString, lowerAbbreviation, index + 1, abbreviationIndex + 2);
      if (transposedScore * SCORE_TRANSPOSITION > score) {
        score = transposedScore * SCORE_TRANSPOSITION;
      }
    }

    if (score > highScore) highScore = score;
    index = lowerString.indexOf(abbreviationChar, index + 1);
  }

  return highScore;
}

function commandScore(string, abbreviation) {
  return commandScoreInner(string, abbreviation, string.toLowerCase(), abbreviation.toLowerCase(), 0, 0);
}
// ---------------------------------------------------------------------------

const rootOf = (el) => el.closest("[data-gsxui-command]");
const itemsOf = (root) => [...root.querySelectorAll("[data-gsxui-command-item]")];
const valueOf = (item) => item.dataset.value || item.textContent.trim();
const disabled = (item) =>
  item.getAttribute("aria-disabled") === "true" || "disabled" in item.dataset;

let uid = 0;

function select(root, item) {
  const input = root.querySelector("[data-gsxui-command-input]");
  for (const other of itemsOf(root)) {
    if (other === item) continue;
    delete other.dataset.selected;
    other.setAttribute("aria-selected", "false");
  }
  if (!item) {
    input?.removeAttribute("aria-activedescendant");
    return;
  }
  item.dataset.selected = "true";
  item.setAttribute("aria-selected", "true");
  if (!item.id) item.id = `gsxui-command-item-${++uid}`;
  input?.setAttribute("aria-activedescendant", item.id);
  item.scrollIntoView({ block: "nearest" });
}

const visibleItems = (root) => itemsOf(root).filter((i) => !i.hidden && !disabled(i));
const selectedOf = (root) => root.querySelector('[data-gsxui-command-item][data-selected="true"]');

// Filter + rank: score every item, hide zero-scores, reorder items within
// their group by score (cmdk reorders the DOM too — element identity is
// preserved, so listeners and state survive appendChild moves), reorder
// groups by their best item, hide empty groups, hide separators while a
// query is active (cmdk's behavior), and toggle the empty state. Ties keep
// source order via the one-time data-gsxui-index stamp.
function filter(root) {
  const input = root.querySelector("[data-gsxui-command-input]");
  const query = (input?.value ?? "").trim();
  const list = root.querySelector("[data-gsxui-command-list]") ?? root;

  // One-time source-order stamps (items AND the list's top-level children:
  // groups, separators, the empty element, ungrouped items). Ranking ties
  // resolve to this order, which is what restores the original DOM layout
  // when the query clears — every score is 1 then, so the stable sort below
  // IS the undo of all previous reorders.
  const items = itemsOf(root);
  items.forEach((item, i) => {
    if (!("gsxuiIndex" in item.dataset)) item.dataset.gsxuiIndex = String(i);
  });
  [...list.children].forEach((child, i) => {
    if (!("gsxuiIndex" in child.dataset)) child.dataset.gsxuiIndex = String(i);
  });

  const scores = new Map();
  let any = false;
  for (const item of items) {
    const score = query ? commandScore(valueOf(item), query) : 1;
    scores.set(item, score);
    item.hidden = score === 0;
    if (!item.hidden) any = true;
  }

  const sourceOrder = (a, b) => Number(a.dataset.gsxuiIndex) - Number(b.dataset.gsxuiIndex);

  // Items reorder within their group every pass (appendChild preserves
  // element identity, so listeners/state survive the moves).
  for (const group of root.querySelectorAll('[data-slot="command-group"]')) {
    const inGroup = items.filter((i) => group.contains(i));
    const ranked = [...inGroup].sort((a, b) => scores.get(b) - scores.get(a) || sourceOrder(a, b));
    for (const item of ranked) group.appendChild(item);
    group.hidden = !inGroup.some((i) => !i.hidden);
  }

  // The list's top level reorders by best contained score (a bare item is
  // its own score; a group is its best item; separators/empty score 0 and
  // are hidden or inert while a query is active anyway).
  const rank = (el) =>
    scores.get(el) ?? Math.max(0, ...items.filter((i) => el.contains(i)).map((i) => scores.get(i)));
  const kids = [...list.children].sort((a, b) => rank(b) - rank(a) || sourceOrder(a, b));
  for (const kid of kids) list.appendChild(kid);

  for (const sep of root.querySelectorAll('[data-slot="command-separator"]')) {
    sep.hidden = query !== "";
  }
  const empty = root.querySelector('[data-slot="command-empty"]');
  if (empty) empty.hidden = any;

  // Selection follows the ranking: top visible item after every keystroke.
  select(root, visibleItems(root)[0] ?? null);
}

function activate(item) {
  if (!item || disabled(item)) return;
  emit(item, "gsxui:select");
  if (item.dataset.href) {
    const dialog = item.closest("dialog[data-gsxui-dialog-content]");
    if (dialog) requestClose(dialog);
    window.location.assign(item.dataset.href);
  }
}

on("input", "[data-gsxui-command-input]", (_e, input) => {
  const root = rootOf(input);
  if (root) filter(root);
});

on("keydown", "[data-gsxui-command-input]", (e, input) => {
  const root = rootOf(input);
  if (!root) return;
  if (e.key === "Enter") {
    e.preventDefault();
    activate(selectedOf(root));
    return;
  }
  const dir = { ArrowDown: 1, ArrowUp: -1 }[e.key];
  if (!dir) return;
  e.preventDefault();
  const items = visibleItems(root);
  if (!items.length) return;
  const i = items.indexOf(selectedOf(root));
  select(root, items[(i + dir + items.length) % items.length]);
});

on("click", "[data-gsxui-command-item]", (_e, item) => activate(item));

// Selection follows the pointer (cmdk hover model — hover IS selection, no
// separate hover style), same rationale as dropdown.js's pointerover focus.
on("pointerover", "[data-gsxui-command-item]", (_e, item) => {
  if (disabled(item) || item.dataset.selected === "true") return;
  const root = rootOf(item);
  if (root) select(root, item);
});

// Initial state: rank/selection for palettes rendered without a query
// (server renders every item visible; this stamps the first selection).
for (const root of document.querySelectorAll("[data-gsxui-command]")) filter(root);

// ⌘K / Ctrl-K toggles the first command dialog on the page. Opening mirrors
// dialog.js's trigger path (stamp data-state BEFORE showModal — the queued
// toggle task alone paints a closed-state frame, the flash class of bugs);
// closing rides requestClose so the exit animation runs.
addEventListener("keydown", (e) => {
  if (e.key.toLowerCase() !== "k" || !(e.metaKey || e.ctrlKey)) return;
  const dialog = document.querySelector("dialog[data-gsxui-command-dialog]");
  if (!dialog) return;
  e.preventDefault();
  if (dialog.open) {
    requestClose(dialog);
  } else {
    dialog.dataset.state = "open";
    dialog.showModal();
  }
});

// A fresh open starts a fresh search: clear the query and re-rank when the
// dialog machinery announces the open (dialog.js emits gsxui:open on the
// <dialog> for every open path — trigger click, hotkey, programmatic).
on(
  "gsxui:open",
  "dialog[data-gsxui-command-dialog]",
  (_e, dialog) => {
    const input = dialog.querySelector("[data-gsxui-command-input]");
    if (input) input.value = "";
    const root = rootOf(input ?? dialog) ?? dialog.querySelector("[data-gsxui-command]");
    if (root) filter(root);
    input?.focus();
  },
  { capture: true },
);
