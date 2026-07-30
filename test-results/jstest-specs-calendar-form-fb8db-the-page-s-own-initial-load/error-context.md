# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> form reset restores a calendar added after the page's own initial load
- Location: jstest/specs/calendar.spec.ts:673:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/form", waiting until "load"

```

# Test source

```ts
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
> 676 |   await page.goto(FORM);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  716 | // them. dayFor is the button half (the one that's actually focusable), used
  717 | // throughout in place of the brief's bare selector.
  718 | test("arrows move by day and by week", async ({ page }) => {
  719 |   await page.goto(BASIC);
  720 |   await dayFor(page, "2026-01-15").click();
  721 |   await dayFor(page, "2026-01-15").focus();
  722 | 
  723 |   await page.keyboard.press("ArrowRight");
  724 |   expect(await focusedDate(page)).toBe("2026-01-16");
  725 | 
  726 |   await page.keyboard.press("ArrowLeft");
  727 |   expect(await focusedDate(page)).toBe("2026-01-15");
  728 | 
  729 |   await page.keyboard.press("ArrowDown");
  730 |   expect(await focusedDate(page)).toBe("2026-01-22");
  731 | 
  732 |   await page.keyboard.press("ArrowUp");
  733 |   expect(await focusedDate(page)).toBe("2026-01-15");
  734 | });
  735 | 
  736 | test("Home and End go to the bounds of the focused week", async ({ page }) => {
  737 |   await page.goto(BASIC);
  738 |   await dayFor(page, "2026-01-15").focus(); // a Thursday
  739 | 
  740 |   await page.keyboard.press("Home");
  741 |   expect(await focusedDate(page)).toBe("2026-01-11"); // Sunday
  742 | 
  743 |   await page.keyboard.press("End");
  744 |   expect(await focusedDate(page)).toBe("2026-01-17"); // Saturday
  745 | });
  746 | 
  747 | test("PageUp and PageDown move by month, with Shift by year", async ({ page }) => {
  748 |   await page.goto(BASIC);
  749 |   await dayFor(page, "2026-01-15").focus();
  750 | 
  751 |   await page.keyboard.press("PageDown");
  752 |   expect(await focusedDate(page)).toBe("2026-02-15");
  753 | 
  754 |   await page.keyboard.press("PageUp");
  755 |   expect(await focusedDate(page)).toBe("2026-01-15");
  756 | 
  757 |   await page.keyboard.press("Shift+PageDown");
  758 |   expect(await focusedDate(page)).toBe("2027-01-15");
  759 | });
  760 | 
  761 | // Task 6 review, Minor 5: the original version of this test moved
  762 | // ArrowRight from 2026-01-31 to 2026-02-01 — a date already present in
  763 | // BASIC's own January grid as trailing outside-day padding (the grid spans
  764 | // 2025-12-28..2026-02-07). Its focusedDate assertion had no teeth: a
  765 | // mutation that skips the goTo() navigation entirely still let .focus()
  766 | // find and land on that already-rendered padding-day button, so only the
  767 | // SECOND assertion (the month attribute) ever caught anything — disclosed
  768 | // honestly in the Task 6 report's own mutation-4 write-up. Retargeted at
  769 | // Shift+PageDown from 2026-01-31 to 2027-01-31 instead: that date is
  770 | // nowhere in the 42-cell window at the moment the key is pressed, so
  771 | // .focus() can only succeed AFTER a real navigation actually repaints a
  772 | // button carrying that date into existence — the first assertion now means
  773 | // what the test's name claims.
  774 | test("arrowing past the month edge navigates and keeps focus on the target", async ({ page }) => {
  775 |   await page.goto(BASIC);
  776 |   await dayFor(page, "2026-01-31").focus();
```