# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> the grid's accessible name follows the displayed month
- Location: jstest/specs/calendar.spec.ts:320:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/basic", waiting until "load"

```

# Test source

```ts
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
> 321 |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  334 | // finding 1: calendar.js never recomputed the nav bounds after a
  335 | // client-side navigation, so 11 clicks on next from 2026-01 reached
  336 | // 2026-12 (where Go WOULD mark next aria-disabled), and a 12th click landed
  337 | // on 2027-01 — a month outside the declared range that the server can
  338 | // never render for this component, and (in captionLayout="dropdown") a
  339 | // year <select> holding a value with no matching <option>.
  340 | test("next stops at the declared navigation bound instead of crossing it", async ({ page }) => {
  341 |   await page.goto(BOUNDED);
  342 |   const root = page.locator("[data-gsxui-calendar]");
  343 |   const next = page.locator("[data-gsxui-calendar-next]");
  344 |   const yearSelect = page.locator("[data-gsxui-calendar-year-select]");
  345 | 
  346 |   for (let i = 0; i < 11; i++) await next.click();
  347 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  348 |   await expect(next).toHaveAttribute("aria-disabled", "true");
  349 |   await expect(yearSelect).toHaveValue("2026");
  350 | 
  351 |   // The 12th click must be a no-op: no server render exists for 2027-01,
  352 |   // and the year <select>'s own options never include it either (they run
  353 |   // navFromYear..navToYear, both 2026 here). force: true — Playwright's own
  354 |   // actionability check treats aria-disabled="true" as non-clickable and
  355 |   // would otherwise wait out the full test timeout for a state that, by
  356 |   // design, never changes; the whole point of this assertion is to attempt
  357 |   // the click anyway and confirm calendar.js's own handler guard (which
  358 |   // checks aria-disabled itself, not Playwright's actionability rules) is
  359 |   // what makes it a no-op.
  360 |   await next.click({ force: true });
  361 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  362 |   await expect(yearSelect).toHaveValue("2026");
  363 | });
  364 | 
  365 | // The lower bound is reachable from first paint (fromYear=toYear=2026, and
  366 | // Bounded's initial month is already 2026-01) — the server itself already
  367 | // marks prev aria-disabled before any JS runs; this checks that state
  368 | // survives being clicked.
  369 | test("prev is disabled at the lower bound from the start and stays a no-op", async ({
  370 |   page,
  371 | }) => {
  372 |   await page.goto(BOUNDED);
  373 |   const root = page.locator("[data-gsxui-calendar]");
  374 |   const prev = page.locator("[data-gsxui-calendar-prev]");
  375 | 
  376 |   await expect(prev).toHaveAttribute("aria-disabled", "true");
  377 |   // force: true — see the equivalent comment on the "next stops at..." test
  378 |   // above: Playwright's own actionability check would otherwise wait out
  379 |   // the whole test timeout for aria-disabled to clear, which by design
  380 |   // never happens.
  381 |   await prev.click({ force: true });
  382 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-01");
  383 | });
  384 | 
  385 | // --- Task 5: selection behavior ---------------------------------------------
  386 | //
  387 | // Both the <td role="gridcell"> and its <button data-gsxui-calendar-day>
  388 | // carry their own data-date (calendar.gsx and calendar.js's own repaint
  389 | // agree on that split — see gridCells() above), so a bare
  390 | // `[data-date="…"]` locator resolves to TWO elements and trips Playwright's
  391 | // strict mode the moment an assertion runs against it. cellFor/dayFor scope
  392 | // to whichever one actually owns the attribute under test: aria-selected/
  393 | // data-selected live on the cell, everything else (click target,
  394 | // data-range-*, data-selected-single, disabled) lives on the button.
  395 | function cellFor(page: import("@playwright/test").Page, iso: string) {
  396 |   return page.locator(`td[data-date="${iso}"]`);
  397 | }
  398 | function dayFor(page: import("@playwright/test").Page, iso: string) {
  399 |   return page.locator(`[data-gsxui-calendar-day][data-date="${iso}"]`);
  400 | }
  401 | 
  402 | // listenForChanges wires up the same window.__changes collection every test
  403 | // below reads from — a bare gsxui:change listener on the document, since the
  404 | // event bubbles (ui/gsxui.js's own emit()) all the way up from the root.
  405 | async function listenForChanges(page: import("@playwright/test").Page) {
  406 |   await page.evaluate(() => {
  407 |     (window as any).__changes = [];
  408 |     document.addEventListener("gsxui:change", (e) =>
  409 |       (window as any).__changes.push((e as CustomEvent).detail),
  410 |     );
  411 |   });
  412 | }
  413 | 
  414 | async function changes(page: import("@playwright/test").Page) {
  415 |   return page.evaluate(() => (window as any).__changes);
  416 | }
  417 | 
  418 | test("single mode replaces the selection and emits once", async ({ page }) => {
  419 |   await page.goto(BASIC);
  420 |   await listenForChanges(page);
  421 | 
```