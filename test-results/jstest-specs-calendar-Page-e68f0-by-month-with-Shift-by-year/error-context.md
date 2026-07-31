# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> PageUp and PageDown move by month, with Shift by year
- Location: jstest/specs/calendar.spec.ts:747:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/basic", waiting until "load"

```

# Test source

```ts
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
> 748 |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  777 | 
  778 |   await page.keyboard.press("Shift+PageDown");
  779 |   expect(await focusedDate(page)).toBe("2027-01-31");
  780 |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  781 |     "data-gsxui-calendar-month",
  782 |     "2027-01",
  783 |   );
  784 | });
  785 | 
  786 | // Task 6 review, Important: the mouse-path twin of this ("next stops at the
  787 | // declared navigation bound instead of crossing it", above) already
  788 | // enforces fromYear/toYear against prev/next clicks; keyboard month/year
  789 | // moves (PageUp/PageDown, Shift+Page) had no equivalent clamp, so a single
  790 | // Shift+PageDown from Bounded's own January 2026 (fromYear=toYear=2026)
  791 | // could reach 2027-01 in one keystroke — a month the server can never
  792 | // render for this component, and (captionLayout="dropdown" here too) a
  793 | // year <select> left holding a value with no matching <option>, verbatim
  794 | // the defect the mouse-path test above exists to prevent.
  795 | test("keyboard month/year moves stop at the declared navigation bound instead of crossing it", async ({
  796 |   page,
  797 | }) => {
  798 |   await page.goto(BOUNDED);
  799 |   const root = page.locator("[data-gsxui-calendar]");
  800 |   const next = page.locator("[data-gsxui-calendar-next]");
  801 |   const yearSelect = page.locator("[data-gsxui-calendar-year-select]");
  802 | 
  803 |   await dayFor(page, "2026-01-15").focus();
  804 |   await page.keyboard.press("Shift+PageDown");
  805 | 
  806 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  807 |   await expect(next).toHaveAttribute("aria-disabled", "true");
  808 |   await expect(yearSelect).toHaveValue("2026");
  809 |   // Clamped to the BOUND ITSELF — the last day of toYear's December —
  810 |   // not to the out-of-range target's day-of-month transplanted into that
  811 |   // month. Upstream's getFocusableDate ends with `min([navEnd,
  812 |   // focusableDate])`, which returns navEnd, and the final review found the
  813 |   // earlier clamp doing the transplant instead (it landed on 2026-12-15
  814 |   // here, which happened to look plausible only because 15 is a valid
  815 |   // December day — the "arrow keys at a bound" test below is where the same
  816 |   // bug jumped focus 24 days in the WRONG DIRECTION).
  817 |   expect(await focusedDate(page)).toBe("2026-12-31");
  818 | 
  819 |   // Already sitting at the bound — one more month-forward press must stay
  820 |   // a no-op, the same "attempt it anyway" discipline the mouse-path test
  821 |   // above uses for its own 12th click.
  822 |   await page.keyboard.press("PageDown");
  823 |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  824 |   await expect(yearSelect).toHaveValue("2026");
  825 |   expect(await focusedDate(page)).toBe("2026-12-31");
  826 | });
  827 | 
  828 | // Final review, Important 4: clampToNavBounds jumped focus forward instead of
  829 | // clamping. It landed on `Date.UTC(fromYear, 0, min(date.getUTCDate(), 31))`
  830 | // — the OUT-OF-RANGE TARGET's day-of-month — while its own comment claimed
  831 | // January 1st of fromYear. On this very page (fromYear=toYear=2026, January
  832 | // displayed): focus 2026-01-01, press ArrowUp, raw target 2025-12-25, clamped
  833 | // to 2026-01-25 — focus jumped 24 days FORWARD out of a backward move.
  834 | //
  835 | // The Shift+PageDown test above never caught it because at the upper bound
  836 | // its own day-of-month (15) happened to survive the transplant unchanged.
  837 | // Arrow keys at a bound are where the two behaviors visibly diverge, so this
  838 | // test presses one at each bound.
  839 | //
  840 | // One-line mutation that turns this red: restore the day-of-month transplant
  841 | // in clampToNavBounds (see the report for the captured failure).
  842 | test("arrow keys at a navigation bound clamp to the bound, not past it", async ({ page }) => {
  843 |   await page.goto(BOUNDED);
  844 |   const root = page.locator("[data-gsxui-calendar]");
  845 | 
  846 |   // Lower bound. A week backward from 2026-01-01 leaves fromYear entirely.
  847 |   await dayFor(page, "2026-01-01").focus();
  848 |   await page.keyboard.press("ArrowUp");
```