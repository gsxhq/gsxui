// Calendar behavior: month navigation with in-place cell updates.
//
// The 42 day cells are rendered by calendar.gsx and never recreated — this
// module only writes textContent and data-* on the button/cell pair that
// already exists at every one of the 42 grid positions. That is what keeps
// every class string in the .gsx (Tailwind scans .gsx and .js both, so this
// is not a styling constraint — it is a single-authority one) and what
// preserves focus across navigation: the exact same button element that had
// focus before a click on prev/next still has it after, just repainted to a
// different day.
//
// monthGrid below is the twin of calendar.gsx's own monthGrid — same
// contract, same six-rows-of-seven fixed grid, UTC throughout. The two must
// agree; jstest/specs/calendar.spec.ts's "Go and JS agree on …" tests are
// exactly that agreement, checked by navigating client-side to a month and
// diffing against the server's own render of that same month.
import { on, emit } from "./gsxui.js";

const ROOT = "[data-gsxui-calendar]";

// monthGrid is the twin of calendar.gsx's monthGrid: 42 dates, six rows of
// seven, starting at the week start on or before the 1st. month is 0-based
// (JS Date convention — January is 0), matching every other date this module
// touches.
function monthGrid(year, month, weekStartsOn) {
  const first = new Date(Date.UTC(year, month, 1));
  const offset = (first.getUTCDay() - weekStartsOn + 7) % 7;
  const grid = [];
  for (let i = 0; i < 42; i++) {
    grid.push(new Date(Date.UTC(year, month, 1 - offset + i)));
  }
  return grid;
}

function iso(date) {
  return date.toISOString().slice(0, 10);
}

// parseISO parses an ISO "YYYY-MM-DD" string into a UTC-midnight Date, or
// null for an empty/malformed value — the same null-for-unset contract
// data-gsxui-calendar-from/-to/-disabled-before/-disabled-after all need,
// since each is omitted entirely rather than emitted empty.
function parseISO(value) {
  if (!value) return null;
  const [y, m, d] = value.split("-").map(Number);
  if (!y || !m || !d) return null;
  const date = new Date(Date.UTC(y, m - 1, d));
  return date.getUTCFullYear() === y && date.getUTCMonth() === m - 1 && date.getUTCDate() === d
    ? date
    : null;
}

// parseMonth turns the root's data-gsxui-calendar-month ("YYYY-MM") into a
// { year, month } pair — month 0-based, matching monthGrid's own argument.
function parseMonth(value) {
  const [y, m] = value.split("-").map(Number);
  return { year: y, month: m - 1 };
}

function formatMonth(year, month) {
  return `${String(year).padStart(4, "0")}-${String(month + 1).padStart(2, "0")}`;
}

const WEEKDAY_NAMES = [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];
const MONTH_NAMES = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

// ariaLabel matches calendar.gsx's own d.Format("Monday, January 2, 2006") —
// weekday-first, so hiding the abbreviated weekday header row (Task 3, still
// aria-hidden after navigation) stays safe: every day's accessible name
// leads with the weekday regardless of which month is currently painted.
function ariaLabel(date) {
  const weekday = WEEKDAY_NAMES[date.getUTCDay()];
  const month = MONTH_NAMES[date.getUTCMonth()];
  return `${weekday}, ${month} ${date.getUTCDate()}, ${date.getUTCFullYear()}`;
}

function captionText(year, month) {
  return `${MONTH_NAMES[month]} ${year}`;
}

// commaList reads a comma-separated data attribute into an array of
// non-empty strings — an absent attribute (calendar.gsx omits each disabled
// rule entirely when unset) yields [], never [""].
function commaList(value) {
  if (!value) return [];
  return value
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
}

function sameUTCDay(a, b) {
  return a !== null && b !== null && a.getTime() === b.getTime();
}

// navBounds reads the root's resolved nav-year bounds (calendar.gsx's own
// calendarNavBounds, always present — see that file's doc comment on why
// this pair is unconditional rather than omit-when-unset like the other
// root attributes).
function navBounds(root) {
  return {
    fromYear: Number(root.dataset.gsxuiCalendarNavFromYear),
    toYear: Number(root.dataset.gsxuiCalendarNavToYear),
  };
}

