import { expect, test } from "../support/fixtures";

const BASIC = "/x/calendar/basic";
const BOUNDED = "/x/calendar/bounded";
const LOADED = "/x/calendar/loaded";
const LOADED_RANGE = "/x/calendar/loadedrange";
const RANGE = "/x/calendar/range";
const MULTIPLE = "/x/calendar/multiple";
const FORM = "/x/calendar/form";

// The grid's 42 data-date values, in DOM order.
async function gridDates(page: import("@playwright/test").Page) {
  return page.$$eval("[data-gsxui-calendar-day]", (els) =>
    els.map((el) => el.getAttribute("data-date")),
  );
}

// gridCells reads the FULL per-day attribute tuple calendar.js's repaint
// writes, not just data-date — Task 4 review's Important finding 2: a diff
// scoped to data-date alone can't catch a bug in dayDisabled, ariaLabel, the
// selection/range flags, or the outside/today computation, since none of
// those feed into data-date at all. Reads from both the button (data-date,
// aria-label, tabindex, data-selected-single, the four data-range-*, and the
// native disabled property) and its enclosing <td> (data-outside/-today/
// -disabled presence, data-selected's literal value, aria-selected) — the
// same split calendar.gsx and calendar.js both use.
async function gridCells(page: import("@playwright/test").Page) {
  return page.$$eval("[data-gsxui-calendar-day]", (buttons) =>
    buttons.map((btn) => {
      const button = btn as HTMLButtonElement;
      const cell = button.closest("td")!;
      return {
        date: button.getAttribute("data-date"),
        outside: cell.hasAttribute("data-outside"),
        today: cell.hasAttribute("data-today"),
        cellDisabled: cell.hasAttribute("data-disabled"),
        cellSelected: cell.getAttribute("data-selected"),
        ariaSelected: cell.getAttribute("aria-selected"),
        ariaLabel: button.getAttribute("aria-label"),
        tabindex: button.getAttribute("tabindex"),
        selectedSingle: button.getAttribute("data-selected-single"),
        rangeStart: button.getAttribute("data-range-start"),
        rangeMiddle: button.getAttribute("data-range-middle"),
        rangeEnd: button.getAttribute("data-range-end"),
        disabled: button.disabled,
      };
    }),
  );
}

test("the server renders a complete grid before any JS runs", async ({ page }) => {
  await page.goto(BASIC);
  const dates = await gridDates(page);
  expect(dates).toHaveLength(42);
  expect(dates[0]).toBe("2025-12-28");
  expect(dates[41]).toBe("2026-02-07");
});

test("next and previous navigate a month at a time", async ({ page }) => {
  await page.goto(BASIC);

  await page.click("[data-gsxui-calendar-next]");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-02",
  );

  await page.click("[data-gsxui-calendar-prev]");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-01",
  );
});

test("navigation never creates or destroys a cell", async ({ page }) => {
  await page.goto(BASIC);
  const before = await page.$$eval("[data-gsxui-calendar-day]", (els) => els.length);

  const first = await page.$("[data-gsxui-calendar-day]");
  await page.click("[data-gsxui-calendar-next]");

  // A real navigation happened — without this, a completely broken click
  // handler (a no-op) would pass the rest of this test too: the cell count
  // and the clicked element's identity are both unchanged by doing nothing
  // at all. Task 4 review's Minor finding 3, confirmed by the Step-4
  // red-validation control: with calendar.js's import removed entirely,
  // this test was one of the two that still passed.
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-02",
  );

  const after = await page.$$eval("[data-gsxui-calendar-day]", (els) => els.length);
  expect(after).toBe(before);

  // The same element object is still there — proof it was updated, not replaced.
  const stillAttached = await first!.evaluate((el) => el.isConnected);
  expect(stillAttached).toBe(true);
});

