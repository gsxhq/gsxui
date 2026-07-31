# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> next and previous navigate a month at a time
- Location: jstest/specs/calendar.spec.ts:71:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/basic", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | const BASIC = "/x/calendar/basic";
  4   | const BOUNDED = "/x/calendar/bounded";
  5   | const LOADED = "/x/calendar/loaded";
  6   | const LOADED_RANGE = "/x/calendar/loadedrange";
  7   | const RANGE = "/x/calendar/range";
  8   | const MULTIPLE = "/x/calendar/multiple";
  9   | const FORM = "/x/calendar/form";
  10  | // The only example passing showOutsideDays={false} — see
  11  | // site/examples/calendar/outside.gsx.
  12  | const HIDDEN_OUTSIDE = "/x/calendar/hiddenoutside";
  13  | 
  14  | // The grid's 42 data-date values, in DOM order.
  15  | async function gridDates(page: import("@playwright/test").Page) {
  16  |   return page.$$eval("[data-gsxui-calendar-day]", (els) =>
  17  |     els.map((el) => el.getAttribute("data-date")),
  18  |   );
  19  | }
  20  | 
  21  | // gridCells reads the FULL per-day attribute tuple calendar.js's repaint
  22  | // writes, not just data-date — Task 4 review's Important finding 2: a diff
  23  | // scoped to data-date alone can't catch a bug in dayDisabled, ariaLabel, the
  24  | // selection/range flags, or the outside/today computation, since none of
  25  | // those feed into data-date at all. Reads from both the button (data-date,
  26  | // aria-label, tabindex, data-selected-single, the four data-range-*, and the
  27  | // native disabled property) and its enclosing <td> (data-outside/-today/
  28  | // -disabled presence, data-selected's literal value, aria-selected) — the
  29  | // same split calendar.gsx and calendar.js both use.
  30  | //
  31  | // The final review added three more fields, all of them showOutsideDays'
  32  | // (Important 1): the cell's own `data-hidden`, the button's `aria-hidden`,
  33  | // and the button's TEXT — a hidden day is blanked rather than removed, and
  34  | // none of the pre-existing fields would notice if calendar.js forgot any part
  35  | // of that on a client-navigated month.
  36  | async function gridCells(page: import("@playwright/test").Page) {
  37  |   return page.$$eval("[data-gsxui-calendar-day]", (buttons) =>
  38  |     buttons.map((btn) => {
  39  |       const button = btn as HTMLButtonElement;
  40  |       const cell = button.closest("td")!;
  41  |       return {
  42  |         date: button.getAttribute("data-date"),
  43  |         text: button.textContent,
  44  |         outside: cell.hasAttribute("data-outside"),
  45  |         hidden: cell.hasAttribute("data-hidden"),
  46  |         ariaHidden: button.getAttribute("aria-hidden"),
  47  |         today: cell.hasAttribute("data-today"),
  48  |         cellDisabled: cell.hasAttribute("data-disabled"),
  49  |         cellSelected: cell.getAttribute("data-selected"),
  50  |         ariaSelected: cell.getAttribute("aria-selected"),
  51  |         ariaLabel: button.getAttribute("aria-label"),
  52  |         tabindex: button.getAttribute("tabindex"),
  53  |         selectedSingle: button.getAttribute("data-selected-single"),
  54  |         rangeStart: button.getAttribute("data-range-start"),
  55  |         rangeMiddle: button.getAttribute("data-range-middle"),
  56  |         rangeEnd: button.getAttribute("data-range-end"),
  57  |         disabled: button.disabled,
  58  |       };
  59  |     }),
  60  |   );
  61  | }
  62  | 
  63  | test("the server renders a complete grid before any JS runs", async ({ page }) => {
  64  |   await page.goto(BASIC);
  65  |   const dates = await gridDates(page);
  66  |   expect(dates).toHaveLength(42);
  67  |   expect(dates[0]).toBe("2025-12-28");
  68  |   expect(dates[41]).toBe("2026-02-07");
  69  | });
  70  | 
  71  | test("next and previous navigate a month at a time", async ({ page }) => {
> 72  |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  73  | 
  74  |   await page.click("[data-gsxui-calendar-next]");
  75  |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  76  |     "data-gsxui-calendar-month",
  77  |     "2026-02",
  78  |   );
  79  | 
  80  |   await page.click("[data-gsxui-calendar-prev]");
  81  |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  82  |     "data-gsxui-calendar-month",
  83  |     "2026-01",
  84  |   );
  85  | });
  86  | 
  87  | test("navigation never creates or destroys a cell", async ({ page }) => {
  88  |   await page.goto(BASIC);
  89  |   const before = await page.$$eval("[data-gsxui-calendar-day]", (els) => els.length);
  90  | 
  91  |   const first = await page.$("[data-gsxui-calendar-day]");
  92  |   await page.click("[data-gsxui-calendar-next]");
  93  | 
  94  |   // A real navigation happened — without this, a completely broken click
  95  |   // handler (a no-op) would pass the rest of this test too: the cell count
  96  |   // and the clicked element's identity are both unchanged by doing nothing
  97  |   // at all. Task 4 review's Minor finding 3, confirmed by the Step-4
  98  |   // red-validation control: with calendar.js's import removed entirely,
  99  |   // this test was one of the two that still passed.
  100 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  101 |     "data-gsxui-calendar-month",
  102 |     "2026-02",
  103 |   );
  104 | 
  105 |   const after = await page.$$eval("[data-gsxui-calendar-day]", (els) => els.length);
  106 |   expect(after).toBe(before);
  107 | 
  108 |   // The same element object is still there — proof it was updated, not replaced.
  109 |   const stillAttached = await first!.evaluate((el) => el.isConnected);
  110 |   expect(stillAttached).toBe(true);
  111 | });
  112 | 
  113 | // The single largest risk in the design: monthGrid exists in Go and in JS and
  114 | // the two must agree. Navigate client-side to a month, then compare against
  115 | // what the server renders for that month directly. Both pages come from the
  116 | // same harness, so this is a real cross-implementation diff.
  117 | //
  118 | // Compares the FULL per-cell attribute tuple (gridCells), not just
  119 | // data-date — Task 4 review's Important finding 2.
  120 | const AGREEMENT_MONTHS = [
  121 |   { clicks: 1, month: "2026-02", why: "28-day February" },
  122 |   { clicks: 11, month: "2026-12", why: "year boundary ahead" },
  123 |   { clicks: -1, month: "2025-12", why: "year boundary behind" },
  124 |   { clicks: -23, month: "2024-02", why: "leap February" },
  125 | ];
  126 | 
  127 | for (const { clicks, month, why } of AGREEMENT_MONTHS) {
  128 |   test(`Go and JS agree on ${month} (${why})`, async ({ page }) => {
  129 |     await page.goto(BASIC);
  130 |     const button = clicks > 0 ? "[data-gsxui-calendar-next]" : "[data-gsxui-calendar-prev]";
  131 |     for (let i = 0; i < Math.abs(clicks); i++) await page.click(button);
  132 | 
  133 |     await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  134 |       "data-gsxui-calendar-month",
  135 |       month,
  136 |     );
  137 |     const clientCells = await gridCells(page);
  138 | 
  139 |     // Same month, rendered by Go.
  140 |     await page.goto(`${BASIC}?month=${month}`);
  141 |     const serverCells = await gridCells(page);
  142 | 
  143 |     expect(clientCells).toEqual(serverCells);
  144 |   });
  145 | }
  146 | 
  147 | for (const year of ["0000", "0099"]) {
  148 |   test(`Go and JS agree for exact year ${year}`, async ({ page }) => {
  149 |     await page.goto(`${BASIC}?month=${year}-01`);
  150 |     await page.click("[data-gsxui-calendar-next]");
  151 | 
  152 |     await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  153 |       "data-gsxui-calendar-month",
  154 |       `${year}-02`,
  155 |     );
  156 |     const clientCells = await gridCells(page);
  157 | 
  158 |     await page.goto(`${BASIC}?month=${year}-02`);
  159 |     const serverCells = await gridCells(page);
  160 | 
  161 |     expect(clientCells).toEqual(serverCells);
  162 |   });
  163 | }
  164 | 
  165 | // Loaded carries a real selection and all four disabled rules, in single
  166 | // mode, with a Monday week start — basic.gsx has none of these (empty
  167 | // selection, no disabled rules, Sunday week start), so every one of
  168 | // calendar.js's ported dayDisabled/daySelected/ariaLabel/outside
  169 | // computations was zero-coverage until now. Task 4 review's Important
  170 | // finding 2 and Minor finding 4: a non-Sunday week start is what makes the
  171 | // JS twin's `(x - weekStartsOn + 7) % 7` wrap-around term reachable at
  172 | // all — with weekStartsOn=0 (basic.gsx) it's dead code, since x - 0 is
```