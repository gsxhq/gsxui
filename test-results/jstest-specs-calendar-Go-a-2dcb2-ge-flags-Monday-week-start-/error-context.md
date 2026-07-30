# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> Go and JS agree on loaded-range 2026-01 (range flags, Monday week start)
- Location: jstest/specs/calendar.spec.ts:230:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/loadedrange", waiting until "load"

```

# Test source

```ts
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
  230 | test("Go and JS agree on loaded-range 2026-01 (range flags, Monday week start)", async ({
  231 |   page,
  232 | }) => {
> 233 |   await page.goto(LOADED_RANGE);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  299 |   await page.goto(HIDDEN_OUTSIDE);
  300 |   await page.click("[data-gsxui-calendar-next]");
  301 | 
  302 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  303 |     "data-gsxui-calendar-month",
  304 |     "2026-02",
  305 |   );
  306 |   const clientCells = await gridCells(page);
  307 |   // Not a vacuous diff: February 2026 really does have padding days to hide.
  308 |   expect(clientCells.filter((c) => c.hidden).length).toBeGreaterThan(0);
  309 | 
  310 |   await page.goto(`${HIDDEN_OUTSIDE}?month=2026-02`);
  311 |   const serverCells = await gridCells(page);
  312 | 
  313 |   expect(clientCells).toEqual(serverCells);
  314 | });
  315 | 
  316 | // role="grid" suppresses a <table>'s implicit naming, so the grid needs an
  317 | // explicit accessible name — upstream's labelGrid, defaulting to the
  318 | // formatted month/year. calendar.js's repaint has to keep it in step with
  319 | // navigation or it keeps announcing the server-rendered month forever.
  320 | test("the grid's accessible name follows the displayed month", async ({ page }) => {
  321 |   await page.goto(BASIC);
  322 |   const grid = page.locator("[data-gsxui-calendar-grid]");
  323 |   await expect(grid).toHaveAttribute("aria-label", "January 2026");
  324 | 
  325 |   await page.click("[data-gsxui-calendar-next]");
  326 |   await expect(grid).toHaveAttribute("aria-label", "February 2026");
  327 | 
  328 |   await page.click("[data-gsxui-calendar-prev]");
  329 |   await expect(grid).toHaveAttribute("aria-label", "January 2026");
  330 | });
  331 | 
  332 | // Bounded pins fromYear=toYear=2026 — the narrowest possible navigation
  333 | // range, reachable in a handful of clicks. Task 4 review's Important
```