// The single largest risk in the design: monthGrid exists in Go and in JS and
// the two must agree. Navigate client-side to a month, then compare against
// what the server renders for that month directly. Both pages come from the
// same harness, so this is a real cross-implementation diff.
//
// Compares the FULL per-cell attribute tuple (gridCells), not just
// data-date — Task 4 review's Important finding 2.
const AGREEMENT_MONTHS = [
  { clicks: 1, month: "2026-02", why: "28-day February" },
  { clicks: 11, month: "2026-12", why: "year boundary ahead" },
  { clicks: -1, month: "2025-12", why: "year boundary behind" },
  { clicks: -23, month: "2024-02", why: "leap February" },
];

for (const { clicks, month, why } of AGREEMENT_MONTHS) {
  test(`Go and JS agree on ${month} (${why})`, async ({ page }) => {
    await page.goto(BASIC);
    const button = clicks > 0 ? "[data-gsxui-calendar-next]" : "[data-gsxui-calendar-prev]";
    for (let i = 0; i < Math.abs(clicks); i++) await page.click(button);

    await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
      "data-gsxui-calendar-month",
      month,
    );
    const clientCells = await gridCells(page);

    // Same month, rendered by Go.
    await page.goto(`${BASIC}?month=${month}`);
    const serverCells = await gridCells(page);

    expect(clientCells).toEqual(serverCells);
  });
}

// Loaded carries a real selection and all four disabled rules, in single
// mode, with a Monday week start — basic.gsx has none of these (empty
// selection, no disabled rules, Sunday week start), so every one of
// calendar.js's ported dayDisabled/daySelected/ariaLabel/outside
// computations was zero-coverage until now. Task 4 review's Important
// finding 2 and Minor finding 4: a non-Sunday week start is what makes the
// JS twin's `(x - weekStartsOn + 7) % 7` wrap-around term reachable at
// all — with weekStartsOn=0 (basic.gsx) it's dead code, since x - 0 is
// already non-negative and `% 7` never changes it.
test("Go and JS agree on loaded 2026-02 (selection, disabled rules, Monday week start)", async ({
  page,
}) => {
  await page.goto(LOADED);
  await page.click("[data-gsxui-calendar-next]");

  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-02",
  );
  const clientCells = await gridCells(page);

  await page.goto(`${LOADED}?month=2026-02`);
  const serverCells = await gridCells(page);

  expect(clientCells).toEqual(serverCells);
});

// A next-then-prev round trip forces TWO repaints (two live calendar.js
// executions) and lands back on the exact month the server rendered first —
// a stronger check than comparing a fresh, un-navigated load against itself
// (which would exercise no JS at all and pass even if calendar.js were
// entirely broken).
test("loaded survives a next/prev round trip back to its own initial month", async ({ page }) => {
  await page.goto(LOADED);
  const serverCells = await gridCells(page);

  await page.click("[data-gsxui-calendar-next]");
  await page.click("[data-gsxui-calendar-prev]");

  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-01",
  );
  const clientCells = await gridCells(page);

  expect(clientCells).toEqual(serverCells);
});

// LoadedRange exercises mode="range" specifically — rangeStart/rangeMiddle/
// rangeEnd and the cell's own cellSelected (true across the WHOLE range,
// not just the two ends) are otherwise never reached client-side, since
// basic.gsx and loaded.gsx above are both single-mode. Task 4 review's
// Important finding 2.
//
// Compares JANUARY 2026, not February: with weekStartsOn=Monday,
// monthGrid(2026, February, Monday) spans 2026-01-26..2026-03-08, which
// does not contain LoadedRange's own from=2026-01-09/to=2026-01-14 at
// all — an earlier version of this test compared February and passed
// vacuously, every range-related field uniformly false/absent on both
// sides (Task 4 review round 3's finding: the diff compared nothing to
// nothing). monthGrid(2026, January, Monday) spans
// 2025-12-29..2026-02-08, which does contain the whole range. A
// next-then-prev round trip (not a bare page load) still forces two live
// repaints, the same "exercise real JS, not just the server's initial
// render" discipline the round-trip test above uses for loaded.gsx.
test("Go and JS agree on loaded-range 2026-01 (range flags, Monday week start)", async ({
  page,
}) => {
  await page.goto(LOADED_RANGE);

  await page.click("[data-gsxui-calendar-next]");
  await page.click("[data-gsxui-calendar-prev]");

  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-01",
  );
  const clientCells = await gridCells(page);

  await page.goto(`${LOADED_RANGE}?month=2026-01`);
  const serverCells = await gridCells(page);

  expect(clientCells).toEqual(serverCells);
});