// prevDisabledAt/nextDisabledAt port calendar.gsx's own
// calendarPrevDisabled/calendarNextDisabled. month here is 0-based (JS
// Date convention, January = 0), unlike Go's 1-based time.Month — so the
// bound compares against 0/11 instead of Go's time.January/time.December.
function prevDisabledAt(year, month, fromYear) {
  return year < fromYear || (year === fromYear && month <= 0);
}
function nextDisabledAt(year, month, toYear) {
  return year > toYear || (year === toYear && month >= 11);
}

// setNavDisabled mirrors calendar.gsx's own `{ if prevDisabled/nextDisabled
// { aria-disabled="true" tabindex="-1" } }` — both attributes present only
// while at the bound, both entirely ABSENT otherwise (never a native
// disabled attribute, never aria-disabled="false"/tabindex="0" — that is
// not what the server ever writes, and the button must stay reachable by
// Tab at a bound, per calendar.gsx's own doc comment on this exact point).
function setNavDisabled(button, atBound) {
  if (!button) return;
  if (atBound) {
    button.setAttribute("aria-disabled", "true");
    button.setAttribute("tabindex", "-1");
  } else {
    button.removeAttribute("aria-disabled");
    button.removeAttribute("tabindex");
  }
}

// syncNavButtons re-evaluates the prev/next buttons' aria-disabled/tabindex
// for the just-navigated-to (year, month) against the root's own nav-bounds
// attributes. Task 3 rendered this correctly for the server's own initial
// month only; without recomputing it here, a caller with a narrow
// fromYear/toYear range (e.g. both 2026) could click straight past the
// declared bound into a month calendar.gsx can never render for this
// component — Task 4 review's Important finding 1.
function syncNavButtons(root, year, month) {
  const { fromYear, toYear } = navBounds(root);
  setNavDisabled(root.querySelector("[data-gsxui-calendar-prev]"), prevDisabledAt(year, month, fromYear));
  setNavDisabled(root.querySelector("[data-gsxui-calendar-next]"), nextDisabledAt(year, month, toYear));
}

// disabledRules reads the root's four disabled-rule attributes (calendar.gsx,
// Task 4) into the shape dayDisabled below expects.
function disabledRules(root) {
  return {
    before: parseISO(root.dataset.gsxuiCalendarDisabledBefore),
    after: parseISO(root.dataset.gsxuiCalendarDisabledAfter),
    dates: commaList(root.dataset.gsxuiCalendarDisabledDates)
      .map(parseISO)
      .filter((d) => d !== null),
    weekdays: commaList(root.dataset.gsxuiCalendarDisabledWeekdays).map(Number),
  };
}

// dayDisabled is calendar.gsx's own dayDisabled ported to JS: strictly
// before the before-bound, strictly after the after-bound, an exact match
// in the disabled-dates list, or a weekday in the disabled-weekdays list.
// All comparisons are UTC-midnight Date comparisons, mirroring calendar.gsx's
// own dayOnly/sameDay discipline.
function dayDisabled(date, rules) {
  if (rules.before && date < rules.before) return true;
  if (rules.after && date > rules.after) return true;
  if (rules.dates.some((d) => sameUTCDay(d, date))) return true;
  if (rules.weekdays.includes(date.getUTCDay())) return true;
  return false;
}

// selection reads the root's server-carried selection (calendar.gsx, Tasks
// 2-3) plus Task 5's own client-only hover-preview attribute. Every
// navigated month still has to reflect whatever was selected before the
// click (Task 4), and every repaint after a click/hover has to reflect the
// selection as it stands right now (Task 5's job).
function selection(root) {
  return {
    selected: commaList(root.dataset.gsxuiCalendarSelected)
      .map(parseISO)
      .filter((d) => d !== null),
    from: parseISO(root.dataset.gsxuiCalendarFrom),
    to: parseISO(root.dataset.gsxuiCalendarTo),
    // hover is set only while range mode has a from and no to yet (the
    // in-progress half of the two-click machine) — see the mouseover/
    // mouseleave handlers below. Never server-rendered, never present
    // once a real `to` exists.
    hover: parseISO(root.dataset.gsxuiCalendarHover),
  };
}

// toggleAttr sets a bare (presence-only) attribute, matching gsx.Toggle's own
// "bare name when true, absent when false" semantics for data-outside/
// data-today/data-disabled.
function toggleAttr(el, name, present) {
  if (present) el.setAttribute(name, "");
  else el.removeAttribute(name);
}

