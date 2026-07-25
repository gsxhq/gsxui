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
import { on } from "./gsxui.js";

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
// 2-3) — Task 4 only navigates months, it never changes the selection
// itself (Task 5's job), but every navigated month still has to reflect
// whatever was selected before the click.
function selection(root) {
  return {
    selected: commaList(root.dataset.gsxuiCalendarSelected)
      .map(parseISO)
      .filter((d) => d !== null),
    from: parseISO(root.dataset.gsxuiCalendarFrom),
    to: parseISO(root.dataset.gsxuiCalendarTo),
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
  const { selected, from, to } = selection(root);

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
    const rangeMiddle = mode === "range" && from !== null && to !== null && date > from && date < to;
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
// write the new month onto the root, repaint the 42 cells, and sync the
// caption/dropdowns to match — the same four things calendar.gsx's own
// render computes from a fresh (year, month), just applied in place instead
// of rebuilt.
function goTo(root, year, month) {
  root.dataset.gsxuiCalendarMonth = formatMonth(year, month);
  repaint(root, year, month);
  updateCaption(root, year, month);
  syncDropdowns(root, year, month);
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