// Bounded pins fromYear=toYear=2026 — the narrowest possible navigation
// range, reachable in a handful of clicks. Task 4 review's Important
// finding 1: calendar.js never recomputed the nav bounds after a
// client-side navigation, so 11 clicks on next from 2026-01 reached
// 2026-12 (where Go WOULD mark next aria-disabled), and a 12th click landed
// on 2027-01 — a month outside the declared range that the server can
// never render for this component, and (in captionLayout="dropdown") a
// year <select> holding a value with no matching <option>.
test("next stops at the declared navigation bound instead of crossing it", async ({ page }) => {
  await page.goto(BOUNDED);
  const root = page.locator("[data-gsxui-calendar]");
  const next = page.locator("[data-gsxui-calendar-next]");
  const yearSelect = page.locator("[data-gsxui-calendar-year-select]");

  for (let i = 0; i < 11; i++) await next.click();
  await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  await expect(next).toHaveAttribute("aria-disabled", "true");
  await expect(yearSelect).toHaveValue("2026");

  // The 12th click must be a no-op: no server render exists for 2027-01,
  // and the year <select>'s own options never include it either (they run
  // navFromYear..navToYear, both 2026 here). force: true — Playwright's own
  // actionability check treats aria-disabled="true" as non-clickable and
  // would otherwise wait out the full test timeout for a state that, by
  // design, never changes; the whole point of this assertion is to attempt
  // the click anyway and confirm calendar.js's own handler guard (which
  // checks aria-disabled itself, not Playwright's actionability rules) is
  // what makes it a no-op.
  await next.click({ force: true });
  await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  await expect(yearSelect).toHaveValue("2026");
});

// The lower bound is reachable from first paint (fromYear=toYear=2026, and
// Bounded's initial month is already 2026-01) — the server itself already
// marks prev aria-disabled before any JS runs; this checks that state
// survives being clicked.
test("prev is disabled at the lower bound from the start and stays a no-op", async ({
  page,
}) => {
  await page.goto(BOUNDED);
  const root = page.locator("[data-gsxui-calendar]");
  const prev = page.locator("[data-gsxui-calendar-prev]");

  await expect(prev).toHaveAttribute("aria-disabled", "true");
  // force: true — see the equivalent comment on the "next stops at..." test
  // above: Playwright's own actionability check would otherwise wait out
  // the whole test timeout for aria-disabled to clear, which by design
  // never happens.
  await prev.click({ force: true });
  await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-01");
});

// --- Task 5: selection behavior ---------------------------------------------
//
// Both the <td role="gridcell"> and its <button data-gsxui-calendar-day>
// carry their own data-date (calendar.gsx and calendar.js's own repaint
// agree on that split — see gridCells() above), so a bare
// `[data-date="…"]` locator resolves to TWO elements and trips Playwright's
// strict mode the moment an assertion runs against it. cellFor/dayFor scope
// to whichever one actually owns the attribute under test: aria-selected/
// data-selected live on the cell, everything else (click target,
// data-range-*, data-selected-single, disabled) lives on the button.
function cellFor(page: import("@playwright/test").Page, iso: string) {
  return page.locator(`td[data-date="${iso}"]`);
}
function dayFor(page: import("@playwright/test").Page, iso: string) {
  return page.locator(`[data-gsxui-calendar-day][data-date="${iso}"]`);
}