// repaint updates every one of the 42 day cells/buttons in place for
// (year, month) — never adding or removing a cell, never touching a class
// string. Mirrors calendar.gsx's own per-day computation (outside/today/
// disabled/selected/range flags) exactly, reading the rules and the current
// selection off the root's own data attributes rather than re-deriving them
// some other way.
function repaint(root, year, month) {
  const weekStartsOn = Number(root.dataset.gsxuiCalendarWeekStart ?? "0");
  const mode = root.dataset.gsxuiCalendarMode || "single";
  const grid = monthGrid(year, month, weekStartsOn);
  const rules = disabledRules(root);
  const { selected, from, to, hover } = selection(root);

  const now = new Date();
  const today = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate()));

  const cells = root.querySelectorAll('td[role="gridcell"]');
  const buttons = root.querySelectorAll("[data-gsxui-calendar-day]");

  let focusIdx = -1;
  for (let i = 0; i < grid.length; i++) {
    if (grid[i].getUTCFullYear() === year && grid[i].getUTCMonth() === month) {
      focusIdx = i;
      break;
    }
  }

  for (let i = 0; i < 42; i++) {
    const date = grid[i];
    const cell = cells[i];
    const button = buttons[i];
    const dateISO = iso(date);

    const outside = date.getUTCFullYear() !== year || date.getUTCMonth() !== month;
    const isToday = sameUTCDay(date, today);
    const disabled = dayDisabled(date, rules);

    const daySelected =
      (mode === "single" || mode === "multiple") && selected.some((s) => sameUTCDay(s, date));
    const rangeStart = mode === "range" && sameUTCDay(from, date);
    const rangeEnd = mode === "range" && sameUTCDay(to, date);
    // rangeMiddle uses the real `to` once it's set; while only `from` is set,
    // it falls back to the hover preview (`other`) so hovering ahead of or
    // behind `from` both preview correctly — sorted into lo/hi since a hover
    // can land on either side, unlike the committed from/to pair (always
    // from <= to by the time both are set — the click handler's own swap).
    // rangeStart/rangeEnd deliberately stay tied to the real from/to only:
    // the brief's own instruction is to recompute the MIDDLE markers on
    // hover, not to relabel the hovered day itself as a start/end.
    const other = to !== null ? to : hover;
    const rangeMiddle =
      mode === "range" &&
      from !== null &&
      other !== null &&
      date > (from < other ? from : other) &&
      date < (from < other ? other : from);
    const cellSelected = daySelected || rangeStart || rangeMiddle || rangeEnd;
    const selectedSingle = daySelected && !rangeStart && !rangeMiddle && !rangeEnd;

    cell.dataset.date = dateISO;
    toggleAttr(cell, "data-outside", outside);
    toggleAttr(cell, "data-today", isToday);
    toggleAttr(cell, "data-disabled", disabled);
    if (cellSelected) cell.setAttribute("data-selected", "true");
    else cell.removeAttribute("data-selected");
    cell.setAttribute("aria-selected", cellSelected ? "true" : "false");

    button.dataset.date = dateISO;
    button.textContent = String(date.getUTCDate());
    button.setAttribute("aria-label", ariaLabel(date));
    button.setAttribute("tabindex", i === focusIdx ? "0" : "-1");
    button.setAttribute("data-selected-single", selectedSingle ? "true" : "false");
    button.setAttribute("data-range-start", rangeStart ? "true" : "false");
    button.setAttribute("data-range-middle", rangeMiddle ? "true" : "false");
    button.setAttribute("data-range-end", rangeEnd ? "true" : "false");
    button.disabled = disabled;
  }
}

// updateCaption writes the same text to every element carrying
// data-gsxui-calendar-caption — one in "label" layout (the visible span),
// two in "dropdown" layout (the sr-only live-region span alongside the two
// visible selects), per calendar.gsx's own doc comment on why dropdown
// layout needs a second, textual announcement target.
function updateCaption(root, year, month) {
  const text = captionText(year, month);
  for (const el of root.querySelectorAll("[data-gsxui-calendar-caption]")) {
    el.textContent = text;
  }
}

