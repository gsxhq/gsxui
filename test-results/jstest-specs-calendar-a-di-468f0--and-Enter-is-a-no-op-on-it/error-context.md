# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> a disabled day keeps focus (aria-disabled, not native) and Enter is a no-op on it
- Location: jstest/specs/calendar.spec.ts:908:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/loaded", waiting until "load"

```

# Test source

```ts
  811  |   // month. Upstream's getFocusableDate ends with `min([navEnd,
  812  |   // focusableDate])`, which returns navEnd, and the final review found the
  813  |   // earlier clamp doing the transplant instead (it landed on 2026-12-15
  814  |   // here, which happened to look plausible only because 15 is a valid
  815  |   // December day — the "arrow keys at a bound" test below is where the same
  816  |   // bug jumped focus 24 days in the WRONG DIRECTION).
  817  |   expect(await focusedDate(page)).toBe("2026-12-31");
  818  | 
  819  |   // Already sitting at the bound — one more month-forward press must stay
  820  |   // a no-op, the same "attempt it anyway" discipline the mouse-path test
  821  |   // above uses for its own 12th click.
  822  |   await page.keyboard.press("PageDown");
  823  |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  824  |   await expect(yearSelect).toHaveValue("2026");
  825  |   expect(await focusedDate(page)).toBe("2026-12-31");
  826  | });
  827  | 
  828  | // Final review, Important 4: clampToNavBounds jumped focus forward instead of
  829  | // clamping. It landed on `Date.UTC(fromYear, 0, min(date.getUTCDate(), 31))`
  830  | // — the OUT-OF-RANGE TARGET's day-of-month — while its own comment claimed
  831  | // January 1st of fromYear. On this very page (fromYear=toYear=2026, January
  832  | // displayed): focus 2026-01-01, press ArrowUp, raw target 2025-12-25, clamped
  833  | // to 2026-01-25 — focus jumped 24 days FORWARD out of a backward move.
  834  | //
  835  | // The Shift+PageDown test above never caught it because at the upper bound
  836  | // its own day-of-month (15) happened to survive the transplant unchanged.
  837  | // Arrow keys at a bound are where the two behaviors visibly diverge, so this
  838  | // test presses one at each bound.
  839  | //
  840  | // One-line mutation that turns this red: restore the day-of-month transplant
  841  | // in clampToNavBounds (see the report for the captured failure).
  842  | test("arrow keys at a navigation bound clamp to the bound, not past it", async ({ page }) => {
  843  |   await page.goto(BOUNDED);
  844  |   const root = page.locator("[data-gsxui-calendar]");
  845  | 
  846  |   // Lower bound. A week backward from 2026-01-01 leaves fromYear entirely.
  847  |   await dayFor(page, "2026-01-01").focus();
  848  |   await page.keyboard.press("ArrowUp");
  849  |   expect(await focusedDate(page)).toBe("2026-01-01");
  850  |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-01");
  851  | 
  852  |   // Upper bound, the mirror image: a week forward from 2026-12-31 leaves
  853  |   // toYear. Reached through the public nav control, same as the mouse-path
  854  |   // bound test above.
  855  |   const next = page.locator("[data-gsxui-calendar-next]");
  856  |   for (let i = 0; i < 11; i++) await next.click();
  857  |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  858  | 
  859  |   await dayFor(page, "2026-12-31").focus();
  860  |   await page.keyboard.press("ArrowDown");
  861  |   expect(await focusedDate(page)).toBe("2026-12-31");
  862  |   await expect(root).toHaveAttribute("data-gsxui-calendar-month", "2026-12");
  863  | });
  864  | 
  865  | test("Enter and Space select the focused day", async ({ page }) => {
  866  |   await page.goto(BASIC);
  867  |   await dayFor(page, "2026-01-15").focus();
  868  |   await page.keyboard.press("Enter");
  869  |   await expect(cellFor(page, "2026-01-15")).toHaveAttribute("aria-selected", "true");
  870  | 
  871  |   await dayFor(page, "2026-01-16").focus();
  872  |   await page.keyboard.press(" ");
  873  |   await expect(cellFor(page, "2026-01-16")).toHaveAttribute("aria-selected", "true");
  874  | });
  875  | 
  876  | test("exactly one day is tabbable at any time", async ({ page }) => {
  877  |   await page.goto(BASIC);
  878  |   const count = () => page.$$eval('[data-gsxui-calendar-day][tabindex="0"]', (e) => e.length);
  879  |   expect(await count()).toBe(1);
  880  | 
  881  |   await dayFor(page, "2026-01-15").focus();
  882  |   await page.keyboard.press("ArrowRight");
  883  |   expect(await count()).toBe(1);
  884  | });
  885  | 
  886  | test("Tab enters the grid when its roving stop is disabled", async ({ page }) => {
  887  |   await page.goto(LOADED);
  888  |   await page.locator("[data-gsxui-calendar-next]").focus();
  889  |   await page.keyboard.press("Tab");
  890  | 
  891  |   expect(await focusedDate(page)).toBe("2026-01-01");
  892  |   await expect(dayFor(page, "2026-01-01")).toHaveJSProperty("disabled", false);
  893  |   await expect(dayFor(page, "2026-01-01")).toHaveAttribute("aria-disabled", "true");
  894  | });
  895  | 
  896  | // Source map §7.2/§8 finding 5, the first of the two behaviors the Task 6
  897  | // brief flags as easy to miss: a disabled day carries the NATIVE `disabled`
  898  | // attribute only while it is not the live-focused day; the moment it becomes
  899  | // the focused day it degrades to aria-disabled instead, so focus is never
  900  | // yanked back out of the grid. Loaded's own 2026-01-20 is disabled
  901  | // server-side (disabledDates, calendar/loaded.gsx) — arrowing onto it from
  902  | // the adjacent, non-disabled 2026-01-19 (a Monday; 2026-01-20's own
  903  | // Saturday-weekday rule and disabledDates entry don't touch it) is what
  904  | // actually puts the split under test, not fabricating disabled state via
  905  | // evaluate(). The round trip before focusing is the same "force one real
  906  | // repaint" discipline the pre-existing "a disabled day cannot be selected"
  907  | // test above already uses.
  908  | test("a disabled day keeps focus (aria-disabled, not native) and Enter is a no-op on it", async ({
  909  |   page,
  910  | }) => {
> 911  |   await page.goto(LOADED);
       |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  912  |   await page.click("[data-gsxui-calendar-next]");
  913  |   await page.click("[data-gsxui-calendar-prev]");
  914  |   await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
  915  |     "data-gsxui-calendar-month",
  916  |     "2026-01",
  917  |   );
  918  |   await listenForChanges(page);
  919  | 
  920  |   await dayFor(page, "2026-01-19").focus();
  921  |   await page.keyboard.press("ArrowRight");
  922  |   expect(await focusedDate(page)).toBe("2026-01-20");
  923  | 
  924  |   // Not .not.toBeDisabled() — Playwright's own accessibility-tree-backed
  925  |   // "disabled" state treats aria-disabled="true" as disabled too (correctly:
  926  |   // that is the entire point of the attribute), so that matcher can't tell
  927  |   // native disabled apart from the aria-disabled degradation this test
  928  |   // exists to prove. toHaveJSProperty reads the raw DOM `disabled` property
  929  |   // directly instead.
  930  |   const day = dayFor(page, "2026-01-20");
  931  |   await expect(day).toHaveJSProperty("disabled", false);
  932  |   await expect(day).toHaveAttribute("aria-disabled", "true");
  933  | 
  934  |   await page.keyboard.press("Enter");
  935  |   await expect(cellFor(page, "2026-01-20")).toHaveAttribute("aria-selected", "false");
  936  |   expect(await changes(page)).toHaveLength(0);
  937  | });
  938  | 
  939  | // The second of the two behaviors the Task 6 brief flags as easy to miss:
  940  | // upstream lands focus imperatively (a ref.current?.focus() call, source map
  941  | // §3/§7.3) — moving tabindex alone would leave the browser's own focus
  942  | // wherever it already was. basic.gsx is pinned to January 2026 (Basic's own
  943  | // DefaultMonth, calendar/basic.gsx); navigating to the client's real current
  944  | // month (whatever that is when this test runs) through the public prev/next
  945  | // controls only, per the Task 6 brief's own ambiguity resolution, rather
  946  | // than reaching for a private hook — is what actually exercises the client
  947  | // path, since basic.gsx's own server render can never contain the real
  948  | // "today" for a fixed, unrelated month.
  949  | test("today is marked from the client's date", async ({ page }) => {
  950  |   await page.goto(BASIC);
  951  | 
  952  |   const { today, delta } = await page.evaluate(() => {
  953  |     const now = new Date();
  954  |     const iso = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
  955  |     // Months from January 2026 to the client's current month.
  956  |     const months = (now.getFullYear() - 2026) * 12 + now.getMonth();
  957  |     return { today: iso, delta: months };
  958  |   });
  959  | 
  960  |   const button = delta >= 0 ? "[data-gsxui-calendar-next]" : "[data-gsxui-calendar-prev]";
  961  |   for (let i = 0; i < Math.abs(delta); i++) await page.click(button);
  962  | 
  963  |   const marked = await page.$$eval("[data-today]", (els) =>
  964  |     els.map((el) => el.getAttribute("data-date")),
  965  |   );
  966  |   expect(marked).toEqual([today]);
  967  | });
  968  | 
  969  | // The dropdown caption is the one place this port renders a real, visible
  970  | // <select> where upstream renders an invisible one under a styled label. Two
  971  | // things went wrong there and neither is visible to any attribute assertion:
  972  | // ui.NativeSelect's own border painted a second concentric rounded rect inside
  973  | // dropdown_root's, and its hardcoded h-8 overflowed the caption row. A third
  974  | // came from `relative` living on the padded root, so `nav`'s `top-0` resolved
  975  | // against the padding box and pinned the buttons above the row they label.
  976  | //
  977  | // Geometry is the only way to catch any of them, so assert it.
  978  | test("dropdown caption is one row, one border", async ({ page }) => {
  979  |   await page.goto("/x/calendar/dropdown");
  980  | 
  981  |   const geom = await page.evaluate(() => {
  982  |     const mid = (el: Element | null) =>
  983  |       el ? Math.round(el.getBoundingClientRect().top + el.getBoundingClientRect().height / 2) : null;
  984  |     const select = document.querySelector("[data-gsxui-calendar-month-select]")!;
  985  |     const wrapper = select.closest('[data-gsxui-slot-native-select-wrapper]')!;
  986  |     const row = wrapper.parentElement!;
  987  |     return {
  988  |       rowMid: mid(row),
  989  |       wrapperMid: mid(wrapper),
  990  |       prevMid: mid(document.querySelector("[data-gsxui-calendar-prev]")),
  991  |       nextMid: mid(document.querySelector("[data-gsxui-calendar-next]")),
  992  |       selectBorder: getComputedStyle(select).borderTopWidth,
  993  |       wrapperBorder: getComputedStyle(wrapper).borderTopWidth,
  994  |       selectHeight: Math.round(select.getBoundingClientRect().height),
  995  |       rowHeight: Math.round(row.getBoundingClientRect().height),
  996  |     };
  997  |   });
  998  | 
  999  |   // Everything in the caption shares one vertical centre line.
  1000 |   expect(geom.wrapperMid).toBe(geom.rowMid);
  1001 |   expect(geom.prevMid).toBe(geom.rowMid);
  1002 |   expect(geom.nextMid).toBe(geom.rowMid);
  1003 | 
  1004 |   // Exactly one border: dropdown_root's, on the wrapper.
  1005 |   expect(geom.wrapperBorder).toBe("1px");
  1006 |   expect(geom.selectBorder).toBe("0px");
  1007 | 
  1008 |   // The control fits the row rather than overflowing it.
  1009 |   expect(geom.selectHeight).toBeLessThanOrEqual(geom.rowHeight);
  1010 | });
  1011 | 
```