// listenForChanges wires up the same window.__changes collection every test
// below reads from — a bare gsxui:change listener on the document, since the
// event bubbles (ui/gsxui.js's own emit()) all the way up from the root.
async function listenForChanges(page: import("@playwright/test").Page) {
  await page.evaluate(() => {
    (window as any).__changes = [];
    document.addEventListener("gsxui:change", (e) =>
      (window as any).__changes.push((e as CustomEvent).detail),
    );
  });
}

async function changes(page: import("@playwright/test").Page) {
  return page.evaluate(() => (window as any).__changes);
}

test("single mode replaces the selection and emits once", async ({ page }) => {
  await page.goto(BASIC);
  await listenForChanges(page);

  await dayFor(page, "2026-01-15").click();
  await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "true");

  await dayFor(page, "2026-01-20").click();
  await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "false");
  await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "true");

  const events = await changes(page);
  expect(events).toHaveLength(2);
  expect(events[1]).toEqual({ mode: "single", selected: ["2026-01-20"] });

  // Re-clicking the now-selected day clears it outright (single mode's own
  // "clicking the already-selected day clears the selection" rule) — a
  // third click the brief's own test doesn't cover, added here because
  // without it a commitSingle that always REPLACED (never cleared) would
  // still satisfy every assertion above.
  await dayFor(page, "2026-01-20").click();
  await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  expect(await changes(page)).toHaveLength(3);
  expect((await changes(page))[2]).toEqual({ mode: "single", selected: [] });
});

test("multiple mode toggles days", async ({ page }) => {
  await page.goto(MULTIPLE);
  await listenForChanges(page);

  await dayFor(page, "2026-01-05").click();
  await dayFor(page, "2026-01-09").click();
  expect(await page.$$eval('[aria-selected="true"]', (e) => e.length)).toBe(2);

  await dayFor(page, "2026-01-05").click();
  expect(await page.$$eval('[aria-selected="true"]', (e) => e.length)).toBe(1);
  await expect(cellFor(page, "2026-01-09")).toHaveAttribute("aria-selected", "true");

  // The toggle-off left the OTHER day selected, sorted, not just "count
  // dropped by one" — a commitMultiple that removed the wrong entry (e.g.
  // always dropping the first list element) would still pass the length
  // assertions above.
  const events = await changes(page);
  expect(events).toHaveLength(3);
  expect(events[2]).toEqual({ mode: "multiple", selected: ["2026-01-09"] });
});

test("range takes two clicks and swaps when the second precedes the first", async ({ page }) => {
  await page.goto(RANGE);
  const root = page.locator("[data-gsxui-calendar]");
  // range.gsx renders name="stay" specifically so this test also exercises
  // syncHiddenInputs's range branch — Task 5 review, Minor 1 verification:
  // the "-to" input has to be told apart from the "from" input by its own
  // data-gsxui-calendar-hidden-to marker, not a name-suffix guess.
  const from = page.locator('input[name="stay"]');
  const to = page.locator('input[name="stay-to"]');
  await expect(from).toHaveValue("");
  await expect(to).toHaveValue("");

  await dayFor(page, "2026-01-20").click();
  await expect(root).toHaveAttribute("data-gsxui-calendar-from", "2026-01-20");
  await expect(root).not.toHaveAttribute("data-gsxui-calendar-to", /.*/);
  await expect(from).toHaveValue("2026-01-20");
  await expect(to).toHaveValue("");

  await dayFor(page, "2026-01-15").click();
  await expect(root).toHaveAttribute("data-gsxui-calendar-from", "2026-01-15");
  await expect(root).toHaveAttribute("data-gsxui-calendar-to", "2026-01-20");
  await expect(dayFor(page, "2026-01-17")).toHaveAttribute("data-range-middle", "true");
  // The two endpoints themselves are NOT the middle — rangeStart/rangeEnd
  // own the boundary days instead; a rangeMiddle computed as a closed
  // interval (>= / <=, not > / <) would also mark these true.
  await expect(dayFor(page, "2026-01-15")).toHaveAttribute("data-range-middle", "false");
  await expect(dayFor(page, "2026-01-20")).toHaveAttribute("data-range-middle", "false");
  await expect(dayFor(page, "2026-01-15")).toHaveAttribute("data-range-start", "true");
  await expect(dayFor(page, "2026-01-20")).toHaveAttribute("data-range-end", "true");
  // Swapped: the FROM input carries the earlier date even though it was
  // clicked SECOND — proof the swap flows into the hidden inputs too, not
  // just the root's own data-* attributes.
  await expect(from).toHaveValue("2026-01-15");
  await expect(to).toHaveValue("2026-01-20");
});