function syncDropdowns(root, year, month) {
  const monthSelect = root.querySelector("[data-gsxui-calendar-month-select]");
  if (monthSelect) monthSelect.value = String(month);
  const yearSelect = root.querySelector("[data-gsxui-calendar-year-select]");
  if (yearSelect) yearSelect.value = String(year);
}

// goTo is the one path every navigation handler below funnels through:
// write the new month onto the root, repaint the 42 cells, sync the
// caption/dropdowns, and re-evaluate the nav buttons' own bound state —
// the same things calendar.gsx's own render computes from a fresh
// (year, month), just applied in place instead of rebuilt.
function goTo(root, year, month) {
  root.dataset.gsxuiCalendarMonth = formatMonth(year, month);
  // A hover preview left over from before the navigation would otherwise
  // bleed into the newly-painted month's rangeMiddle computation — clear it
  // defensively; the mouseleave handler below normally does this, but a
  // click on prev/next while a preview is live doesn't fire mouseleave on
  // the root at all (the pointer never actually left it).
  delete root.dataset.gsxuiCalendarHover;
  repaint(root, year, month);
  updateCaption(root, year, month);
  syncDropdowns(root, year, month);
  syncNavButtons(root, year, month);
}

on("click", "[data-gsxui-calendar-prev]", (_event, el) => {
  // Nav buttons never take a native disabled attribute (calendar.gsx's own
  // doc comment on why — only aria-disabled + tabindex="-1", so they stay
  // reachable by Tab at a navigation bound). aria-disabled doesn't block a
  // click on its own, so this handler has to check it explicitly.
  if (el.getAttribute("aria-disabled") === "true") return;
  const root = el.closest(ROOT);
  if (!root) return;
  const { year, month } = parseMonth(root.dataset.gsxuiCalendarMonth);
  const prevMonth = month === 0 ? 11 : month - 1;
  const prevYear = month === 0 ? year - 1 : year;
  goTo(root, prevYear, prevMonth);
});

on("click", "[data-gsxui-calendar-next]", (_event, el) => {
  if (el.getAttribute("aria-disabled") === "true") return;
  const root = el.closest(ROOT);
  if (!root) return;
  const { year, month } = parseMonth(root.dataset.gsxuiCalendarMonth);
  const nextMonth = month === 11 ? 0 : month + 1;
  const nextYear = month === 11 ? year + 1 : year;
  goTo(root, nextYear, nextMonth);
});

on("change", "[data-gsxui-calendar-month-select]", (_event, el) => {
  const root = el.closest(ROOT);
  if (!root) return;
  const { year } = parseMonth(root.dataset.gsxuiCalendarMonth);
  goTo(root, year, Number(el.value));
});

on("change", "[data-gsxui-calendar-year-select]", (_event, el) => {
  const root = el.closest(ROOT);
  if (!root) return;
  const { month } = parseMonth(root.dataset.gsxuiCalendarMonth);
  goTo(root, Number(el.value), month);
});

// --- selection (Task 5) ----------------------------------------------------
//
// Three independent state machines, one per mode (source map §7.1's own
// useSingle/useMulti/useRange, simplified to what the brief's Step 4 spells
// out verbatim rather than every upstream resetOnSelect/excludeDisabled
// nuance — those are real props react-day-picker exposes that gsxui's own
// ui.Calendar never surfaces, so there is no caller-visible knob to wire):
//
// - single: replace outright; re-clicking the selected day clears it.
// - multiple: toggle the clicked day into/out of a sorted list.
// - range: no from -> set from, clear to; from set, no to -> set to
//   (swapping if the click precedes from); both set -> start over with a
//   new from.
//
// Each mutator writes ISO strings straight onto the root's own
// data-gsxui-calendar-selected/-from/-to — the same attributes repaint()
// already reads via selection() above, so committing a selection is just
// "write the attribute, then repaint" (Task 4's own goTo shape, minus the
// month/caption/dropdown/nav-button steps that only apply to navigation).
// Absent, never empty-string, matching calendar.gsx's own
// omit-when-unset convention for these same three attributes.

function setSelected(root, list) {
  if (list.length) root.dataset.gsxuiCalendarSelected = list.join(",");
  else delete root.dataset.gsxuiCalendarSelected;
}

// commitSingle: clicking the already-selected day clears the selection;
// clicking any other day replaces it outright. Returns the new selected list
// (0 or 1 ISO string) for the gsxui:change detail.
function commitSingle(root, dateISO) {
  const current = commaList(root.dataset.gsxuiCalendarSelected);
  const next = current.length === 1 && current[0] === dateISO ? [] : [dateISO];
  setSelected(root, next);
  return next;
}

