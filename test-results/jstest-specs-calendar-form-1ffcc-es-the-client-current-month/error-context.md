# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> form reset restores the client-current month
- Location: jstest/specs/calendar.spec.ts:614:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/form", waiting until "load"

```

# Test source

```ts
  515 | });
  516 | 
  517 | test("range previews on hover while only the start is set", async ({ page }) => {
  518 |   await page.goto(RANGE);
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
> 615 |   await page.goto(FORM);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  616 |   const root = page.locator("[data-gsxui-calendar]");
  617 |   const currentMonth = await page.evaluate(() => {
  618 |     const now = new Date();
  619 |     return `${String(now.getFullYear()).padStart(4, "0")}-${String(now.getMonth() + 1).padStart(2, "0")}`;
  620 |   });
  621 | 
  622 |   await page.click("[data-gsxui-calendar-next]");
  623 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-02");
  624 | 
  625 |   await page.locator('button[type="reset"]').click();
  626 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", currentMonth);
  627 | });
  628 | 
  629 | test("one form reset restores both calendar and combobox state", async ({ page }) => {
  630 |   await page.goto(FORM);
  631 |   const calendarInput = page.locator('input[name="date"]');
  632 |   const comboInput = page.locator("[data-gsxui-combobox-input]");
  633 |   const comboBridge = page.locator("[data-gsxui-combobox-bridge]");
  634 |   const comboItem = page.locator('[data-gsxui-combobox-item][data-value="sveltekit"]');
  635 | 
  636 |   await dayFor(page, "2026-01-15").click();
  637 |   await comboInput.click();
  638 |   await comboItem.click();
  639 | 
  640 |   await expect(calendarInput).toHaveValue("2026-01-15");
  641 |   await expect(comboInput).toHaveValue("SvelteKit");
  642 |   await expect(comboBridge).toHaveValue("sveltekit");
  643 |   await expect(comboItem).toHaveAttribute("aria-selected", "true");
  644 | 
  645 |   await page.locator('button[type="reset"]').click();
  646 | 
  647 |   await expect(calendarInput).toHaveValue("");
  648 |   await expect(page.locator("[data-gsxui-calendar]")).not.toHaveAttribute(
  649 |     "data-gsxui-calendar-selected",
  650 |     /.*/,
  651 |   );
  652 |   expect(await page.$$eval('[data-gsxui-calendar] [aria-selected="true"]', (els) => els.length)).toBe(
  653 |     0,
  654 |   );
  655 |   await expect(comboInput).toHaveValue("");
  656 |   await expect(comboBridge).toHaveValue("");
  657 |   await expect(comboItem).toHaveAttribute("aria-selected", "false");
  658 | });
  659 | 
  660 | // Task 5 review, Important 3: captureDefaults used to run exactly once, at
  661 | // module load, over whatever calendar roots already existed — a root an
  662 | // HTMX swap (or any other post-load insertion) adds afterward had no
  663 | // snapshot at all, so resetting ITS form silently left the user's
  664 | // selection and hidden input untouched instead of reverting them.
  665 | // ui/gsxui.js's own header promises elements added later just work; this
  666 | // is the one place calendar.js broke that promise. Simulated here by
  667 | // cloning the FORM example's own pristine calendar root into a brand-new
  668 | // <form> well after calendar.js's module-load sweep has already run — the
  669 | // clone is a distinct node instance the WeakMap has never seen. One-line
  670 | // mutation that reproduces the bug: delete the `captureDefaults(root);`
  671 | // call at the top of the day-click handler (see the report for the
  672 | // captured failure).
  673 | test("form reset restores a calendar added after the page's own initial load", async ({
  674 |   page,
  675 | }) => {
  676 |   await page.goto(FORM);
  677 | 
  678 |   await page.evaluate(() => {
  679 |     const original = document.querySelector("[data-gsxui-calendar]")!;
  680 |     const clone = original.cloneNode(true) as HTMLElement;
  681 |     const form = document.createElement("form");
  682 |     form.appendChild(clone);
  683 |     document.body.appendChild(form);
  684 |   });
  685 | 
  686 |   const added = page.locator("[data-gsxui-calendar]").nth(1);
  687 |   const addedDay = added.locator('[data-gsxui-calendar-day][data-date="2026-01-15"]');
  688 |   const addedInput = added.locator('input[name="date"]');
  689 | 
  690 |   await addedDay.click();
  691 |   await expect(added.locator('td[data-date="2026-01-15"]')).toHaveAttribute(
  692 |     "aria-selected",
  693 |     "true",
  694 |   );
  695 |   await expect(addedInput).toHaveValue("2026-01-15");
  696 | 
  697 |   await added.evaluate((el) => (el.closest("form") as HTMLFormElement).reset());
  698 |   await expect(addedInput).toHaveValue("");
  699 |   await expect(added).not.toHaveAttribute("data-gsxui-calendar-selected", /.*/);
  700 |   expect(await added.locator('[aria-selected="true"]').count()).toBe(0);
  701 | });
  702 | 
  703 | // --- Task 6: keyboard grid and today reconciliation -------------------------
  704 | //
  705 | // document.activeElement is the one thing no locator search can stand in
  706 | // for: several tests below exist specifically to discover WHICH day ended
  707 | // up holding focus after a key press, not to assert a guessed one already
  708 | // has it.
  709 | async function focusedDate(page: import("@playwright/test").Page) {
  710 |   return page.evaluate(() => document.activeElement?.getAttribute("data-date") ?? null);
  711 | }
  712 | 
  713 | // The brief's own selectors here are bare `[data-date="…"]`, which — same
  714 | // strict-mode ambiguity as every Task 5 test above — resolve to both a
  715 | // `<td>` and a `<button>` and throw the moment Playwright tries to act on
```