test("range previews on hover while only the start is set", async ({ page }) => {
  await page.goto(RANGE);

  await dayFor(page, "2026-01-10").click();
  await dayFor(page, "2026-01-14").hover();
  await expect(dayFor(page, "2026-01-12")).toHaveAttribute("data-range-middle", "true");

  // Hovering elsewhere moves the preview rather than accumulating it.
  await dayFor(page, "2026-01-11").hover();
  await expect(dayFor(page, "2026-01-12")).toHaveAttribute("data-range-middle", "false");

  // Leaving the whole calendar clears the preview entirely — re-hovering
  // 2026-01-14 after a mouseleave should reproduce the same preview as the
  // very first hover, not leave 2026-01-12 stuck in whatever state the
  // in-between hover left it in.
  await dayFor(page, "2026-01-14").hover();
  await page.mouse.move(0, 0);
  await expect(dayFor(page, "2026-01-12")).toHaveAttribute("data-range-middle", "false");

  // Hovering never emits gsxui:change — only a committed click does.
  await listenForChanges(page);
  await dayFor(page, "2026-01-14").hover();
  await dayFor(page, "2026-01-11").hover();
  expect(await changes(page)).toHaveLength(0);
});

// Ambiguity resolution (per the task-5 brief's own note): rather than
// fabricating a disabled day on BASIC via evaluate() and clicking it with
// force:true, this exercises a day LOADED already disables server-side
// (2026-01-20 is in Loaded's own disabledDates — see calendar/loaded.gsx) —
// a test that fabricates the state it then checks proves less than one that
// exercises the real render.
//
// The next/prev round trip before clicking is load-bearing, not decoration:
// on a bare page load, 2026-01-20's disabled <button> attribute comes
// straight from calendar.gsx's own server render, so a real browser refuses
// to dispatch "click" on it at all regardless of what ui/calendar.js does —
// a click attempted against that untouched server HTML can't tell a correct
// repaint() apart from one that forgot `button.disabled = disabled`, since
// repaint() never even runs. Forcing one live repaint (Task 4's own
// "exercise real JS, not just the server's initial render" discipline) means
// the disabled state under test is JS-computed, so the assertion actually
// depends on ui/calendar.js getting it right.
test("a disabled day cannot be selected", async ({ page }) => {
  await page.goto(LOADED);
  await page.click("[data-gsxui-calendar-next]");
  await page.click("[data-gsxui-calendar-prev]");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-01",
  );
  await listenForChanges(page);

  await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  await expect(dayFor(page, "2026-01-20")).toBeDisabled();

  // force: true — Playwright's own actionability check refuses to click a
  // genuinely disabled button; the point of this test is that even an
  // attempted click changes nothing, not that Playwright politely declines.
  await dayFor(page, "2026-01-20").click({ force: true });

  await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  expect(await changes(page)).toHaveLength(0);
});