// commitMultiple: toggle the clicked day in the comma-separated list, kept
// sorted — ISO "YYYY-MM-DD" strings sort lexically in date order, so a
// plain Array#sort is exact date order too, no Date parsing needed.
function commitMultiple(root, dateISO) {
  const current = commaList(root.dataset.gsxuiCalendarSelected);
  const idx = current.indexOf(dateISO);
  const next = idx === -1 ? [...current, dateISO] : current.filter((d) => d !== dateISO);
  next.sort();
  setSelected(root, next);
  return next;
}

// commitRange: the two-click machine. Returns { from, to } (to is null while
// only the first click has landed) for the gsxui:change detail.
function commitRange(root, dateISO) {
  const from = root.dataset.gsxuiCalendarFrom || null;
  const to = root.dataset.gsxuiCalendarTo || null;

  let nextFrom, nextTo;
  if (from === null) {
    nextFrom = dateISO;
    nextTo = null;
  } else if (to === null) {
    // Swap if the click precedes `from` — ISO strings compare lexically in
    // date order, same reasoning as commitMultiple's sort above.
    if (dateISO < from) {
      nextFrom = dateISO;
      nextTo = from;
    } else {
      nextFrom = from;
      nextTo = dateISO;
    }
  } else {
    // Both already set: start over with a new from.
    nextFrom = dateISO;
    nextTo = null;
  }

  root.dataset.gsxuiCalendarFrom = nextFrom;
  if (nextTo !== null) root.dataset.gsxuiCalendarTo = nextTo;
  else delete root.dataset.gsxuiCalendarTo;
  return { from: nextFrom, to: nextTo };
}

// syncHiddenInputs mirrors calendar.gsx's own showHiddenSingle/-From/-To:
// a plain <input type="hidden" name={name}> for single/multiple (carrying
// the FIRST selected date, same as the server — calendar.gsx's own
// showHiddenSingle reads selected[0], not the whole list, even in multiple
// mode), and a name/name+"-to" pair for range. Each input only exists in the
// DOM at all when the caller passed `name` (and, per calendar.gsx's own
// showHidden* guards, only once there's something to put in it) — Task 5
// never creates one; if none was server-rendered there is nothing to sync,
// same "never creates or destroys" discipline the day cells themselves are
// held to.
function syncHiddenInputs(root, mode, selected, from, to) {
  const inputs = root.querySelectorAll('input[type="hidden"]');
  if (!inputs.length) return;
  if (mode === "range") {
    for (const input of inputs) {
      if (input.name.endsWith("-to")) input.value = to ?? "";
      else input.value = from ?? "";
    }
  } else {
    inputs[0].value = selected[0] ?? "";
  }
}

function currentMonth(root) {
  return parseMonth(root.dataset.gsxuiCalendarMonth);
}

on("click", "[data-gsxui-calendar-day]", (_event, el) => {
  // Disabled days ignore the click entirely and emit nothing. Real browsers
  // never dispatch "click" on a native-disabled button in the first place
  // (repaint() sets button.disabled = disabled for exactly this reason) —
  // this check is the defense-in-depth twin of that, not the only thing
  // standing between a disabled day and a selection.
  if (el.disabled) return;
  const root = el.closest(ROOT);
  if (!root) return;
  const mode = root.dataset.gsxuiCalendarMode || "single";
  const dateISO = el.dataset.date;

  let detail;
  if (mode === "multiple") {
    const selected = commitMultiple(root, dateISO);
    syncHiddenInputs(root, mode, selected, null, null);
    detail = { mode: "multiple", selected };
  } else if (mode === "range") {
    const { from, to } = commitRange(root, dateISO);
    syncHiddenInputs(root, mode, [], from, to);
    detail = { mode: "range", from, to };
  } else {
    const selected = commitSingle(root, dateISO);
    syncHiddenInputs(root, mode, selected, null, null);
    detail = { mode: "single", selected };
  }

  // A committed click always supersedes whatever hover preview was live —
  // clear it before repainting so the just-committed from/to (not a stale
  // preview) drives this repaint's rangeMiddle.
  delete root.dataset.gsxuiCalendarHover;
  const { year, month } = currentMonth(root);
  repaint(root, year, month);
  emit(root, "gsxui:change", detail);
});

