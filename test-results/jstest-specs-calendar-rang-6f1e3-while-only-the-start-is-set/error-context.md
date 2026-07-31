# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> range previews on hover while only the start is set
- Location: jstest/specs/calendar.spec.ts:517:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/range", waiting until "load"

```

# Test source

```ts
  418 | test("single mode replaces the selection and emits once", async ({ page }) => {
  419 |   await page.goto(BASIC);
  420 |   await listenForChanges(page);
  421 | 
  422 |   await dayFor(page, "2026-01-15").click();
  423 |   await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "true");
  424 | 
  425 |   await dayFor(page, "2026-01-20").click();
  426 |   await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "false");
  427 |   await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "true");
  428 | 
  429 |   const events = await changes(page);
  430 |   expect(events).toHaveLength(2);
  431 |   expect(events[1]).toEqual({ mode: "single", selected: ["2026-01-20"] });
  432 | 
  433 |   // Re-clicking the now-selected day clears it outright (single mode's own
  434 |   // "clicking the already-selected day clears the selection" rule) — a
  435 |   // third click the brief's own test doesn't cover, added here because
  436 |   // without it a commitSingle that always REPLACED (never cleared) would
  437 |   // still satisfy every assertion above.
  438 |   await dayFor(page, "2026-01-20").click();
  439 |   await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  440 |   expect(await changes(page)).toHaveLength(3);
  441 |   expect((await changes(page))[2]).toEqual({ mode: "single", selected: [] });
  442 | });
  443 | 
  444 | test("multiple mode toggles days", async ({ page }) => {
  445 |   await page.goto(MULTIPLE);
  446 |   await listenForChanges(page);
  447 | 
  448 |   await dayFor(page, "2026-01-05").click();
  449 |   await dayFor(page, "2026-01-09").click();
  450 |   expect(await page.$$eval('[aria-selected="true"]', (e) => e.length)).toBe(2);
  451 | 
  452 |   await dayFor(page, "2026-01-05").click();
  453 |   expect(await page.$$eval('[aria-selected="true"]', (e) => e.length)).toBe(1);
  454 |   await expect(cellFor(page, "2026-01-09")).toHaveAttribute("aria-selected", "true");
  455 | 
  456 |   // The toggle-off left the OTHER day selected, sorted, not just "count
  457 |   // dropped by one" — a commitMultiple that removed the wrong entry (e.g.
  458 |   // always dropping the first list element) would still pass the length
  459 |   // assertions above.
  460 |   const events = await changes(page);
  461 |   expect(events).toHaveLength(3);
  462 |   expect(events[2]).toEqual({ mode: "multiple", selected: ["2026-01-09"] });
  463 | });
  464 | 
  465 | test("multiple mode submits every selected date as a repeated form value", async ({ page }) => {
  466 |   await page.goto(MULTIPLE);
  467 |   const form = page.locator("form");
  468 | 
  469 |   expect(
  470 |     await form.evaluate((el) => new FormData(el as HTMLFormElement).getAll("dates")),
  471 |   ).toEqual([""]);
  472 | 
  473 |   await dayFor(page, "2026-01-05").click();
  474 |   await dayFor(page, "2026-01-09").click();
  475 | 
  476 |   expect(
  477 |     await form.evaluate((el) => new FormData(el as HTMLFormElement).getAll("dates")),
  478 |   ).toEqual(["2026-01-05", "2026-01-09"]);
  479 | });
  480 | 
  481 | test("range takes two clicks and swaps when the second precedes the first", async ({ page }) => {
  482 |   await page.goto(RANGE);
  483 |   const root = page.locator("[data-gsxui-calendar]");
  484 |   // range.gsx renders name="stay" specifically so this test also exercises
  485 |   // syncHiddenInputs's range branch — Task 5 review, Minor 1 verification:
  486 |   // the "-to" input has to be told apart from the "from" input by its own
  487 |   // data-gsxui-calendar-hidden-to marker, not a name-suffix guess.
  488 |   const from = page.locator('input[name="stay"]');
  489 |   const to = page.locator('input[name="stay-to"]');
  490 |   await expect(from).toHaveValue("");
  491 |   await expect(to).toHaveValue("");
  492 | 
  493 |   await dayFor(page, "2026-01-20").click();
  494 |   await expect(root).toHaveAttribute("data-gsxui-calendar-from", "2026-01-20");
  495 |   await expect(root).not.toHaveAttribute("data-gsxui-calendar-to", /.*/);
  496 |   await expect(from).toHaveValue("2026-01-20");
  497 |   await expect(to).toHaveValue("");
  498 | 
  499 |   await dayFor(page, "2026-01-15").click();
  500 |   await expect(root).toHaveAttribute("data-gsxui-calendar-from", "2026-01-15");
  501 |   await expect(root).toHaveAttribute("data-gsxui-calendar-to", "2026-01-20");
  502 |   await expect(dayFor(page, "2026-01-17")).toHaveAttribute("data-range-middle", "true");
  503 |   // The two endpoints themselves are NOT the middle — rangeStart/rangeEnd
  504 |   // own the boundary days instead; a rangeMiddle computed as a closed
  505 |   // interval (>= / <=, not > / <) would also mark these true.
  506 |   await expect(dayFor(page, "2026-01-15")).toHaveAttribute("data-range-middle", "false");
  507 |   await expect(dayFor(page, "2026-01-20")).toHaveAttribute("data-range-middle", "false");
  508 |   await expect(dayFor(page, "2026-01-15")).toHaveAttribute("data-range-start", "true");
  509 |   await expect(dayFor(page, "2026-01-20")).toHaveAttribute("data-range-end", "true");
  510 |   // Swapped: the FROM input carries the earlier date even though it was
  511 |   // clicked SECOND — proof the swap flows into the hidden inputs too, not
  512 |   // just the root's own data-* attributes.
  513 |   await expect(from).toHaveValue("2026-01-15");
  514 |   await expect(to).toHaveValue("2026-01-20");
  515 | });
  516 | 
  517 | test("range previews on hover while only the start is set", async ({ page }) => {
> 518 |   await page.goto(RANGE);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  519 | 
  520 |   await dayFor(page, "2026-01-10").click();
  521 |   await dayFor(page, "2026-01-14").hover();
  522 |   await expect(dayFor(page, "2026-01-12")).toHaveAttribute("data-range-middle", "true");
  523 | 
  524 |   // Hovering elsewhere moves the preview rather than accumulating it.
  525 |   await dayFor(page, "2026-01-11").hover();
  526 |   await expect(dayFor(page, "2026-01-12")).toHaveAttribute("data-range-middle", "false");
  527 | 
  528 |   // Leaving the whole calendar clears the preview entirely — re-hovering
  529 |   // 2026-01-14 after a mouseleave should reproduce the same preview as the
  530 |   // very first hover, not leave 2026-01-12 stuck in whatever state the
  531 |   // in-between hover left it in.
  532 |   await dayFor(page, "2026-01-14").hover();
  533 |   await page.mouse.move(0, 0);
  534 |   await expect(dayFor(page, "2026-01-12")).toHaveAttribute("data-range-middle", "false");
  535 | 
  536 |   // Hovering never emits gsxui:change — only a committed click does.
  537 |   await listenForChanges(page);
  538 |   await dayFor(page, "2026-01-14").hover();
  539 |   await dayFor(page, "2026-01-11").hover();
  540 |   expect(await changes(page)).toHaveLength(0);
  541 | });
  542 | 
  543 | // Ambiguity resolution (per the task-5 brief's own note): rather than
  544 | // fabricating a disabled day on BASIC via evaluate() and clicking it with
  545 | // force:true, this exercises a day LOADED already disables server-side
  546 | // (2026-01-20 is in Loaded's own disabledDates — see calendar/loaded.gsx) —
  547 | // a test that fabricates the state it then checks proves less than one that
  548 | // exercises the real render.
  549 | //
  550 | // The next/prev round trip before clicking is load-bearing, not decoration:
  551 | // on a bare page load, 2026-01-20's disabled <button> attribute comes
  552 | // straight from calendar.gsx's own server render, so a real browser refuses
  553 | // to dispatch "click" on it at all regardless of what ui/calendar.js does —
  554 | // a click attempted against that untouched server HTML can't tell a correct
  555 | // repaint() apart from one that forgot `button.disabled = disabled`, since
  556 | // repaint() never even runs. Forcing one live repaint (Task 4's own
  557 | // "exercise real JS, not just the server's initial render" discipline) means
  558 | // the disabled state under test is JS-computed, so the assertion actually
  559 | // depends on ui/calendar.js getting it right.
  560 | test("a disabled day cannot be selected", async ({ page }) => {
  561 |   await page.goto(LOADED);
  562 |   await page.click("[data-gsxui-calendar-next]");
  563 |   await page.click("[data-gsxui-calendar-prev]");
  564 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  565 |     "data-gsxui-calendar-month",
  566 |     "2026-01",
  567 |   );
  568 |   await listenForChanges(page);
  569 | 
  570 |   await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  571 |   await expect(dayFor(page, "2026-01-20")).toBeDisabled();
  572 | 
  573 |   // force: true — Playwright's own actionability check refuses to click a
  574 |   // genuinely disabled button; the point of this test is that even an
  575 |   // attempted click changes nothing, not that Playwright politely declines.
  576 |   await dayFor(page, "2026-01-20").click({ force: true });
  577 | 
  578 |   await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  579 |   expect(await changes(page)).toHaveLength(0);
  580 | });
  581 | 
  582 | // calendar/form.gsx (Task 5 review, Important 1) renders name="date" with
  583 | // NOTHING preselected — real corpus coverage for the Critical fix: the
  584 | // hidden input exists (empty) from first paint, and clicking a day has to
  585 | // actually flow into it. Asserting the input's own value (not just
  586 | // aria-selected) is Important 2 — the test used to be named "…and the
  587 | // hidden input" while never once looking at one, and with the input
  588 | // existing conditionally on a selection (the pre-fix behavior), the
  589 | // BASIC-based version of this test never had one to check in the first
  590 | // place. One-line mutation that turns this from a real assertion back into
  591 | // exactly that no-op: replace syncHiddenInputs's body with a bare `return;`
  592 | // — every other assertion in this test still passes (see the report for
  593 | // the captured failure).
  594 | test("form reset clears the selection and the hidden input", async ({ page }) => {
  595 |   await page.goto(FORM);
  596 |   const input = page.locator('input[name="date"]');
  597 |   await expect(input).toHaveValue("");
  598 | 
  599 |   await dayFor(page, "2026-01-15").click();
  600 |   await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "true");
  601 |   await expect(input).toHaveValue("2026-01-15");
  602 | 
  603 |   await page.locator('button[type="reset"]').click();
  604 |   await expect(input).toHaveValue("");
  605 |   await expect(page.locator("[data-gsxui-calendar]")).not.toHaveAttribute(
  606 |     "data-gsxui-calendar-selected",
  607 |     /.*/,
  608 |   );
  609 |   expect(await page.$$eval('[data-gsxui-calendar] [aria-selected="true"]', (els) => els.length)).toBe(
  610 |     0,
  611 |   );
  612 | });
  613 | 
  614 | test("form reset restores the client-current month", async ({ page }) => {
  615 |   await page.goto(FORM);
  616 |   const root = page.locator("[data-gsxui-calendar]");
  617 |   const currentMonth = await page.evaluate(() => {
  618 |     const now = new Date();
```