// calendar/form.gsx (Task 5 review, Important 1) renders name="date" with
// NOTHING preselected — real corpus coverage for the Critical fix: the
// hidden input exists (empty) from first paint, and clicking a day has to
// actually flow into it. Asserting the input's own value (not just
// aria-selected) is Important 2 — the test used to be named "…and the
// hidden input" while never once looking at one, and with the input
// existing conditionally on a selection (the pre-fix behavior), the
// BASIC-based version of this test never had one to check in the first
// place. One-line mutation that turns this from a real assertion back into
// exactly that no-op: replace syncHiddenInputs's body with a bare `return;`
// — every other assertion in this test still passes (see the report for
// the captured failure).
test("form reset clears the selection and the hidden input", async ({ page }) => {
  await page.goto(FORM);
  const input = page.locator('input[name="date"]');
  await expect(input).toHaveValue("");

  await dayFor(page, "2026-01-15").click();
  await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "true");
  await expect(input).toHaveValue("2026-01-15");

  await page.locator('button[type="reset"]').click();
  await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "false");
  await expect(input).toHaveValue("");
  await expect(page.locator("[data-gsxui-calendar]")).not.toHaveAttribute(
    "data-gsxui-calendar-selected",
    /.*/,
  );
});

// Task 5 review, Important 3: captureDefaults used to run exactly once, at
// module load, over whatever calendar roots already existed — a root an
// HTMX swap (or any other post-load insertion) adds afterward had no
// snapshot at all, so resetting ITS form silently left the user's
// selection and hidden input untouched instead of reverting them.
// ui/gsxui.js's own header promises elements added later just work; this
// is the one place calendar.js broke that promise. Simulated here by
// cloning the FORM example's own pristine calendar root into a brand-new
// <form> well after calendar.js's module-load sweep has already run — the
// clone is a distinct node instance the WeakMap has never seen. One-line
// mutation that reproduces the bug: delete the `captureDefaults(root);`
// call at the top of the day-click handler (see the report for the
// captured failure).
test("form reset restores a calendar added after the page's own initial load", async ({
  page,
}) => {
  await page.goto(FORM);

  await page.evaluate(() => {
    const original = document.querySelector("[data-gsxui-calendar]")!;
    const clone = original.cloneNode(true) as HTMLElement;
    const form = document.createElement("form");
    form.appendChild(clone);
    document.body.appendChild(form);
  });

  const added = page.locator("[data-gsxui-calendar]").nth(1);
  const addedCell = added.locator('td[data-date="2026-01-15"]');
  const addedDay = added.locator('[data-gsxui-calendar-day][data-date="2026-01-15"]');
  const addedInput = added.locator('input[name="date"]');

  await addedDay.click();
  await expect(addedCell).toHaveAttribute("aria-selected", "true");
  await expect(addedInput).toHaveValue("2026-01-15");

  await added.evaluate((el) => (el.closest("form") as HTMLFormElement).reset());
  await expect(addedCell).toHaveAttribute("aria-selected", "false");
  await expect(addedInput).toHaveValue("");
});

// --- Task 6: keyboard grid and today reconciliation -------------------------
//
// document.activeElement is the one thing no locator search can stand in
// for: several tests below exist specifically to discover WHICH day ended
// up holding focus after a key press, not to assert a guessed one already
// has it.
async function focusedDate(page: import("@playwright/test").Page) {
  return page.evaluate(() => document.activeElement?.getAttribute("data-date") ?? null);
}

// The brief's own selectors here are bare `[data-date="…"]`, which — same
// strict-mode ambiguity as every Task 5 test above — resolve to both a
// `<td>` and a `<button>` and throw the moment Playwright tries to act on
// them. dayFor is the button half (the one that's actually focusable), used
// throughout in place of the brief's bare selector.
test("arrows move by day and by week", async ({ page }) => {
  await page.goto(BASIC);
  await dayFor(page, "2026-01-15").click();
  await dayFor(page, "2026-01-15").focus();

  await page.keyboard.press("ArrowRight");
  expect(await focusedDate(page)).toBe("2026-01-16");

  await page.keyboard.press("ArrowLeft");
  expect(await focusedDate(page)).toBe("2026-01-15");

  await page.keyboard.press("ArrowDown");
  expect(await focusedDate(page)).toBe("2026-01-22");

  await page.keyboard.press("ArrowUp");
  expect(await focusedDate(page)).toBe("2026-01-15");
});

