# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> dropdown caption is one row, one border
- Location: jstest/specs/calendar.spec.ts:978:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/dropdown", waiting until "load"

```

# Test source

```ts
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
  911  |   await page.goto(LOADED);
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
> 979  |   await page.goto("/x/calendar/dropdown");
       |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  1012 | // ui.NativeSelect ships py-1 for its h-8 standalone form-control shape. At the
  1013 | // caption's h-6 that leaves a 16px content box for a 20px line — the line
  1014 | // cannot fit, so no browser can centre it and the text renders visibly low.
  1015 | // The invariant is that the line fits and the padding is symmetric; the
  1016 | // browser centres it from there.
  1017 | test("dropdown caption text can be vertically centred", async ({ page }) => {
  1018 |   await page.goto("/x/calendar/dropdown");
  1019 | 
  1020 |   const m = await page.evaluate(() => {
  1021 |     const sel = document.querySelector("[data-gsxui-calendar-month-select]")!;
  1022 |     const cs = getComputedStyle(sel);
  1023 |     const px = (v: string) => parseFloat(v);
  1024 |     return {
  1025 |       height: sel.getBoundingClientRect().height,
  1026 |       padTop: px(cs.paddingTop),
  1027 |       padBottom: px(cs.paddingBottom),
  1028 |       lineHeight: px(cs.lineHeight),
  1029 |       borderTop: px(cs.borderTopWidth),
  1030 |       borderBottom: px(cs.borderBottomWidth),
  1031 |     };
  1032 |   });
  1033 | 
  1034 |   expect(m.padTop).toBe(m.padBottom);
  1035 |   const contentBox = m.height - m.padTop - m.padBottom - m.borderTop - m.borderBottom;
  1036 |   expect(contentBox).toBeGreaterThanOrEqual(m.lineHeight);
  1037 | });
  1038 | 
  1039 | // The chevron is absolutely positioned inside the select, so its clearance is
  1040 | // reserved with right padding rather than by flex flow the way upstream's
  1041 | // caption label does it. Over-reserve and the controls crowd the nav buttons
  1042 | // until they nearly touch — 2px of gap at pr-6, against the 6px that separates
  1043 | // the two controls from each other. Assert the caption breathes evenly.
  1044 | test("dropdown caption does not crowd the nav buttons", async ({ page }) => {
  1045 |   await page.goto("/x/calendar/dropdown");
  1046 | 
  1047 |   const gaps = await page.evaluate(() => {
  1048 |     const r = (el: Element) => el.getBoundingClientRect();
  1049 |     const wraps = [
  1050 |       ...document.querySelectorAll('[data-gsxui-slot-native-select-wrapper]'),
  1051 |     ];
  1052 |     const prev = document.querySelector("[data-gsxui-calendar-prev]")!;
  1053 |     const next = document.querySelector("[data-gsxui-calendar-next]")!;
  1054 |     return {
  1055 |       beforeFirst: Math.round(r(wraps[0]).left - r(prev).right),
  1056 |       between: Math.round(r(wraps[1]).left - r(wraps[0]).right),
  1057 |       afterLast: Math.round(r(next).left - r(wraps[wraps.length - 1]).right),
  1058 |     };
  1059 |   });
  1060 | 
  1061 |   // gap-1.5 between the controls is the reference spacing; the outer gaps must
  1062 |   // be at least as generous, and equal to each other.
  1063 |   expect(gaps.beforeFirst).toBeGreaterThanOrEqual(gaps.between);
  1064 |   expect(gaps.afterLast).toBeGreaterThanOrEqual(gaps.between);
  1065 |   expect(gaps.beforeFirst).toBe(gaps.afterLast);
  1066 | });
  1067 | 
```