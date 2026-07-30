# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> Go and JS agree on 2026-12 (year boundary ahead)
- Location: jstest/specs/calendar.spec.ts:128:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/basic", waiting until "load"

```

# Test source

```ts
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
  72  |   await page.goto(BASIC);
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
> 129 |     await page.goto(BASIC);
      |                ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  173 | // already non-negative and `% 7` never changes it.
  174 | test("Go and JS agree on loaded 2026-02 (selection, disabled rules, Monday week start)", async ({
  175 |   page,
  176 | }) => {
  177 |   await page.goto(LOADED);
  178 |   await page.click("[data-gsxui-calendar-next]");
  179 | 
  180 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  181 |     "data-gsxui-calendar-month",
  182 |     "2026-02",
  183 |   );
  184 |   const clientCells = await gridCells(page);
  185 | 
  186 |   await page.goto(`${LOADED}?month=2026-02`);
  187 |   const serverCells = await gridCells(page);
  188 | 
  189 |   expect(clientCells).toEqual(serverCells);
  190 | });
  191 | 
  192 | // A next-then-prev round trip forces TWO repaints (two live calendar.js
  193 | // executions) and lands back on the exact month the server rendered first —
  194 | // a stronger check than comparing a fresh, un-navigated load against itself
  195 | // (which would exercise no JS at all and pass even if calendar.js were
  196 | // entirely broken).
  197 | test("loaded survives a next/prev round trip back to its own initial month", async ({ page }) => {
  198 |   await page.goto(LOADED);
  199 |   const serverCells = await gridCells(page);
  200 | 
  201 |   await page.click("[data-gsxui-calendar-next]");
  202 |   await page.click("[data-gsxui-calendar-prev]");
  203 | 
  204 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  205 |     "data-gsxui-calendar-month",
  206 |     "2026-01",
  207 |   );
  208 |   const clientCells = await gridCells(page);
  209 | 
  210 |   expect(clientCells).toEqual(serverCells);
  211 | });
  212 | 
  213 | // LoadedRange exercises mode="range" specifically — rangeStart/rangeMiddle/
  214 | // rangeEnd and the cell's own cellSelected (true across the WHOLE range,
  215 | // not just the two ends) are otherwise never reached client-side, since
  216 | // basic.gsx and loaded.gsx above are both single-mode. Task 4 review's
  217 | // Important finding 2.
  218 | //
  219 | // Compares JANUARY 2026, not February: with weekStartsOn=Monday,
  220 | // monthGrid(2026, February, Monday) spans 2026-01-26..2026-03-08, which
  221 | // does not contain LoadedRange's own from=2026-01-09/to=2026-01-14 at
  222 | // all — an earlier version of this test compared February and passed
  223 | // vacuously, every range-related field uniformly false/absent on both
  224 | // sides (Task 4 review round 3's finding: the diff compared nothing to
  225 | // nothing). monthGrid(2026, January, Monday) spans
  226 | // 2025-12-29..2026-02-08, which does contain the whole range. A
  227 | // next-then-prev round trip (not a bare page load) still forces two live
  228 | // repaints, the same "exercise real JS, not just the server's initial
  229 | // render" discipline the round-trip test above uses for loaded.gsx.
```