// --- range hover preview (range mode only, while from is set and to is not)

on("mouseover", "[data-gsxui-calendar-day]", (_event, el) => {
  const root = el.closest(ROOT);
  if (!root) return;
  const mode = root.dataset.gsxuiCalendarMode || "single";
  if (mode !== "range") return;
  const from = root.dataset.gsxuiCalendarFrom;
  const to = root.dataset.gsxuiCalendarTo;
  if (!from || to) return;
  const dateISO = el.dataset.date;
  if (root.dataset.gsxuiCalendarHover === dateISO) return;
  root.dataset.gsxuiCalendarHover = dateISO;
  const { year, month } = currentMonth(root);
  repaint(root, year, month);
});

// mouseleave doesn't bubble — registered with { capture: true } per
// ui/gsxui.js's own header comment on exactly this case. Delegation via
// closest() still matches every mouseleave fired on any descendant of the
// root (each day button fires its own mouseleave when the pointer moves to
// a DIFFERENT day, since it leaves that button specifically even though it
// never leaves the root) — event.target !== root filters those out, the
// same guard ui/dialog.js's own backdrop-click handler uses for an
// analogous "only when it's really this element" check.
on(
  "mouseleave",
  ROOT,
  (event, root) => {
    if (event.target !== root) return;
    if (!("gsxuiCalendarHover" in root.dataset)) return;
    delete root.dataset.gsxuiCalendarHover;
    const { year, month } = currentMonth(root);
    repaint(root, year, month);
  },
  { capture: true },
);

// --- form reset (ui/combobox.js's own reflectFromBridge shape) -------------
//
// A native <form>.reset() has nothing to revert here — data-gsxui-calendar-*
// attributes aren't form-control state a browser knows how to snapshot the
// way it does an <input>'s defaultValue, so Task 5 keeps its own snapshot:
// captureDefaults below runs once per root at module load (the same "walk
// every root once at import time" shape ui/combobox.js's own init() loop
// uses), before any click has had a chance to mutate it.
const defaults = new WeakMap();

function captureDefaults(root) {
  if (defaults.has(root)) return;
  defaults.set(root, {
    selected: root.dataset.gsxuiCalendarSelected || "",
    from: root.dataset.gsxuiCalendarFrom || "",
    to: root.dataset.gsxuiCalendarTo || "",
  });
}

for (const root of document.querySelectorAll(ROOT)) captureDefaults(root);

// Scoped to a form that actually contains a calendar (`:has()`), not a bare
// "form" — ui/combobox.js's own reset handler already claims plain "form"
// for this exact (type, capture) pair, and gsxui.js's registry dispatches
// to EVERY handler whose selector matches a given element, so two modules
// both claiming literal "form" would double-fire on ANY form's reset,
// combobox or not (jstest/specs/invariants.spec.ts's own Invariant 4 — the
// same hook-prefix-collision defect class that shipped in Tier 4 Batch B
// between dropdown.js and context-menu.js, caught here for "reset" instead
// of "click"). Scoping to :has([data-gsxui-calendar]) keeps the two modules
// disjoint on every existing example page — none combine a calendar and a
// combobox/input inside one shared <form> — while still matching the form
// this handler actually needs.
on("reset", "form:has([data-gsxui-calendar])", (_event, form) => {
  for (const root of form.querySelectorAll(ROOT)) {
    const snapshot = defaults.get(root);
    if (!snapshot) continue;

    if (snapshot.selected) root.dataset.gsxuiCalendarSelected = snapshot.selected;
    else delete root.dataset.gsxuiCalendarSelected;
    if (snapshot.from) root.dataset.gsxuiCalendarFrom = snapshot.from;
    else delete root.dataset.gsxuiCalendarFrom;
    if (snapshot.to) root.dataset.gsxuiCalendarTo = snapshot.to;
    else delete root.dataset.gsxuiCalendarTo;
    delete root.dataset.gsxuiCalendarHover;

    const mode = root.dataset.gsxuiCalendarMode || "single";
    syncHiddenInputs(root, mode, commaList(snapshot.selected), snapshot.from || null, snapshot.to || null);
    const { year, month } = currentMonth(root);
    repaint(root, year, month);
  }
});
