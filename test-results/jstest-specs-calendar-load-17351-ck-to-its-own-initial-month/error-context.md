# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> loaded survives a next/prev round trip back to its own initial month
- Location: jstest/specs/calendar.spec.ts:197:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/loaded", waiting until "load"

```

# Test source

```ts
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
> 198 |   await page.goto(LOADED);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  230 | test("Go and JS agree on loaded-range 2026-01 (range flags, Monday week start)", async ({
  231 |   page,
  232 | }) => {
  233 |   await page.goto(LOADED_RANGE);
  234 | 
  235 |   await page.click("[data-gsxui-calendar-next]");
  236 |   await page.click("[data-gsxui-calendar-prev]");
  237 | 
  238 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  239 |     "data-gsxui-calendar-month",
  240 |     "2026-01",
  241 |   );
  242 |   const clientCells = await gridCells(page);
  243 | 
  244 |   await page.goto(`${LOADED_RANGE}?month=2026-01`);
  245 |   const serverCells = await gridCells(page);
  246 | 
  247 |   expect(clientCells).toEqual(serverCells);
  248 | });
  249 | 
  250 | // --- final review, Important 1: showOutsideDays -----------------------------
  251 | //
  252 | // The parameter was declared and read nowhere, so passing false had no
  253 | // observable effect. Upstream's semantics are not "drop the cell"
  254 | // (react-day-picker 9.14.0, helpers/createGetModifiers.js): showOutsideDays
  255 | // false sets modifiers.hidden, the <td> stays for layout, classNames.hidden
  256 | // is "invisible", and no DayButton renders inside it. gsxui keeps the button
  257 | // (calendar.js may never create or destroy one) and blanks it instead.
  258 | test("outside days are hidden without their cells leaving the grid", async ({ page }) => {
  259 |   await page.goto(HIDDEN_OUTSIDE);
  260 | 
  261 |   const cells = await gridCells(page);
  262 |   expect(cells).toHaveLength(42);
  263 | 
  264 |   // January 2026, Sunday week start: 2025-12-28..2026-02-07, so four leading
  265 |   // and seven trailing padding days.
  266 |   const hidden = cells.filter((c) => c.hidden);
  267 |   expect(hidden).toHaveLength(11);
  268 |   // Hidden and outside coincide exactly when showOutsideDays is false.
  269 |   expect(cells.filter((c) => c.outside).map((c) => c.date)).toEqual(hidden.map((c) => c.date));
  270 | 
  271 |   // Each hidden day is blanked: no text, no tab stop, hidden from assistive
  272 |   // tech, not activatable.
  273 |   for (const c of hidden) {
  274 |     expect(c.text).toBe("");
  275 |     expect(c.ariaHidden).toBe("true");
  276 |     expect(c.tabindex).toBe("-1");
  277 |     expect(c.disabled).toBe(true);
  278 |   }
  279 | 
  280 |   // In-month days are untouched, and the roving tab stop is still exactly
  281 |   // one, on an in-month day — hidden days are OUT of the tab sequence, the
  282 |   // opposite of the deliberate call made for disabled days.
  283 |   const shown = cells.filter((c) => !c.hidden);
  284 |   expect(shown.every((c) => c.text !== "" && c.ariaHidden === null)).toBe(true);
  285 |   const stops = cells.filter((c) => c.tabindex === "0");
  286 |   expect(stops.map((c) => c.date)).toEqual(["2026-01-01"]);
  287 | });
  288 | 
  289 | // The client half of the same fix: calendar.js re-derives `hidden` from the
  290 | // root's own data-gsxui-calendar-show-outside-days for every month the server
  291 | // never rendered, so it needs the same agreement diff every other example
  292 | // gets. A next click lands on February 2026, whose own padding set is
  293 | // completely different from January's.
  294 | //
  295 | // One-line mutation that turns this red: delete the
  296 | // `toggleAttr(cell, "data-hidden", hidden)` line from calendar.js's repaint
  297 | // (see the report for the captured failure).
  298 | test("Go and JS agree on hidden-outside 2026-02 (showOutsideDays=false)", async ({ page }) => {
```