test("Home and End go to the bounds of the focused week", async ({ page }) => {
  await page.goto(BASIC);
  await dayFor(page, "2026-01-15").focus(); // a Thursday

  await page.keyboard.press("Home");
  expect(await focusedDate(page)).toBe("2026-01-11"); // Sunday

  await page.keyboard.press("End");
  expect(await focusedDate(page)).toBe("2026-01-17"); // Saturday
});

test("PageUp and PageDown move by month, with Shift by year", async ({ page }) => {
  await page.goto(BASIC);
  await dayFor(page, "2026-01-15").focus();

  await page.keyboard.press("PageDown");
  expect(await focusedDate(page)).toBe("2026-02-15");

  await page.keyboard.press("PageUp");
  expect(await focusedDate(page)).toBe("2026-01-15");

  await page.keyboard.press("Shift+PageDown");
  expect(await focusedDate(page)).toBe("2027-01-15");
});

// Task 6 review, Minor 5: the original version of this test moved
// ArrowRight from 2026-01-31 to 2026-02-01 — a date already present in
// BASIC's own January grid as trailing outside-day padding (the grid spans
// 2025-12-28..2026-02-07). Its focusedDate assertion had no teeth: a
// mutation that skips the goTo() navigation entirely still let .focus()
// find and land on that already-rendered padding-day button, so only the
// SECOND assertion (the month attribute) ever caught anything — disclosed
// honestly in the Task 6 report's own mutation-4 write-up. Retargeted at
// Shift+PageDown from 2026-01-31 to 2027-01-31 instead: that date is
// nowhere in the 42-cell window at the moment the key is pressed, so
// .focus() can only succeed AFTER a real navigation actually repaints a
// button carrying that date into existence — the first assertion now means
// what the test's name claims.
test("arrowing past the month edge navigates and keeps focus on the target", async ({ page }) => {
  await page.goto(BASIC);
  await dayFor(page, "2026-01-31").focus();

  await page.keyboard.press("Shift+PageDown");
  expect(await focusedDate(page)).toBe("2027-01-31");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2027-01",
  );
});

// Task 6 review, Important: the mouse-path twin of this ("next stops at the
// declared navigation bound instead of crossing it", above) already
// enforces fromYear/toYear against prev/next clicks; keyboard month/year
// moves (PageUp/PageDown, Shift+Page) had no equivalent clamp, so a single
// Shift+PageDown from Bounded's own January 2026 (fromYear=toYear=2026)
// could reach 2027-01 in one keystroke — a month the server can never
// render for this component, and (captionLayout="dropdown" here too) a
// year <select> left holding a value with no matching <option>, verbatim
// the defect the mouse-path test above exists to prevent.
test("keyboard month/year moves stop at the declared navigation bound instead of crossing it", async ({
  page,
}) => {
  await page.goto(BOUNDED);
  const root = page.locator("[data-gsxui-calendar]");
  const next = page.locator("[data-gsxui-calendar-next]");
  const yearSelect = page.locator("[data-gsxui-calendar-year-select]");

  await dayFor(page, "2026-01-15").focus();
  await page.keyboard.press("Shift+PageDown");

  await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  await expect(next).toHaveAttribute("aria-disabled", "true");
  await expect(yearSelect).toHaveValue("2026");
  // Clamped, not stranded: the day-of-month survives the clamp (15 is a
  // valid December day) and focus lands INSIDE the clamped month, not left
  // behind in January.
  expect(await focusedDate(page)).toBe("2026-12-15");

  // Already sitting at the bound — one more month-forward press must stay
  // a no-op, the same "attempt it anyway" discipline the mouse-path test
  // above uses for its own 12th click.
  await page.keyboard.press("PageDown");
  await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  await expect(yearSelect).toHaveValue("2026");
  expect(await focusedDate(page)).toBe("2026-12-15");
});

