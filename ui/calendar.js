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

const ROOT = "[data-gsxui-slot-calendar]";

// Date.UTC treats years 0..99 as 1900..1999. Setting the full year on an
// existing UTC date preserves the caller's exact year while retaining the
// useful month/day overflow normalization calendar arithmetic relies on.
function utcDate(year, month, day) {
  const date = new Date(0);
  date.setUTCHours(0, 0, 0, 0);
  date.setUTCFullYear(year, month, day);
  return date;
}

// monthGrid is the twin of calendar.gsx's monthGrid: 42 dates, six rows of
// seven, starting at the week start on or before the 1st. month is 0-based
// (JS Date convention — January is 0), matching every other date this module
// touches.
function monthGrid(year, month, weekStartsOn) {
  const first = utcDate(year, month, 1);
  const offset = (first.getUTCDay() - weekStartsOn + 7) % 7;
  const grid = [];
  for (let i = 0; i < 42; i++) {
    grid.push(utcDate(year, month, 1 - offset + i));
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
  const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(value);
  if (!match) return null;
  const [, year, month, day] = match;
  const [y, m, d] = [year, month, day].map(Number);
  const date = utcDate(y, m - 1, d);
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
  const year = String(date.getUTCFullYear()).padStart(4, "0");
  return `${weekday}, ${month} ${date.getUTCDate()}, ${year}`;
}

function captionText(year, month) {
  return `${MONTH_NAMES[month]} ${String(year).padStart(4, "0")}`;
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
  setNavDisabled(root.querySelector("[data-gsxui-slot-calendar-previous]"), prevDisabledAt(year, month, fromYear));
  setNavDisabled(root.querySelector("[data-gsxui-slot-calendar-next]"), nextDisabledAt(year, month, toYear));
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

// toggleAttr sets a bare (presence-only) attribute, matching the server's
// data-outside/data-today/data-disabled rendering: a bool on a data-* name
// always renders as presence — bare when true, absent when false.
function toggleAttr(el, name, present) {
  if (present) el.setAttribute(name, "");
  else el.removeAttribute(name);
}

// clientToday returns the CLIENT's own local calendar date, as a
// UTC-midnight Date so it compares directly against the UTC-midnight grid
// dates monthGrid produces (sameUTCDay below). Deliberately reads the local
// getters (getFullYear/getMonth/getDate), not the UTC ones.
//
// Task 6 review, Minor 3: the server's own "today" (calendar.gsx:557,
// `time.Now().UTC()`) is the UTC calendar date, not "the server's own
// timezone" — the server has no timezone of its own in this computation at
// all, it's pinned to UTC unconditionally. The two implementations can
// still legitimately disagree, just not for the reason an earlier version
// of this comment claimed: any CLIENT whose local zone isn't UTC can
// already be on a different calendar date than the UTC one the instant the
// page loads (e.g. a client west of UTC, still on "yesterday" at 8pm local
// while UTC has already rolled over to today). Reading
// getUTCFullYear/getUTCMonth/getUTCDate off `new Date()` here would
// silently recompute the SAME UTC date the server already used — not a
// bug exactly, just a no-op that could never reconcile anything.
function clientToday() {
  const now = new Date();
  return utcDate(now.getFullYear(), now.getMonth(), now.getDate());
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
  // calendar.gsx serializes the resolved flag unconditionally as the literal
  // "true"/"false" (its own boolStr), so an exact === "true" test is the
  // faithful twin — no `|| default` fallback, because there is no
  // omit-when-unset case to fall back FROM.
  const showOutsideDays = root.dataset.gsxuiCalendarShowOutsideDays === "true";
  const grid = monthGrid(year, month, weekStartsOn);
  const rules = disabledRules(root);
  const { selected, from, to, hover } = selection(root);

  const today = clientToday();
  const focusedISO = currentFocusedDate(root);

  const cells = root.querySelectorAll('td[role="gridcell"]');
  const buttons = root.querySelectorAll("[data-gsxui-slot-calendar-day-button]");

  // The roving tabindex's single tab stop (source map §7.3's
  // isFocusTarget): the sticky tabStop day if one has ever been set for
  // this root (ui/toggle-group.js's own setActiveTabStop discipline — stays
  // where the user left it, even across blur), else calendar.gsx's own
  // "first day of the displayed month" fallback (firstFocusableIndex),
  // unchanged from Tasks 1-5.
  //
  // Task 6 review, Minor 4: the sticky day is only honored while it's
  // still INSIDE the displayed month. A day that was the tab stop before a
  // navigation can end up rendered as an OUTSIDE padding cell of the newly
  // displayed month (e.g. focus 2026-01-31, then click next on a Monday-
  // week-start calendar — February's own grid still contains 2026-01-31 as
  // a greyed leading day) — calendar.gsx's own firstFocusableIndex would
  // never place tabindex="0" on an outside day for a fresh render of that
  // month, so honoring a stale sticky value there is a real Go/JS
  // divergence, not a cosmetic one (tabindex is part of the agreement
  // tuple gridCells() reads). Deliberately NOT also rejecting a sticky day
  // for being disabled: calendar.gsx's own firstFocusableIndex doesn't
  // filter disabled either (Tasks 1-5's already-narrower-than-upstream
  // scope), and the "a disabled day keeps focus" test above depends on a
  // day the user just arrowed onto staying the sticky tab stop even though
  // dayDisabled(day, rules) is true for it.
  const stickyISO = currentTabStop(root);
  let focusIdx = -1;
  if (stickyISO) {
    const idx = grid.findIndex((d) => iso(d) === stickyISO);
    if (idx !== -1 && grid[idx].getUTCFullYear() === year && grid[idx].getUTCMonth() === month) {
      focusIdx = idx;
    }
  }
  if (focusIdx === -1) {
    for (let i = 0; i < grid.length; i++) {
      if (grid[i].getUTCFullYear() === year && grid[i].getUTCMonth() === month) {
        focusIdx = i;
        break;
      }
    }
  }

  for (let i = 0; i < 42; i++) {
    const date = grid[i];
    const cell = cells[i];
    const button = buttons[i];
    const dateISO = iso(date);

    const outside = date.getUTCFullYear() !== year || date.getUTCMonth() !== month;
    // calendar.gsx's own `hiddenDay := outside && !showOutsideDays`.
    const hidden = outside && !showOutsideDays;
    const isToday = sameUTCDay(date, today);
    const disabled = dayDisabled(date, rules);
    const isFocused = focusedISO !== null && dateISO === focusedISO;

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
    toggleAttr(cell, "data-hidden", hidden);
    toggleAttr(cell, "data-today", isToday);
    toggleAttr(cell, "data-disabled", disabled);
    // data-focused drives the button's own group-data-[focused=true]/day:
    // ring styling (source map §3) via the cell's group/day class — set
    // only while this day holds LIVE DOM focus (cleared on blur, see the
    // focusin/focusout handlers below), not the sticky tab-stop day.
    toggleAttr(cell, "data-focused", isFocused);
    if (cellSelected) cell.setAttribute("data-selected", "true");
    else cell.removeAttribute("data-selected");
    cell.setAttribute("aria-selected", cellSelected ? "true" : "false");

    button.dataset.date = dateISO;
    // A hidden day's button is blanked, not removed (calendar.js never
    // creates or destroys an element; upstream simply skips rendering the
    // DayButton, which this port cannot do). Empty text + aria-hidden +
    // native disabled + tabindex="-1" is the closest reachable equivalent:
    // nothing to read, nothing to click, and — unlike DISABLED days, which
    // deliberately stay in the roving sequence — nothing to tab onto either.
    // An invisible focusable button is a screen-reader trap.
    button.textContent = hidden ? "" : String(date.getUTCDate());
    button.setAttribute("aria-label", ariaLabel(date));
    if (hidden) button.setAttribute("aria-hidden", "true");
    else button.removeAttribute("aria-hidden");
    const isRovingStop = i === focusIdx && !hidden;
    button.setAttribute("tabindex", isRovingStop ? "0" : "-1");
    button.setAttribute("data-selected-single", selectedSingle ? "true" : "false");
    button.setAttribute("data-range-start", rangeStart ? "true" : "false");
    button.setAttribute("data-range-middle", rangeMiddle ? "true" : "false");
    button.setAttribute("data-range-end", rangeEnd ? "true" : "false");
    // Source map §7.2/§8 finding 5: a disabled day carries the native
    // `disabled` attribute only while it is NOT the live-focused day. A
    // disabled day that IS currently focused degrades to aria-disabled
    // instead, so focus is never yanked out of the grid — the click
    // handler below checks aria-disabled explicitly for exactly this case,
    // since a native-disabled=false button dispatches "click" normally.
    //
    // A HIDDEN day always takes the native attribute regardless (mirroring
    // calendar.gsx's own `disabled={ dayDis || hiddenDay }`): the
    // degrade-to-aria-disabled branch exists so a focused day never loses
    // focus, and a hidden day can never be the focused one — it is out of
    // the roving sequence entirely, by exactly the reasoning above.
    if (disabled && (isFocused || isRovingStop) && !hidden) {
      button.disabled = false;
      button.setAttribute("aria-disabled", "true");
    } else {
      button.disabled = disabled || hidden;
      button.removeAttribute("aria-disabled");
    }
  }

  // Upstream sets the grid's accessible name from labelGrid (the formatted
  // month/year) — calendar.gsx renders it as aria-label={captionText}, and
  // it has to follow a client-side navigation the same way the caption text
  // itself does, or the grid keeps announcing the month it was
  // server-rendered for.
  const gridEl = root.querySelector("[data-gsxui-slot-calendar-grid]");
  if (gridEl) gridEl.setAttribute("aria-label", captionText(year, month));
}

// updateCaption writes the same text to every element carrying
// data-gsxui-slot-calendar-caption — one in "label" layout (the visible span),
// two in "dropdown" layout (the sr-only live-region span alongside the two
// visible selects), per calendar.gsx's own doc comment on why dropdown
// layout needs a second, textual announcement target.
function updateCaption(root, year, month) {
  const text = captionText(year, month);
  for (const el of root.querySelectorAll("[data-gsxui-slot-calendar-caption]")) {
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

on("click", "[data-gsxui-slot-calendar-previous]", (_event, el) => {
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

on("click", "[data-gsxui-slot-calendar-next]", (_event, el) => {
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

// syncHiddenInputs mirrors calendar.gsx's hidden form bridge: one input for
// single, repeated same-name inputs for every multiple selection, and a
// name/name+"-to" pair for range. Each mode keeps an empty placeholder when
// unselected so the caller's field remains present in FormData.
// Each input renders whenever
// the caller passed `name` — unconditionally, value empty until there's a
// real selection (Task 5 review, Critical: an earlier revision on both the
// Go and JS sides only synced/rendered once something was already selected,
// so `<ui.Calendar mode="single" name="date"/>` with nothing preselected —
// the single most likely way a caller wires this into a form — submitted
// with no "date" field at all; range mode had the same hole for "-to"
// whenever only `from` was set). Task 5 still never CREATES an input if the
// caller never passed `name` at all — same "never creates or destroys"
// discipline the day cells themselves are held to — this only stopped
// requiring a pre-existing selection to sync the one(s) that do exist.
//
// The range pair is told apart by the `to` input's own
// data-gsxui-calendar-hidden-to marker, not by `input.name.endsWith("-to")`
// (Task 5 review, Minor 1): a caller's own `name` can legitimately end in
// "-to" (e.g. name="valid-to"), and a suffix check would then misidentify
// that caller's FROM input as the TO input.
function syncHiddenInputs(root, mode, selected, from, to) {
  if (mode === "range") {
    const inputs = root.querySelectorAll('input[type="hidden"]');
    for (const input of inputs) {
      if ("gsxuiCalendarHiddenTo" in input.dataset) input.value = to ?? "";
      else input.value = from ?? "";
    }
    return;
  }
  if (mode === "multiple") {
    const inputs = [...root.querySelectorAll("[data-gsxui-calendar-hidden-multiple]")];
    if (!inputs.length) return;
    const values = selected.length ? selected : [""];
    const name = inputs[0].name;
    while (inputs.length < values.length) {
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = name;
      input.setAttribute("data-gsxui-calendar-hidden-multiple", "");
      root.append(input);
      inputs.push(input);
    }
    for (let i = 0; i < values.length; i++) inputs[i].value = values[i];
    for (let i = values.length; i < inputs.length; i++) inputs[i].remove();
    return;
  }
  const input = root.querySelector('input[type="hidden"]');
  if (input) input.value = selected[0] ?? "";
}

function currentMonth(root) {
  return parseMonth(root.dataset.gsxuiCalendarMonth);
}

on("click", "[data-gsxui-slot-calendar-day-button]", (_event, el) => {
  // Disabled days ignore the click entirely and emit nothing. Real browsers
  // never dispatch "click" on a native-disabled button in the first place
  // (repaint() sets button.disabled = disabled for exactly this reason) —
  // this check is the defense-in-depth twin of that, not the only thing
  // standing between a disabled day and a selection. The aria-disabled
  // check covers the one case el.disabled can't: a disabled day that
  // currently holds live focus degrades to aria-disabled instead of
  // native disabled (source map §7.2/§8 finding 5, repaint() above), so a
  // real click (or the Enter/Space activation the native <button> element
  // synthesizes into one) DOES reach this handler for that day and has to
  // be turned away explicitly.
  if (el.disabled || el.getAttribute("aria-disabled") === "true") return;
  const root = el.closest(ROOT);
  if (!root) return;
  // A no-op the first time any given root is clicked (captureDefaults
  // itself guards on the WeakMap already having an entry) — but the FIRST
  // time matters for a root that didn't exist at module load, e.g. one an
  // HTMX swap inserted after the page's own initial captureDefaults sweep
  // (Task 5 review, Important 3: ui/gsxui.js's own header promises elements
  // added later just work; a snapshot taken once at import time and never
  // revisited broke that promise for form reset specifically — a calendar
  // added post-load had no snapshot at all, so resetting its form silently
  // left the user's selection and hidden input(s) untouched instead of
  // reverting them). Capturing here, before commitSingle/commitMultiple/
  // commitRange mutate anything below, guarantees the snapshot is always
  // this root's PRISTINE state, never a post-click one.
  captureDefaults(root);
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

on("mouseover", "[data-gsxui-slot-calendar-day-button]", (_event, el) => {
  // The same guard the day-click handler opens with, for the same reason:
  // a day that cannot be SELECTED must not paint a preview of a range that
  // ending on it would produce. Without this, hovering a disabled (or
  // hidden) day dragged a data-range-middle band right through it, promising
  // a selection the click path would then refuse. Native-disabled buttons
  // don't dispatch mouse events in most browsers, but the aria-disabled
  // degradation (a disabled day holding focus) and the hidden days both
  // reach here otherwise.
  if (el.disabled || el.getAttribute("aria-disabled") === "true") return;
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

// --- keyboard grid (Task 6) -------------------------------------------------
//
// Roving tabindex + arrow-key navigation, ported from source map §7.3
// (`useFocus.js`, `helpers/{calculateFocusTarget,getFocusableDate,
// getNextFocus}.js`, `DayPicker.js`'s own `handleDayKeyDown`).
//
// Two distinct per-root flags, easy to conflate (source map §7.3's own
// warning) — kept as two separate WeakMaps rather than one, precisely
// because they react to different events and drive different attributes:
//
// - liveFocused: which day currently holds REAL DOM focus right now,
//   cleared the moment it blurs. Drives the cell's data-focused (hence the
//   button's own group-data-[focused=true]/day: ring, source map §3) and
//   the disabled/aria-disabled split in repaint() above (source map §7.2/
//   §8 finding 5).
// - tabStop: the roving-tabindex's single tabindex="0" day. Sticky across
//   blur — once a day has held focus it stays the sole tab stop until
//   another one takes over, the same "stays where the user left it"
//   discipline ui/toggle-group.js's own setActiveTabStop already uses, so
//   re-tabbing into the grid lands back where the user was, not always the
//   first day of the month. This is source map §7.3's own `lastFocused`,
//   simplified: gsxui exposes no controlled `modifiers` prop, so the
//   "explicit modifiers.focused" priority level above lastFocused never
//   applies here, and this port does not implement the "selected day" /
//   "today" fallback levels below lastFocused either — no test exercises
//   them, and calendar.gsx's own firstFocusableIndex (Tasks 1-5, still the
//   fallback once tabStop is unset) was already scoped to "first day of
//   the displayed month" rather than the full upstream priority chain.
const liveFocused = new WeakMap();
const tabStop = new WeakMap();

function currentFocusedDate(root) {
  return liveFocused.get(root) ?? null;
}
function currentTabStop(root) {
  return tabStop.get(root) ?? null;
}

// addDaysUTC/addMonthsUTC/addYearsUTC port getFocusableDate's own date
// arithmetic (source map §7.3): day/week moves are exact ±N-day shifts;
// month/year moves clamp the resulting day-of-month to the last day of the
// target month (date-fns' own addMonths/addYears semantics — the dateLib
// getFocusableDate itself delegates to) rather than letting it roll over
// the way plain Date.UTC(year, month + delta, day) arithmetic would (Jan 31
// PageDown would otherwise land on Mar 3, not Feb 28).
function addDaysUTC(date, delta) {
  return utcDate(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate() + delta);
}

// lastDayOfMonthUTC returns the number of days in (year, month) — day 0 of
// the FOLLOWING month is always the last day of THIS one, in any calendar
// arithmetic that normalizes automatically (JS Date.UTC does).
function lastDayOfMonthUTC(year, month) {
  return utcDate(year, month + 1, 0).getUTCDate();
}

function addMonthsUTC(date, delta) {
  const targetMonthIndex = date.getUTCMonth() + delta;
  const lastDay = lastDayOfMonthUTC(date.getUTCFullYear(), targetMonthIndex);
  return utcDate(
    date.getUTCFullYear(),
    targetMonthIndex,
    Math.min(date.getUTCDate(), lastDay),
  );
}

function addYearsUTC(date, delta) {
  return addMonthsUTC(date, delta * 12);
}

function startOfWeekUTC(date, weekStartsOn) {
  const diff = (date.getUTCDay() - weekStartsOn + 7) % 7;
  return addDaysUTC(date, -diff);
}

function endOfWeekUTC(date, weekStartsOn) {
  return addDaysUTC(startOfWeekUTC(date, weekStartsOn), 6);
}

// clampToNavBounds ports getFocusableDate's own clamp-to-navStart/navEnd
// behavior (source map §7.3: "clamped to navStart/navEnd if the move would
// cross them") — Task 6 review's Important finding: without this, a
// keyboard move could carry the displayed month straight past a caller's
// declared fromYear/toYear range (calendar.gsx's own calendarNavBounds,
// mirrored client-side by navBounds() above), exactly the defect Task 4's
// review fixed for the mouse path (syncNavButtons/prev/next's own
// aria-disabled guard).
//
// This is a CLAMP, not a re-projection: upstream's getFocusableDate ends with
// `focusableDate = max([navStart, focusableDate])` (moving backward) or
// `min([navEnd, focusableDate])` (moving forward) — it returns the BOUND
// ITSELF, not the out-of-range target's day-of-month transplanted into the
// bound month. The final review caught the earlier version doing the latter:
// it landed on `Date.UTC(fromYear, 0, min(date.getUTCDate(), 31))`, so on the
// shipped /x/calendar/bounded (fromYear=toYear=2026, January displayed),
// focusing 2026-01-01 and pressing ArrowUp produced a raw 2025-12-25 that
// "clamped" to 2026-01-25 — jumping focus 24 days FORWARD out of a backward
// move, and past the very bound it was meant to hold. Same shape at the
// upper bound. Only the YEAR can go out of range here (month is normalized
// to 0-11 by every KEY_MOVES entry), so the two bounds are the first day of
// fromYear's January and the last day of toYear's December.
function clampToNavBounds(date, fromYear, toYear) {
  const year = date.getUTCFullYear();
  if (year < fromYear) return utcDate(fromYear, 0, 1);
  if (year > toYear) return utcDate(toYear, 11, lastDayOfMonthUTC(toYear, 11));
  return date;
}

// KEY_MOVES ports handleDayKeyDown's own keyMap verbatim (source map §7.3's
// table). gsxui exposes no rtl/dir prop, so ArrowLeft/ArrowRight always take
// the table's non-rtl branch. Every key NOT in this map is left completely
// untouched — no preventDefault, no handling — matching upstream's own `if
// (keyMap[e.key])` guard; this is also what lets Enter/Space fall through
// to the native <button> element's own default activation (which
// synthesizes a "click", handled by the day-click handler above) rather
// than being reimplemented here.
const KEY_MOVES = {
  ArrowLeft: (date, shift) => (shift ? addMonthsUTC(date, -1) : addDaysUTC(date, -1)),
  ArrowRight: (date, shift) => (shift ? addMonthsUTC(date, 1) : addDaysUTC(date, 1)),
  ArrowDown: (date, shift) => (shift ? addYearsUTC(date, 1) : addDaysUTC(date, 7)),
  ArrowUp: (date, shift) => (shift ? addYearsUTC(date, -1) : addDaysUTC(date, -7)),
  PageDown: (date, shift) => (shift ? addYearsUTC(date, 1) : addMonthsUTC(date, 1)),
  PageUp: (date, shift) => (shift ? addYearsUTC(date, -1) : addMonthsUTC(date, -1)),
  Home: (date, _shift, weekStartsOn) => startOfWeekUTC(date, weekStartsOn),
  End: (date, _shift, weekStartsOn) => endOfWeekUTC(date, weekStartsOn),
};

// Disabled days deliberately stay IN the roving sequence (ambiguity
// resolved by the Task 6 brief, a deviation from upstream's own
// getNextFocus, which recurses past them): arrow keys move onto a disabled
// day same as any other, and the day-click handler's own aria-disabled
// check is what keeps Enter/Space on one a no-op. Skipping them here would
// strand a keyboard user on whichever side of a disabled span they
// approached from.
//
// HIDDEN days (showOutsideDays=false) are the opposite call, and never come
// up here anyway: a hidden day is by definition an outside day, so any move
// that resolves onto one crosses a month boundary and goTo repaints it as an
// ordinary in-month day of the newly displayed month before focus lands. The
// roving sequence they are excluded from is the TAB one (tabindex="0"), in
// repaint above — not this one.
on("keydown", "[data-gsxui-slot-calendar-day-button]", (event, el) => {
  const move = KEY_MOVES[event.key];
  if (!move) return;
  const root = el.closest(ROOT);
  if (!root) return;
  const current = parseISO(el.dataset.date);
  if (!current) return;
  event.preventDefault();

  const weekStartsOn = Number(root.dataset.gsxuiCalendarWeekStart ?? "0");
  const raw = move(current, event.shiftKey, weekStartsOn);
  // Task 6 review, Important: clamp to the SAME nav bounds the prev/next
  // buttons already enforce (navBounds() reads calendar.gsx's own resolved
  // data-gsxui-calendar-nav-from-year/-to-year) before this move is allowed
  // to change the displayed month at all — PageDown/Shift+PageDown could
  // otherwise carry a caller's narrow fromYear/toYear range (bounded.gsx:
  // fromYear=toYear=2026) straight past its declared bound in one keystroke.
  const { fromYear, toYear } = navBounds(root);
  const target = clampToNavBounds(raw, fromYear, toYear);
  const targetISO = iso(target);

  // Pre-register the target as this root's live-focused/tab-stop day BEFORE
  // the repaint below runs, not after — a target that's a disabled day
  // needs its OWN repaint to already see itself as focused, so it takes the
  // aria-disabled branch (source map §7.2/§8 finding 5) instead of the
  // native-disabled one. A native-disabled button cannot receive .focus()
  // at all; calling it while still native-disabled would silently strand
  // focus on the day the user started the move from instead of landing it
  // on the target.
  liveFocused.set(root, targetISO);
  tabStop.set(root, targetISO);

  // Crossing a month boundary navigates the grid FIRST (source map §7.3:
  // moveFocus also calls calendar.goToDay(nextFocus) whenever the resolved
  // target falls outside the currently-shown month) — goTo's own repaint
  // both paints the target's own month (so its button exists at all, since
  // the target's data-date must actually appear among the 42 cells before
  // the imperative focus() call below can find it) and, since liveFocused/
  // tabStop are already set above, gets the target's focused/tabindex/
  // disabled state right in that same pass. Staying within the same month
  // still needs its own repaint for that second part, hence the else.
  const { year, month } = currentMonth(root);
  if (target.getUTCFullYear() !== year || target.getUTCMonth() !== month) {
    goTo(root, target.getUTCFullYear(), target.getUTCMonth());
  } else {
    repaint(root, year, month);
  }

  // Focus lands imperatively (source map §3/§7.3: upstream calls
  // ref.current?.focus() whenever a day becomes the focused one — moving
  // tabindex alone does not move focus).
  const targetButton = root.querySelector(`[data-gsxui-slot-calendar-day-button][data-date="${targetISO}"]`);
  if (targetButton) targetButton.focus();
});

on("focusin", "[data-gsxui-slot-calendar-day-button]", (_event, el) => {
  const root = el.closest(ROOT);
  if (!root) return;
  const dateISO = el.dataset.date;
  liveFocused.set(root, dateISO);
  tabStop.set(root, dateISO);
  const { year, month } = currentMonth(root);
  repaint(root, year, month);
});

on("focusout", "[data-gsxui-slot-calendar-day-button]", (_event, el) => {
  const root = el.closest(ROOT);
  if (!root) return;
  // A stale event: this root's live-focused day has already moved on (the
  // keydown handler's own goTo repaint, when a move crosses a month
  // boundary, rewrites every button's data-date — including this one —
  // before the new target is ever focused, so by the time this button's
  // blur fires its data-date no longer matches what liveFocused still
  // holds). Nothing to clear in that case; the focusin firing right after
  // this one, for the new target, will set liveFocused correctly.
  if (liveFocused.get(root) !== el.dataset.date) return;
  liveFocused.delete(root);
  const { year, month } = currentMonth(root);
  repaint(root, year, month);
});

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
// "form". A form may also contain a combobox and then intentionally matches
// both modules for reset:false: each handler restores disjoint descendants.
// That exact pair is reviewed in selector-allowlist.ts and exercised by the
// mixed calendar/form example; all other forms remain module-scoped.
on("reset", "form:has([data-gsxui-slot-calendar])", (_event, form) => {
  // The reset event and its microtask checkpoint both run before the
  // browser's default reset action. Reconcile in the next task, once native
  // controls and this component's custom DOM can reach one final state.
  setTimeout(() => {
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
      liveFocused.delete(root);
      tabStop.delete(root);

      const mode = root.dataset.gsxuiCalendarMode || "single";
      syncHiddenInputs(
        root,
        mode,
        commaList(snapshot.selected),
        snapshot.from || null,
        snapshot.to || null,
      );
      const today = clientToday();
      goTo(root, today.getUTCFullYear(), today.getUTCMonth());
    }
  }, 0);
});

// --- today reconciliation (Task 6) ------------------------------------------
//
// repaint() above already recomputes "today" from clientToday() (the
// client's own local date, Minor 3's corrected doc comment on clientToday()
// explains what can actually disagree) on every call it makes — every
// subsequent navigation, click, hover, or focus change already
// self-corrects. This loop is what makes that correction land on the FIRST
// render too, before any of those have to happen.
//
// Task 6 review, CRITICAL finding: an earlier version of this loop called
// the FULL repaint(root, year, month) here, not just a data-today
// reconciliation — which repaints every field (outside/selected/range/
// disabled/tabindex), not just today, for every root on EVERY page at
// module load, before any test ever interacts with the page. That broke
// jstest/specs/calendar.spec.ts's own "Go and JS agree" tests without
// making any of them fail: gridCells() reads only fields repaint() writes,
// so once the page underneath every "serverCells" snapshot was ALREADY a
// JS repaint (not the server's own untouched render), every one of those
// tests degenerated into comparing repaint(root, y, m) against itself on
// identical inputs — a tautology that can no longer catch monthGrid's two
// twins (this file's own and calendar.gsx's) diverging, which is this
// design's single largest named risk. (The safety argument that version's
// comment made — citing the "loaded survives a next/prev round trip" test
// as proof repaint() reproduces the server's own render exactly — was
// circular for the same reason: that test's own "untouched server render"
// baseline was itself already repainted by this loop.)
//
// reconcileToday fixes this by touching ONLY data-today, on the 42 cells
// of whatever month a root already has server-rendered, leaving every
// other field (outside/selected/range/disabled/tabindex) exactly as the
// server wrote it — genuinely untouched, so the agreement tests stay real
// diffs against a real server render.
function reconcileToday(root, year, month) {
  const weekStartsOn = Number(root.dataset.gsxuiCalendarWeekStart ?? "0");
  const grid = monthGrid(year, month, weekStartsOn);
  const today = clientToday();
  const cells = root.querySelectorAll('td[role="gridcell"]');
  for (let i = 0; i < 42; i++) {
    toggleAttr(cells[i], "data-today", sameUTCDay(grid[i], today));
  }
}

for (const root of document.querySelectorAll(ROOT)) {
  const { year, month } = currentMonth(root);
  reconcileToday(root, year, month);
}
