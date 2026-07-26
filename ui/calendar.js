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

// clientToday returns the CLIENT's own local calendar date, as a
// UTC-midnight Date so it compares directly against the UTC-midnight grid
// dates monthGrid produces (sameUTCDay below). Deliberately reads the local
// getters (getFullYear/getMonth/getDate), not the UTC ones — calendar.gsx's
// own doc comment on the Calendar component is explicit that its "today" is
// computed in the SERVER's timezone, and this is "the one place the two
// implementations may legitimately disagree" (Task 6 brief): a client in a
// different zone can already be on a different calendar date the instant
// the page loads. Reading getUTCFullYear/getUTCMonth/getUTCDate off `new
// Date()` here would silently compute today's date AT UTC, not the client's
// own local date — a different, also-wrong value that just happens to
// coincide with the client's local date for clients already in UTC.
function clientToday() {
  const now = new Date();
  return new Date(Date.UTC(now.getFullYear(), now.getMonth(), now.getDate()));
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

  const today = clientToday();
  const focusedISO = currentFocusedDate(root);

  const cells = root.querySelectorAll('td[role="gridcell"]');
  const buttons = root.querySelectorAll("[data-gsxui-calendar-day]");

  // The roving tabindex's single tab stop (source map §7.3's
  // isFocusTarget): the sticky tabStop day if one has ever been set for
  // this root (ui/toggle-group.js's own setActiveTabStop discipline — stays
  // where the user left it, even across blur), else calendar.gsx's own
  // "first day of the displayed month" fallback (firstFocusableIndex),
  // unchanged from Tasks 1-5.
  const stickyISO = currentTabStop(root);
  let focusIdx = stickyISO ? grid.findIndex((d) => iso(d) === stickyISO) : -1;
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
    button.textContent = String(date.getUTCDate());
    button.setAttribute("aria-label", ariaLabel(date));
    button.setAttribute("tabindex", i === focusIdx ? "0" : "-1");
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
    if (disabled && isFocused) {
      button.disabled = false;
      button.setAttribute("aria-disabled", "true");
    } else {
      button.disabled = disabled;
      button.removeAttribute("aria-disabled");
    }
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
// mode), and a name/name+"-to" pair for range. Each input renders whenever
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
  const inputs = root.querySelectorAll('input[type="hidden"]');
  if (!inputs.length) return;
  if (mode === "range") {
    for (const input of inputs) {
      if ("gsxuiCalendarHiddenTo" in input.dataset) input.value = to ?? "";
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
  return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate() + delta));
}

function addMonthsUTC(date, delta) {
  const targetMonthIndex = date.getUTCMonth() + delta;
  const lastDay = new Date(Date.UTC(date.getUTCFullYear(), targetMonthIndex + 1, 0)).getUTCDate();
  return new Date(
    Date.UTC(date.getUTCFullYear(), targetMonthIndex, Math.min(date.getUTCDate(), lastDay)),
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
on("keydown", "[data-gsxui-calendar-day]", (event, el) => {
  const move = KEY_MOVES[event.key];
  if (!move) return;
  const root = el.closest(ROOT);
  if (!root) return;
  const current = parseISO(el.dataset.date);
  if (!current) return;
  event.preventDefault();

  const weekStartsOn = Number(root.dataset.gsxuiCalendarWeekStart ?? "0");
  const target = move(current, event.shiftKey, weekStartsOn);
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
  const targetButton = root.querySelector(`[data-gsxui-calendar-day][data-date="${targetISO}"]`);
  if (targetButton) targetButton.focus();
});

on("focusin", "[data-gsxui-calendar-day]", (_event, el) => {
  const root = el.closest(ROOT);
  if (!root) return;
  const dateISO = el.dataset.date;
  liveFocused.set(root, dateISO);
  tabStop.set(root, dateISO);
  const { year, month } = currentMonth(root);
  repaint(root, year, month);
});

on("focusout", "[data-gsxui-calendar-day]", (_event, el) => {
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

// --- today reconciliation (Task 6) ------------------------------------------
//
// calendar.gsx's own doc comment: the server computes "today" in the
// SERVER's timezone; nothing about a client's own local date is knowable
// there. repaint() above already recomputes "today" from clientToday() (the
// client's own clock) on every call it makes — every subsequent navigation,
// click, hover, or focus change already self-corrects. This loop is what
// makes that correction land on the FIRST render too, before any of those
// have to happen: every root already on the page at module load gets one
// repaint of whatever month it's already showing, so a client already on a
// different calendar date than the server sees the right "today" marker
// immediately, not only after its first interaction.
//
// Safe to run unconditionally alongside the server's own initial render:
// every OTHER field repaint() computes (outside/selected/range/disabled/
// tabindex) is re-derived from the same root data attributes the server
// itself wrote, and reproduces the server's own values exactly — proven
// already by jstest/specs/calendar.spec.ts's "loaded survives a next/prev
// round trip back to its own initial month" test, which diffs a
// repainted-twice grid against the untouched server render and requires
// them equal.
for (const root of document.querySelectorAll(ROOT)) {
  const { year, month } = currentMonth(root);
  repaint(root, year, month);
}
