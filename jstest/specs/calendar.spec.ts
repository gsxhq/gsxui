import { expect, test } from "../support/fixtures";

const BASIC = "/x/calendar/basic";
const BOUNDED = "/x/calendar/bounded";
const LOADED = "/x/calendar/loaded";
const LOADED_RANGE = "/x/calendar/loadedrange";

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
test("Go and JS agree on loaded-range 2026-02 (range flags, Monday week start)", async ({
  page,
}) => {
  await page.goto(LOADED_RANGE);
  await page.click("[data-gsxui-calendar-next]");

  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-02",
  );
  const clientCells = await gridCells(page);

  await page.goto(`${LOADED_RANGE}?month=2026-02`);
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