test("Enter and Space select the focused day", async ({ page }) => {
  await page.goto(BASIC);
  await dayFor(page, "2026-01-15").focus();
  await page.keyboard.press("Enter");
  await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "true");

  await dayFor(page, "2026-01-16").focus();
  await page.keyboard.press(" ");
  await expect(cellFor(page, "2026-01-16")).toHaveAttribute("aria-selected", "true");
});

test("exactly one day is tabbable at any time", async ({ page }) => {
  await page.goto(BASIC);
  const count = () => page.$$eval('[data-gsxui-calendar-day][tabindex="0"]', (e) => e.length);
  expect(await count()).toBe(1);

  await dayFor(page, "2026-01-15").focus();
  await page.keyboard.press("ArrowRight");
  expect(await count()).toBe(1);
});

// Source map §7.2/§8 finding 5, the first of the two behaviors the Task 6
// brief flags as easy to miss: a disabled day carries the NATIVE `disabled`
// attribute only while it is not the live-focused day; the moment it becomes
// the focused day it degrades to aria-disabled instead, so focus is never
// yanked back out of the grid. Loaded's own 2026-01-20 is disabled
// server-side (disabledDates, calendar/loaded.gsx) — arrowing onto it from
// the adjacent, non-disabled 2026-01-19 (a Monday; 2026-01-20's own
// Saturday-weekday rule and disabledDates entry don't touch it) is what
// actually puts the split under test, not fabricating disabled state via
// evaluate(). The round trip before focusing is the same "force one real
// repaint" discipline the pre-existing "a disabled day cannot be selected"
// test above already uses.
test("a disabled day keeps focus (aria-disabled, not native) and Enter is a no-op on it", async ({
  page,
}) => {
  await page.goto(LOADED);
  await page.click("[data-gsxui-calendar-next]");
  await page.click("[data-gsxui-calendar-prev]");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-01",
  );
  await listenForChanges(page);

  await dayFor(page, "2026-01-19").focus();
  await page.keyboard.press("ArrowRight");
  expect(await focusedDate(page)).toBe("2026-01-20");

  // Not .not.toBeDisabled() — Playwright's own accessibility-tree-backed
  // "disabled" state treats aria-disabled="true" as disabled too (correctly:
  // that is the entire point of the attribute), so that matcher can't tell
  // native disabled apart from the aria-disabled degradation this test
  // exists to prove. toHaveJSProperty reads the raw DOM `disabled` property
  // directly instead.
  const day = dayFor(page, "2026-01-20");
  await expect(day).toHaveJSProperty("disabled", false);
  await expect(day).toHaveAttribute("aria-disabled", "true");

  await page.keyboard.press("Enter");
  await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  expect(await changes(page)).toHaveLength(0);
});

// The second of the two behaviors the Task 6 brief flags as easy to miss:
// upstream lands focus imperatively (a ref.current?.focus() call, source map
// §3/§7.3) — moving tabindex alone would leave the browser's own focus
// wherever it already was. basic.gsx is pinned to January 2026 (Basic's own
// DefaultMonth, calendar/basic.gsx); navigating to the client's real current
// month (whatever that is when this test runs) through the public prev/next
// controls only, per the Task 6 brief's own ambiguity resolution, rather
// than reaching for a private hook — is what actually exercises the client
// path, since basic.gsx's own server render can never contain the real
// "today" for a fixed, unrelated month.
test("today is marked from the client's date", async ({ page }) => {
  await page.goto(BASIC);

  const { today, delta } = await page.evaluate(() => {
    const now = new Date();
    const iso = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
    // Months from January 2026 to the client's current month.
    const months = (now.getFullYear() - 2026) * 12 + now.getMonth();
    return { today: iso, delta: months };
  });

  const button = delta >= 0 ? "[data-gsxui-calendar-next]" : "[data-gsxui-calendar-prev]";
  for (let i = 0; i < Math.abs(delta); i++) await page.click(button);

  const marked = await page.$$eval("[data-today]", (els) =>
    els.map((el) => el.getAttribute("data-date")),
  );
  expect(marked).toEqual([today]);
});
